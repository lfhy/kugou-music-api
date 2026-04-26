package sdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/core/config"
	"github.com/lfhy/kugou-music-api/core/kugou"
	"github.com/lfhy/kugou-music-api/core/util"
)

// User listening and video metadata endpoints keep their signature logic here.
func (c *Client) UserCloudUrl(ctx context.Context, req UserCloudUrlRequest) (*UserCloudUrlResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	hash := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", firstAny(params["hash"], req.Hash))))
	pid := "20026"
	paramsMap := map[string]any{
		"hash":           hash,
		"ssa_flag":       "is_fromtrack",
		"version":        "20102",
		"ssl":            0,
		"album_audio_id": toInt(firstAny(params["album_audio_id"], req.AlbumAudioId), 0),
		"pid":            20026,
		"audio_id":       toInt(firstAny(params["audio_id"], req.AudioId), 0),
		"kv_id":          2,
		"key":            signCloudKey(hash, pid),
		"bucket":         "musicclound",
		"name":           firstNonEmpty(fmt.Sprintf("%v", firstAny(params["name"], req.Name)), ""),
		"with_res_tag":   0,
	}
	resp, err := c.Call(ctx, RouteUserCloudUrl, Request{
		Method:      "GET",
		URL:         "/bsstrackercdngz/v2/query_musicclound_url",
		Params:      paramsMap,
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := UserCloudUrlResponse(*resp)
	return &out, nil
}

func (c *Client) UserFollow(ctx context.Context, req UserFollowRequest) (*UserFollowResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	if merged, ok := c.ensureLoginValid(ctx, cookies); ok {
		cookies = merged
	}
	token := firstNonEmpty(fmt.Sprintf("%v", req.Token), cookies["token"], "")
	userid := firstNonEmpty(fmt.Sprintf("%v", req.Userid), cookies["userid"], "0")
	dateTime := time.Now().Unix()
	pub := kugou.PublicRASKey
	if c.isLite {
		pub = kugou.PublicLiteRASKey
	}
	p, err := kugou.CryptoRSAEncryptRawHex(map[string]any{"clienttime": dateTime, "token": token}, pub)
	if err != nil {
		return nil, err
	}
	resp, err := c.Call(ctx, RouteUserFollow, Request{
		Method:      "POST",
		URL:         "/v4/follow_list",
		Data:        map[string]any{"merge": 2, "need_iden_type": 1, "ext_params": "k_pic,jumptype,singerid,score", "userid": userid, "type": 0, "id_type": 0, "p": strings.ToUpper(p)},
		Params:      map[string]any{"plat": 1},
		Cookie:      cookies,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "relationuser.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := UserFollowResponse(*resp)
	return &out, nil
}

func (c *Client) UserListen(ctx context.Context, req UserListenRequest) (*UserListenResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	if merged, ok := c.ensureLoginValid(ctx, cookies); ok {
		cookies = merged
	}
	token := firstNonEmpty(fmt.Sprintf("%v", req.Token), cookies["token"], "")
	userid := firstNonEmpty(fmt.Sprintf("%v", req.Userid), cookies["userid"], "0")
	clienttime := time.Now().Unix()
	pub := kugou.PublicRASKey
	if c.isLite {
		pub = kugou.PublicLiteRASKey
	}
	p, err := kugou.CryptoRSAEncryptRawHex(map[string]any{"clienttime": clienttime, "token": token}, pub)
	if err != nil {
		return nil, err
	}
	resp, err := c.Call(ctx, RouteUserListen, Request{
		Method:      "POST",
		BaseURL:     "https://listenservice.kugou.com",
		URL:         "/v2/get_list",
		Data:        map[string]any{"t_userid": userid, "userid": userid, "list_type": toInt(req.Type, 0), "area_code": 1, "cover": 2, "p": strings.ToUpper(p)},
		Params:      map[string]any{"clienttime": clienttime, "plat": 0},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := UserListenResponse(*resp)
	return &out, nil
}

func (c *Client) VideoDetail(ctx context.Context, req VideoDetailRequest) (*VideoDetailResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	dfid := firstNonEmpty(cookies["dfid"], "-")
	mid := cookies["KUGOU_API_MID"]
	uuid := util.MD5Hex(dfid + mid)
	token := firstNonEmpty(fmt.Sprintf("%v", req.Token), cookies["token"], "")
	clienttime := time.Now().Unix()
	resource := make([]map[string]any, 0)
	for _, s := range splitCSV(fmt.Sprintf("%v", req.Id)) {
		resource = append(resource, map[string]any{"video_id": s})
	}
	data := map[string]any{
		"appid":           appid,
		"clientver":       clientver,
		"clienttime":      clienttime,
		"mid":             mid,
		"uuid":            uuid,
		"dfid":            dfid,
		"token":           token,
		"key":             signParamsKey(strconv.FormatInt(clienttime, 10), c.isLite),
		"show_resolution": 1,
		"data":            resource,
	}
	resp, err := c.Call(ctx, RouteVideoDetail, Request{
		Method:             "POST",
		URL:                "/v1/video",
		Data:               data,
		Cookie:             cookies,
		EncryptType:        "android",
		ClearDefaultParams: boolPtr(true),
		Headers:            map[string]string{"x-router": "kmr.service.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := VideoDetailResponse(*resp)
	return &out, nil
}

func (c *Client) VideoPrivilege(ctx context.Context, req VideoPrivilegeRequest) (*VideoPrivilegeResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	resource := make([]map[string]any, 0)
	for _, s := range splitCSV(fmt.Sprintf("%v", req.Hash)) {
		resource = append(resource, map[string]any{"hash": s, "id": 0, "name": ""})
	}
	data := map[string]any{
		"appid":     appid,
		"area_code": 1,
		"behavior":  "play",
		"clientver": clientver,
		"dfid":      firstNonEmpty(cookies["dfid"], "-"),
		"mid":       cookies["KUGOU_API_MID"],
		"resource":  resource,
		"token":     firstNonEmpty(cookies["token"], ""),
		"userid":    firstNonEmpty(cookies["userid"], "0"),
		"vip":       firstNonEmpty(cookies["vip_type"], "0"),
	}
	resp, err := c.Call(ctx, RouteVideoPrivilege, Request{
		Method:      "POST",
		URL:         "/v1/get_video_privilege",
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "media.store.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := VideoPrivilegeResponse(*resp)
	return &out, nil
}
