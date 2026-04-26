package sdk

import (
	"context"
	"fmt"

	"github.com/lfhy/kugou-music-api/core/config"
)

// Scene and search endpoints need manual parameter normalization.

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
