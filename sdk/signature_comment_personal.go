package sdk

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lfhy/kugou-music-api/core/config"
)

// Comment and personal FM endpoints stay grouped by shared request signing rules.
func (c *Client) CommentAlbum(ctx context.Context, req CommentAlbumRequest) (*CommentAlbumResponse, error) {
	resp, err := c.commentCommon(ctx, RouteCommentAlbum, map[string]any{
		"childrenid":        req.Id,
		"need_show_image":   1,
		"p":                 toInt(req.Page, 1),
		"pagesize":          toInt(req.Pagesize, 30),
		"show_classify":     toInt(req.ShowClassify, 1),
		"show_hotword_list": toInt(req.ShowHotwordList, 1),
		"code":              "94f1792ced1df89aa68a7939eaf2efca",
	}, req.Cookie)
	if err != nil {
		return nil, err
	}
	out := CommentAlbumResponse(*resp)
	return &out, nil
}

func (c *Client) CommentFloor(ctx context.Context, req CommentFloorRequest) (*CommentFloorResponse, error) {
	resp, err := c.commentCommon(ctx, RouteCommentFloor, map[string]any{
		"childrenid":        req.SpecialId,
		"mixsongid":         req.Mixsongid,
		"need_show_image":   1,
		"p":                 toInt(req.Page, 1),
		"pagesize":          toInt(req.Pagesize, 30),
		"show_classify":     toInt(req.ShowClassify, 1),
		"show_hotword_list": toInt(req.ShowHotwordList, 1),
		"code":              "fc4be23b4e972707f36b8a828a93ba8a",
		"tid":               req.Tid,
	}, req.Cookie)
	if err != nil {
		return nil, err
	}
	out := CommentFloorResponse(*resp)
	return &out, nil
}

func (c *Client) CommentMusic(ctx context.Context, req CommentMusicRequest) (*CommentMusicResponse, error) {
	resp, err := c.commentCommon(ctx, RouteCommentMusic, map[string]any{
		"mixsongid":         req.Mixsongid,
		"need_show_image":   1,
		"p":                 toInt(req.Page, 1),
		"pagesize":          toInt(req.Pagesize, 30),
		"show_classify":     toInt(req.ShowClassify, 1),
		"show_hotword_list": toInt(req.ShowHotwordList, 1),
		"extdata":           "0",
		"code":              "fc4be23b4e972707f36b8a828a93ba8a",
	}, req.Cookie)
	if err != nil {
		return nil, err
	}
	out := CommentMusicResponse(*resp)
	return &out, nil
}

func (c *Client) CommentMusicHotword(ctx context.Context, req CommentMusicHotwordRequest) (*CommentMusicHotwordResponse, error) {
	resp, err := c.commentCommon(ctx, RouteCommentMusicHotword, map[string]any{
		"mixsongid":       req.Mixsongid,
		"need_show_image": 1,
		"p":               toInt(req.Page, 1),
		"pagesize":        toInt(req.Pagesize, 30),
		"hot_word":        req.HotWord,
		"extdata":         "0",
		"code":            "fc4be23b4e972707f36b8a828a93ba8a",
	}, req.Cookie)
	if err != nil {
		return nil, err
	}
	out := CommentMusicHotwordResponse(*resp)
	return &out, nil
}

func (c *Client) CommentPlaylist(ctx context.Context, req CommentPlaylistRequest) (*CommentPlaylistResponse, error) {
	resp, err := c.commentCommon(ctx, RouteCommentPlaylist, map[string]any{
		"childrenid":        req.Id,
		"need_show_image":   1,
		"p":                 toInt(req.Page, 1),
		"pagesize":          toInt(req.Pagesize, 30),
		"show_classify":     toInt(req.ShowClassify, 1),
		"show_hotword_list": toInt(req.ShowHotwordList, 1),
		"code":              "ca53b96fe5a1d9c22d71c8f522ef7c4f",
		"content_type":      0,
		"tag":               5,
	}, req.Cookie)
	if err != nil {
		return nil, err
	}
	out := CommentPlaylistResponse(*resp)
	return &out, nil
}

func (c *Client) PersonalFm(ctx context.Context, req PersonalFmRequest) (*PersonalFmResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	if merged, ok := c.ensureLoginValid(ctx, cookies); ok {
		cookies = merged
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	dateTime := time.Now().UnixMilli()
	userid := firstNonEmpty(cookies["userid"], fmt.Sprintf("%v", req.Userid), "0")
	token := firstNonEmpty(cookies["token"], fmt.Sprintf("%v", req.Token), "")
	vipType := firstNonEmpty(cookies["vip_type"], fmt.Sprintf("%v", req.VipType), "0")
	data := map[string]any{
		"appid":                   appid,
		"clienttime":              dateTime,
		"mid":                     cookies["KUGOU_API_MID"],
		"action":                  firstNonEmpty(fmt.Sprintf("%v", req.Action), "play"),
		"recommend_source_locked": 0,
		"song_pool_id":            toInt(req.SongPoolId, 0),
		"callerid":                0,
		"m_type":                  1,
		"platform":                firstNonEmpty(fmt.Sprintf("%v", req.Platform), "ios"),
		"area_code":               1,
		"remain_songcnt":          toInt(req.RemainSongcnt, 0),
		"clientver":               clientver,
		"is_overplay":             ternaryInt(toBool(req.IsOverplay, false), 1, 0),
		"mode":                    firstNonEmpty(fmt.Sprintf("%v", req.Mode), "normal"),
		"fakem":                   "ca981cfc583a4c37f28d2d49000013c16a0a",
		"key":                     signParamsKey(strconv.FormatInt(dateTime, 10), c.isLite),
	}
	if userid != "0" {
		data["userid"] = userid
		data["kguid"] = userid
	}
	if token != "" && token != "0" {
		data["token"] = token
	}
	if vipType != "0" {
		data["vip_type"] = vipType
	}
	if fmt.Sprintf("%v", req.Hash) != "" {
		data["hash"] = req.Hash
	}
	if fmt.Sprintf("%v", req.Songid) != "" {
		data["songid"] = req.Songid
	}
	if fmt.Sprintf("%v", req.Playtime) != "" {
		data["playtime"] = req.Playtime
	}
	resp, err := c.Call(ctx, RoutePersonalFm, Request{
		Method:      "POST",
		URL:         "/v2/personal_recommend",
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "persnfm.service.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := PersonalFmResponse(*resp)
	return &out, nil
}
