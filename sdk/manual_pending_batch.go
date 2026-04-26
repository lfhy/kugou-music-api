package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/lfhy/kugou-music-api/core/config"
)

func normalizeOptionalLogin(ctx context.Context, c *Client, cookie map[string]string) map[string]string {
	token := strings.TrimSpace(cookie["token"])
	userid := strings.TrimSpace(cookie["userid"])
	if token == "" || userid == "" || userid == "0" {
		return cookie
	}
	if merged, ok := c.ensureLoginValid(ctx, cookie); ok {
		return merged
	}
	delete(cookie, "token")
	delete(cookie, "userid")
	delete(cookie, "vip_token")
	delete(cookie, "vip_type")
	return cookie
}

func (c *Client) AlbumDetail(ctx context.Context, req AlbumDetailRequest) (*AlbumDetailResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteAlbumDetail, Request{
		Method: "POST",
		URL:    "/kmr/v2/albums",
		Data: map[string]any{
			"data": []map[string]any{{
				"album_id": firstAny(params["id"], params["album_id"]),
			}},
			"is_buy": toInt(firstAny(params["is_buy"], nil), 0),
			"fields": "album_id,album_name,publish_date,sizable_cover,intro,language,is_publish,heat,type,quality,authors,exclusive,author_name,trans_param",
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "openapi.kugou.com", "kg-tid": "255"},
	})
	if err != nil {
		return nil, err
	}
	out := AlbumDetailResponse(*resp)
	return &out, nil
}

func (c *Client) AlbumShop(ctx context.Context, req AlbumShopRequest) (*AlbumShopResponse, error) {
	resp, err := c.Call(ctx, RouteAlbumShop, Request{
		Method:      "GET",
		URL:         "/zhuanjidata/v3/album_shop_v2/get_classify_data",
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := AlbumShopResponse(*resp)
	return &out, nil
}

func (c *Client) ArtistDetail(ctx context.Context, req ArtistDetailRequest) (*ArtistDetailResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteArtistDetail, Request{
		Method: "POST",
		URL:    "/kmr/v3/author",
		Data: map[string]any{
			"author_id": firstAny(params["id"], params["author_id"]),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "openapi.kugou.com", "kg-tid": "36"},
	})
	if err != nil {
		return nil, err
	}
	out := ArtistDetailResponse(*resp)
	return &out, nil
}

func (c *Client) ArtistFollowNewsongs(ctx context.Context, req ArtistFollowNewsongsRequest) (*ArtistFollowNewsongsResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, requireLoginCookie(cookies)
	}
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	lastAlbumID := toInt(firstAny(params["last_album_id"], nil), 0)
	resp, err := c.Call(ctx, RouteArtistFollowNewsongs, Request{
		Method: "POST",
		URL:    "/feed/v1/follow/newsong_album_list",
		Params: map[string]any{
			"last_album_id": lastAlbumID,
			"page_size":     toInt(firstAny(params["pagesize"], params["page_size"]), 30),
			"opt_sort":      ternaryInt(toInt(firstAny(params["opt_sort"], nil), 1) == 2, 2, 1),
		},
		Data: map[string]any{
			"last_album_id": lastAlbumID,
		},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := ArtistFollowNewsongsResponse(*resp)
	return &out, nil
}

func (c *Client) ArtistHonour(ctx context.Context, req ArtistHonourRequest) (*ArtistHonourResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteArtistHonour, Request{
		Method:  "POST",
		BaseURL: "http://h5activity.kugou.com",
		URL:     "/v1/query_singer_honour_detail",
		Params: map[string]any{
			"singer_id": firstAny(params["id"], params["singer_id"]),
			"pagesize":  toInt(firstAny(params["pagesize"], nil), 30),
			"page":      toInt(firstAny(params["page"], nil), 1),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := ArtistHonourResponse(*resp)
	return &out, nil
}

func (c *Client) EverydayFriend(ctx context.Context, req EverydayFriendRequest) (*EverydayFriendResponse, error) {
	resp, err := c.Call(ctx, RouteEverydayFriend, Request{
		Method:  "POST",
		BaseURL: "https://acsing.service.kugou.com",
		URL:     "/sing7/relation/json/v3/friend_rec_by_using_song_list",
		Data: map[string]any{
			"list": []any{
				map[string]any{
					"user_id": 853927886,
					"mixsong_ids": []any{
						290083753, 251724346, 571554587, 250126644, 208831644, 40328518, 250504076, 581706850, 318347675, 585258401,
						288481998, 407414475, 28239430, 280584633, 291957521, 64556644, 243149863, 488725103, 32114153, 39951172,
						29019580, 40397606, 327507651, 32029382, 32218359, 340353127, 276448762, 177071956, 100031397, 249251602,
					},
				},
			},
		},
		Params:      map[string]any{"channel": 130, "isteen": 0, "platform": 2, "usemkv": 1},
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"pid": "126556797"},
	})
	if err != nil {
		return nil, err
	}
	out := EverydayFriendResponse(*resp)
	return &out, nil
}

func (c *Client) EverydayHistory(ctx context.Context, req EverydayHistoryRequest) (*EverydayHistoryResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	paramsMap := map[string]any{
		"mode":     firstNonEmpty(fmt.Sprintf("%v", firstAny(params["mode"], nil)), "list"),
		"platform": firstNonEmpty(fmt.Sprintf("%v", firstAny(params["platform"], nil)), "ios"),
	}
	if v := strings.TrimSpace(fmt.Sprintf("%v", firstAny(params["history_name"], nil))); v != "" && v != "<nil>" {
		paramsMap["history_name"] = v
	}
	if v := strings.TrimSpace(fmt.Sprintf("%v", firstAny(params["date"], nil))); v != "" && v != "<nil>" {
		paramsMap["date"] = v
	}
	resp, err := c.Call(ctx, RouteEverydayHistory, Request{
		Method:      "POST",
		URL:         "/everyday/api/v1/get_history",
		Params:      paramsMap,
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "everydayrec.service.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := EverydayHistoryResponse(*resp)
	return &out, nil
}

func (c *Client) EverydayStyleRecommend(ctx context.Context, req EverydayStyleRecommendRequest) (*EverydayStyleRecommendResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteEverydayStyleRecommend, Request{
		Method: "POST",
		URL:    "/everydayrec.service/everyday_style_recommend",
		Params: map[string]any{
			"tagids": firstNonEmpty(fmt.Sprintf("%v", firstAny(params["tagids"], nil)), ""),
		},
		Data:        map[string]any{},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := EverydayStyleRecommendResponse(*resp)
	return &out, nil
}

func (c *Client) FavoriteCount(ctx context.Context, req FavoriteCountRequest) (*FavoriteCountResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteFavoriteCount, Request{
		Method: "GET",
		URL:    "/count/v1/audio/mget_collect",
		Params: map[string]any{
			"mixsongids": firstAny(params["mixsongids"], nil),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := FavoriteCountResponse(*resp)
	return &out, nil
}

func (c *Client) IpZone(ctx context.Context, req IpZoneRequest) (*IpZoneResponse, error) {
	resp, err := c.Call(ctx, RouteIpZone, Request{
		Method:      "GET",
		URL:         "/v1/zone/index",
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "yuekucategory.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		if status, _ := resp.Body["status"].(float64); int(status) == 1 {
			if dm, ok := resp.Body["data"].(map[string]any); ok {
				if list, ok := dm["list"].([]any); ok {
					for i, item := range list {
						im, ok := item.(map[string]any)
						if !ok {
							continue
						}
						rawLink := strings.TrimSpace(fmt.Sprintf("%v", im["special_link"]))
						if rawLink == "" || rawLink == "<nil>" {
							continue
						}
						linkQS, err := url.ParseQuery(rawLink)
						if err != nil {
							continue
						}
						pathRaw := strings.TrimSpace(linkQS.Get("path"))
						if pathRaw == "" {
							continue
						}
						pathQS, err := url.ParseQuery(pathRaw)
						if err != nil {
							continue
						}
						ipID := strings.TrimSpace(pathQS.Get("ip_id"))
						if ipID == "" {
							continue
						}
						if n, err := strconv.Atoi(ipID); err == nil {
							im["ip_id"] = n
						}
						list[i] = im
					}
					dm["list"] = list
					resp.Body["data"] = dm
				}
			}
		}
	}
	out := IpZoneResponse(*resp)
	return &out, nil
}

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

func (c *Client) SceneAudioList(ctx context.Context, req SceneAudioListRequest) (*SceneAudioListResponse, error) {
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
	appid, clientver := config.PlatformConfig(c.isLite)
	userid := firstNonEmpty(asIntString(firstAny(params["userid"], nil)), cookies["userid"], "0")
	token := firstNonEmpty(fmt.Sprintf("%v", firstAny(params["token"], nil)), cookies["token"], "")
	resp, err := c.Call(ctx, RouteSceneAudioList, Request{
		Method: "POST",
		URL:    "/scene/v1/scene/audio_list",
		Params: map[string]any{
			"scene_id":  firstAny(params["id"], params["scene_id"]),
			"module_id": firstAny(params["module_id"], nil),
			"tag":       firstAny(params["tag"], nil),
			"page":      toInt(firstAny(params["page"], nil), 1),
			"page_size": toInt(firstAny(params["pagesize"], params["page_size"]), 30),
		},
		Data: map[string]any{
			"appid":     appid,
			"clientver": clientver,
			"token":     token,
			"userid":    userid,
		},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SceneAudioListResponse(*resp)
	return &out, nil
}

func (c *Client) SceneCollectionList(ctx context.Context, req SceneCollectionListRequest) (*SceneCollectionListResponse, error) {
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
	appid, clientver := config.PlatformConfig(c.isLite)
	userid := firstNonEmpty(asIntString(firstAny(params["userid"], nil)), cookies["userid"], "0")
	token := firstNonEmpty(fmt.Sprintf("%v", firstAny(params["token"], nil)), cookies["token"], "")
	resp, err := c.Call(ctx, RouteSceneCollectionList, Request{
		Method: "POST",
		URL:    "/scene/v1/distribution/collection_list",
		Data: map[string]any{
			"appid":        appid,
			"clientver":    clientver,
			"token":        token,
			"userid":       userid,
			"tag_id":       firstAny(params["tag_id"], nil),
			"page":         toInt(firstAny(params["page"], nil), 1),
			"page_size":    toInt(firstAny(params["pagesize"], params["page_size"]), 30),
			"exposed_data": []any{},
		},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SceneCollectionListResponse(*resp)
	return &out, nil
}

func (c *Client) SceneLists(ctx context.Context, req SceneListsRequest) (*SceneListsResponse, error) {
	resp, err := c.Call(ctx, RouteSceneLists, Request{
		Method:      "GET",
		URL:         "/scene/v1/scene/list",
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SceneListsResponse(*resp)
	return &out, nil
}

func (c *Client) SceneListsV2(ctx context.Context, req SceneListsV2Request) (*SceneListsV2Response, error) {
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
	sortType := map[string]int{"rec": 1, "hot": 2, "new": 3}
	sortKey := firstNonEmpty(fmt.Sprintf("%v", firstAny(params["sort"], nil)), "rec")
	sortVal := sortType[sortKey]
	if sortVal == 0 {
		sortVal = 1
	}
	userid := firstNonEmpty(asIntString(firstAny(params["userid"], params["kugouid"])), cookies["userid"], "0")
	resp, err := c.Call(ctx, RouteSceneListsV2, Request{
		Method: "POST",
		URL:    "/scene/v1/scene/list_v2",
		Params: map[string]any{
			"scene_id":  firstAny(params["id"], params["scene_id"]),
			"page":      toInt(firstAny(params["page"], nil), 1),
			"pagesize":  toInt(firstAny(params["pagesize"], params["page_size"]), 30),
			"sort_type": sortVal,
			"kugouid":   userid,
		},
		Data:        map[string]any{"exposure": []any{}},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SceneListsV2Response(*resp)
	return &out, nil
}

func (c *Client) SceneModule(ctx context.Context, req SceneModuleRequest) (*SceneModuleResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteSceneModule, Request{
		Method: "POST",
		URL:    "/scene/v1/scene/module",
		Params: map[string]any{
			"scene_id": firstAny(params["id"], params["scene_id"]),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SceneModuleResponse(*resp)
	return &out, nil
}

func (c *Client) SceneModuleInfo(ctx context.Context, req SceneModuleInfoRequest) (*SceneModuleInfoResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteSceneModuleInfo, Request{
		Method: "GET",
		URL:    "/scene/v1/scene/module_info",
		Params: map[string]any{
			"scene_id":  firstAny(params["id"], params["scene_id"]),
			"module_id": firstAny(params["module_id"], nil),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SceneModuleInfoResponse(*resp)
	return &out, nil
}

func (c *Client) SceneMusic(ctx context.Context, req SceneMusicRequest) (*SceneMusicResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteSceneMusic, Request{
		Method: "POST",
		URL:    "/genesisapi/v1/scene_music/rec_music",
		Params: map[string]any{
			"scene_id": firstAny(params["id"], params["scene_id"]),
			"page":     toInt(firstAny(params["page"], nil), 1),
			"pagesize": toInt(firstAny(params["pagesize"], params["page_size"]), 30),
		},
		Data:        map[string]any{"exposure": []any{}},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SceneMusicResponse(*resp)
	return &out, nil
}

func (c *Client) SceneVideoList(ctx context.Context, req SceneVideoListRequest) (*SceneVideoListResponse, error) {
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
	appid, clientver := config.PlatformConfig(c.isLite)
	userid := firstNonEmpty(asIntString(firstAny(params["userid"], nil)), cookies["userid"], "0")
	token := firstNonEmpty(fmt.Sprintf("%v", firstAny(params["token"], nil)), cookies["token"], "")
	resp, err := c.Call(ctx, RouteSceneVideoList, Request{
		Method: "POST",
		URL:    "/scene/v1/distribution/video_list",
		Data: map[string]any{
			"appid":        appid,
			"clientver":    clientver,
			"token":        token,
			"userid":       userid,
			"tag_id":       firstAny(params["tag_id"], nil),
			"page":         toInt(firstAny(params["page"], nil), 1),
			"page_size":    toInt(firstAny(params["pagesize"], params["page_size"]), 30),
			"exposed_data": []any{},
		},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SceneVideoListResponse(*resp)
	return &out, nil
}

func (c *Client) SearchSuggest(ctx context.Context, req SearchSuggestRequest) (*SearchSuggestResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteSearchSuggest, Request{
		Method: "GET",
		URL:    "/v2/getSearchTip",
		Params: map[string]any{
			"keyword":         firstAny(params["keywords"], params["keyword"]),
			"AlbumTipCount":   toInt(firstAny(params["albumTipCount"], params["album_tip_count"]), 10),
			"CorrectTipCount": toInt(firstAny(params["correctTipCount"], params["correct_tip_count"]), 10),
			"MVTipCount":      toInt(firstAny(params["mvTipCount"], params["mv_tip_count"]), 10),
			"MusicTipCount":   toInt(firstAny(params["musicTipCount"], params["music_tip_count"]), 10),
			"radiotip":        1,
		},
		Headers: map[string]string{"x-router": "searchtip.kugou.com"},
		Cookie:  req.Cookie,
	})
	if err != nil {
		return nil, err
	}
	out := SearchSuggestResponse(*resp)
	return &out, nil
}

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

func (c *Client) UserVipDetail(ctx context.Context, req UserVipDetailRequest) (*UserVipDetailResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, requireLoginCookie(cookies)
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
		return nil, requireLoginCookie(cookies)
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
		return nil, requireLoginCookie(cookies)
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
		return nil, requireLoginCookie(cookies)
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
		return nil, requireLoginCookie(cookies)
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
		return nil, requireLoginCookie(cookies)
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
