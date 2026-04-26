package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// Miscellaneous server, sheet, singer, and song ranking endpoints.
func (c *Client) ServerNow(ctx context.Context, req ServerNowRequest) (*ServerNowResponse, error) {
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
	userid := firstNonEmpty(asIntString(firstAny(params["userid"], nil)), cookies["userid"], "0")
	token := firstNonEmpty(fmt.Sprintf("%v", firstAny(params["token"], nil)), cookies["token"], "")
	resp, err := c.Call(ctx, RouteServerNow, Request{
		Method: "POST",
		URL:    "/v1/server_now",
		Params: map[string]any{"plat": 3},
		Data: map[string]any{
			"token":  token,
			"userid": userid,
		},
		Cookie:      cookies,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "usercenter.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := ServerNowResponse(*resp)
	return &out, nil
}

func (c *Client) SheetHot(ctx context.Context, req SheetHotRequest) (*SheetHotResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteSheetHot, Request{
		Method: "GET",
		URL:    "/miniyueku/v1/opern_square/get_home_hot_opern",
		Params: map[string]any{
			"srcappid":   2919,
			"opern_type": toInt(firstAny(params["opern_type"], nil), 1),
		},
		Cookie:      req.Cookie,
		EncryptType: "web",
	})
	if err != nil {
		return nil, err
	}
	out := SheetHotResponse(*resp)
	return &out, nil
}

func (c *Client) SingerList(ctx context.Context, req SingerListRequest) (*SingerListResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteSingerList, Request{
		Method: "GET",
		URL:    "/ocean/v6/singer/list",
		Params: map[string]any{
			"hotsize":  toInt(firstAny(params["hotsize"], nil), 200),
			"musician": 0,
			"sextype":  toInt(firstAny(params["sextype"], nil), 0),
			"showtype": 2,
			"type":     toInt(firstAny(params["type"], nil), 0),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SingerListResponse(*resp)
	return &out, nil
}

func (c *Client) SongClimax(ctx context.Context, req SongClimaxRequest) (*SongClimaxResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	hashes := splitCSV(firstNonEmpty(fmt.Sprintf("%v", firstAny(params["hash"], nil)), ""))
	items := make([]map[string]string, 0, len(hashes))
	for _, h := range hashes {
		if strings.TrimSpace(h) == "" {
			continue
		}
		items = append(items, map[string]string{"hash": h})
	}
	dataJSON, _ := json.Marshal(items)
	resp, err := c.Call(ctx, RouteSongClimax, Request{
		Method:  "GET",
		BaseURL: "https://expendablekmrcdn.kugou.com",
		URL:     "/v1/audio_climax/audio",
		Params: map[string]any{
			"data": string(dataJSON),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SongClimaxResponse(*resp)
	return &out, nil
}

func (c *Client) SongRanking(ctx context.Context, req SongRankingRequest) (*SongRankingResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteSongRanking, Request{
		Method: "GET",
		URL:    "/grow/v1/song_ranking/play_page/ranking_info",
		Params: map[string]any{
			"album_audio_id": firstAny(params["album_audio_id"], params["albumAudioId"]),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SongRankingResponse(*resp)
	return &out, nil
}

func (c *Client) SongRankingFilter(ctx context.Context, req SongRankingFilterRequest) (*SongRankingFilterResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteSongRankingFilter, Request{
		Method: "GET",
		URL:    "/grow/v1/song_ranking/unlock/v2/ranking_filter",
		Params: map[string]any{
			"album_audio_id": firstAny(params["album_audio_id"], params["albumAudioId"]),
			"page":           toInt(firstAny(params["page"], nil), 1),
			"pagesize":       toInt(firstAny(params["pagesize"], nil), 30),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SongRankingFilterResponse(*resp)
	return &out, nil
}
