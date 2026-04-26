package sdk

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lfhy/kugou-music-api/core/kugou"
)

// finalizeLoginResponse normalizes login payloads and writes refreshed auth cookies back to the pool.
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
