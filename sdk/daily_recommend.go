package sdk

import (
	"context"
	"fmt"
	"strings"
)

// GetDailyRecommendGuest fetches daily recommend list for guest mode.
// It first calls EverydayRecommend, and if business code 200103 is returned,
// it falls back to RecommendSongs(platform=android, userid=0).
func (c *Client) GetDailyRecommendGuest(ctx context.Context, cookie map[string]string) (*Response, error) {
	effective := map[string]string{}
	for k, v := range cookie {
		effective[k] = v
	}
	if merged, ok := c.ensureLoginValid(ctx, effective); ok {
		effective = merged
	} else {
		delete(effective, "token")
		delete(effective, "userid")
		delete(effective, "vip_type")
		delete(effective, "vip_token")
	}

	resp, err := c.EverydayRecommend(ctx, EverydayRecommendRequest{
		Platform: "ios",
		Cookie:   effective,
	})
	if err != nil {
		return resp, err
	}

	if isBizCode(resp, "200103") {
		return c.RecommendSongs(ctx, RecommendSongsRequest{
			Platform: "android",
			Userid:   "0",
			Cookie:   effective,
		})
	}

	return resp, nil
}

// EverydayRecommend returns daily recommendation with login-aware cookie handling.
func (c *Client) EverydayRecommend(ctx context.Context, req EverydayRecommendRequest) (*EverydayRecommendResponse, error) {
	effective := c.Cookie()
	for k, v := range req.Cookie {
		effective[k] = v
	}
	if merged, ok := c.ensureLoginValid(ctx, effective); ok {
		effective = merged
	} else {
		delete(effective, "token")
		delete(effective, "userid")
		delete(effective, "vip_type")
		delete(effective, "vip_token")
	}

	platform := strings.TrimSpace(fmt.Sprintf("%v", req.Platform))
	if platform == "" {
		platform = "ios"
	}
	resp, err := c.Call(ctx, RouteEverydayRecommend, Request{
		Method:      "POST",
		URL:         "/everyday_song_recommend",
		Params:      map[string]any{"platform": platform},
		Cookie:      effective,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "everydayrec.service.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := EverydayRecommendResponse(*resp)
	return &out, nil
}

// RecommendSongs returns recommend songs. If login is invalid, it falls back to guest userid=0.
func (c *Client) RecommendSongs(ctx context.Context, req RecommendSongsRequest) (*RecommendSongsResponse, error) {
	effective := c.Cookie()
	for k, v := range req.Cookie {
		effective[k] = v
	}
	if merged, ok := c.ensureLoginValid(ctx, effective); ok {
		effective = merged
	}

	platform := strings.TrimSpace(fmt.Sprintf("%v", req.Platform))
	if platform == "" {
		platform = "android"
	}
	userid := strings.TrimSpace(fmt.Sprintf("%v", req.Userid))
	if userid == "" || userid == "<nil>" {
		userid = strings.TrimSpace(effective["userid"])
	}
	if userid == "" {
		userid = "0"
	}

	resp, err := c.Call(ctx, RouteRecommendSongs, Request{
		Method:      "POST",
		URL:         "/everyday_song_recommend",
		Data:        map[string]any{"platform": platform, "userid": userid},
		Cookie:      effective,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "everydayrec.service.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := RecommendSongsResponse(*resp)
	return &out, nil
}

func isBizCode(resp *Response, code string) bool {
	if resp == nil || resp.Body == nil {
		return false
	}
	v, ok := resp.Body["error_code"]
	if !ok {
		return false
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v)) == code
}
