package sdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/core/config"
	"github.com/lfhy/kugou-music-api/core/kugou"
)

// Catalog and artist endpoints need hand-built signature payloads.
func (c *Client) AiRecommend(ctx context.Context, req AiRecommendRequest) (*AiRecommendResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	clienttime := time.Now().UnixMilli()
	userid := firstNonEmpty(fmt.Sprintf("%v", req.Userid), cookies["userid"], "0")
	recommendSource := make([]map[string]any, 0)
	for _, s := range splitCSV(fmt.Sprintf("%v", req.AlbumAudioId)) {
		recommendSource = append(recommendSource, map[string]any{"ID": toInt(s, 0)})
	}
	data := map[string]any{
		"platform":         "ios",
		"clientver":        clientver,
		"clienttime":       clienttime,
		"userid":           userid,
		"client_playlist":  []any{},
		"source_type":      2,
		"playlist_ver":     2,
		"area_code":        1,
		"appid":            appid,
		"key":              signParamsKey(strconv.FormatInt(clienttime, 10), c.isLite),
		"mid":              cookies["KUGOU_API_MID"],
		"recommend_source": recommendSource,
	}
	resp, err := c.Call(ctx, RouteAiRecommend, Request{
		Method:             "POST",
		URL:                "/recommend",
		Data:               data,
		Cookie:             cookies,
		EncryptType:        "android",
		ClearDefaultParams: boolPtr(true),
		Headers:            map[string]string{"x-router": "songlistairec.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := AiRecommendResponse(*resp)
	return &out, nil
}

func (c *Client) Album(ctx context.Context, req AlbumRequest) (*AlbumResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	dateTime := time.Now().UnixMilli()
	userid := firstNonEmpty(cookies["userid"], fmt.Sprintf("%v", req.Userid), "0")
	token := firstNonEmpty(cookies["token"], fmt.Sprintf("%v", req.Token), "")
	dataList := make([]map[string]any, 0)
	for _, s := range splitCSV(fmt.Sprintf("%v", req.AlbumId)) {
		dataList = append(dataList, map[string]any{"album_id": s, "album_name": "", "author_name": ""})
	}
	data := map[string]any{
		"appid":      appid,
		"clienttime": dateTime,
		"clientver":  clientver,
		"data":       dataList,
		"dfid":       firstNonEmpty(cookies["dfid"], req.Dfid, "-"),
		"fields":     firstNonEmpty(req.Fields, ""),
		"key":        signParamsKey(strconv.FormatInt(dateTime, 10), c.isLite),
		"mid":        cookies["KUGOU_API_MID"],
	}
	if strings.TrimSpace(token) != "" && token != "0" {
		data["token"] = token
	}
	if strings.TrimSpace(userid) != "" && userid != "0" {
		data["userid"] = userid
	}
	resp, err := c.Call(ctx, RouteAlbum, Request{
		Method:      "POST",
		BaseURL:     "http://kmr.service.kugou.com",
		URL:         "/v1/album",
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
		Headers: map[string]string{
			"x-router":     "kmr.service.kugou.com",
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return nil, err
	}
	out := AlbumResponse(*resp)
	return &out, nil
}

func (c *Client) ArtistAudios(ctx context.Context, req ArtistAudiosRequest) (*ArtistAudiosResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	clienttime := time.Now().Unix()
	data := map[string]any{
		"appid":      appid,
		"clientver":  clientver,
		"mid":        cookies["KUGOU_API_MID"],
		"clienttime": clienttime,
		"key":        signParamsKey(strconv.FormatInt(clienttime, 10), c.isLite),
		"author_id":  req.Id,
		"pagesize":   toInt(req.Pagesize, 30),
		"page":       toInt(req.Page, 1),
		"sort":       ternaryInt(strings.TrimSpace(fmt.Sprintf("%v", req.Sort)) == "hot", 1, 2),
		"area_code":  "all",
	}
	resp, err := c.Call(ctx, RouteArtistAudios, Request{
		Method:      "POST",
		BaseURL:     "https://openapi.kugou.com",
		URL:         "/kmr/v1/audio_group/author",
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "openapi.kugou.com", "kg-tid": "220"},
	})
	if err != nil {
		return nil, err
	}
	out := ArtistAudiosResponse(*resp)
	return &out, nil
}

func (c *Client) ArtistFollow(ctx context.Context, req ArtistFollowRequest) (*ArtistFollowResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	if merged, ok := c.ensureLoginValid(ctx, cookies); ok {
		cookies = merged
	}
	token := firstNonEmpty(fmt.Sprintf("%v", req.Token), cookies["token"], "")
	userid := toInt(firstAny(req.Userid, cookies["userid"]), 0)
	singerid := toInt(req.Id, 0)
	clienttime := time.Now().Unix()
	enc, err := kugou.CryptoAesEncrypt(map[string]any{"singerid": singerid, "token": token}, nil)
	if err != nil {
		return nil, err
	}
	pub := kugou.PublicRASKey
	if c.isLite {
		pub = kugou.PublicLiteRASKey
	}
	p, err := kugou.CryptoRSAEncryptPKCS1Hex(map[string]any{"clienttime": clienttime, "key": enc.Key}, pub)
	if err != nil {
		return nil, err
	}
	resp, err := c.Call(ctx, RouteArtistFollow, Request{
		Method:      "POST",
		URL:         "/followservice/v3/follow_singer",
		Params:      map[string]any{"clienttime": clienttime},
		Data:        map[string]any{"plat": 0, "userid": userid, "singerid": singerid, "source": 7, "p": p, "params": enc.Str},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := ArtistFollowResponse(*resp)
	return &out, nil
}

func (c *Client) ArtistUnfollow(ctx context.Context, req ArtistUnfollowRequest) (*ArtistUnfollowResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	if merged, ok := c.ensureLoginValid(ctx, cookies); ok {
		cookies = merged
	}
	token := firstNonEmpty(fmt.Sprintf("%v", req.Token), cookies["token"], "")
	userid := firstNonEmpty(fmt.Sprintf("%v", req.Userid), cookies["userid"], "0")
	singerid := fmt.Sprintf("%v", req.Id)
	clienttime := time.Now().Unix()
	enc, err := kugou.CryptoAesEncrypt(map[string]any{"singerid": singerid, "token": token}, nil)
	if err != nil {
		return nil, err
	}
	pub := kugou.PublicRASKey
	if c.isLite {
		pub = kugou.PublicLiteRASKey
	}
	p, err := kugou.CryptoRSAEncryptPKCS1Hex(map[string]any{"clienttime": clienttime, "key": enc.Key}, pub)
	if err != nil {
		return nil, err
	}
	resp, err := c.Call(ctx, RouteArtistUnfollow, Request{
		Method:      "POST",
		URL:         "/followservice/v3/unfollow_singer",
		Data:        map[string]any{"plat": 0, "userid": userid, "singerid": singerid, "source": 7, "p": p, "params": enc.Str},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := ArtistUnfollowResponse(*resp)
	return &out, nil
}

func (c *Client) Audio(ctx context.Context, req AudioRequest) (*AudioResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	dateTime := time.Now().UnixMilli()
	userid := firstNonEmpty(cookies["userid"], fmt.Sprintf("%v", req.Userid), "0")
	token := firstNonEmpty(cookies["token"], fmt.Sprintf("%v", req.Token), "")
	dataList := make([]map[string]any, 0)
	for _, s := range splitCSV(fmt.Sprintf("%v", req.Hash)) {
		dataList = append(dataList, map[string]any{"hash": s, "audio_id": 0})
	}
	data := map[string]any{
		"appid":      appid,
		"clienttime": dateTime,
		"clientver":  clientver,
		"data":       dataList,
		"dfid":       firstNonEmpty(cookies["dfid"], req.Dfid, "-"),
		"key":        signParamsKey(strconv.FormatInt(dateTime, 10), c.isLite),
		"mid":        cookies["KUGOU_API_MID"],
	}
	if strings.TrimSpace(token) != "" && token != "0" {
		data["token"] = token
	}
	if strings.TrimSpace(userid) != "" && userid != "0" {
		data["userid"] = userid
	}
	resp, err := c.Call(ctx, RouteAudio, Request{
		Method:      "POST",
		BaseURL:     "http://kmr.service.kugou.com",
		URL:         "/v1/audio/audio",
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
		Headers: map[string]string{
			"x-router":     "kmr.service.kugou.com",
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return nil, err
	}
	out := AudioResponse(*resp)
	return &out, nil
}
