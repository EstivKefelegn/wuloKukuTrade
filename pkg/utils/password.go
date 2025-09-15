package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)


func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrorHandler(errors.New("password can't be null"), "password can't be null")
	}

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	
	if err != nil {
		return "", ErrorHandler(err, "failed to generate the salt")
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hastBasse64 := base64.StdEncoding.EncodeToString(hash)
	
	encodedhash := fmt.Sprintf("%s.%s", saltBase64, &hastBasse64)


	return encodedhash, nil
}