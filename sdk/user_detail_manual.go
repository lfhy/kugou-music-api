package sdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/core/kugou"
)

// UserDetail overrides generated behavior with JS-compatible signing for /user/detail.
func (c *Client) UserDetail(ctx context.Context, req UserDetailRequest) (*UserDetailResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}

	for k, v := range req.Extra {
		if sv, ok := v.(string); ok {
			cookies[k] = sv
		}
	}

	token := strings.TrimSpace(fmt.Sprintf("%v", req.Token))
	if token == "" || token == "<nil>" {
		token = strings.TrimSpace(cookies["token"])
	}

	userID := req.Userid
	if userID == 0 {
		if uidStr := strings.TrimSpace(cookies["userid"]); uidStr != "" {
			if v, err := strconv.Atoi(uidStr); err == nil {
				userID = v
			} else if fv, ferr := strconv.ParseFloat(uidStr, 64); ferr == nil {
				userID = int(fv)
			}
		}
	}
	if userID == 0 {
		if uidStr := strings.TrimSpace(fmt.Sprintf("%v", req.Extra["userid"])); uidStr != "" && uidStr != "<nil>" {
			if v, err := strconv.Atoi(uidStr); err == nil {
				userID = v
			} else if fv, ferr := strconv.ParseFloat(uidStr, 64); ferr == nil {
				userID = int(fv)
			}
		}
	}

	clientTime := time.Now().Unix()
	pubKey := kugou.PublicRASKey
	if c.isLite {
		pubKey = kugou.PublicLiteRASKey
	}
	p, err := kugou.CryptoRSAEncryptRawHex(map[string]any{
		"token":      token,
		"clienttime": clientTime,
	}, pubKey)
	if err != nil {
		return nil, err
	}

	resp, err := c.Call(ctx, RouteUserDetail, Request{
		Params: map[string]any{"plat": 1},
		Data: map[string]any{
			"visit_time": clientTime,
			"usertype":   1,
			"p":          strings.ToUpper(p),
			"userid":     userID,
		},
		Cookie: cookies,
	})
	if err != nil {
		return nil, err
	}
	out := UserDetailResponse(*resp)
	return &out, nil
}
