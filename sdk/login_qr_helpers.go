package sdk

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/skip2/go-qrcode"
)

const loginQRCodeURLPrefix = "https://h5.kugou.com/apps/loginQRCode/html/index.html?qrcode="

func (r *LoginQrKeyResponse) QRCodeKey() string {
	if r == nil {
		return ""
	}
	return pickStringDeepMap(r.Body, "qrcode", "key", "qrkey")
}

func (r *LoginQrKeyResponse) QRCodeURL() string {
	return buildLoginQRCodeURL(r.QRCodeKey())
}

func (r *LoginQrKeyResponse) QRCodeImageURL() string {
	return r.QRCodeURL()
}

func (r *LoginQrCheckResponse) StatusCode() int {
	if r == nil {
		return -1
	}
	if code := pickNestedStatus(r.Body); code != -1 {
		return code
	}
	return pickIntDeepMap(r.Body, "status")
}

func (r *LoginQrCheckResponse) Token() string {
	if r == nil {
		return ""
	}
	return pickStringDeepMap(r.Body, "token")
}

func (r *LoginQrCheckResponse) UserID() string {
	if r == nil {
		return ""
	}
	return pickStringDeepMap(r.Body, "userid", "user_id", "uid")
}

func (r *LoginQrCreateResponse) URL() string {
	if r == nil {
		return ""
	}
	return pickStringDeepMap(r.Body, "url")
}

func (r *LoginQrCreateResponse) Base64() string {
	if r == nil {
		return ""
	}
	return pickStringDeepMap(r.Body, "base64")
}

func finalizeLoginQrCheckResponse(c *Client, resp *Response) *LoginQrCheckResponse {
	if resp == nil {
		return nil
	}
	out := LoginQrCheckResponse(*resp)
	if out.StatusCode() == 4 {
		token := out.Token()
		userid := out.UserID()
		if token != "" {
			c.SetCookie("token", token)
			resp.Cookie = append(resp.Cookie, "token="+token)
		}
		if userid != "" {
			c.SetCookie("userid", userid)
			resp.Cookie = append(resp.Cookie, "userid="+userid)
		}
	}
	return &out
}

func buildLoginQrCreateResponse(req LoginQrCreateRequest) (*LoginQrCreateResponse, error) {
	key := strings.TrimSpace(fmt.Sprintf("%v", req.Key))
	if key == "" || key == "<nil>" {
		return nil, fmt.Errorf("login qr create: missing key")
	}

	data := map[string]any{
		"url":    buildLoginQRCodeURL(key),
		"base64": "",
	}
	if toBool(req.Qrimg, false) {
		png, err := qrcode.Encode(data["url"].(string), qrcode.Medium, 256)
		if err != nil {
			return nil, err
		}
		data["base64"] = "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
	}

	body := map[string]any{
		"code": 200,
		"data": data,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	out := LoginQrCreateResponse{
		Status:  200,
		RawBody: raw,
		Body:    body,
		Headers: map[string]string{},
		Cookie:  []string{},
	}
	return &out, nil
}

func buildLoginQRCodeURL(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	return loginQRCodeURLPrefix + url.QueryEscape(key)
}

func pickStringDeepMap(v any, keys ...string) string {
	m, ok := v.(map[string]any)
	if !ok || m == nil {
		return ""
	}
	for _, k := range keys {
		if x, ok := m[k]; ok {
			s := strings.TrimSpace(fmt.Sprintf("%v", x))
			if s != "" && s != "<nil>" {
				return s
			}
		}
	}
	for _, x := range m {
		if mm, ok := x.(map[string]any); ok {
			if s := pickStringDeepMap(mm, keys...); s != "" {
				return s
			}
		}
	}
	return ""
}

func pickIntDeepMap(v any, key string) int {
	m, ok := v.(map[string]any)
	if !ok || m == nil {
		return -1
	}
	if x, ok := m[key]; ok {
		var n int
		fmt.Sscanf(fmt.Sprintf("%v", x), "%d", &n)
		return n
	}
	for _, x := range m {
		if mm, ok := x.(map[string]any); ok {
			n := pickIntDeepMap(mm, key)
			if n != -1 {
				return n
			}
		}
	}
	return -1
}

func pickNestedStatus(v any) int {
	root, ok := v.(map[string]any)
	if !ok || root == nil {
		return -1
	}
	data, ok := root["data"].(map[string]any)
	if !ok || data == nil {
		return -1
	}
	if x, ok := data["status"]; ok {
		var n int
		fmt.Sscanf(fmt.Sprintf("%v", x), "%d", &n)
		return n
	}
	return -1
}
