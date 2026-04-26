package sdk

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/lfhy/kugou-music-api/core/kugou"
)

type SendCaptchaRequest struct {
	Mobile string
	Cookie map[string]string
}

// SendCaptcha sends login SMS code with JS-compatible payload.
// Equivalent to module/captcha_sent.js: businessid=5, plat=3.
func (c *Client) SendCaptcha(ctx context.Context, req SendCaptchaRequest) (*Response, error) {
	merged := c.Cookie()
	for k, v := range req.Cookie {
		merged[k] = v
	}

	// Keep JS behavior: only send mid cookie to avoid stale auth cookies causing 20018.
	mid := strings.TrimSpace(firstNonEmpty(merged["mid"], merged["KUGOU_API_MID"]))
	reqCookie := map[string]string{}
	if mid != "" {
		reqCookie["mid"] = mid
		reqCookie["KUGOU_API_MID"] = mid
	}

	raw, err := c.core.CreateRequest(ctx, kugou.RequestConfig{
		Method:      "POST",
		BaseURL:     "http://login.user.kugou.com",
		URL:         "/v7/send_mobile_code",
		Data:        map[string]any{"businessid": 5, "mobile": req.Mobile, "plat": 3},
		EncryptType: "android",
		Cookie:      reqCookie,
	})

	if len(raw.Cookie) > 0 {
		c.updateCookiePool(raw.Cookie)
	}

	out := &Response{Status: raw.Status, RawBody: raw.Body, Headers: raw.Headers, Cookie: raw.Cookie}
	var body map[string]any
	if json.Unmarshal(raw.Body, &body) == nil {
		out.Body = body
	}
	return out, err
}
