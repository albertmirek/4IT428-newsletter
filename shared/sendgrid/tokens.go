package sendgrid

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/pbkdf2"
	"os"
	"strings"
)

type ValueFromToken struct {
	UserID       string
	UserEmail    string
	NewsletterID string
}

const (
	iterations = 4096
	keyLen     = 32 // AES-256
)

func EncryptUserNewsletterToken(userID, userEmail string, newsletterID string) (string, error) {
	values := userID + ":" + newsletterID + ":" + userEmail
	password := []byte(os.Getenv("ENCRYPTION_PASSWORD"))
	salt := []byte(os.Getenv("SALT"))

	key := pbkdf2.Key(password, salt, iterations, keyLen, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(values))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(values))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decryptUserNewsletterToken(encryptedText string) (string, error) {
	password := []byte(os.Getenv("ENCRYPTION_PASSWORD"))
	salt := []byte(os.Getenv("SALT"))
	key := pbkdf2.Key(password, salt, iterations, keyLen, sha256.New)

	ciphertext, err := base64.URLEncoding.DecodeString(encryptedText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("Cipher text is too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func GetValuesFromEncryptedToken(token string) (*ValueFromToken, error) {
	decryptedText, err := decryptUserNewsletterToken(token)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(decryptedText, ":")
	if len(parts) != 3 {
		return nil, err
	}

	values := &ValueFromToken{
		UserID:       parts[0],
		UserEmail:    parts[2],
		NewsletterID: parts[1],
	}

	return values, nil
}
