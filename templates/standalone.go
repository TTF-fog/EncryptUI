// template file modify here
package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ncruces/zenity"
	"io/ioutil"
	"os"
	"strings"
)

const MARKER = "===DATA_START==="

func main() {
	s := string(extractData())
	_, password, _ := zenity.Password(zenity.Title("Enter Password"))
	reader := bytes.NewReader([]byte(func(password string) string {
		data, err := DecryptAES([]byte(password), s)
		if err != nil {
			zenity.Error(err.Error())
			panic(err)
			return ""
		}
		return data
	}(password)))
	gzreader, e1 := gzip.NewReader(reader)
	if e1 != nil {
		fmt.Println(e1)
	}
	output, e2 := ioutil.ReadAll(gzreader)
	if e2 != nil {
		fmt.Println(e2)
	}
	result := string(output)
	zenity.Info(result)
}
func hashPassword(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
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

func extractData() []byte {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	markerBytes := []byte(MARKER)
	idx := bytes.Index(data, markerBytes)
	if idx == -1 {
		fmt.Println("No embedded data found")
		panic(errors.New("No embedded data found"))
	}
	embeddedData := data[idx+len(MARKER):]
	config := strings.TrimSpace(string(embeddedData))

	fmt.Printf("Embedded config: %s\n", config)
	return []byte(config)
}
