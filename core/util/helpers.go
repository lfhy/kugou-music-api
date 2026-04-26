package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/big"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	alphaNumUpper = []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	cookieCleaner = regexp.MustCompile(`\\s*(Domain|domain|path|expires)=[^(;|$)]+;*`)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func MD5Hex(text string) string {
	sum := md5.Sum([]byte(text))
	return hex.EncodeToString(sum[:])
}

func RandomString(n int) string {
	if n <= 0 {
		n = 16
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = alphaNumUpper[rand.Intn(len(alphaNumUpper))]
	}
	return string(b)
}

func ParseCookieString(cookie string) string {
	t := cookieCleaner.ReplaceAllString(cookie, "")
	return strings.ReplaceAll(t, ";HttpOnly", "")
}

func CookieToMap(cookie string) map[string]string {
	out := map[string]string{}
	if cookie == "" {
		return out
	}
	parts := strings.Split(cookie, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key, err1 := url.QueryUnescape(strings.TrimSpace(kv[0]))
		val, err2 := url.QueryUnescape(strings.TrimSpace(kv[1]))
		if err1 != nil {
			key = strings.TrimSpace(kv[0])
		}
		if err2 != nil {
			val = strings.TrimSpace(kv[1])
		}
		out[key] = val
	}
	return out
}

func CalculateMid(seed string) string {
	digest := MD5Hex(seed)
	n := new(big.Int)
	n.SetString(digest, 16)
	return n.String()
}

func BoolFromString(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}
