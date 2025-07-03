package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"
)

type file struct {
	name     string
	text     string
	password string
}

func hashPassword(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

const BlockSize = 16

func EncryptAES(key []byte, plaintext string) string {
	key = hashPassword(string(key))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := strings.Repeat(string(byte(padding)), padding)
	plaintextPadded := []byte(plaintext + padtext)

	ciphertext := make([]byte, aes.BlockSize+len(plaintextPadded))
	iv := ciphertext[:aes.BlockSize]

	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintextPadded)

	return hex.EncodeToString(ciphertext)
}
func DecryptAES(key []byte, hexCipher string) string {
	key = hashPassword(string(key))
	ciphertext, err := hex.DecodeString(hexCipher)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	padding := int(ciphertext[len(ciphertext)-1])
	if padding > aes.BlockSize || padding == 0 {

		return "Invalid Password"
	}
	for i := len(ciphertext) - padding; i < len(ciphertext); i++ {
		if ciphertext[i] != byte(padding) {

			return "Incorrect Password"
		}
	}
	return string(ciphertext[:len(ciphertext)-padding])
}
