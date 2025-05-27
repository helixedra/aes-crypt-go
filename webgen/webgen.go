package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	
	http.HandleFunc("/", handleQR)
	fmt.Println("Server started at:")
	fmt.Println(" - http://localhost:8080")
	fmt.Println(" - http://" + getLocalIP() + ":8080")

	http.ListenAndServe("0.0.0.0:8080", nil)
}

func handleQR(w http.ResponseWriter, r *http.Request) {
	qrParam := r.URL.RawQuery
	m, _ := url.ParseQuery(qrParam)
	qr := m.Get("qr")
	qr = strings.ReplaceAll(qr, " ", "+")
	qr = strings.TrimSpace(qr)
	if qr == "" {
		http.Error(w, "Parameter 'qr' is required", http.StatusBadRequest)
		return
	}

	// decode Base64 (get seed)
	seed, err := base64.StdEncoding.DecodeString(qr)
	if err != nil || len(seed) != 32 {
		http.Error(w, "Invalid base64 or seed length (must be 32 bytes)", http.StatusBadRequest)
		return
	}

	// calculate SHA-256 from seed
	key := sha256.Sum256(seed) // key — это [32]byte
	keyBytes := key[:]         // convert to []byte

	// save key (SHA-256) to temporary file
	tmpDir := os.TempDir()
	keyFile := filepath.Join(tmpDir, "download.key")
	err = os.WriteFile(keyFile, keyBytes, 0644)
	if err != nil {
		http.Error(w, "Error writing key file", http.StatusInternalServerError)
		return
	}

	// send download link to user
	fmt.Fprintf(w, `
		<html>
		<head><title>🔐 Download key</title></head>
		<body>
			<h2>Your key is ready!</h2>
			<p>Download and use for decrypt files</p>
			<p><a href="/download.key" download>📥 Download key</a></p>
		</body>
		</html>
	`)
}

// download handler
func init() {
	http.HandleFunc("/download.key", func(w http.ResponseWriter, r *http.Request) {
		tmpFile := filepath.Join(os.TempDir(), "download.key")
		f, err := os.Open(tmpFile)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\"mykey.key\"")
		io.Copy(w, f)
	})
}

// get Local IP
func getLocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "unknown"
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			if strings.HasPrefix(ip.String(), "192.168.") || strings.HasPrefix(ip.String(), "10.") {
				return ip.String()
			}
		}
	}
	return "unknown"
}