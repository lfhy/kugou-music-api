package sdk

import (
	"context"
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rushingrain/kugou-music-api/core/config"
	"github.com/rushingrain/kugou-music-api/core/kugou"
)

// Cloud playlist mutations send raw AES ciphertext with an RSA session portrait.
// The gateway's ordinary Android signature and default query are not used.
func (c *Client) callPlaylistCloud(ctx context.Context, endpoint string, data map[string]any, cookie map[string]string) (*Response, error) {
	cookies := c.mergeRequestCookies(cookie, false)
	enc, err := playlistAesEncrypt(data)
	if err != nil {
		return nil, err
	}
	pub := kugou.PublicRASKey
	if c.isLite {
		pub = kugou.PublicLiteRASKey
	}
	var userid any = cookies["userid"]
	if userid == "" {
		userid = 0
	}
	portrait, err := kugou.CryptoRSAEncryptPKCS1Hex(map[string]any{
		"aes": enc.Key, "uid": userid, "token": cookies["token"],
	}, pub)
	if err != nil {
		return nil, err
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	clienttime := strconv.FormatInt(time.Now().Unix(), 10)
	mid := firstNonEmpty(cookies["KUGOU_API_MID"], cookies["mid"])
	cookies["KUGOU_API_MID"] = mid
	body, err := base64.StdEncoding.DecodeString(enc.CipherBase64)
	if err != nil {
		return nil, err
	}
	raw, requestErr := c.doRequest(ctx, kugou.RequestConfig{
		Method: "POST", URL: endpoint,
		Params: map[string]any{
			"appid": appid, "clientver": clientver, "mid": mid,
			"clienttime": clienttime, "key": signParamsKey(clienttime, c.isLite),
			"dfid": firstNonEmpty(cookies["dfid"], "-"), "p": strings.ToUpper(portrait),
		},
		Data: body, Cookie: cookies, ClearDefaultParams: true, NotSignature: true,
		Headers: map[string]string{
			"x-router":     "cloudlist.service.kugou.com",
			"Content-Type": "application/json;charset=utf-8",
		},
	})
	if len(raw.Cookie) > 0 {
		c.updateCookiePool(raw.Cookie)
	}
	out := &Response{Status: raw.Status, RawBody: raw.Body, Headers: raw.Headers, Cookie: raw.Cookie}
	// Plain JSON errors may have any length; CBC decryption requires full blocks.
	if len(raw.Body) > 0 && len(raw.Body)%aes.BlockSize == 0 {
		if decoded, err := playlistAesDecryptFromRaw(raw.Body, enc.Key); err == nil {
			out.Body = decoded
			out.RawBody, _ = json.Marshal(decoded)
		}
	}
	if out.Body == nil {
		out.Body = decodeJSONBody(raw.Body)
	}
	if requestErr != nil {
		out.Status = 502
		return out, requestErr
	}
	if out.Body == nil {
		out.Status = 500
		out.Body = map[string]any{"status": 0, "msg": string(raw.Body)}
		return out, fmt.Errorf("playlist cloud: invalid JSON response")
	}
	status, hasStatus := out.Body["status"].(float64)
	errorCode := out.Body["error_code"]
	if (hasStatus && status == 0) || (errorCode != nil && errorCode != float64(0) && errorCode != "" && errorCode != false) {
		out.Status = 502
		return out, fmt.Errorf("playlist cloud: %v", out.Body)
	}
	out.Status = 200
	return out, nil
}
