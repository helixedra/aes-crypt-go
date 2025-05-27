package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("--------------------")
	fmt.Println("Key Generator v0.1.1")
	fmt.Println("--------------------")

	// Generate key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("Error generating key:", err)
		return
	}

	// Choose output format
	fmt.Print("\nChoose output format:\n1) .key file\n2) Base64\n3) Both\n\nYour choice (1-3): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		saveKeyFile(key)
	case "2":
		showBase64(key)
	case "3":
		saveKeyFile(key)
		showBase64(key)
	default:
		fmt.Println("Error: Invalid choice")
	}
	fmt.Println("\nPress Enter to exit...")
	reader.ReadString('\n')
}

func saveKeyFile(key []byte) {
	fmt.Print("\nEnter filename (default: secret.key): ")
	filename, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	filename = strings.TrimSpace(filename)

	if filename == "" {
		filename = "secret.key"
	}

	if !strings.HasSuffix(filename, ".key") {
		filename += ".key"
	}

	err := os.WriteFile(filename, key, 0600)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Printf("Key saved to %s\n", filename)
}

func showBase64(key []byte) {
	keyB64 := base64.StdEncoding.EncodeToString(key)
	fmt.Printf("\nBase64 key: %s\n", keyB64)
	fmt.Println("\nKeep this key secret! Copy it to a secure place.\nNow you can run Cryptor and use this key to encrypt/decrypt files.")
}
