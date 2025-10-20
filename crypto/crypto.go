package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha3"
	"errors"
	"hash"

	"golang.org/x/crypto/chacha20poly1305"
)

// GenerateNonce generates a 24-byte nonce for XChaCha20.
func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, chacha20poly1305.NonceSizeX) // 24 bytes
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	return nonce, nil
}

// EncryptXChaCha20 encrypts plaintext using key and nonce.
func EncryptXChaCha20(key, nonce, plaintext []byte) ([]byte, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, errors.New("invalid key length")
	}
	if len(nonce) != chacha20poly1305.NonceSizeX {
		return nil, errors.New("invalid nonce length")
	}
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	return aead.Seal(nil, nonce, plaintext, nil), nil
}

// DecryptXChaCha20 decrypts ciphertext using key and nonce.
func DecryptXChaCha20(key, nonce, ciphertext []byte) ([]byte, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, errors.New("invalid key length")
	}
	if len(nonce) != chacha20poly1305.NonceSizeX {
		return nil, errors.New("invalid nonce length")
	}
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	return aead.Open(nil, nonce, ciphertext, nil)
}

// HashSHA3 returns the SHA3-256 hash of input data.
func HashSHA3(data []byte) []byte {
	h := sha3.New256()
	h.Write(data)
	return h.Sum(nil)
}

func HMAC_SHA3(key, data []byte) []byte {
	h := hmac.New(func() hash.Hash { return sha3.New256() }, key)
	h.Write(data)
	return h.Sum(nil)
}
