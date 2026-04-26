package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/core/kugou"
)

type userDetailSignature struct {
	Token      string `json:"token"`
	ClientTime int64  `json:"clienttime"`
}

type userDetailPayload struct {
	VisitTime int64  `json:"visit_time"`
	UserType  int    `json:"usertype"`
	P         string `json:"p"`
	UserID    int    `json:"userid"`
}

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
	p, err := kugou.CryptoRSAEncryptRawHex(userDetailSignature{
		Token:      token,
		ClientTime: clientTime,
	}, pubKey)
	if err != nil {
		return nil, err
	}

	raw, err := c.core.CreateRequest(ctx, kugou.RequestConfig{
		Method:      "POST",
		URL:         "/v3/get_my_info",
		Params:      map[string]any{"plat": 1},
		Data:        userDetailPayload{VisitTime: clientTime, UserType: 1, P: strings.ToUpper(p), UserID: userID},
		Cookie:      cookies,
		EncryptType: "android",
		Headers: map[string]string{
			"x-router": "usercenter.kugou.com",
		},
	})
	if len(raw.Cookie) > 0 {
		c.updateCookiePool(raw.Cookie)
	}
	out := UserDetailResponse{
		Status:  raw.Status,
		RawBody: raw.Body,
		Headers: raw.Headers,
		Cookie:  raw.Cookie,
	}
	if err == nil {
		out.Body = map[string]any{}
		_ = json.Unmarshal(raw.Body, &out.Body)
	}
	return &out, nil
}
