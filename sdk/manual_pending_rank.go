package sdk

import (
	"context"
	"fmt"
)

// Ranking endpoints keep custom legacy parameter names for compatibility.
func (c *Client) RankInfo(ctx context.Context, req RankInfoRequest) (*RankInfoResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteRankInfo, Request{
		Method: "GET",
		URL:    "/ocean/v6/rank/info",
		Params: map[string]any{
			"rank_cid":       toInt(firstAny(params["rank_cid"], params["rankCid"]), 0),
			"rankid":         firstAny(params["rankid"], params["id"]),
			"with_album_img": toInt(firstAny(params["album_img"], params["with_album_img"]), 1),
			"zone":           firstNonEmpty(fmt.Sprintf("%v", firstAny(params["zone"], nil)), ""),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := RankInfoResponse(*resp)
	return &out, nil
}

func (c *Client) RankList(ctx context.Context, req RankListRequest) (*RankListResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteRankList, Request{
		Method: "GET",
		URL:    "/ocean/v6/rank/list",
		Params: map[string]any{
			"plat":     2,
			"withsong": toInt(firstAny(params["withsong"], nil), 1),
			"parentid": 0,
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := RankListResponse(*resp)
	return &out, nil
}

func (c *Client) RankTop(ctx context.Context, req RankTopRequest) (*RankTopResponse, error) {
	resp, err := c.Call(ctx, RouteRankTop, Request{
		Method:      "GET",
		URL:         "/mobileservice/api/v5/rank/rec_rank_list",
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := RankTopResponse(*resp)
	return &out, nil
}

func (c *Client) RankVol(ctx context.Context, req RankVolRequest) (*RankVolResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteRankVol, Request{
		Method: "GET",
		URL:    "/ocean/v6/rank/vol",
		Params: map[string]any{
			"rank_cid": toInt(firstAny(params["rank_cid"], params["rankCid"]), 0),
			"rankid":   firstAny(params["rankid"], params["id"]),
			"ranktype": 1,
			"type":     0,
			"plat":     2,
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := RankVolResponse(*resp)
	return &out, nil
}
