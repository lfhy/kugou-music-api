package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func saveQRCodePNG(key, dataURI string) (string, error) {
	key = sanitizeQRCodeFileKey(key)
	if key == "" {
		return "", nil
	}
	dataURI = strings.TrimSpace(dataURI)
	if dataURI == "" {
		return "", nil
	}

	raw := dataURI
	if idx := strings.Index(raw, ","); idx >= 0 {
		raw = raw[idx+1:]
	}

	png, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return "", fmt.Errorf("decode png base64: %w", err)
	}

	path := filepath.Join(os.TempDir(), "kugou-login-qrcode-"+key+".png")
	if err := os.WriteFile(path, png, 0o600); err != nil {
		return "", fmt.Errorf("write png: %w", err)
	}
	return path, nil
}

func sanitizeQRCodeFileKey(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	var out strings.Builder
	for _, r := range key {
		switch {
		case r >= 'a' && r <= 'z':
			out.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			out.WriteRune(r)
		case r >= '0' && r <= '9':
			out.WriteRune(r)
		case r == '-' || r == '_':
			out.WriteRune(r)
		}
		if out.Len() >= 48 {
			break
		}
	}
	return out.String()
}
