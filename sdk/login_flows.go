package sdk

import (
	"context"

	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/core/kugou"
	"github.com/lfhy/kugou-music-api/core/util"
)

const (
	loginT1 = "562a6f12a6e803453647d16a08f5f0c2ff7eee692cba2ab74cc4c8ab47fc467561a7c6b586ce7dc46a63613b246737c03a1dc8f8d162d8ce1d2c71893d19f1d4b797685a4c6d3d81341cbde65e488c4829a9b4d42ef2df470eb102979fa5adcdd9b4eecfea8b909ff7599abeb49867640f10c3c70fc444effca9d15db44a9a6c907731e2bb0f22cd9b3536380169995693e5f0e2424e3378097d3813186e3fe96bbe7023808a0981b4e2b6135a76faac"
	loginT2 = "31c4daf4cf480169ccea1cb7d4a209295865a9d2b788510301694db229b87807469ea0d41b4d4b9173c2151da7294aeebfc9738df154bbdf11a4e117bb5dff6a3af8ce5ce333e681c1f29a44038f27567d58992eb81283e080778ac77db1400fdf49b7cf7e26be2e5af4da7830cc3be4"
	loginT3 = "MCwwLDAsMCwwLDAsMCwwLDA="
)

var (
	liteT2Key = "fd14b35e3f81af3817a20ae7adae7020"
	liteT2Iv  = "17a20ae7adae7020"
	liteT1Key = "5e4ef500e9597fe004bd09a46d8add98"
	liteT1Iv  = "04bd09a46d8add98"

	refreshKey     = "90b8382a1bb4ccdcf063102053fd75b8"
	refreshIv      = "f063102053fd75b8"
	refreshLiteKey = "c24f74ca2820225badc01946dba4fdf7"
	refreshLiteIv  = "adc01946dba4fdf7"
)

type PasswordLoginRequest struct {
	Username string
	Password string
	Cookie   map[string]string
}

type CellphoneLoginRequest struct {
	Mobile string
	Code   string
	UserID string
	Cookie map[string]string
}

type TokenLoginRequest struct {
	Token  string
	UserID string
	Cookie map[string]string
}

type loginPKPayload struct {
	ClientTimeMS int64  `json:"clienttime_ms"`
	Key          string `json:"key"`
}

type tokenRefreshPayload struct {
	Dfid         string `json:"dfid"`
	P3           string `json:"p3"`
	Plat         int    `json:"plat"`
	T1           int    `json:"t1"`
	T2           int    `json:"t2"`
	T3           string `json:"t3"`
	PK           string `json:"pk"`
	Params       string `json:"params"`
	UserID       string `json:"userid"`
	ClientTimeMS int64  `json:"clienttime_ms"`
}

func (c *Client) LoginByPassword(ctx context.Context, req PasswordLoginRequest) (*Response, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}

	nowMs := time.Now().UnixMilli()
	encrypt, err := kugou.CryptoAesEncrypt(map[string]any{
		"pwd":           req.Password,
		"code":          "",
		"clienttime_ms": nowMs,
	}, nil)
	if err != nil {
		return nil, err
	}

	pubKey := kugou.PublicRASKey
	if c.isLite {
		pubKey = kugou.PublicLiteRASKey
	}
	pk, err := kugou.CryptoRSAEncryptRawHex(loginPKPayload{
		ClientTimeMS: nowMs,
		Key:          encrypt.Key,
	}, pubKey)
	if err != nil {
		return nil, err
	}

	data := map[string]any{
		"plat":          1,
		"support_multi": 1,
		"clienttime_ms": nowMs,
		"t1":            loginT1,
		"t2":            loginT2,
		"t3":            loginT3,
		"username":      req.Username,
		"params":        encrypt.Str,
		"pk":            strings.ToUpper(pk),
	}

	raw, err := c.doRequest(ctx, kugou.RequestConfig{
		Method:      "POST",
		URL:         "/v9/login_by_pwd",
		Data:        data,
		EncryptType: "android",
		Cookie:      cookies,
		Headers: map[string]string{
			"x-router": "login.user.kugou.com",
		},
	})
	return c.finalizeLoginResponse(raw, encrypt.Key, err)
}

func (c *Client) LoginByCellphone(ctx context.Context, req CellphoneLoginRequest) (*Response, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}

	nowMs := time.Now().UnixMilli()
	encrypt, err := kugou.CryptoAesEncrypt(map[string]any{"mobile": req.Mobile, "code": req.Code}, nil)
	if err != nil {
		return nil, err
	}

	pubKey := kugou.PublicRASKey
	if c.isLite {
		pubKey = kugou.PublicLiteRASKey
	}
	pk, err := kugou.CryptoRSAEncryptRawHex(loginPKPayload{
		ClientTimeMS: nowMs,
		Key:          encrypt.Key,
	}, pubKey)
	if err != nil {
		return nil, err
	}

	t2, _ := kugou.CryptoAesEncrypt(fmt.Sprintf("%s|0f607264fc6318a92b9e13c65db7cd3c|%s|%s|%d", cookies["KUGOU_API_GUID"], cookies["KUGOU_API_MAC"], cookies["KUGOU_API_DEV"], nowMs), &kugou.AesOpt{Key: liteT2Key, IV: liteT2Iv})
	t1, _ := kugou.CryptoAesEncrypt(fmt.Sprintf("|%d", nowMs), &kugou.AesOpt{Key: liteT1Key, IV: liteT1Iv})

	data := map[string]any{
		"plat":          1,
		"support_multi": 1,
		"t1":            0,
		"t2":            0,
		"clienttime_ms": nowMs,
		"mobile":        maskMobile(req.Mobile),
		"key":           signParamsKey(strconv.FormatInt(nowMs, 10), c.isLite),
		"pk":            strings.ToUpper(pk),
		"params":        encrypt.Str,
	}
	if strings.TrimSpace(req.UserID) != "" {
		data["userid"] = req.UserID
	}
	if c.isLite {
		data["t1"] = t1.Str
		data["t2"] = t2.Str
		data["dfid"] = firstNonEmpty(cookies["dfid"], util.RandomString(24))
		data["dev"] = cookies["KUGOU_API_DEV"]
		data["gitversion"] = "5f0b7c4"
	} else {
		data["t3"] = loginT3
	}

	raw, err := c.doRequest(ctx, kugou.RequestConfig{
		Method:      "POST",
		BaseURL:     "https://loginserviceretry.kugou.com",
		URL:         "/v7/login_by_verifycode",
		Data:        data,
		EncryptType: "android",
		Cookie:      cookies,
		Headers: map[string]string{
			"support-calm": "1",
			"User-Agent":   "Android16-1070-11440-130-0-LOGIN-wifi",
		},
	})
	return c.finalizeLoginResponse(raw, encrypt.Key, err)
}

func (c *Client) LoginByToken(ctx context.Context, req TokenLoginRequest) (*Response, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}

	nowMs := time.Now().UnixMilli()
	token := firstNonEmpty(req.Token, cookies["token"])
	userid := firstNonEmpty(req.UserID, cookies["userid"], "0")

	encKey := refreshKey
	encIv := refreshIv
	if c.isLite {
		encKey = refreshLiteKey
		encIv = refreshLiteIv
	}

	p3, err := kugou.CryptoAesEncrypt(map[string]any{"clienttime": nowMs / 1000, "token": token}, &kugou.AesOpt{Key: encKey, IV: encIv})
	if err != nil {
		return nil, err
	}
	encryptParams, err := kugou.CryptoAesEncrypt(map[string]any{}, nil)
	if err != nil {
		return nil, err
	}

	pubKey := kugou.PublicRASKey
	if c.isLite {
		pubKey = kugou.PublicLiteRASKey
	}
	pk, err := kugou.CryptoRSAEncryptRawHex(loginPKPayload{
		ClientTimeMS: nowMs,
		Key:          encryptParams.Key,
	}, pubKey)
	if err != nil {
		return nil, err
	}

	data := tokenRefreshPayload{
		Dfid:         firstNonEmpty(cookies["dfid"], "-"),
		P3:           p3.Str,
		Plat:         1,
		T1:           0,
		T2:           0,
		T3:           loginT3,
		PK:           pk,
		Params:       encryptParams.Str,
		UserID:       userid,
		ClientTimeMS: nowMs,
	}

	raw, err := c.doRequest(ctx, kugou.RequestConfig{
		Method:      "POST",
		BaseURL:     "http://login.user.kugou.com",
		URL:         ternaryString(c.isLite, "/v4/login_by_token", "/v5/login_by_token"),
		Data:        data,
		EncryptType: "android",
		Cookie:      cookies,
		Headers: map[string]string{
			"x-router": "login.user.kugou.com",
		},
	})
	return c.finalizeLoginResponse(raw, encryptParams.Key, err)
}
