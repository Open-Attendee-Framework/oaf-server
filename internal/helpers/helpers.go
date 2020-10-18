package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"

	"golang.org/x/crypto/scrypt"
)

const (
	saltSize = 32
)

//FileExists check if a file exists
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

//GenerateSalt generates a random Salt for a Passowrd
func GenerateSalt() (string, error) {
	buf := make([]byte, saltSize)
	_, err := io.ReadFull(rand.Reader, buf)
	return string(buf), err
}

//GeneratePasswordHash takes a password and a salt to generate a SCRYPT hashed password
func GeneratePasswordHash(p, s string) (string, error) {
	dk, err := scrypt.Key([]byte(p), []byte(s), 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}
	h := hex.EncodeToString(dk)
	return h, err
}

//CopyIfNotEmpty takes two strings an copies the second one over the first if its not empty
func CopyIfNotEmpty(str1, str2 string) string {
	if str2 != "" {
		return str2
	}
	return str1

}
