// template file modify here
package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ncruces/zenity"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const MARKER = "===DATA_START==="

type Config struct {
	Destruct    int    `json:"destruct"`
	ExecuteMode int    `json:"is_executable"`
	Single_use  int    `json:"single_use"`
	Data        string `json:"data"`
}

func main() {
	//TODO: encrypt JSON to prevent data from being modified
	config := extractData()
	zenity.Info(decrypt(config))
	go destruct()
	zenity.Info("This message has self-destructed")
	return
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

func extractData() Config {
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
	var config_json Config
	json.Unmarshal([]byte(config), &config_json)
	fmt.Printf("Embedded config: %s\n", config)
	return config_json
}

func decrypt(config Config) string {
	path, _ := os.Executable()
	_, password, _ := zenity.Password(zenity.Title("Enter Password"))
	decryptedData, err := DecryptAES([]byte(password), config.Data)

	if err != nil {
		if err.Error() == "incorrect password" {
			//TODO: prevent restarting from clearing counter (maybe write to the end?)
			if config.Destruct > 0 {
				config.Destruct -= 1
				if config.Destruct == 0 {
					os.Remove(path)
					return ""
				}
				zenity.Error(fmt.Sprintf("%d Tries Left", config.Destruct))
				return decrypt(config)
			}
			return ""
		}
		zenity.Error(err.Error())
		panic(err)
	}

	reader := bytes.NewReader([]byte(decryptedData))
	gzreader, e1 := gzip.NewReader(reader)
	if e1 != nil {
		zenity.Error("Failed to create gzip reader: " + e1.Error())
		panic(e1)
	}

	data, e2 := ioutil.ReadAll(gzreader)
	if e2 != nil {
		zenity.Error("Failed to decompress data: " + e2.Error())
		panic(e2)
	}

	return string(data)
}

func destruct() {
	time.Sleep(2 * time.Second)
	fmt.Println("Destruct executed, goodbye")
	path, _ := os.Executable()
	os.Remove(path)
}
