package sdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/core/config"
)

// Signature-based audio recommendation endpoints keep their custom payloads here.
func (c *Client) AudioRelated(ctx context.Context, req AudioRelatedRequest) (*AudioRelatedResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}

	dataMap := map[string]any{
		"album_audio_id": toInt(params["album_audio_id"], req.AlbumAudioId),
		"appid":          1005,
		"area_code":      1,
		"clientver":      12329,
	}

	showDetail := false
	if v, ok := params["show_detail"]; ok {
		showDetail = toInt(v, 1) == 0
	}

	endpoint := "/v2/audio_related/total"
	if !showDetail {
		endpoint = "/v3/album_audio/related"
		dataMap["page"] = toInt(params["page"], 1)
		dataMap["pagesize"] = toInt(params["pagesize"], 30)
		dataMap["show_input"] = 1
		dataMap["show_type"] = toInt(params["show_type"], 0)
		dataMap["sort"] = audioRelatedSort(params["sort"])
		dataMap["type"] = toInt(params["type"], 0)
	}
	dataMap["version"] = 1
	dataMap["signature"] = md5SortedWithKey(dataMap, "OIlwieks28dk2k092lksi2UIkp")

	resp, err := c.Call(ctx, RouteAudioRelated, Request{
		Method:             "GET",
		BaseURL:            "https://listkmrp3cdnretry.kugou.com",
		URL:                endpoint,
		Params:             dataMap,
		Cookie:             req.Cookie,
		EncryptType:        "android",
		ClearDefaultParams: boolPtr(true),
	})
	if err != nil {
		return nil, err
	}
	out := AudioRelatedResponse(*resp)
	return &out, nil
}

func (c *Client) AudioAccompanyMatching(ctx context.Context, req AudioAccompanyMatchingRequest) (*AudioAccompanyMatchingResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	appid, _ := config.PlatformConfig(c.isLite)
	dataMap := map[string]any{
		"isteen":   0,
		"mixId":    toInt(firstAny(params["mixId"], req.MixId), 0),
		"usemkv":   1,
		"platform": 2,
		"fileName": firstNonEmpty(fmt.Sprintf("%v", firstAny(params["fileName"], req.FileName)), ""),
		"hash":     firstNonEmpty(fmt.Sprintf("%v", firstAny(params["hash"], req.Hash)), ""),
		"version":  12375,
		"appid":    appid,
	}
	dataMap["sign"] = md5AmpersandSign(dataMap)

	resp, err := c.Call(ctx, RouteAudioAccompanyMatching, Request{
		Method:             "GET",
		BaseURL:            "https://nsongacsing.kugou.com",
		URL:                "/sing7/accompanywan/json/v2/cdn/optimal_matching_accompany_2_listen.do",
		Params:             dataMap,
		Cookie:             req.Cookie,
		EncryptType:        "android",
		ClearDefaultParams: boolPtr(true),
		NotSignature:       boolPtr(true),
	})
	if err != nil {
		return nil, err
	}
	out := AudioAccompanyMatchingResponse(*resp)
	return &out, nil
}

func (c *Client) AudioKtvTotal(ctx context.Context, req AudioKtvTotalRequest) (*AudioKtvTotalResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	appid, _ := config.PlatformConfig(c.isLite)
	dataMap := map[string]any{
		"isteen":     0,
		"songId":     toInt(firstAny(params["songId"], req.SongId), 0),
		"usemkv":     1,
		"platform":   2,
		"singerName": firstNonEmpty(fmt.Sprintf("%v", firstAny(params["singerName"], req.SingerName)), ""),
		"songHash":   firstNonEmpty(fmt.Sprintf("%v", firstAny(params["songHash"], req.SongHash)), ""),
		"version":    12375,
		"appid":      appid,
	}
	dataMap["sign"] = md5AmpersandSign(dataMap)

	resp, err := c.Call(ctx, RouteAudioKtvTotal, Request{
		Method:             "GET",
		BaseURL:            "https://acsing.service.kugou.com",
		URL:                "/sing7/listenguide/json/v2/cdn/listenguide/get_total_opus_num_v02.do",
		Params:             dataMap,
		Cookie:             req.Cookie,
		EncryptType:        "android",
		ClearDefaultParams: boolPtr(true),
		NotSignature:       boolPtr(true),
	})
	if err != nil {
		return nil, err
	}
	out := AudioKtvTotalResponse(*resp)
	return &out, nil
}

func (c *Client) Brush(ctx context.Context, req BrushRequest) (*BrushResponse, error) {
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

	appid, _ := config.PlatformConfig(c.isLite)
	dateTime := time.Now().UnixMilli()
	userid := strings.TrimSpace(firstNonEmpty(cookies["userid"], asIntString(firstAny(params["userid"], req.Userid)), "0"))
	vipType := strings.TrimSpace(firstNonEmpty(cookies["vip_type"], asIntString(firstAny(params["vipType"], req.VipType)), "0"))
	mode := strings.TrimSpace(firstNonEmpty(fmt.Sprintf("%v", firstAny(params["mode"], req.Mode)), "normal"))
	if mode == "" || mode == "<nil>" {
		mode = "normal"
	}

	personalRecommend := map[string]any{
		"userid":                  userid,
		"appid":                   appid,
		"playlist_ver":            2,
		"clienttime":              dateTime,
		"mid":                     cookies["KUGOU_API_MID"],
		"new_sync_point":          dateTime,
		"module_id":               1,
		"action":                  "login",
		"vip_type":                vipType,
		"vip_flags":               3,
		"recommend_source_locked": 0,
		"song_pool_id":            toInt(firstAny(params["song_pool_id"], req.SongPoolId), 0),
		"callerid":                0,
		"m_type":                  1,
		"kguid":                   userid,
		"platform":                "ios",
		"area_code":               1,
		"fakem":                   "ca981cfc583a4c37f28d2d49000013c16a0a",
		"clientver":               11850,
		"mode":                    mode,
		"active_swtich":           "on",
		"key":                     signParamsKey(strconv.FormatInt(dateTime, 10), c.isLite),
	}

	data := map[string]any{
		"behaviors": []any{},
		"abtest": map[string]any{
			"abtest": map[string]any{"shuashua": map[string]any{"commentcard": 2}},
		},
		"personal_recommend_params": personalRecommend,
	}

	resp, err := c.Call(ctx, RouteBrush, Request{
		Method:      "POST",
		URL:         "/genesisapi/v1/newepoch_song_rec/feed",
		Params:      map[string]any{"sort_type": 1, "platform": "ios", "page": 1, "content_ver": 4, "clientver": 11850},
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := BrushResponse(*resp)
	return &out, nil
}
