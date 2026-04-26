package sdk

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/core/config"
	"github.com/lfhy/kugou-music-api/core/kugou"
)

// Playlist and cloud endpoints require custom encryption and signature steps.
func (c *Client) PlaylistDel(ctx context.Context, req PlaylistDelRequest) (*PlaylistDelResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	if merged, ok := c.ensureLoginValid(ctx, cookies); ok {
		cookies = merged
	}
	userid := firstNonEmpty(cookies["userid"], fmt.Sprintf("%v", req.Userid), "0")
	token := firstNonEmpty(cookies["token"], fmt.Sprintf("%v", req.Token), "")
	clienttime := time.Now().Unix()
	aesData, err := playlistAesEncrypt(map[string]any{
		"listid":    toInt(req.Listid, 0),
		"total_ver": 0,
		"type":      1,
	})
	if err != nil {
		return nil, err
	}
	pub := kugou.PublicRASKey
	if c.isLite {
		pub = kugou.PublicLiteRASKey
	}
	p, err := kugou.CryptoRSAEncryptPKCS1Hex(map[string]any{"aes": aesData.Key, "uid": userid, "token": token}, pub)
	if err != nil {
		return nil, err
	}
	params := map[string]any{
		"clienttime": clienttime,
		"key":        signParamsKey(strconv.FormatInt(clienttime, 10), c.isLite),
		"last_area":  "gztx",
		"clientver":  currentClientVer(c.isLite),
		"appid":      currentAppid(c.isLite),
		"last_time":  clienttime,
		"p":          strings.ToUpper(p),
	}
	raw, err := c.core.CreateRequest(ctx, kugou.RequestConfig{
		Method:      "POST",
		URL:         "/v2/delete_list",
		BaseURL:     "",
		Params:      params,
		Data:        aesData.CipherBase64,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "cloudlist.service.kugou.com"},
		Cookie:      cookies,
	})
	out := &Response{Status: raw.Status, RawBody: raw.Body, Headers: raw.Headers, Cookie: raw.Cookie}
	if decoded, derr := playlistAesDecryptFromRaw(raw.Body, aesData.Key); derr == nil {
		out.Body = decoded
		if b, jerr := json.Marshal(decoded); jerr == nil {
			out.RawBody = b
		}
	}
	return (*PlaylistDelResponse)(out), err
}

func (c *Client) PlaylistSimilar(ctx context.Context, req PlaylistSimilarRequest) (*PlaylistSimilarResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	clienttime := time.Now().UnixMilli()
	dataList := make([]map[string]any, 0)
	for _, s := range splitCSV(req.Ids) {
		dataList = append(dataList, map[string]any{"global_collection_id": s})
	}
	data := map[string]any{
		"appid":      appid,
		"clientver":  clientver,
		"clienttime": clienttime,
		"key":        signParamsKey(strconv.FormatInt(clienttime, 10), c.isLite),
		"userid":     firstNonEmpty(cookies["userid"], fmt.Sprintf("%v", req.Userid), "0"),
		"ugc":        1,
		"show_list":  1,
		"need_songs": 1,
		"data":       dataList,
	}
	resp, err := c.Call(ctx, RoutePlaylistSimilar, Request{
		Method:      "POST",
		URL:         "/pubsongs/v1/kmr_get_similar_lists",
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := PlaylistSimilarResponse(*resp)
	return &out, nil
}

func (c *Client) TopCard(ctx context.Context, req TopCardRequest) (*TopCardResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	dateTime := time.Now().UnixMilli()
	fakem := "60f7ebf1f812edbac3c63a7310001701760f"
	data := map[string]any{
		"appid":           appid,
		"clientver":       clientver,
		"platform":        "android",
		"clienttime":      dateTime,
		"userid":          firstNonEmpty(cookies["userid"], fmt.Sprintf("%v", req.Userid), "0"),
		"key":             signParamsKey(strconv.FormatInt(dateTime, 10), c.isLite),
		"fakem":           fakem,
		"area_code":       1,
		"mid":             cookies["KUGOU_API_MID"],
		"uuid":            "-",
		"client_playlist": []any{},
		"u_info":          "a0c35cd40af564444b5584c2754dedec",
	}
	resp, err := c.Call(ctx, RouteTopCard, Request{
		Method:      "POST",
		URL:         "/singlecardrec.service/v1/single_card_recommend",
		Data:        data,
		Params:      map[string]any{"card_id": toInt(req.CardId, 1), "fakem": fakem, "area_code": 1, "platform": "ios"},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := TopCardResponse(*resp)
	return &out, nil
}

func (c *Client) TopPlaylist(ctx context.Context, req TopPlaylistRequest) (*TopPlaylistResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	dateTime := time.Now().Unix()
	specialRecommend := map[string]any{
		"withtag":       toInt(req.Withtag, 1),
		"withsong":      toInt(req.Withsong, 1),
		"sort":          toInt(req.Sort, 1),
		"ugc":           1,
		"is_selected":   0,
		"withrecommend": 1,
		"area_code":     1,
		"categoryid":    toInt(req.CategoryId, 0),
	}
	data := map[string]any{
		"appid":               appid,
		"mid":                 cookies["KUGOU_API_MID"],
		"clientver":           clientver,
		"platform":            "android",
		"clienttime":          dateTime,
		"userid":              firstNonEmpty(cookies["userid"], fmt.Sprintf("%v", req.Userid), "0"),
		"module_id":           toInt(req.ModuleId, 1),
		"page":                toInt(req.Page, 1),
		"pagesize":            toInt(req.Pagesize, 30),
		"key":                 signParamsKey(strconv.FormatInt(dateTime, 10), c.isLite),
		"special_recommend":   specialRecommend,
		"req_multi":           1,
		"retrun_min":          5,
		"return_special_falg": 1,
	}
	resp, err := c.Call(ctx, RouteTopPlaylist, Request{
		Method:      "POST",
		URL:         "/v2/special_recommend",
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "specialrec.service.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := TopPlaylistResponse(*resp)
	return &out, nil
}

func (c *Client) UserCloud(ctx context.Context, req UserCloudRequest) (*UserCloudResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	if merged, ok := c.ensureLoginValid(ctx, cookies); ok {
		cookies = merged
	}
	userid := firstNonEmpty(cookies["userid"], fmt.Sprintf("%v", req.Userid), "0")
	token := firstNonEmpty(cookies["token"], fmt.Sprintf("%v", req.Token), "")
	mid := cookies["KUGOU_API_MID"]
	clienttime := time.Now().Unix()
	aesData, err := playlistAesEncrypt(map[string]any{"page": toInt(req.Page, 1), "pagesize": toInt(req.Pagesize, 30), "getkmr": 1})
	if err != nil {
		return nil, err
	}
	pub := kugou.PublicRASKey
	if c.isLite {
		pub = kugou.PublicLiteRASKey
	}
	p, err := kugou.CryptoRSAEncryptPKCS1Hex(map[string]any{"aes": aesData.Key, "uid": userid, "token": token}, pub)
	if err != nil {
		return nil, err
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	params := map[string]any{
		"clienttime": clienttime,
		"mid":        mid,
		"key":        signParamsKey(strconv.FormatInt(clienttime, 10), c.isLite),
		"clientver":  clientver,
		"appid":      appid,
		"p":          strings.ToUpper(p),
	}
	dataRaw, _ := base64.StdEncoding.DecodeString(aesData.CipherBase64)
	raw, err := c.core.CreateRequest(ctx, kugou.RequestConfig{
		Method:             "POST",
		BaseURL:            "https://mcloudservice.kugou.com",
		URL:                "/v1/get_list",
		Params:             params,
		Data:               dataRaw,
		EncryptType:        "android",
		Cookie:             cookies,
		ClearDefaultParams: true,
		NotSignature:       true,
	})
	out := &Response{Status: raw.Status, RawBody: raw.Body, Headers: raw.Headers, Cookie: raw.Cookie}
	if decoded, derr := playlistAesDecryptFromRaw(raw.Body, aesData.Key); derr == nil {
		out.Body = decoded
		if b, jerr := json.Marshal(decoded); jerr == nil {
			out.RawBody = b
		}
	}
	return (*UserCloudResponse)(out), err
}
