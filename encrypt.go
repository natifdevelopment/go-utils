package utils

import "github.com/natifdevelopment/go-types"

func Encrypt(plaintext string, key string) (string, error) {
	return types.Encrypt(plaintext, key)
}

func Decrypt(ciphertextHex string, key string) (string, error) {
	return types.Decrypt(ciphertextHex, key)
}
