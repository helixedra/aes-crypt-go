package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
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
		base64Str   = flag.String("base64", "", "Base64-encoded seed string")
	)
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: no input file specified")
		fmt.Println("Usage: ./cryptor --encrypt --key my.key secret.txt")
		return
	}
	inputPath := args[0]

	var mode string
	if *encryptFlag {
		mode = "encrypt"
	} else if *decryptFlag {
		mode = "decrypt"
	} else {
		mode = ask("Select mode ([1] encrypt / [2] decrypt): ")
		if mode == "1" {
			mode = "encrypt"
		} else if mode == "2" {
			mode = "decrypt"
		} else {
			fmt.Println("Error: unknown mode")
			return
		}
	}

	outputPath := getOutputPath(inputPath, mode)

	var key []byte
	if *keyPath != "" {
		keyData, err := os.ReadFile(*keyPath)
		if err != nil || len(keyData) != 32 {
			fmt.Println("Error: invalid key file or length (32 bytes)")
			return
		}
		key = keyData
	} else if *base64Str != "" {
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(*base64Str))
		if err != nil {
			fmt.Println("Error: invalid base64 string")
			return
		}
		// Always use SHA-256 of decoded data
		hashed := sha256.Sum256(decoded)
		key = hashed[:]
	} else {
		// ask user interactively
		source := ask("Select source of key ([1] passphrase / [2] key file / [3] base64 seed): ")
		switch source {
		case "1":
			phrase := ask("Enter passphrase: ")
			key = deriveKeyFromString(phrase)
		case "2":
			path := ask("Enter path to .key file: ")
			keyData, err := os.ReadFile(path)
			if err != nil || len(keyData) != 32 {
				fmt.Println("Error: invalid key file or length (32 bytes)")
				return
			}
			key = keyData
		case "3":
			b64 := ask("Enter base64 seed: ")
			decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(b64))
			if err != nil {
				fmt.Println("Error: invalid base64 string")
				return
			}
			hashed := sha256.Sum256(decoded)
			key = hashed[:]
		default:
			fmt.Println("Error: unknown source of key")
			return
		}
	}

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Error: reading file:", err)
		return
	}

	var outputData []byte
	if mode == "encrypt" {
		outputData, err = encrypt(inputData, key)
	} else {
		outputData, err = decrypt(inputData, key)
	}
	if err != nil {
		fmt.Println("Error: processing file:", err)
		return
	}

	err = os.WriteFile(outputPath, outputData, 0644)
	if err != nil {
		fmt.Println("Error: writing file:", err)
		return
	}

	fmt.Println("Done. File:", outputPath)
}

func ask(prompt string) string {
	fmt.Print(prompt)
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	return strings.TrimSpace(scan.Text())
}

func getOutputPath(input, mode string) string {
	if mode == "encrypt" {
		return input + ".enc"
	}
	if strings.HasSuffix(input, ".enc") {
		return strings.TrimSuffix(input, ".enc")
	}
	return input + ".dec"
}

func deriveKeyFromString(pass string) []byte {
	h := sha256.Sum256([]byte(pass))
	return h[:]
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