package main

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

// Small formatting, parsing, and coercion helpers are grouped here.
func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func getenv(k, dv string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return dv
}
