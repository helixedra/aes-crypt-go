package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const nonceSize = 12

func main() {
	var (
		encryptFlag = flag.Bool("encrypt", false, "Encrypt file")
		decryptFlag = flag.Bool("decrypt", false, "Decrypt file")
		keyPath     = flag.String("key", "", "Path to .key file")
		base64Key   = flag.String("base64", "", "Base64-encoded key")
	)
	flag.Parse()

	// Check mode (encryption/decryption)
	mode := ""
	switch {
	case *encryptFlag:
		mode = "encrypt"
	case *decryptFlag:
		mode = "decrypt"
	default:
		fmt.Println("Error: Specify mode: --encrypt or --decrypt")
		return
	}

	// Check file path
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Error: Specify file path (e.g. ./cryptor --encrypt --key secret.key data.txt)")
		return
	}
	inputPath := args[0]

	// Get key
	var key []byte
	if *keyPath != "" {
		// Load from .key file
		var err error
		key, err = os.ReadFile(*keyPath)
		if err != nil || len(key) != 32 {
			fmt.Println("Error reading .key file (must be 32 bytes)")
			return
		}
	} else if *base64Key != "" {
		// Decode base64
		var err error
		key, err = base64.StdEncoding.DecodeString(strings.TrimSpace(*base64Key))
		if err != nil || len(key) != 32 {
			fmt.Println("Error: Invalid base64 key (must decode to 32 bytes)")
			return
		}
	} else {
		fmt.Println("Error: Specify key (--key or --base64)")
		return
	}

	// Process file
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Error: reading file:", err)
		return
	}

	var outputData []byte
	switch mode {
	case "encrypt":
		outputData, err = encrypt(key, inputData)
	case "decrypt":
		outputData, err = decrypt(key, inputData)
	}
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Save result
	outputPath := inputPath + ".enc"
	if mode == "decrypt" {
		if strings.HasSuffix(inputPath, ".enc") {
			outputPath = strings.TrimSuffix(inputPath, ".enc")
		} else {
			outputPath = inputPath + ".dec"
		}
	}

	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		fmt.Println("Error: writing file:", err)
		return
	}

	fmt.Println("Done. Result in", outputPath)
}

// Encrypt (AES-256-GCM)
func encrypt(key, plaintext []byte) ([]byte, error) {
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

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt (AES-256-GCM)
func decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("too short ciphertext")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}