package sdk

import "context"

// Yueku endpoints are grouped separately to keep the youth wrapper files small.
func (c *Client) Yueku(ctx context.Context, req YuekuRequest) (*YuekuResponse, error) {
	resp, err := c.Call(ctx, RouteYueku, Request{
		Method: "GET",
		URL:    "/v1/yueku/recommend_v2",
		Params: map[string]any{"operator": 7, "plat": 0, "type": 11, "area_code": 1, "req_multi": 1},
		Cookie: req.Cookie,
		Headers: map[string]string{
			"x-router": "service.mobile.kugou.com",
		},
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := YuekuResponse(*resp)
	return &out, nil
}

func (c *Client) YuekuFm(ctx context.Context, req YuekuFmRequest) (*YuekuFmResponse, error) {
	resp, err := c.Call(ctx, RouteYuekuFm, Request{
		Method: "GET",
		URL:    "/v1/time_fm_info",
		Params: map[string]any{"operator": 7, "plat": 0, "type": 11, "area_code": 1, "req_multi": 1},
		Cookie: req.Cookie,
		Headers: map[string]string{
			"x-router": "fm.service.kugou.com",
		},
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := YuekuFmResponse(*resp)
	return &out, nil
}
