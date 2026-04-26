package sdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/core/config"
	"github.com/lfhy/kugou-music-api/core/util"
)

// Search and FM endpoints use upstream signature flows that differ from generic wrappers.
func (c *Client) SearchMixed(ctx context.Context, req SearchMixedRequest) (*SearchMixedResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	t := time.Now().UnixMilli()
	keyword := firstNonEmpty(fmt.Sprintf("%v", firstAny(params["keyword"], req.Keyword)), "")
	dataMap := map[string]any{
		"ab_tag":           0,
		"ability":          511,
		"albumhide":        0,
		"apiver":           22,
		"area_code":        1,
		"clientver":        20125,
		"cursor":           0,
		"is_gpay":          0,
		"iscorrection":     1,
		"keyword":          keyword,
		"nocollect":        0,
		"osversion":        16.5,
		"platform":         "IOSFilter",
		"recver":           2,
		"req_ai":           1,
		"requestid":        util.MD5Hex("bdaa53d04e7475feb9024164a47032f9"+strconv.FormatInt(t, 10)) + "_0",
		"search_ability":   3,
		"sec_aggre":        1,
		"sec_aggre_bitmap": 0,
		"style_type":       3,
		"tag":              "em",
	}
	resp, err := c.Call(ctx, RouteSearchMixed, Request{
		Method: "GET",
		URL:    "/v3/search/mixed",
		Params: dataMap,
		Headers: map[string]string{
			"x-router":        "complexsearch.kugou.com",
			"kg-clienttimems": strconv.FormatInt(t, 10),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SearchMixedResponse(*resp)
	return &out, nil
}

func (c *Client) FmClass(ctx context.Context, req FmClassRequest) (*FmClassResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	dateTime := time.Now().UnixMilli()
	appid, clientver := config.PlatformConfig(c.isLite)
	userid := firstNonEmpty(cookies["userid"], asIntString(firstAny(req.Userid, 0)), "0")
	data := map[string]any{
		"kguid":       userid,
		"clienttime":  dateTime,
		"mid":         cookies["KUGOU_API_MID"],
		"platform":    "android",
		"clientver":   clientver,
		"uid":         userid,
		"get_tracker": 1,
		"key":         signParamsKey(strconv.FormatInt(dateTime, 10), c.isLite),
		"appid":       appid,
	}
	resp, err := c.Call(ctx, RouteFmClass, Request{
		Method:      "POST",
		URL:         "/v1/class_fm_song",
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "fm.service.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := FmClassResponse(*resp)
	return &out, nil
}

func (c *Client) FmImage(ctx context.Context, req FmImageRequest) (*FmImageResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	dateTime := time.Now().UnixMilli()
	appid, clientver := config.PlatformConfig(c.isLite)
	userid := firstNonEmpty(cookies["userid"], fmt.Sprintf("%v", req.Userid))
	token := firstNonEmpty(cookies["token"], fmt.Sprintf("%v", req.Token))
	dfid := firstNonEmpty(cookies["dfid"], req.Dfid, "-")

	fmData := make([]map[string]any, 0)
	for _, s := range splitCSV(req.Fmid) {
		fmData = append(fmData, map[string]any{"fields": "imgUrl100,imgUrl50", "fmid": s, "fmtype": 2})
	}
	data := map[string]any{
		"appid":      appid,
		"clienttime": dateTime,
		"clientver":  clientver,
		"data":       fmData,
		"dfid":       dfid,
		"key":        signParamsKey(strconv.FormatInt(dateTime, 10), c.isLite),
		"mid":        cookies["KUGOU_API_MID"],
	}
	if strings.TrimSpace(userid) != "" && userid != "<nil>" {
		data["userid"] = userid
	}
	if strings.TrimSpace(token) != "" && token != "<nil>" {
		data["token"] = token
	}
	resp, err := c.Call(ctx, RouteFmImage, Request{
		Method:      "POST",
		URL:         "/v1/fm_info",
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
		Headers: map[string]string{
			"x-router":     "fm.service.kugou.com",
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return nil, err
	}
	out := FmImageResponse(*resp)
	return &out, nil
}

func (c *Client) FmRecommend(ctx context.Context, req FmRecommendRequest) (*FmRecommendResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	dateTime := time.Now().UnixMilli()
	appid, clientver := config.PlatformConfig(c.isLite)
	data := map[string]any{
		"appid":         appid,
		"clientver":     clientver,
		"clienttime":    dateTime,
		"mid":           cookies["KUGOU_API_MID"],
		"key":           signParamsKey(strconv.FormatInt(dateTime, 10), c.isLite),
		"rcmdsongcount": 1,
		"level":         0,
		"area_code":     1,
		"get_tracker":   1,
		"uid":           0,
	}
	resp, err := c.Call(ctx, RouteFmRecommend, Request{
		Method:      "POST",
		URL:         "/v1/rcmd_list",
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "fm.service.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := FmRecommendResponse(*resp)
	return &out, nil
}

func (c *Client) FmSongs(ctx context.Context, req FmSongsRequest) (*FmSongsResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	dateTime := time.Now().UnixMilli()
	appid, clientver := config.PlatformConfig(c.isLite)
	userid := firstNonEmpty(cookies["userid"], fmt.Sprintf("%v", req.Userid))

	fmids := splitCSV(req.Fmid)
	fmtypes := splitCSV(req.Fmtype)
	fmoffsets := splitCSV(req.Fmoffset)
	fmsizes := splitCSV(req.Fmsize)

	fmTypeDef := toInt(req.Type, 2)
	if fmTypeDef == 0 {
		fmTypeDef = 2
	}
	offsetDef := toInt(req.Offset, -1)
	sizeDef := req.Size
	if sizeDef == 0 {
		sizeDef = 20
	}

	fmData := make([]map[string]any, 0, len(fmids))
	for i, id := range fmids {
		item := map[string]any{
			"fmid":       id,
			"fmtype":     fmTypeDef,
			"offset":     offsetDef,
			"size":       sizeDef,
			"singername": "",
		}
		if i < len(fmtypes) && fmtypes[i] != "" {
			item["fmtype"] = fmtypes[i]
		}
		if i < len(fmoffsets) && fmoffsets[i] != "" {
			item["offset"] = fmoffsets[i]
		}
		if i < len(fmsizes) && fmsizes[i] != "" {
			item["size"] = fmsizes[i]
		}
		fmData = append(fmData, item)
	}

	data := map[string]any{
		"appid":       appid,
		"area_code":   1,
		"clienttime":  dateTime,
		"clientver":   clientver,
		"data":        fmData,
		"get_tracker": 1,
		"key":         signParamsKey(strconv.FormatInt(dateTime, 10), c.isLite),
		"mid":         cookies["KUGOU_API_MID"],
		"uid":         userid,
	}
	resp, err := c.Call(ctx, RouteFmSongs, Request{
		Method:      "POST",
		URL:         "/v1/app_song_list_offset",
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
		Headers: map[string]string{
			"x-router":     "fm.service.kugou.com",
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return nil, err
	}
	out := FmSongsResponse(*resp)
	return &out, nil
}
