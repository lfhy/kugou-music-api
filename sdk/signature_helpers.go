package sdk

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/lfhy/kugou-music-api/core/config"
	"github.com/lfhy/kugou-music-api/core/util"
)

// Shared signature helpers are isolated so the per-feature files stay compact.
func (c *Client) commentCommon(ctx context.Context, route string, params map[string]any, cookie map[string]string) (*Response, error) {
	resp, err := c.Call(ctx, route, Request{
		Method:      "POST",
		Params:      params,
		Cookie:      cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func audioRelatedSort(v any) int {
	s := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", v)))
	switch s {
	case "hot":
		return 2
	case "new":
		return 3
	default:
		return 1
	}
}

func md5SortedWithKey(m map[string]any, key string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+sigVal(m[k]))
	}
	return util.MD5Hex(key + strings.Join(parts, "") + key)
}

func md5AmpersandSign(m map[string]any) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+sigVal(m[k]))
	}
	sum := util.MD5Hex(strings.Join(pairs, "&") + "*s&iN#G70*")
	if len(sum) < 24 {
		return sum
	}
	return sum[8:24]
}

func sigVal(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

type playlistEncryptResult struct {
	Key          string
	CipherBase64 string
}

func playlistAesEncrypt(data any) (playlistEncryptResult, error) {
	plain, err := json.Marshal(data)
	if err != nil {
		return playlistEncryptResult{}, err
	}
	keySeed := strings.ToLower(util.RandomString(6))
	md5 := util.MD5Hex(keySeed)
	key := []byte(md5[:16])
	iv := []byte(md5[16:32])

	block, err := aes.NewCipher(key)
	if err != nil {
		return playlistEncryptResult{}, err
	}
	padded := pkcs7Pad(plain, aes.BlockSize)
	out := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out, padded)

	return playlistEncryptResult{Key: keySeed, CipherBase64: base64.StdEncoding.EncodeToString(out)}, nil
}

func playlistAesDecryptFromRaw(raw []byte, keySeed string) (map[string]any, error) {
	md5 := util.MD5Hex(keySeed)
	key := []byte(md5[:16])
	iv := []byte(md5[16:32])

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(raw))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(out, raw)
	plain, err := pkcs7Unpad(out, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(plain, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func toInt(v any, def int) int {
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "" || s == "<nil>" {
		return def
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return int(f)
	}
	return def
}

func toBool(v any, def bool) bool {
	s := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", v)))
	switch s {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	case "", "<nil>":
		return def
	default:
		return def
	}
}

func boolPtr(v bool) *bool { return &v }

func splitCSV(s string) []string {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil
	}
	parts := strings.Split(t, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func signCloudKey(hash, pid string) string {
	return util.MD5Hex("musicclound" + hash + pid + "ebd1ac3134c880bda6a2194537843caa0162e2e7")
}

func currentAppid(isLite bool) string {
	a, _ := config.PlatformConfig(isLite)
	return a
}

func currentClientVer(isLite bool) string {
	_, v := config.PlatformConfig(isLite)
	return v
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	pad := blockSize - len(data)%blockSize
	out := make([]byte, len(data)+pad)
	copy(out, data)
	for i := len(data); i < len(out); i++ {
		out[i] = byte(pad)
	}
	return out
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded data")
	}
	pad := int(data[len(data)-1])
	if pad <= 0 || pad > blockSize || pad > len(data) {
		return nil, fmt.Errorf("invalid pad size")
	}
	for i := len(data) - pad; i < len(data); i++ {
		if int(data[i]) != pad {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:len(data)-pad], nil
}
