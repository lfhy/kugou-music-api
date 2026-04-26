package sdk

import (
	"context"
	"fmt"
)

// Media detail and long-audio endpoints use manual payload shaping.
func (c *Client) KmrAudioMv(ctx context.Context, req KmrAudioMvRequest) (*KmrAudioMvResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	data := make([]map[string]any, 0)
	for _, s := range splitCSV(firstNonEmpty(fmt.Sprintf("%v", firstAny(params["album_audio_id"], nil)), "")) {
		data = append(data, map[string]any{"album_audio_id": s})
	}
	resp, err := c.Call(ctx, RouteKmrAudioMv, Request{
		Method: "POST",
		URL:    "/kmr/v1/audio/mv",
		Data: map[string]any{
			"data":   data,
			"fields": firstNonEmpty(fmt.Sprintf("%v", firstAny(params["fields"], nil)), ""),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "openapi.kugou.com", "KG-TID": "38"},
	})
	if err != nil {
		return nil, err
	}
	out := KmrAudioMvResponse(*resp)
	return &out, nil
}

func (c *Client) KrmAudio(ctx context.Context, req KrmAudioRequest) (*KrmAudioResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	data := make([]map[string]any, 0)
	for _, s := range splitCSV(firstNonEmpty(fmt.Sprintf("%v", firstAny(params["album_audio_id"], nil)), "")) {
		data = append(data, map[string]any{"entity_id": toInt(s, 0)})
	}
	resp, err := c.Call(ctx, RouteKrmAudio, Request{
		Method: "POST",
		URL:    "/kmr/v2/audio",
		Data: map[string]any{
			"data":   data,
			"fields": firstNonEmpty(fmt.Sprintf("%v", firstAny(params["fields"], nil)), "base"),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "openapi.kugou.com", "KG-TID": "238"},
	})
	if err != nil {
		return nil, err
	}
	out := KrmAudioResponse(*resp)
	return &out, nil
}

func (c *Client) LongaudioAlbumAudios(ctx context.Context, req LongaudioAlbumAudiosRequest) (*LongaudioAlbumAudiosResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteLongaudioAlbumAudios, Request{
		Method: "POST",
		URL:    "/longaudio/v2/album_audios",
		Data: map[string]any{
			"album_id":  firstAny(params["album_id"], nil),
			"area_code": 1,
			"tagid":     0,
			"page":      toInt(firstAny(params["page"], nil), 1),
			"pagesize":  toInt(firstAny(params["pagesize"], nil), 30),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "openapi.kugou.com", "KG-TID": "78"},
	})
	if err != nil {
		return nil, err
	}
	out := LongaudioAlbumAudiosResponse(*resp)
	return &out, nil
}

func (c *Client) LongaudioAlbumDetail(ctx context.Context, req LongaudioAlbumDetailRequest) (*LongaudioAlbumDetailResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	list := make([]map[string]any, 0)
	for _, s := range splitCSV(firstNonEmpty(fmt.Sprintf("%v", firstAny(params["album_id"], nil)), "")) {
		list = append(list, map[string]any{"album_id": s})
	}
	resp, err := c.Call(ctx, RouteLongaudioAlbumDetail, Request{
		Method: "POST",
		URL:    "/openapi/v2/broadcast",
		Data: map[string]any{
			"data":           list,
			"show_album_tag": 1,
			"fields":         "album_name,album_id,category,authors,sizable_cover,intro,author_name,trans_param,album_tag,mix_intro,full_intro,is_publish",
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"KG-TID": "78"},
	})
	if err != nil {
		return nil, err
	}
	out := LongaudioAlbumDetailResponse(*resp)
	return &out, nil
}

func (c *Client) LongaudioDailyRecommend(ctx context.Context, req LongaudioDailyRecommendRequest) (*LongaudioDailyRecommendResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteLongaudioDailyRecommend, Request{
		Method: "POST",
		URL:    "/longaudio/v1/home_new/daily_recommend",
		Params: map[string]any{
			"module_id": 1,
			"size":      toInt(firstAny(params["pagesize"], nil), 30),
			"page":      toInt(firstAny(params["page"], nil), 1),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := LongaudioDailyRecommendResponse(*resp)
	return &out, nil
}

func (c *Client) LongaudioRankRecommend(ctx context.Context, req LongaudioRankRecommendRequest) (*LongaudioRankRecommendResponse, error) {
	resp, err := c.Call(ctx, RouteLongaudioRankRecommend, Request{
		Method: "GET",
		URL:    "/longaudio/v1/home_new/rank_card_recommend",
		Params: map[string]any{"platform": "ios"},
		Cookie: req.Cookie,
	})
	if err != nil {
		return nil, err
	}
	out := LongaudioRankRecommendResponse(*resp)
	return &out, nil
}

func (c *Client) LongaudioVipRecommend(ctx context.Context, req LongaudioVipRecommendRequest) (*LongaudioVipRecommendResponse, error) {
	resp, err := c.Call(ctx, RouteLongaudioVipRecommend, Request{
		Method: "POST",
		URL:    "/longaudio/v1/home_new/vip_select_recommend",
		Data:   map[string]any{"album_playlist": []any{}},
		Params: map[string]any{"position": "2", "clientver": 12329},
		Cookie: req.Cookie,
	})
	if err != nil {
		return nil, err
	}
	out := LongaudioVipRecommendResponse(*resp)
	return &out, nil
}

func (c *Client) LongaudioWeekRecommend(ctx context.Context, req LongaudioWeekRecommendRequest) (*LongaudioWeekRecommendResponse, error) {
	resp, err := c.Call(ctx, RouteLongaudioWeekRecommend, Request{
		Method: "POST",
		URL:    "/longaudio/v1/home_new/week_new_albums_recommend",
		Data:   map[string]any{"album_playlist": []any{}},
		Params: map[string]any{"clientver": 12329},
		Cookie: req.Cookie,
	})
	if err != nil {
		return nil, err
	}
	out := LongaudioWeekRecommendResponse(*resp)
	return &out, nil
}
