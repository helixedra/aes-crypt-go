# FILE ENCRYPTION TOOL
# ===================

## KEY GENERATION
$ go run keygen.go

Generates:
- secret.key (32-byte binary key file)
- Base64 encoded key (for manual transfer)

## ENCRYPTION
$ go run cryptor.go --encrypt --key secret.key file.txt
OR
$ go run cryptor.go --encrypt --base64 "your_base64_key" file.txt

Creates: file.txt.enc

## DECRYPTION 
$ go run cryptor.go --decrypt --key secret.key file.txt.enc
OR
$ go run cryptor.go --decrypt --base64 "your_base64_key" file.txt.enc

Restores original file (removes .enc extension)

## SECURITY SPECS
- Algorithm: AES-256-GCM
- Key size: 32 bytes (256-bit)
- Secure random generation (crypto/rand)
- Authenticated encryption

## SAFETY TIPS
1. NEVER share .key files or base64 keys publicly
2. ALWAYS backup your keys - no recovery if lost
3. DELETE temporary key files after use
4. USE secure channels for key transfer

## EXAMPLE WORKFLOW
1. Generate key: go run keygen.go
2. Encrypt: go run cryptor.go --encrypt --key secret.txt document.pdf
3. Securely transfer:
   - document.pdf.enc (safe to share)
   - secret.key (keep private!)
4. Decrypt: go run cryptor.go --decrypt --key secret.key document.pdf.enc

# WARNING
# =======
# LOST KEY = LOST DATA
# Keep your keys secure and backed up!
