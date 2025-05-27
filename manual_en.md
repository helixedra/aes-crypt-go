# User Manual: Key Generator and File Cryptor

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Key Generation](#key-generation)
   - [Saving to a File](#saving-to-a-file)
   - [Displaying as Base64](#displaying-as-base64)
4. [File Encryption](#file-encryption)
5. [File Decryption](#file-decryption)
6. [Security](#security)
7. [Troubleshooting](#troubleshooting)

## Introduction

Key Generator and File Cryptor are utilities for secure file encryption and decryption using the AES-256 algorithm.

## Installation

1. Ensure you have Go (version 1.16 or higher) installed
2. Clone the repository or copy the files
3. Build the executable files:
   ```bash
   go build keygen.go
   go build cryptor.go
   ```

## Key Generation

Run the `keygen` program:

```bash
./keygen
```

### Saving to a File

1. Select option `1) .key file`
2. Enter a filename (default: `secret.key`)
3. The key will be saved to the specified file

### Displaying as Base64

1. Select option `2) Base64`
2. The key will be displayed in the console in Base64 format
3. Copy and save this key in a secure location

## File Encryption

1. Run `cryptor`
2. Select `1) Encrypt`
3. Enter the path to the file you want to encrypt
4. Choose the key source:
   - `1) .key file` - specify the path to the key file
   - `2) Base64 key` - paste the key in Base64 format
5. The encrypted file will be saved with a `.enc` extension

## File Decryption

1. Run `cryptor`
2. Select `2) Decrypt`
3. Enter the path to the encrypted file (with `.enc` extension)
4. Provide the key (from file or enter as Base64)
5. The decrypted file will be saved with the original filename (without `.enc`)

## Security

- Always store your keys in a secure location
- Never transmit keys over unsecured channels
- It's recommended to use `.key` files instead of Base64, as they don't remain in command line history
- Keep backup copies of your encryption keys
- Losing the key means losing access to your encrypted data

## Troubleshooting

### "Invalid key length" Error

Make sure that:

- The key file contains exactly 32 bytes
- The Base64 key is properly encoded and decodes to 32 bytes

### File Read Error

Check:

- If the specified file exists
- If you have read/write permissions
- If the file path is correct

### Decryption Error

- Ensure you're using the same key that was used for encryption
- Check if the encrypted file is not corrupted

### General Tips

- The program will overwrite files without warning
- One key can be used for multiple files, but using different keys for different files is more secure
- The `.enc` extension is added automatically during encryption

---

For additional assistance, please contact support.
