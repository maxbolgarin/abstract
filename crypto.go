// Copy paste from https://github.com/gtank/cryptopasta/tree/master
package abstract

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
)

// NewEncryptionKey generates a random 256-bit key for Encrypt() and
// Decrypt(). It panics if the source of randomness fails.
func NewEncryptionKey() *[32]byte {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return &key
}

// EncryptAES encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func EncryptAES(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
	if plaintext == nil {
		return nil, errors.New("plaintext is nil")
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// DecryptAES decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func DecryptAES(ciphertext []byte, key *[32]byte) (plaintext []byte, err error) {
	if ciphertext == nil {
		return nil, errors.New("ciphertext is nil")
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

// HashHMAC generates a hash of data using HMAC-SHA-512/256. The tag is intended to
// be a natural-language string describing the purpose of the hash, such as
// "hash file for lookup key" or "master secret to client secret".  It serves
// as an HMAC "key" and ensures that different purposes will have different
// hash output. This function is NOT suitable for hashing passwords.
func HashHMAC(tag string, data []byte) []byte {
	if len(data) == 0 {
		return nil
	}

	h := hmac.New(sha512.New512_256, []byte(tag))
	h.Write(data)
	return h.Sum(nil)
}

// DecodePublicKey decodes a PEM-encoded ECDSA public key.
func DecodePublicKey(encodedKey []byte) (*ecdsa.PublicKey, error) {
	if len(encodedKey) == 0 {
		return nil, errors.New("encoded key is empty")
	}

	block, _ := pem.Decode(encodedKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("marshal: could not decode PEM block or not a PUBLIC KEY")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("marshal: data was not an ECDSA public key")
	}

	return ecdsaPub, nil
}

// EncodePublicKey encodes an ECDSA public key to PEM format.
func EncodePublicKey(key *ecdsa.PublicKey) ([]byte, error) {
	if key == nil {
		return nil, errors.New("key is nil")
	}

	derBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	}

	return pem.EncodeToMemory(block), nil
}

// DecodePrivateKey decodes a PEM-encoded ECDSA private key.
func DecodePrivateKey(encodedKey []byte) (*ecdsa.PrivateKey, error) {
	if len(encodedKey) == 0 {
		return nil, errors.New("encoded key is empty")
	}

	var skippedTypes []string
	var block *pem.Block

	for {
		block, encodedKey = pem.Decode(encodedKey)

		if block == nil {
			return nil, fmt.Errorf("failed to find EC PRIVATE KEY in PEM data after skipping types %v", skippedTypes)
		}

		if block.Type == "EC PRIVATE KEY" {
			break
		} else {
			skippedTypes = append(skippedTypes, block.Type)
			continue
		}
	}

	privKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

// EncodePrivateKey encodes an ECDSA private key to PEM format.
func EncodePrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	if key == nil {
		return nil, errors.New("key is nil")
	}

	derKey, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, err
	}

	keyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: derKey,
	}

	return pem.EncodeToMemory(keyBlock), nil
}

// Encodes an ECDSA signature according to
// https://tools.ietf.org/html/rfc7515#appendix-A.3.1
func EncodeSignatureJWT(sig []byte) string {
	if len(sig) == 0 {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(sig)
}

// Decodes an ECDSA signature according to
// https://tools.ietf.org/html/rfc7515#appendix-A.3.1
func DecodeSignatureJWT(b64sig string) ([]byte, error) {
	if b64sig == "" {
		return nil, errors.New("empty signature")
	}
	return base64.RawURLEncoding.DecodeString(b64sig)
}

// NewHMACKey generates a random 256-bit secret key for HMAC use.
// Because key generation is critical, it panics if the source of randomness fails.
func NewHMACKey() *[32]byte {
	key := &[32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return key
}

// GenerateHMAC produces a symmetric signature using a shared secret key.
func GenerateHMAC(data []byte, key *[32]byte) []byte {
	if len(data) == 0 || key == nil {
		return nil
	}

	h := hmac.New(sha512.New512_256, key[:])
	h.Write(data)
	return h.Sum(nil)
}

// CheckHMAC securely checks the supplied MAC against a message using the shared secret key.
func CheckHMAC(data, suppliedMAC []byte, key *[32]byte) bool {
	if len(data) == 0 || len(suppliedMAC) == 0 || key == nil {
		return false
	}

	expectedMAC := GenerateHMAC(data, key)
	return subtle.ConstantTimeCompare(expectedMAC, suppliedMAC) == 1
}

// NewSigningKey generates a random P-256 ECDSA private key.
func NewSigningKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return key, err
}

// SignData signs arbitrary data using ECDSA.
func SignData(data []byte, privkey *ecdsa.PrivateKey) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	if privkey == nil {
		return nil, errors.New("private key is nil")
	}

	// hash message
	digest := sha256.Sum256(data)

	// sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privkey, digest[:])
	if err != nil {
		return nil, err
	}

	// ensure s is in lower half of curve order
	// this protects against signature malleability
	halfOrder := new(big.Int).Rsh(privkey.Curve.Params().N, 1)
	if s.Cmp(halfOrder) > 0 {
		s.Sub(privkey.Curve.Params().N, s)
	}

	// encode the signature {R, S}
	// big.Int.Bytes() will need padding in the case of leading zero bytes
	params := privkey.Curve.Params()
	curveOrderByteSize := params.P.BitLen() / 8
	rBytes, sBytes := r.Bytes(), s.Bytes()
	signature := make([]byte, curveOrderByteSize*2)
	copy(signature[curveOrderByteSize-len(rBytes):], rBytes)
	copy(signature[curveOrderByteSize*2-len(sBytes):], sBytes)

	return signature, nil
}

// VerifySign checks a raw ECDSA signature.
// Returns true if it's valid and false if not.
func VerifySign(data, signature []byte, pubkey *ecdsa.PublicKey) bool {
	if len(data) == 0 || len(signature) == 0 || pubkey == nil {
		return false
	}

	// hash message
	digest := sha256.Sum256(data)

	curveOrderByteSize := pubkey.Curve.Params().P.BitLen() / 8

	if len(signature) < curveOrderByteSize*2 {
		return false
	}

	r, s := new(big.Int), new(big.Int)
	r.SetBytes(signature[:curveOrderByteSize])
	s.SetBytes(signature[curveOrderByteSize:])

	// Verify s is in the lower half of the curve order
	// This protects against signature malleability
	halfOrder := new(big.Int).Rsh(pubkey.Curve.Params().N, 1)
	if s.Cmp(halfOrder) > 0 {
		return false
	}

	return ecdsa.Verify(pubkey, digest[:], r, s)
}
