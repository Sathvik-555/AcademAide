package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/pbkdf2"
)

// GenerateWallet creates a new Ethereum wallet and returns the private key (hex) and address (hex).
func GenerateWallet() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	return privateKeyHex, address, nil
}

// DeriveKey derives a 32-byte key from the password using PBKDF2.
// In a real production system, use a random salt stored with the data.
// For this implementation, we'll use a deterministic salt or include it in the output.
// To keep it simple and stateless for this lab, we will use a static salt or just the address as salt if we had it.
// Let's use a fixed salt for simplicity in this lab environment, or better, generate a salt and prepend it.
func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)
}

// EncryptPrivateKey encrypts the private key using the password.
// Returns hex encoded string: salt + nonce + ciphertext
func EncryptPrivateKey(privateKeyHex string, password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(privateKeyHex), nil)

	// Result: salt + ciphertext (nonce is inside ciphertext usually? No, Seal appends to dst. We passed nonce as dst)
	// gcm.Seal(dst, nonce, plaintext, data) -> dst + result.
	// So if we pass nonce as dst, it returns nonce + ciphertext + tag.
	// We need to store salt too.

	final := append(salt, ciphertext...)
	return hex.EncodeToString(final), nil
}

// DecryptPrivateKey decrypts the encrypted private key using the password.
func DecryptPrivateKey(encryptedHex string, password string) (string, error) {
	data, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}

	if len(data) < 16+12 { // salt(16) + nonce(12) + min_tag(16)
		return "", errors.New("invalid ciphertext length")
	}

	salt := data[:16]
	ciphertextWithNonce := data[16:]

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertextWithNonce) < gcm.NonceSize() {
		return "", errors.New("malformed ciphertext")
	}

	nonce := ciphertextWithNonce[:gcm.NonceSize()]
	actualCiphertext := ciphertextWithNonce[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
