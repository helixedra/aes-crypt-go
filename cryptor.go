package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const nonceSize = 12

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("-------------------")
	fmt.Println("File Cryptor v0.1.1")
	fmt.Println("-------------------")

	// Choose mode
	fmt.Print("\nChoose mode:\n1) Encrypt\n2) Decrypt\n\nYour choice (1-2): ")
	mode, _ := reader.ReadString('\n')
	mode = strings.TrimSpace(mode)

	// Choose file
	fmt.Print("\nEnter file path: ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Choose key source
	fmt.Print("\nChoose key source:\n1) .key file\n2) Base64 key\n\nYour choice (1-2): ")
	keySource, _ := reader.ReadString('\n')
	keySource = strings.TrimSpace(keySource)

	var key []byte

	switch keySource {
	case "1":
		fmt.Print("Enter .key file path: ")
		keyPath, _ := reader.ReadString('\n')
		keyPath = strings.TrimSpace(keyPath)

		key, err = os.ReadFile(keyPath)
		if err != nil {
			fmt.Println("Error reading key file:", err)
			return
		}

		if len(key) != 32 {
			fmt.Println("Invalid key length (must be 32 bytes)")
			return
		}

	case "2":
		fmt.Print("Enter Base64 key: ")
		keyB64, _ := reader.ReadString('\n')
		keyB64 = strings.TrimSpace(keyB64)

		key, err = base64.StdEncoding.DecodeString(keyB64)
		if err != nil {
			fmt.Println("Error decoding Base64 key:", err)
			return
		}

		if len(key) != 32 {
			fmt.Println("Invalid key length (must decode to 32 bytes)")
			return
		}

	default:
		fmt.Println("Error: Invalid choice")
		return
	}

	// Process file
	var result []byte
	outputPath := getOutputPath(filePath, mode == "1")

	switch mode {
	case "1": // Encrypt
		result, err = encrypt(key, fileData)
		if err != nil {
			fmt.Println("Encryption failed:", err)
			return
		}
	case "2": // Decrypt
		result, err = decrypt(key, fileData)
		if err != nil {
			fmt.Println("Decryption failed:", err)
			return
		}
	default:
		fmt.Println("Error: Invalid mode selected")
		return
	}

	// Save result
	err = os.WriteFile(outputPath, result, 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Printf("\nSuccess! File saved to:\n%s\n", outputPath)
	fmt.Println("\nPress Enter to exit...")
	reader.ReadString('\n')
}

func getOutputPath(input string, isEncrypt bool) string {
	if isEncrypt {
		return input + ".enc"
	}

	if strings.HasSuffix(input, ".enc") {
		return strings.TrimSuffix(input, ".enc")
	}

	return input + ".dec"
}

func encrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
