package main

import (
	"bytes"
	compress "compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

//go:embed templates/standalone.go
var standalone string

func hashPassword(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func EncryptAES(key []byte, plaintext []byte) (string, error) {
	key = hashPassword(string(key))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	plaintextPadded := append(plaintext, padtext...)

	ciphertext := make([]byte, aes.BlockSize+len(plaintextPadded))
	iv := ciphertext[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintextPadded)

	return hex.EncodeToString(ciphertext), nil
}

func DecryptAES(key []byte, hexCipher string) (string, error) {
	key = hashPassword(string(key))
	ciphertext, err := hex.DecodeString(hexCipher)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	padding := int(ciphertext[len(ciphertext)-1])
	if padding > aes.BlockSize || padding == 0 {
		return "", errors.New("incorrect password")
	}
	for i := len(ciphertext) - padding; i < len(ciphertext); i++ {
		if ciphertext[i] != byte(padding) {
			return "", errors.New("incorrect password")
		}
	}
	return string(ciphertext[:len(ciphertext)-padding]), nil
}

type Config struct {
	Destruct    int    `json:"destruct"`
	ExecuteMode bool   `json:"is___executable"`
	Single_use  bool   `json:"single___use"`
	Data        string `json:"data"`
}

func makeNewStandalone(password string, binaryPath string, config Config) {
	var b bytes.Buffer
	w := compress.NewWriter(&b)
	w.Write([]byte(config.Data))
	w.Close()
	compressed := b.Bytes()

	encrypted, err := EncryptAES([]byte(password), compressed)
	if err != nil {
		fmt.Printf("Error encrypting data: %v\n", err)
		return
	}
	config.Data = encrypted
	configJSON, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("Error marshalling config to JSON: %v\n", err)
		return
	}

	_, err = copy("templates/templates", binaryPath+"/standalone")
	if err != nil {
		fmt.Printf("Error copying template: %v\n", err)
		return
	}

	binaryData, err := os.ReadFile(binaryPath + "/standalone")
	if err != nil {
		fmt.Printf("Error reading binary: %v\n", err)
		return
	}

	marker := "===DATA_START==="

	var finalBinary bytes.Buffer
	finalBinary.Write(binaryData)
	finalBinary.WriteString(marker)
	finalBinary.Write(configJSON)

	err = os.WriteFile(binaryPath+"/standalone", finalBinary.Bytes(), 0755)
	if err != nil {
		fmt.Printf("Error writing final binary: %v\n", err)
		return
	}
}
func copy(src, dst string) (int64, error) {

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
