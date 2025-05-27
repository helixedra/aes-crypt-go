package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/skip2/go-qrcode"
)

const (
	randomSeedLength = 32
	keyLength        = 32
	baseName         = "secret"
	qrPrefixURL      = "https://mycode.com?qr="
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// flag to generate QR
	autoQR := len(os.Args) > 1 && os.Args[1] == "--qr"

	// generate random key
	seed := make([]byte, randomSeedLength)
	_, err := rand.Read(seed)
	if err != nil {
		fmt.Println("Error generating seed:", err)
		return
	}

	// encode to Base64
	keyB64 := base64.StdEncoding.EncodeToString(seed)

	// ask for output key name
	fmt.Print("Enter filename for key: ")
	var keyName string
	keyName, _ = reader.ReadString('\n')
	keyName = sanitize(keyName)
	if keyName == "" {
		keyName = baseName
	}

	// Get SHA-256 from seed
	finalKey := sha256.Sum256(seed)
	keyBytes := finalKey[:]

	// Save binary key to file
	keyFile := keyName + ".key"
	err = os.WriteFile(keyFile, keyBytes, 0644)
	if err != nil {
		fmt.Println("Error saving key:", err)
		return
	}


	// generate QR code by flag or question
	if autoQR || askYesNo("Create QR code with link? [y/n]: ", reader) {
		qrContent := qrPrefixURL + keyB64
		qrFile := keyName + ".png"
		err := qrcode.WriteFile(qrContent, qrcode.Medium, 256, qrFile)
		if err != nil {
			fmt.Println("Error generating QR code:", err)
			return
		}
		fmt.Println("QR code saved to file:", qrFile)
	}

	fmt.Println("Key saved to file:", keyFile)
fmt.Println("")
fmt.Println("Key*:", keyB64)
fmt.Println("\n*you can use this key to decrypt files as passphrase or \ngenerate .key file by online key generator")
fmt.Println("")
fmt.Println("Press Enter to exit...")
fmt.Scanln()
}

func sanitize(s string) string {
	return strings.TrimSpace(s)
}

func askYesNo(prompt string, reader *bufio.Reader) bool {
	fmt.Print(prompt)
	var response string
	response, _ = reader.ReadString('\n')
	response = sanitize(strings.ToLower(response))
	return response == "y" || response == "yes"
}