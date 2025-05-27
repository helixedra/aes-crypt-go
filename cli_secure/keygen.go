package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	// 1. Generate random 32-byte key
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		fmt.Println("Error: generating key:", err)
		return
	}

	// 2. Save to binary file
	keyFile := "secret.key"
	if err := os.WriteFile(keyFile, key, 0600); err != nil {
		fmt.Println("Error: writing key to file:", err)
		return
	}

	// 3. Output in Base64 (for alternative usage)
	keyB64 := base64.StdEncoding.EncodeToString(key)
	fmt.Printf("Key saved to file: %s\n", keyFile)
	fmt.Printf("Base64-key: %s\n", keyB64)
}