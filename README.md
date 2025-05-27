# AES File Cryptor

AES File Cryptor is a set of command-line utilities for secure file encryption and decryption using the AES-256-GCM algorithm. The project consists of two main tools:

- **Key Generator (`keygen`)** - Generates secure cryptographic keys
- **File Cryptor (`cryptor`)** - Encrypts and decrypts files using AES-256

## Features

- 🔒 Strong AES-256-GCM encryption
- 🔑 Secure key generation
- 📁 File encryption with `.enc` extension
- 🔄 Simple and intuitive CLI interface
- 🔐 Support for both file-based and Base64-encoded keys
- 🚀 Written in Go for cross-platform compatibility

## Prerequisites

- Go 1.16 or higher
- Basic command-line knowledge

## Installation

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd AES-crypt
   ```

2. Build the executables:
   ```bash
   go build keygen.go
   go build cryptor.go
   ```

## Quick Start

### 1. Generate a Key

```bash
# Generate a new key
./keygen

# Choose option 1 to save as .key file
# Or option 2 to display as Base64
# Or option 3 for both
```

### 2. Encrypt a File

```bash
./cryptor
# Choose 1) Encrypt
# Enter file path
# Select key source (file or Base64)
```

### 3. Decrypt a File

```bash
./cryptor
# Choose 2) Decrypt
# Enter .enc file path
# Provide the key
```

## Security Notes

- Always keep your encryption keys secure and never share them
- The `.key` files contain raw binary data - handle with care
- For maximum security, use file-based keys as they don't leave traces in shell history
- The encrypted files are saved with `.enc` extension by default

## License

This project is open source and available under the [MIT License](LICENSE).
