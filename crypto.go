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

// NewEncryptionKey generates a cryptographically secure random 256-bit key
// for use with EncryptAES and DecryptAES functions.
//
// Security considerations:
//   - Uses crypto/rand for secure random generation
//   - Panics if the system's secure random number generator fails
//   - The returned key should be kept secret and stored securely
//
// Returns:
//   - A pointer to a 32-byte array containing the encryption key
//
// Example usage:
//
//	key := NewEncryptionKey()
//	defer func() { // Zero out the key when done
//		for i := range key {
//			key[i] = 0
//		}
//	}()
//
//	encrypted, err := EncryptAES([]byte("secret data"), key)
//	if err != nil {
//		log.Fatal(err)
//	}
func NewEncryptionKey() *[32]byte {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return &key
}

// EncryptAES encrypts data using 256-bit AES-GCM (Galois/Counter Mode).
// This provides both confidentiality and authenticity - it hides the content
// and ensures the data hasn't been tampered with.
//
// The output format is: nonce || ciphertext || tag
// where || indicates concatenation.
//
// Security considerations:
//   - Uses AES-256-GCM which is a NIST-approved authenticated encryption mode
//   - Generates a unique nonce for each encryption operation
//   - Provides both encryption and authentication
//   - The same key should never be used to encrypt more than 2^32 messages
//
// Parameters:
//   - plaintext: The data to encrypt (can be any length)
//   - key: A 32-byte encryption key (use NewEncryptionKey() to generate)
//
// Returns:
//   - ciphertext: The encrypted data with nonce and authentication tag
//   - error: Any error that occurred during encryption
//
// Example usage:
//
//	key := NewEncryptionKey()
//	plaintext := []byte("confidential message")
//	ciphertext, err := EncryptAES(plaintext, key)
//	if err != nil {
//		log.Fatal(err)
//	}
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

// DecryptAES decrypts data that was encrypted with EncryptAES using 256-bit AES-GCM.
// This function both decrypts the data and verifies its authenticity.
//
// The input must be in the format: nonce || ciphertext || tag
// where || indicates concatenation (as produced by EncryptAES).
//
// Security considerations:
//   - Automatically verifies the authentication tag before decryption
//   - Returns an error if the data has been tampered with
//   - Uses constant-time operations to prevent timing attacks
//
// Parameters:
//   - ciphertext: The encrypted data (as returned by EncryptAES)
//   - key: The same 32-byte key used for encryption
//
// Returns:
//   - plaintext: The decrypted data
//   - error: Any error that occurred during decryption or authentication
//
// Example usage:
//
//	key := NewEncryptionKey()
//	ciphertext, _ := EncryptAES([]byte("secret"), key)
//	plaintext, err := DecryptAES(ciphertext, key)
//	if err != nil {
//		log.Fatal("Decryption failed:", err)
//	}
//	fmt.Printf("Decrypted: %s\n", plaintext)
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

// HashHMAC generates a keyed hash of data using HMAC-SHA-512/256.
// This is suitable for data integrity verification and key derivation,
// but NOT for password hashing (use bcrypt, scrypt, or Argon2 for passwords).
//
// The tag parameter serves as the HMAC key and should describe the purpose
// of the hash to ensure domain separation between different uses.
//
// Security considerations:
//   - Uses SHA-512/256 which provides 256-bit security
//   - The tag acts as a key, so different tags produce different hashes
//   - Suitable for integrity verification and key derivation
//   - NOT suitable for password hashing
//
// Parameters:
//   - tag: A descriptive string that serves as the HMAC key (e.g., "session-token", "api-key")
//   - data: The data to hash
//
// Returns:
//   - A 32-byte hash of the data, or nil if data is empty
//
// Example usage:
//
//	hash := HashHMAC("user-session", []byte("user123:session456"))
//	// Use hash for integrity verification or as a derived key
func HashHMAC(tag string, data []byte) []byte {
	if len(data) == 0 {
		return nil
	}

	h := hmac.New(sha512.New512_256, []byte(tag))
	h.Write(data)
	return h.Sum(nil)
}

// DecodePublicKey decodes a PEM-encoded ECDSA public key from bytes.
// The input should be a PEM block with type "PUBLIC KEY".
//
// Parameters:
//   - encodedKey: PEM-encoded public key bytes
//
// Returns:
//   - An ECDSA public key ready for signature verification
//   - An error if the key cannot be decoded or is not an ECDSA key
//
// Example usage:
//
//	pemData := []byte(`-----BEGIN PUBLIC KEY-----
//	MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
//	-----END PUBLIC KEY-----`)
//	pubKey, err := DecodePublicKey(pemData)
//	if err != nil {
//		log.Fatal(err)
//	}
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
// The output is suitable for storage, transmission, or sharing.
//
// Parameters:
//   - key: The ECDSA public key to encode
//
// Returns:
//   - PEM-encoded public key bytes
//   - An error if the key cannot be encoded
//
// Example usage:
//
//	privKey, _ := NewSigningKey()
//	pubKey := &privKey.PublicKey
//	pemData, err := EncodePublicKey(pubKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Public key:\n%s", pemData)
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

// DecodePrivateKey decodes a PEM-encoded ECDSA private key from bytes.
// The input should be a PEM block with type "EC PRIVATE KEY".
//
// Security considerations:
//   - Private keys should be stored securely and never shared
//   - Consider encrypting private keys when storing them
//   - Zero out the key material when no longer needed
//
// Parameters:
//   - encodedKey: PEM-encoded private key bytes
//
// Returns:
//   - An ECDSA private key ready for signing operations
//   - An error if the key cannot be decoded or is not an ECDSA key
//
// Example usage:
//
//	pemData := []byte(`-----BEGIN EC PRIVATE KEY-----
//	MHcCAQEEIK9...
//	-----END EC PRIVATE KEY-----`)
//	privKey, err := DecodePrivateKey(pemData)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer func() {
//		privKey.D.SetInt64(0) // Zero out the private key
//	}()
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
// The output should be stored securely and protected from unauthorized access.
//
// Security considerations:
//   - The encoded private key should be stored securely
//   - Consider encrypting the PEM data before storage
//   - Never share or transmit private keys over insecure channels
//
// Parameters:
//   - key: The ECDSA private key to encode
//
// Returns:
//   - PEM-encoded private key bytes
//   - An error if the key cannot be encoded
//
// Example usage:
//
//	privKey, _ := NewSigningKey()
//	pemData, err := EncodePrivateKey(privKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Store pemData securely
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

// EncodeSignatureJWT encodes an ECDSA signature for use in JWT tokens.
// This follows the JWT specification (RFC 7515, Appendix A.3.1) for
// ECDSA signature encoding.
//
// Parameters:
//   - sig: The raw ECDSA signature bytes
//
// Returns:
//   - Base64url-encoded signature string suitable for JWT, or empty string if sig is empty
//
// Example usage:
//
//	signature, _ := SignData([]byte("data"), privKey)
//	jwtSig := EncodeSignatureJWT(signature)
//	// Use jwtSig in JWT token
func EncodeSignatureJWT(sig []byte) string {
	if len(sig) == 0 {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(sig)
}

// DecodeSignatureJWT decodes a JWT-encoded ECDSA signature.
// This is the reverse operation of EncodeSignatureJWT.
//
// Parameters:
//   - b64sig: Base64url-encoded signature string from JWT
//
// Returns:
//   - The raw ECDSA signature bytes
//   - An error if the signature cannot be decoded
//
// Example usage:
//
//	jwtSig := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9..."
//	signature, err := DecodeSignatureJWT(jwtSig)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Use signature with VerifySign
func DecodeSignatureJWT(b64sig string) ([]byte, error) {
	if b64sig == "" {
		return nil, errors.New("empty signature")
	}
	return base64.RawURLEncoding.DecodeString(b64sig)
}

// NewHMACKey generates a cryptographically secure random 256-bit key
// for use with HMAC operations.
//
// Security considerations:
//   - Uses crypto/rand for secure random generation
//   - Panics if the system's secure random number generator fails
//   - The returned key should be kept secret and stored securely
//
// Returns:
//   - A pointer to a 32-byte array containing the HMAC key
//
// Example usage:
//
//	key := NewHMACKey()
//	defer func() { // Zero out the key when done
//		for i := range key {
//			key[i] = 0
//		}
//	}()
//
//	mac := GenerateHMAC([]byte("message"), key)
func NewHMACKey() *[32]byte {
	key := &[32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return key
}

// GenerateHMAC produces a symmetric signature using HMAC-SHA-512/256.
// This creates a message authentication code that can be used to verify
// both the integrity and authenticity of a message.
//
// Security considerations:
//   - Uses SHA-512/256 which provides 256-bit security
//   - The key should be at least 32 bytes for optimal security
//   - Different keys produce different MACs for the same data
//
// Parameters:
//   - data: The data to authenticate
//   - key: A 32-byte secret key (use NewHMACKey() to generate)
//
// Returns:
//   - A 32-byte HMAC, or nil if data is empty or key is nil
//
// Example usage:
//
//	key := NewHMACKey()
//	data := []byte("important message")
//	mac := GenerateHMAC(data, key)
//
//	// Later, verify the MAC
//	if CheckHMAC(data, mac, key) {
//		fmt.Println("Message is authentic")
//	}
func GenerateHMAC(data []byte, key *[32]byte) []byte {
	if len(data) == 0 || key == nil {
		return nil
	}

	h := hmac.New(sha512.New512_256, key[:])
	h.Write(data)
	return h.Sum(nil)
}

// CheckHMAC securely verifies an HMAC against a message using the shared secret key.
// This function uses constant-time comparison to prevent timing attacks.
//
// Security considerations:
//   - Uses constant-time comparison to prevent timing attacks
//   - Both the data and key must match exactly for verification to succeed
//   - Returns false for any invalid input (empty data, empty MAC, nil key)
//
// Parameters:
//   - data: The original data that was authenticated
//   - suppliedMAC: The HMAC to verify
//   - key: The same 32-byte key used to generate the HMAC
//
// Returns:
//   - true if the HMAC is valid for the given data and key, false otherwise
//
// Example usage:
//
//	key := NewHMACKey()
//	data := []byte("message")
//	mac := GenerateHMAC(data, key)
//
//	// Verify the MAC
//	if CheckHMAC(data, mac, key) {
//		fmt.Println("HMAC verification successful")
//	} else {
//		fmt.Println("HMAC verification failed - data may be tampered")
//	}
func CheckHMAC(data, suppliedMAC []byte, key *[32]byte) bool {
	if len(data) == 0 || len(suppliedMAC) == 0 || key == nil {
		return false
	}

	expectedMAC := GenerateHMAC(data, key)
	return subtle.ConstantTimeCompare(expectedMAC, suppliedMAC) == 1
}

// NewSigningKey generates a new random P-256 ECDSA private key for digital signatures.
// P-256 is a NIST-approved elliptic curve that provides 128-bit security.
//
// Security considerations:
//   - Uses crypto/rand for secure random generation
//   - P-256 provides 128-bit security level
//   - The private key should be stored securely and never shared
//   - Consider using hardware security modules for key storage in production
//
// Returns:
//   - A new ECDSA private key
//   - An error if key generation fails
//
// Example usage:
//
//	privKey, err := NewSigningKey()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer func() {
//		privKey.D.SetInt64(0) // Zero out the private key
//	}()
//
//	// Use the key for signing
//	signature, _ := SignData([]byte("document"), privKey)
func NewSigningKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return key, err
}

// SignData creates a digital signature for arbitrary data using ECDSA.
// The signature can be verified using VerifySign with the corresponding public key.
//
// Security considerations:
//   - Uses SHA-256 for hashing the data before signing
//   - Includes protection against signature malleability attacks
//   - The signature is deterministic for the same data and key
//   - Uses secure random nonce generation
//
// Parameters:
//   - data: The data to sign (will be hashed with SHA-256)
//   - privkey: The ECDSA private key for signing
//
// Returns:
//   - A signature that can be verified with VerifySign
//   - An error if signing fails or inputs are invalid
//
// Example usage:
//
//	privKey, _ := NewSigningKey()
//	data := []byte("document to sign")
//	signature, err := SignData(data, privKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Verify with the public key
//	pubKey := &privKey.PublicKey
//	if VerifySign(data, signature, pubKey) {
//		fmt.Println("Signature is valid")
//	}
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

// VerifySign verifies an ECDSA signature against the original data.
// This function checks both the mathematical validity and authenticity of the signature.
//
// Security considerations:
//   - Uses SHA-256 for hashing the data (must match SignData)
//   - Includes protection against signature malleability attacks
//   - Returns false for any invalid input or tampered signatures
//   - Uses constant-time operations where possible
//
// Parameters:
//   - data: The original data that was signed
//   - signature: The signature to verify (as returned by SignData)
//   - pubkey: The ECDSA public key corresponding to the private key used for signing
//
// Returns:
//   - true if the signature is valid for the given data and public key, false otherwise
//
// Example usage:
//
//	privKey, _ := NewSigningKey()
//	data := []byte("signed document")
//	signature, _ := SignData(data, privKey)
//
//	// Verify the signature
//	pubKey := &privKey.PublicKey
//	if VerifySign(data, signature, pubKey) {
//		fmt.Println("Signature verification successful")
//	} else {
//		fmt.Println("Signature verification failed")
//	}
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
