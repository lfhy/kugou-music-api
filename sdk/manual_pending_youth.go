package sdk

import (
	"context"
	"fmt"
)

// Youth and yueku endpoints are grouped here to keep manual wrappers small.
func (c *Client) UserVipDetail(ctx context.Context, req UserVipDetailRequest) (*UserVipDetailResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, c.loginStateError(cookies)
	}
	resp, err := c.Call(ctx, RouteUserVipDetail, Request{
		Method:  "GET",
		BaseURL: "https://kugouvip.kugou.com",
		URL:     "/v1/get_union_vip",
		Params:  map[string]any{"busi_type": "concept"},
		Cookie:  cookies,
	})
	if err != nil {
		return nil, err
	}
	out := UserVipDetailResponse(*resp)
	return &out, nil
}

func (c *Client) YouthChannelAll(ctx context.Context, req YouthChannelAllRequest) (*YouthChannelAllResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteYouthChannelAll, Request{
		Method: "GET",
		URL:    "/youth/v2/channel/channel_all_list",
		Params: map[string]any{
			"page":     toInt(firstAny(params["page"], nil), 1),
			"pagesize": toInt(firstAny(params["pagesize"], nil), 30),
			"type":     1,
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := YouthChannelAllResponse(*resp)
	return &out, nil
}

func (c *Client) YouthChannelAmway(ctx context.Context, req YouthChannelAmwayRequest) (*YouthChannelAmwayResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteYouthChannelAmway, Request{
		Method: "GET",
		URL:    "/youth/api/amway/v2/index",
		Params: map[string]any{
			"global_collection_id": firstAny(params["global_collection_id"], nil),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := YouthChannelAmwayResponse(*resp)
	return &out, nil
}

func (c *Client) YouthChannelDetail(ctx context.Context, req YouthChannelDetailRequest) (*YouthChannelDetailResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	data := make([]map[string]any, 0)
	for _, s := range splitCSV(firstNonEmpty(fmt.Sprintf("%v", firstAny(params["global_collection_id"], nil)), "")) {
		data = append(data, map[string]any{"global_collection_id": s})
	}
	resp, err := c.Call(ctx, RouteYouthChannelDetail, Request{
		Method:      "POST",
		URL:         "/youth/api/channel/v1/channel_list_by_id",
		Data:        map[string]any{"data": data},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := YouthChannelDetailResponse(*resp)
	return &out, nil
}

func (c *Client) YouthChannelSimilar(ctx context.Context, req YouthChannelSimilarRequest) (*YouthChannelSimilarResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	cookies = normalizeOptionalLogin(ctx, c, cookies)
	vipType := firstNonEmpty(asIntString(firstAny(params["vip_type"], nil)), cookies["vip_type"], "0")
	resp, err := c.Call(ctx, RouteYouthChannelSimilar, Request{
		Method: "POST",
		URL:    "/youth/v1/channel/get_friendly_channel",
		Params: map[string]any{
			"channel_id": firstAny(params["channel_id"], nil),
		},
		Data: map[string]any{
			"area_code":    1,
			"playlist_ver": 2,
			"vip_type":     vipType,
			"platform":     "ios",
		},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := YouthChannelSimilarResponse(*resp)
	return &out, nil
}

func (c *Client) YouthChannelSub(ctx context.Context, req YouthChannelSubRequest) (*YouthChannelSubResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, c.loginStateError(cookies)
	}
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	t := ternaryInt(toInt(firstAny(params["t"], nil), 1) == 0, 0, 1)
	method := "POST"
	urlPath := "/youth/v1/channel_subscribe"
	if t == 0 {
		method = "DELETE"
		urlPath = "/youth/v1/channel_un_subscribe"
	}
	resp, err := c.Call(ctx, RouteYouthChannelSub, Request{
		Method: method,
		URL:    urlPath,
		Params: map[string]any{
			"global_collection_id": firstAny(params["global_collection_id"], nil),
			"source":               1,
		},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := YouthChannelSubResponse(*resp)
	return &out, nil
}

func (c *Client) YouthDayVip(ctx context.Context, req YouthDayVipRequest) (*YouthDayVipResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, c.loginStateError(cookies)
	}
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteYouthDayVip, Request{
		Method: "POST",
		URL:    "/youth/v1/recharge/receive_vip_listen_song",
		Params: map[string]any{
			"source_id":   90139,
			"receive_day": firstAny(params["receive_day"], nil),
		},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := YouthDayVipResponse(*resp)
	return &out, nil
}

func (c *Client) YouthDynamic(ctx context.Context, req YouthDynamicRequest) (*YouthDynamicResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, c.loginStateError(cookies)
	}
	resp, err := c.Call(ctx, RouteYouthDynamic, Request{
		Method:      "GET",
		URL:         "/youth/v3/user/get_dynamic",
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := YouthDynamicResponse(*resp)
	return &out, nil
}

func (c *Client) YouthDynamicRecent(ctx context.Context, req YouthDynamicRecentRequest) (*YouthDynamicRecentResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, c.loginStateError(cookies)
	}
	resp, err := c.Call(ctx, RouteYouthDynamicRecent, Request{
		Method:      "GET",
		URL:         "/youth/v3/user/recent_dynamic",
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := YouthDynamicRecentResponse(*resp)
	return &out, nil
}

func (c *Client) YouthMonthVipRecord(ctx context.Context, req YouthMonthVipRecordRequest) (*YouthMonthVipRecordResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, c.loginStateError(cookies)
	}
	resp, err := c.Call(ctx, RouteYouthMonthVipRecord, Request{
		Method: "GET",
		URL:    "/youth/v1/activity/get_month_vip_record",
		Params: map[string]any{"latest_limit": 100},
		Cookie: cookies,
	})
	if err != nil {
		return nil, err
	}
	out := YouthMonthVipRecordResponse(*resp)
	return &out, nil
}
