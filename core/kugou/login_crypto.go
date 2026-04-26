package kugou

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
	"strings"

	"github.com/lfhy/kugou-music-api/core/util"
)

const (
	PublicRASKey     = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDIAG7QOELSYoIJvTFJhMpe1s/gbjDJX51HBNnEl5HXqTW6lQ7LC8jr9fWZTwusknp+sVGzwd40MwP6U5yDE27M/X1+UR4tvOGOqp94TJtQ1EPnWGWXngpeIW5GxoQGao1rmYWAu6oi1z9XkChrsUdC6DJE5E221wf/4WLFxwAtRQIDAQAB\n-----END PUBLIC KEY-----"
	PublicLiteRASKey = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDECi0Np2UR87scwrvTr72L6oO01rBbbBPriSDFPxr3Z5syug0O24QyQO8bg27+0+4kBzTBTBOZ/WWU0WryL1JSXRTXLgFVxtzIY41Pe7lPOgsfTCn5kZcvKhYKJesKnnJDNr5/abvTGf+rHG3YRwsCHcQ08/q6ifSioBszvb3QiwIDAQAB\n-----END PUBLIC KEY-----"
)

type AesOpt struct {
	Key string
	IV  string
}

type AesEncryptResult struct {
	Str string
	Key string
}

func CryptoAesEncrypt(data any, opt *AesOpt) (AesEncryptResult, error) {
	plain, err := normalizeDataString(data)
	if err != nil {
		return AesEncryptResult{}, err
	}

	var key string
	var iv string
	out := AesEncryptResult{}

	if opt != nil && opt.Key != "" && opt.IV != "" {
		key = opt.Key
		iv = opt.IV
	} else {
		tempKey := strings.ToLower(util.RandomString(16))
		md5Key := md5HexLogin(tempKey)
		key = md5Key
		iv = md5Key[len(md5Key)-16:]
		out.Key = tempKey
	}

	cipherHex, err := aesCBCEncryptHex([]byte(plain), []byte(key), []byte(iv))
	if err != nil {
		return AesEncryptResult{}, err
	}
	out.Str = cipherHex
	return out, nil
}

func CryptoAesDecryptHex(dataHex, key, iv string) (any, error) {
	if iv == "" {
		key = md5HexLogin(key)
		iv = key[len(key)-16:]
	}

	plain, err := aesCBCDecryptHex(dataHex, []byte(key), []byte(iv))
	if err != nil {
		return nil, err
	}

	var obj map[string]any
	if json.Unmarshal(plain, &obj) == nil {
		return obj, nil
	}
	return string(plain), nil
}

func CryptoRSAEncryptRawHex(data any, publicKeyPEM string) (string, error) {
	buf, err := normalizeDataBytes(data)
	if err != nil {
		return "", err
	}

	pub, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return "", err
	}

	k := (pub.N.BitLen() + 7) / 8
	if len(buf) > k {
		return "", errors.New("data length exceeds key size")
	}
	if len(buf) < k {
		padded := make([]byte, k)
		copy(padded, buf)
		buf = padded
	}

	m := new(big.Int).SetBytes(buf)
	e := big.NewInt(int64(pub.E))
	c := new(big.Int).Exp(m, e, pub.N)
	h := c.Text(16)
	if len(h) < k*2 {
		h = strings.Repeat("0", k*2-len(h)) + h
	}
	return h, nil
}

func CryptoRSAEncryptPKCS1Hex(data any, publicKeyPEM string) (string, error) {
	buf, err := normalizeDataBytes(data)
	if err != nil {
		return "", err
	}
	pub, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return "", err
	}
	encrypted, err := rsa.EncryptPKCS1v15(nil, pub, buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(encrypted), nil
}

func normalizeDataString(data any) (string, error) {
	switch v := data.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}

func normalizeDataBytes(data any) ([]byte, error) {
	s, err := normalizeDataString(data)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

func parsePublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("invalid public key pem")
	}
	k, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := k.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not rsa public key")
	}
	return pub, nil
}

func md5HexLogin(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func aesCBCEncryptHex(plain, key, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(iv) != aes.BlockSize {
		return "", errors.New("invalid iv length")
	}
	plain = pkcs7Pad(plain, aes.BlockSize)
	out := make([]byte, len(plain))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out, plain)
	return hex.EncodeToString(out), nil
}

func aesCBCDecryptHex(dataHex string, key, iv []byte) ([]byte, error) {
	raw, err := hex.DecodeString(strings.TrimSpace(dataHex))
	if err != nil {
		return nil, err
	}
	if len(raw)%aes.BlockSize != 0 {
		return nil, errors.New("invalid ciphertext size")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != aes.BlockSize {
		return nil, errors.New("invalid iv length")
	}
	out := make([]byte, len(raw))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(out, raw)
	return pkcs7Unpad(out, aes.BlockSize)
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
		return nil, errors.New("invalid padded data")
	}
	pad := int(data[len(data)-1])
	if pad <= 0 || pad > blockSize || pad > len(data) {
		return nil, errors.New("invalid pad size")
	}
	for i := len(data) - pad; i < len(data); i++ {
		if int(data[i]) != pad {
			return nil, errors.New("invalid padding")
		}
	}
	return data[:len(data)-pad], nil
}
