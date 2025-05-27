package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const nonceSize = 12 // GCM standard

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("AES Crypt v0.1")
	fmt.Println("------------------------------------")
	fmt.Println("Select mode:")
	fmt.Println("[1] Encrypt file")
	fmt.Println("[2] Decrypt file")
	fmt.Print("Type 1 or 2: ")
	mode, _ := reader.ReadString('\n')
	mode = sanitize(mode)
	if mode != "1" && mode != "2" {
		fmt.Println("Unknown mode:", mode)
		return
	}

	fmt.Print("Enter input file path: ")
	inputPath, _ := reader.ReadString('\n')
	inputPath = sanitize(inputPath)

	outputPath := inputPath
	if mode == "1" {
		outputPath += ".enc"
	} else if mode == "2" && strings.HasSuffix(inputPath, ".enc") {
		outputPath = strings.TrimSuffix(inputPath, ".enc")
	} else if mode == "2" {
		outputPath += ".dec"
	}

	fmt.Println("Select key source:")
	fmt.Println("[1] Enter passphrase")
	fmt.Println("[2] Use key file (.key)")
	fmt.Print("Type 1 or 2: ")
	keyMode, _ := reader.ReadString('\n')
	keyMode = sanitize(keyMode)

	var key []byte
	switch keyMode {
	case "1":
		fmt.Print("Enter passphrase: ")
		passphrase, _ := reader.ReadString('\n')
		passphrase = sanitize(passphrase)
		key = deriveKeyFromString(passphrase)
	case "2":
		fmt.Print("Enter path to .key file: ")
		keyPath, _ := reader.ReadString('\n')
		keyPath = sanitize(keyPath)
		readKey, err := os.ReadFile(keyPath)
		if err != nil {
			fmt.Println("Error reading key:", err)
			return
		}
		if len(readKey) != 32 {
			fmt.Println("Invalid key length (32 bytes required)")
			return
		}
		key = readKey
	default:
		fmt.Println("Unknown key source")
		return
	}

	// If key was entered manually, offer to save it
	var saveKey string
	if keyMode == "1" {
		fmt.Print("Save key to file? [y/n]: ")
		saveKey, _ = reader.ReadString('\n')
		saveKey = sanitize(saveKey)
		if saveKey == "y" {
			err := os.WriteFile(inputPath+".key", key, 0600)
			if err != nil {
				fmt.Println("Error saving key:", err)
				return
			}
			fmt.Println("Key saved to:", inputPath+".key")
		}
	}

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	var outputData []byte
	if mode == "1" {
		outputData, err = encrypt(inputData, key)
	} else {
		outputData, err = decrypt(inputData, key)
	}
	if err != nil {
		fmt.Println("Processing error:", err)
		return
	}

	err = os.WriteFile(outputPath, outputData, 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}

	fmt.Println(" Done!\nOutput file:", outputPath)
	if saveKey == "y" {
		fmt.Println("Key saved to:", inputPath+".key")
	}
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}

func sanitize(s string) string {
	return strings.TrimSpace(s)
}

func deriveKeyFromString(pass string) []byte {
	hash := sha256.Sum256([]byte(pass))
	return hash[:]
}

func encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	return append(nonce, ciphertext...), nil
}

func decrypt(data, key []byte) ([]byte, error) {
	if len(data) < nonceSize {
		return nil, errors.New("data too short")
	}
	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}
