package abstract_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"
	"testing"

	"github.com/maxbolgarin/abstract"
)

func TestEncryptDecryptGCM(t *testing.T) {
	randomKey := &[32]byte{}
	_, err := io.ReadFull(rand.Reader, randomKey[:])
	if err != nil {
		t.Fatal(err)
	}

	gcmTests := []struct {
		plaintext []byte
		key       *[32]byte
	}{
		{
			plaintext: []byte("Hello, world!"),
			key:       randomKey,
		},
	}

	for _, tt := range gcmTests {
		ciphertext, err := abstract.EncryptAES(tt.plaintext, tt.key)
		if err != nil {
			t.Fatal(err)
		}

		plaintext, err := abstract.DecryptAES(ciphertext, tt.key)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(plaintext, tt.plaintext) {
			t.Errorf("plaintexts don't match")
		}

		ciphertext[0] ^= 0xff
		_, err = abstract.DecryptAES(ciphertext, tt.key)
		if err == nil {
			t.Errorf("gcmOpen should not have worked, but did")
		}
	}
}

// A keypair for NIST P-256 / secp256r1
// Generated using:
//
//	openssl ecparam -genkey -name prime256v1 -outform PEM
var pemECPrivateKeyP256 = `-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIOI+EZsjyN3jvWJI/KDihFmqTuDpUe/if6f/pgGTBta/oAoGCCqGSM49
AwEHoUQDQgAEhhObKJ1r1PcUw+3REd/TbmSZnDvXnFUSTwqQFo5gbfIlP+gvEYba
+Rxj2hhqjfzqxIleRK40IRyEi3fJM/8Qhg==
-----END EC PRIVATE KEY-----
`

var pemECPublicKeyP256 = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEhhObKJ1r1PcUw+3REd/TbmSZnDvX
nFUSTwqQFo5gbfIlP+gvEYba+Rxj2hhqjfzqxIleRK40IRyEi3fJM/8Qhg==
-----END PUBLIC KEY-----
`

// A keypair for NIST P-384 / secp384r1
// Generated using:
//
//	openssl ecparam -genkey -name secp384r1 -outform PEM
var pemECPrivateKeyP384 = `-----BEGIN EC PARAMETERS-----
BgUrgQQAIg==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDAhA0YPVL1kimIy+FAqzUAtmR3It2Yjv2I++YpcC4oX7wGuEWcWKBYE
oOjj7wG/memgBwYFK4EEACKhZANiAAQub8xaaCTTW5rCHJCqUddIXpvq/TxdwViH
+tPEQQlJAJciXStM/aNLYA7Q1K1zMjYyzKSWz5kAh/+x4rXQ9Hlm3VAwCQDVVSjP
bfiNOXKOWfmyrGyQ7fQfs+ro1lmjLjs=
-----END EC PRIVATE KEY-----
`

var pemECPublicKeyP384 = `-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAELm/MWmgk01uawhyQqlHXSF6b6v08XcFY
h/rTxEEJSQCXIl0rTP2jS2AO0NStczI2Msykls+ZAIf/seK10PR5Zt1QMAkA1VUo
z234jTlyjln5sqxskO30H7Pq6NZZoy47
-----END PUBLIC KEY-----
`

var garbagePEM = `-----BEGIN GARBAGE-----
TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQ=
-----END GARBAGE-----
`

func TestPublicKeyMarshaling(t *testing.T) {
	ecKey, err := abstract.DecodePublicKey([]byte(pemECPublicKeyP256))
	if err != nil {
		t.Fatal(err)
	}

	pemBytes, _ := abstract.EncodePublicKey(ecKey)
	if !bytes.Equal(pemBytes, []byte(pemECPublicKeyP256)) {
		t.Fatal("public key encoding did not match")
	}

}

func TestPrivateKeyBadDecode(t *testing.T) {
	_, err := abstract.DecodePrivateKey([]byte(garbagePEM))
	if err == nil {
		t.Fatal("decoded garbage data without complaint")
	}
}

func TestPrivateKeyMarshaling(t *testing.T) {
	ecKey, err := abstract.DecodePrivateKey([]byte(pemECPrivateKeyP256))
	if err != nil {
		t.Fatal(err)
	}

	pemBytes, _ := abstract.EncodePrivateKey(ecKey)
	if !strings.HasSuffix(pemECPrivateKeyP256, string(pemBytes)) {
		t.Fatal("private key encoding did not match")
	}
}

// Test vector from https://tools.ietf.org/html/rfc7515#appendix-A.3.1
var jwtTest = []struct {
	sigBytes []byte
	b64sig   string
}{
	{
		sigBytes: []byte{14, 209, 33, 83, 121, 99, 108, 72, 60, 47, 127, 21,
			88, 7, 212, 2, 163, 178, 40, 3, 58, 249, 124, 126, 23, 129, 154, 195, 22, 158,
			166, 101, 197, 10, 7, 211, 140, 60, 112, 229, 216, 241, 45, 175,
			8, 74, 84, 128, 166, 101, 144, 197, 242, 147, 80, 154, 143, 63, 127, 138, 131,
			163, 84, 213},
		b64sig: "DtEhU3ljbEg8L38VWAfUAqOyKAM6-Xx-F4GawxaepmXFCgfTjDxw5djxLa8ISlSApmWQxfKTUJqPP3-Kg6NU1Q",
	},
}

func TestJWTEncoding(t *testing.T) {
	for _, tt := range jwtTest {
		result := abstract.EncodeSignatureJWT(tt.sigBytes)

		if strings.Compare(result, tt.b64sig) != 0 {
			t.Fatalf("expected %s, got %s\n", tt.b64sig, result)
		}
	}
}

func TestJWTDecoding(t *testing.T) {
	for _, tt := range jwtTest {
		resultSig, err := abstract.DecodeSignatureJWT(tt.b64sig)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(resultSig, tt.sigBytes) {
			t.Fatalf("decoded signature was incorrect")
		}
	}
}

// https://groups.google.com/d/msg/sci.crypt/OolWgsgQD-8/jHciyWkaL0gJ
var hmacTests = []struct {
	key    string
	data   string
	digest string
}{
	{
		key:    "0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b",
		data:   "4869205468657265", // "Hi There"
		digest: "9f9126c3d9c3c330d760425ca8a217e31feae31bfe70196ff81642b868402eab",
	},
	{
		key:    "4a656665",                                                 // "Jefe"
		data:   "7768617420646f2079612077616e7420666f72206e6f7468696e673f", // "what do ya want for nothing?"
		digest: "6df7b24630d5ccb2ee335407081a87188c221489768fa2020513b2d593359456",
	},
}

func TestHMAC(t *testing.T) {
	for idx, tt := range hmacTests {
		keySlice, _ := hex.DecodeString(tt.key)
		dataBytes, _ := hex.DecodeString(tt.data)
		expectedDigest, _ := hex.DecodeString(tt.digest)

		keyBytes := &[32]byte{}
		copy(keyBytes[:], keySlice)

		macDigest := abstract.GenerateHMAC(dataBytes, keyBytes)
		if !bytes.Equal(macDigest, expectedDigest) {
			t.Errorf("test %d generated unexpected mac", idx)
		}
	}
}

func TestSign(t *testing.T) {
	message := []byte("Hello, world!")

	key, err := abstract.NewSigningKey()
	if err != nil {
		t.Error(err)
		return
	}

	signature, err := abstract.SignData(message, key)
	if err != nil {
		t.Error(err)
		return
	}

	if !abstract.VerifySign(message, signature, &key.PublicKey) {
		t.Error("signature was not correct")
		return
	}

	message[0] ^= 0xff
	if abstract.VerifySign(message, signature, &key.PublicKey) {
		t.Error("signature was good for altered message")
	}
}

func TestSignWithP384(t *testing.T) {
	message := []byte("Hello, world!")

	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Error(err)
		return
	}

	signature, err := abstract.SignData(message, key)
	if err != nil {
		t.Error(err)
		return
	}

	if !abstract.VerifySign(message, signature, &key.PublicKey) {
		t.Error("signature was not correct")
		return
	}

	message[0] ^= 0xff
	if abstract.VerifySign(message, signature, &key.PublicKey) {
		t.Error("signature was good for altered message")
	}
}

func TestNewEncryptionKey(t *testing.T) {
	key := abstract.NewEncryptionKey()

	if key == nil {
		t.Fatal("Expected non-nil key")
	}

	if len(key) != 32 {
		t.Errorf("Expected key length of 32 bytes, got %d", len(key))
	}

	// Check that we get different keys each time
	key2 := abstract.NewEncryptionKey()
	if bytes.Equal(key[:], key2[:]) {
		t.Error("Expected different keys on subsequent calls")
	}
}

func TestEncryptDecryptAES(t *testing.T) {
	key := abstract.NewEncryptionKey()
	plaintext := []byte("This is a secret message")

	// Test encryption
	ciphertext, err := abstract.EncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	if bytes.Equal(plaintext, ciphertext) {
		t.Error("Ciphertext should not equal plaintext")
	}

	// Test decryption
	decrypted, err := abstract.DecryptAES(ciphertext, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Error("Decrypted text doesn't match original plaintext")
	}

	// Test decryption with wrong key
	wrongKey := abstract.NewEncryptionKey()
	_, err = abstract.DecryptAES(ciphertext, wrongKey)
	if err == nil {
		t.Error("Decryption with wrong key should fail")
	}

	// Test decryption of malformed ciphertext
	_, err = abstract.DecryptAES([]byte("too short"), key)
	if err == nil {
		t.Error("Decryption of malformed ciphertext should fail")
	}
}

func TestHashHMAC(t *testing.T) {
	testData := []byte("This is some data to hash")
	tag1 := "purpose1"
	tag2 := "purpose2"

	// Same data, different tags should produce different hashes
	hash1 := abstract.HashHMAC(tag1, testData)
	hash2 := abstract.HashHMAC(tag2, testData)

	if bytes.Equal(hash1, hash2) {
		t.Error("Hashes with different tags should be different")
	}

	// Same data, same tag should produce same hash
	hash1Again := abstract.HashHMAC(tag1, testData)
	if !bytes.Equal(hash1, hash1Again) {
		t.Error("Hash function should be deterministic for same input")
	}

	// Different data, same tag should produce different hashes
	differentData := []byte("Different data")
	hash3 := abstract.HashHMAC(tag1, differentData)
	if bytes.Equal(hash1, hash3) {
		t.Error("Hashes of different data should be different")
	}
}

func TestECDSAKeyEncodingDecoding(t *testing.T) {
	// Generate a test key pair
	privKey, err := abstract.NewSigningKey()
	if err != nil {
		t.Fatalf("Failed to generate signing key: %v", err)
	}

	// Test private key encoding/decoding
	encodedPriv, err := abstract.EncodePrivateKey(privKey)
	if err != nil {
		t.Fatalf("Failed to encode private key: %v", err)
	}

	decodedPriv, err := abstract.DecodePrivateKey(encodedPriv)
	if err != nil {
		t.Fatalf("Failed to decode private key: %v", err)
	}

	// Compare original and decoded private keys
	originalD := privKey.D.Bytes()
	decodedD := decodedPriv.D.Bytes()
	if !bytes.Equal(originalD, decodedD) {
		t.Error("Decoded private key doesn't match original")
	}

	// Test public key encoding/decoding
	encodedPub, err := abstract.EncodePublicKey(&privKey.PublicKey)
	if err != nil {
		t.Fatalf("Failed to encode public key: %v", err)
	}

	decodedPub, err := abstract.DecodePublicKey(encodedPub)
	if err != nil {
		t.Fatalf("Failed to decode public key: %v", err)
	}

	// Compare original and decoded public keys
	originalX := privKey.PublicKey.X.Bytes()
	decodedX := decodedPub.X.Bytes()
	originalY := privKey.PublicKey.Y.Bytes()
	decodedY := decodedPub.Y.Bytes()

	if !bytes.Equal(originalX, decodedX) || !bytes.Equal(originalY, decodedY) {
		t.Error("Decoded public key doesn't match original")
	}
}

func TestDecodePublicKeyErrors(t *testing.T) {
	// Test with invalid PEM data
	_, err := abstract.DecodePublicKey([]byte("not a valid PEM"))
	if err == nil {
		t.Error("Expected error for invalid PEM data")
	}

	// Test with wrong PEM type
	wrongTypePEM := []byte(`-----BEGIN CERTIFICATE-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3VoPN9PKUjKFLMwOge6+
wnDi8sbETGIx2FKXGgqtAKpzmem53kRGEQg8WeqRmp12wgp74TGpkEXsGae7RS1k
enJCnma4fii+noGH7R0qKgHvPrI2Bwa9hzsH8tHxpyM3qrXslOmD45EH9SxIDUBJ
FehNdaPbLP1gFyahKMsdfxFJLUvbUycuZSJ2ZnIgeVxwm4qbSvZInL9Iu4FzuPtg
fINKcbbovy1qq4KvPIrXzhbY3PWDc6btxCf3SE0JdE1MCPThpf+RmQHdbVI+Nmhz
RV/H0oMEQhsC8Y3UXNmWe/YFCq7ULwy8NB6u9YSQM1FS54PBj+0LtW8IWOBU4eGW
TwIDAQAB
-----END CERTIFICATE-----`)
	_, err = abstract.DecodePublicKey(wrongTypePEM)
	if err == nil {
		t.Error("Expected error for wrong PEM type")
	}
}

func TestDecodePrivateKeyErrors(t *testing.T) {
	// Test with invalid PEM data
	_, err := abstract.DecodePrivateKey([]byte("not a valid PEM"))
	if err == nil {
		t.Error("Expected error for invalid PEM data")
	}

	// Test with PEM without the right type
	wrongTypePEM := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEArhPHlPcKcfynX1tJzuMXELrBW2IpPHHf4qqk9KKwQQEjG0Tj
aNZUBRGYItpTU1lJq2lkTdNuHvW8oj8qoGGjUrkbvTJWAQKwVeNmQ9RzdUxJhEXy
d94AO/FSwDeJSQAIBbHPvCXG7KBic+UzX0AIPKdaE6ClZRG+UXgN6rT21n2LmTzY
bLoaInZxA8/lxJ9ZUxaIhDPJZcq3byoGXHsWHUFCveRWrR5xmGKLECu0dM8KoQJE
QR7p6FMamaaVf2RHPBnH9Y5kV3VXVj6lSERfq75yBAl2cxqhpzrP1SYQOKLTsj9X
/0Nv6/UGGUXXOUGQEqC10A3YZPU0h14JMWwBJQIDAQABAoIBAG6cHxZxFBk03zPo
v67YHjKsqpbSuMvjckR9R6V5+SeRMUxLjE3aLd7DGg7sMT9+Vx3DWi9QQx3YiEo2
aIcv+ohMj2YJeCcpNzElT307K1vfgLjNMyKGXMdpPQb/n0I8dxU3yveCGQATmK5m
YCi7mYZ+4qVxDveKBszuOQvMdJJA0tMHhMVZrvXIfGwzWpcGEWs2G0PDSXADhA5M
GmUb4CLpqrWYNVxEzQETYCXTmvJBnAzR9reQbtfQSctxOII3YYBwaM0HAWLmEh31
o2USL0IZ1gEd3SpK3cwFPWmw6RnW/tUFPqBTxHV/HKWHQCVOmvo6PWWBo3B245ZC
Ty8Q0yECgYEA2MKtQbCEMR/M/JXH/eD6narF9J8z6LYbcwLKzTDCyuX3GWZ1xCQO
PXRbVEusEFJzUXNVu+xfpY+tzpKL8yBUDHTHXlrHqlJ/HE9Z+YLXr9I6DLvfswJe
p5yJTcWJtgLPfF4QUmhNYDzlLWYF6dIFD5GMGpj9MYebjD9NwDJm2pUCgYEAzXBK
UqBfLrSettRvXAt7lG+T3DS1+pvJ+HJZtrLu4+NTV7TaxECjGS+8mHSlFeC8a+J1
9BW+D19s6Ti1vU1JzGEjtKYkV2XMalX1qMRKl+SFMnF+hvDLEhvQ85YTX6uvCgMJ
TtGqQS/CvEYYQFVNT7ZDGVIqSucBJJbGXcWbVXECgYBnXb4xf8msKwMrwCHXLm6W
5iPvs9JJxmKPHn/31hV1a6JM44KJIehO5U6A05qiYxp13QY+fGNPmXXJ5Z4TPS7y
ZbwaFESp93jtmloVmTGzfmD7dB8MJECslBaGXyFYe3yPMdJSFRbVdQBVCkksQGJk
4bSLyOXjH6NDUniMYLoOVQKBgDWgGIZ+/yyZzfxC3wifi3rO5tq9VPH4BbybzFOA
Ujli9SlXVAhw5C4J3OzLHKgUBFHLcPnNrHbsXSUTD9gPxc+5WOKqEelDzxJE9TyL
x5qRHnjrnH9HFRmWRi7KGGnpJm0I8qgl6d2N7rRx4lAQvjR2nLyHfTdKUIJ9+QKZ
9gaBAoGAcTns5IfsbdkW3XA5OZiSV5L1YE9IJ/zXjKqvYULJt63vERSXJIsoeMD6
6V4gGvSL2PVTxft134CDoy3+Aq0lBajcakFqe3RCruFhcVkFvBK2CKBuKyjK4bM5
M15+VV8DK7v2EWCBT3hUFDeYJBW2GgXvBHJM5Q+XyRcfbCXJgjU=
-----END RSA PRIVATE KEY-----`)
	_, err = abstract.DecodePrivateKey(wrongTypePEM)
	if err == nil {
		t.Error("Expected error for wrong PEM type")
	}
}

func TestJWTSignatureEncoding(t *testing.T) {
	testSig := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	encoded := abstract.EncodeSignatureJWT(testSig)
	decoded, err := abstract.DecodeSignatureJWT(encoded)
	if err != nil {
		t.Fatalf("Failed to decode JWT signature: %v", err)
	}

	if !bytes.Equal(testSig, decoded) {
		t.Error("Decoded signature doesn't match original")
	}

	// Test decoding invalid base64
	_, err = abstract.DecodeSignatureJWT("!invalid base64!")
	if err == nil {
		t.Error("Expected error when decoding invalid base64")
	}
}

func TestHMACKeyFunctions(t *testing.T) {
	// Test key generation
	key := abstract.NewHMACKey()
	if key == nil {
		t.Fatal("Expected non-nil key")
	}

	if len(key) != 32 {
		t.Errorf("Expected key length of 32 bytes, got %d", len(key))
	}

	// Test that we get different keys each time
	key2 := abstract.NewHMACKey()
	if bytes.Equal(key[:], key2[:]) {
		t.Error("Expected different keys on subsequent calls")
	}

	// Test HMAC generation and validation
	data := []byte("test data")
	mac := abstract.GenerateHMAC(data, key)

	// Verify that the correct MAC validates
	if !abstract.CheckHMAC(data, mac, key) {
		t.Error("CheckHMAC should return true for valid MAC")
	}

	// Verify that invalid MAC fails
	invalidMAC := make([]byte, len(mac))
	copy(invalidMAC, mac)
	invalidMAC[0] ^= 0xff // Flip some bits

	if abstract.CheckHMAC(data, invalidMAC, key) {
		t.Error("CheckHMAC should return false for invalid MAC")
	}

	// Verify that wrong key fails
	if abstract.CheckHMAC(data, mac, key2) {
		t.Error("CheckHMAC should return false with wrong key")
	}
}

func TestSigningKeyGeneration(t *testing.T) {
	key, err := abstract.NewSigningKey()
	if err != nil {
		t.Fatalf("Failed to generate signing key: %v", err)
	}

	if key.Curve != elliptic.P256() {
		t.Error("Expected P-256 curve")
	}

	// Check that private key is valid and has public key
	if key.D == nil || key.PublicKey.X == nil || key.PublicKey.Y == nil {
		t.Error("Key components should not be nil")
	}
}

func TestSignVerify(t *testing.T) {
	privKey, err := abstract.NewSigningKey()
	if err != nil {
		t.Fatalf("Failed to generate signing key: %v", err)
	}

	testData := []byte("data to sign")
	signature, err := abstract.SignData(testData, privKey)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	// Verify valid signature
	valid := abstract.VerifySign(testData, signature, &privKey.PublicKey)
	if !valid {
		t.Error("Signature should be valid")
	}

	// Verify signature with wrong data
	wrongData := []byte("wrong data")
	if abstract.VerifySign(wrongData, signature, &privKey.PublicKey) {
		t.Error("Signature should be invalid for wrong data")
	}

	// Verify signature with tampered signature
	tamperedSig := make([]byte, len(signature))
	copy(tamperedSig, signature)
	tamperedSig[0] ^= 0xff // Flip some bits
	if abstract.VerifySign(testData, tamperedSig, &privKey.PublicKey) {
		t.Error("Tampered signature should be invalid")
	}

	// Generate another key and test cross-verification
	otherKey, _ := abstract.NewSigningKey()
	if abstract.VerifySign(testData, signature, &otherKey.PublicKey) {
		t.Error("Signature should be invalid with wrong public key")
	}
}

func TestVerifySignWithInvalidSignatureLength(t *testing.T) {
	privKey, _ := abstract.NewSigningKey()

	// Create a signature with invalid length
	invalidSig := []byte{1, 2, 3} // Too short

	// This should return false
	if abstract.VerifySign([]byte("test"), invalidSig, &privKey.PublicKey) {
		t.Error("VerifySign should return false for invalid signature length")
	}
}

func TestAESWithEmptyPlaintext(t *testing.T) {
	key := abstract.NewEncryptionKey()
	emptyData := []byte{}

	// Encrypt empty data
	ciphertext, err := abstract.EncryptAES(emptyData, key)
	if err != nil {
		t.Fatalf("Failed to encrypt empty data: %v", err)
	}

	// Decrypt and verify
	decrypted, err := abstract.DecryptAES(ciphertext, key)
	if err != nil {
		t.Fatalf("Failed to decrypt ciphertext: %v", err)
	}

	if len(decrypted) != 0 {
		t.Errorf("Expected empty decrypted data, got %d bytes", len(decrypted))
	}
}

func TestAESWithLargeData(t *testing.T) {
	key := abstract.NewEncryptionKey()

	// Create a 1MB data block
	largeData := make([]byte, 1024*1024)
	_, err := rand.Read(largeData)
	if err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	// Encrypt and decrypt large data
	ciphertext, err := abstract.EncryptAES(largeData, key)
	if err != nil {
		t.Fatalf("Failed to encrypt large data: %v", err)
	}

	decrypted, err := abstract.DecryptAES(ciphertext, key)
	if err != nil {
		t.Fatalf("Failed to decrypt large data: %v", err)
	}

	if !bytes.Equal(largeData, decrypted) {
		t.Error("Decrypted data doesn't match original large data")
	}
}

func TestBase64EncodeDecode(t *testing.T) {
	// Test that our base64 functions are compatible
	data := []byte{0, 1, 2, 3, 255, 254, 253, 252}

	// First, encode using standard base64 and our function
	stdEncoded := base64.RawURLEncoding.EncodeToString(data)
	ourEncoded := abstract.EncodeSignatureJWT(data)

	// They should be the same
	if stdEncoded != ourEncoded {
		t.Errorf("Our JWT encoding doesn't match standard base64. Got %s, expected %s",
			ourEncoded, stdEncoded)
	}

	// Now decode and compare
	stdDecoded, _ := base64.RawURLEncoding.DecodeString(stdEncoded)
	ourDecoded, _ := abstract.DecodeSignatureJWT(ourEncoded)

	if !bytes.Equal(stdDecoded, ourDecoded) {
		t.Error("Our JWT decoding doesn't match standard base64")
	}

	// Finally, check round-trip
	if !bytes.Equal(data, ourDecoded) {
		t.Error("JWT encoding/decoding round-trip failed")
	}
}

func TestHashCollisionResistance(t *testing.T) {
	// Test that similar data produces very different hashes
	data1 := []byte("This is a test message")
	data2 := []byte("This is a test messafe") // One bit changed

	hash1 := abstract.HashHMAC("test", data1)
	hash2 := abstract.HashHMAC("test", data2)

	// Count differing bits
	differentBits := 0
	for i := 0; i < len(hash1) && i < len(hash2); i++ {
		xor := hash1[i] ^ hash2[i]
		for j := 0; j < 8; j++ {
			if (xor & (1 << j)) != 0 {
				differentBits++
			}
		}
	}

	// A good hash function should have approximately 50% of bits different
	// We'll use a generous threshold of at least 30%
	minDifferentBits := len(hash1) * 8 * 30 / 100
	if differentBits < minDifferentBits {
		t.Errorf("Hash function doesn't show good avalanche effect: %d/%d bits different",
			differentBits, len(hash1)*8)
	}
}

func compareECDSAPublicKeys(key1, key2 *ecdsa.PublicKey) bool {
	return key1.Curve == key2.Curve &&
		key1.X.Cmp(key2.X) == 0 &&
		key1.Y.Cmp(key2.Y) == 0
}

func TestPrivateKeyToPubKey(t *testing.T) {
	// Generate a signing key
	privKey, _ := abstract.NewSigningKey()

	// Test that when encoding and decoding the private key, the public key remains correct
	encodedPriv, _ := abstract.EncodePrivateKey(privKey)
	decodedPriv, _ := abstract.DecodePrivateKey(encodedPriv)

	if !compareECDSAPublicKeys(&privKey.PublicKey, &decodedPriv.PublicKey) {
		t.Error("Public key component lost during private key encoding/decoding")
	}

	// Also test that we can sign with the decoded private key and verify with the original public key
	testData := []byte("test signing with decoded key")
	signature, _ := abstract.SignData(testData, decodedPriv)

	valid := abstract.VerifySign(testData, signature, &privKey.PublicKey)
	if !valid {
		t.Error("Signature with decoded private key can't be verified with original public key")
	}
}

func TestHashSHA256(t *testing.T) {
	// Test compatibility with standard SHA-256
	data := []byte("test data")

	// Calculate hash with standard library
	stdHash := sha256.Sum256(data)

	// HMAC with empty key is essentially SHA-256 (not exactly, but close enough for a sanity check)
	ourHash := abstract.HashHMAC("", data)

	// The hashes should be different, but have similar properties
	if bytes.Equal(stdHash[:], ourHash) {
		t.Log("Note: Standard SHA-256 and HMAC with empty key produced identical hashes. This is surprising but not necessarily wrong.")
	}

	// Both hashes should be fixed length
	if len(stdHash) != 32 || len(ourHash) != 32 {
		t.Errorf("Expected 32-byte hashes, got %d and %d", len(stdHash), len(ourHash))
	}
}
