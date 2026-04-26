package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/core/config"
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
	pk, err := kugou.CryptoRSAEncryptRawHex(map[string]any{"clienttime_ms": nowMs, "key": encrypt.Key}, pubKey)
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

	raw, err := c.core.CreateRequest(ctx, kugou.RequestConfig{
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
	pk, err := kugou.CryptoRSAEncryptRawHex(map[string]any{"clienttime_ms": nowMs, "key": encrypt.Key}, pubKey)
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

	raw, err := c.core.CreateRequest(ctx, kugou.RequestConfig{
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
	pk, err := kugou.CryptoRSAEncryptRawHex(map[string]any{"clienttime_ms": nowMs, "key": encryptParams.Key}, pubKey)
	if err != nil {
		return nil, err
	}

	t2, _ := kugou.CryptoAesEncrypt(fmt.Sprintf("%s|0f607264fc6318a92b9e13c65db7cd3c|%s|%s|%d", cookies["KUGOU_API_GUID"], cookies["KUGOU_API_MAC"], cookies["KUGOU_API_DEV"], nowMs), &kugou.AesOpt{Key: liteT2Key, IV: liteT2Iv})
	t1Source := fmt.Sprintf("|%d", nowMs)
	if cookies["t1"] != "" {
		t1Source = fmt.Sprintf("%s|%d", cookies["t1"], nowMs)
	}
	t1, _ := kugou.CryptoAesEncrypt(t1Source, &kugou.AesOpt{Key: liteT1Key, IV: liteT1Iv})

	data := map[string]any{
		"dfid":          firstNonEmpty(cookies["dfid"], "-"),
		"p3":            p3.Str,
		"plat":          1,
		"t1":            0,
		"t2":            0,
		"t3":            loginT3,
		"pk":            pk,
		"params":        encryptParams.Str,
		"userid":        userid,
		"clienttime_ms": nowMs,
	}
	if c.isLite {
		data["t1"] = t1.Str
		data["t2"] = t2.Str
		data["dev"] = cookies["KUGOU_API_DEV"]
	}

	raw, err := c.core.CreateRequest(ctx, kugou.RequestConfig{
		Method:      "POST",
		BaseURL:     "http://login.user.kugou.com",
		URL:         "/v5/login_by_token",
		Data:        data,
		EncryptType: "android",
		Cookie:      cookies,
	})
	return c.finalizeLoginResponse(raw, encryptParams.Key, err)
}

func (c *Client) finalizeLoginResponse(raw kugou.Response, aesKey string, reqErr error) (*Response, error) {
	if len(raw.Cookie) > 0 {
		c.updateCookiePool(raw.Cookie)
	}

	out := &Response{Status: raw.Status, RawBody: raw.Body, Headers: raw.Headers, Cookie: raw.Cookie}
	var body map[string]any
	if json.Unmarshal(raw.Body, &body) == nil {
		if status, _ := body["status"].(float64); int(status) == 1 {
			if dataMap, ok := body["data"].(map[string]any); ok {
				if secu, ok := dataMap["secu_params"].(string); ok && strings.TrimSpace(secu) != "" {
					decoded, err := kugou.CryptoAesDecryptHex(secu, aesKey, "")
					if err == nil {
						switch t := decoded.(type) {
						case map[string]any:
							for k, v := range t {
								dataMap[k] = v
								out.Cookie = append(out.Cookie, fmt.Sprintf("%s=%v", k, v))
							}
						default:
							dataMap["token"] = t
							out.Cookie = append(out.Cookie, fmt.Sprintf("token=%v", t))
						}
					}
				}

				if v, ok := dataMap["t1"]; ok {
					out.Cookie = append(out.Cookie, fmt.Sprintf("t1=%v", v))
				}
				// Always overwrite auth cookies after successful login to avoid stale session values.
				out.Cookie = append(out.Cookie, "token="+asString(firstAny(dataMap["token"], "")))
				out.Cookie = append(out.Cookie, "userid="+asIntString(firstAny(dataMap["userid"], 0)))
				out.Cookie = append(out.Cookie, "vip_type="+asIntString(firstAny(dataMap["vip_type"], 0)))
				out.Cookie = append(out.Cookie, "vip_token="+asString(firstAny(dataMap["vip_token"], "")))
			}
		}

		clean := dedupSetCookie(out.Cookie)
		out.Cookie = clean
		c.updateCookiePool(clean)

		if b, err := json.Marshal(body); err == nil {
			out.RawBody = b
		}
		out.Body = body
	}

	if reqErr != nil {
		return out, reqErr
	}
	return out, nil
}

func dedupSetCookie(input []string) []string {
	m := map[string]string{}
	for _, x := range input {
		x = strings.TrimSpace(x)
		if x == "" {
			continue
		}
		kv := strings.SplitN(strings.Split(x, ";")[0], "=", 2)
		if len(kv) != 2 {
			continue
		}
		m[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	out := make([]string, 0, len(m))
	for k, v := range m {
		out = append(out, k+"="+v)
	}
	return out
}

func firstAny(v any, def any) any {
	if v == nil {
		return def
	}
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "" || s == "<nil>" {
		return def
	}
	return v
}

func asString(v any) string {
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "<nil>" {
		return ""
	}
	return s
}

func asIntString(v any) string {
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "" || s == "<nil>" {
		return "0"
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return strconv.FormatInt(i, 10)
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return strconv.FormatInt(int64(f), 10)
	}
	return "0"
}

func maskMobile(mobile string) string {
	m := strings.TrimSpace(mobile)
	if len(m) < 3 {
		return m
	}
	last := m[len(m)-1:]
	if len(m) >= 2 {
		return m[:2] + "*****" + last
	}
	return m
}

func signParamsKey(data string, isLite bool) string {
	appid, clientver := config.PlatformConfig(isLite)
	key := "OIlwieks28dk2k092lksi2UIkp"
	if isLite {
		key = "LnT6xpN3khm36zse0QzvmgTZ3waWdRSA"
	}
	return util.MD5Hex(appid + key + clientver + data)
}
