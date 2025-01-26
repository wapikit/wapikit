package encryption_service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
)

type EncryptionService struct {
	Logger *slog.Logger
	Key    string
}

func NewEncryptionService(
	logger *slog.Logger,
	key string,
) *EncryptionService {
	return &EncryptionService{
		Logger: logger,
		Key:    key,
	}
}

// Encrypt any data and return an encrypted string.
func (es *EncryptionService) EncryptData(data interface{}) (string, error) {
	// Serialize the data into JSON
	plainText, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to serialize data: %w", err)
	}

	fmt.Println("encryption key: ", es.Key)

	// Decode the base64 key
	keyBytes, err := base64.StdEncoding.DecodeString(es.Key)
	if err != nil {
		return "", fmt.Errorf("invalid key: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM (Galois/Counter Mode)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt a string and store the result in the provided pointer.
func (es *EncryptionService) DecryptData(encryptedData string, result interface{}) error {
	// Decode the base64 key and encrypted data
	keyBytes, err := base64.StdEncoding.DecodeString(es.Key)
	if err != nil {
		return fmt.Errorf("invalid key: %w", err)
	}

	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return fmt.Errorf("invalid encrypted data: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM (Galois/Counter Mode)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// Separate the nonce and cipher text
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return errors.New("invalid encrypted data length")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	// Decrypt the data
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	// Deserialize the JSON into the result pointer
	err = json.Unmarshal(plainText, result)
	if err != nil {
		return fmt.Errorf("failed to deserialize data: %w", err)
	}

	return nil
}
