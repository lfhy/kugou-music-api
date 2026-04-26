package sdk

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/core/config"
	"github.com/lfhy/kugou-music-api/core/kugou"
	"github.com/lfhy/kugou-music-api/core/util"
)

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
		return nil, requireLoginCookie(cookies)
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
		"userid":                   userid,
		"appid":                    appid,
		"playlist_ver":             2,
		"clienttime":               dateTime,
		"mid":                      cookies["KUGOU_API_MID"],
		"new_sync_point":           dateTime,
		"module_id":                1,
		"action":                   "login",
		"vip_type":                 vipType,
		"vip_flags":                3,
		"recommend_source_locked":  0,
		"song_pool_id":             toInt(firstAny(params["song_pool_id"], req.SongPoolId), 0),
		"callerid":                 0,
		"m_type":                   1,
		"kguid":                    userid,
		"platform":                 "ios",
		"area_code":                1,
		"fakem":                    "ca981cfc583a4c37f28d2d49000013c16a0a",
		"clientver":                11850,
		"mode":                     mode,
		"active_swtich":            "on",
		"key":                      signParamsKey(strconv.FormatInt(dateTime, 10), c.isLite),
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

func (c *Client) RegisterDev(ctx context.Context, req RegisterDevRequest) (*RegisterDevResponse, error) {
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

	userid := firstNonEmpty(fmt.Sprintf("%v", firstAny(params["userid"], req.Userid)), cookies["userid"], "0")
	token := firstNonEmpty(fmt.Sprintf("%v", firstAny(params["token"], req.Token)), cookies["token"], "")
	guid := firstNonEmpty(cookies["KUGOU_API_GUID"], cookies["mid"])

	dataMap := map[string]any{
		"availableRamSize":   toInt(firstAny(params["availableRamSize"], req.AvailableRamSize), 4983533568),
		"availableRomSize":   toInt(firstAny(params["availableRomSize"], req.AvailableRomSize), 48114719),
		"availableSDSize":    toInt(firstAny(params["availableSDSize"], req.AvailableSDSize), 48114717),
		"basebandVer":        firstNonEmpty(fmt.Sprintf("%v", firstAny(params["basebandVer"], req.BasebandVer)), ""),
		"batteryLevel":       toInt(firstAny(params["batteryLevel"], req.BatteryLevel), 100),
		"batteryStatus":      toInt(firstAny(params["batteryStatus"], req.BatteryStatus), 3),
		"brand":              firstNonEmpty(fmt.Sprintf("%v", firstAny(params["brand"], req.Brand)), "Redmi"),
		"buildSerial":        firstNonEmpty(fmt.Sprintf("%v", firstAny(params["buildSerial"], req.BuildSerial)), "unknown"),
		"device":             firstNonEmpty(fmt.Sprintf("%v", firstAny(params["device"], req.Device)), "marble"),
		"imei":               firstNonEmpty(fmt.Sprintf("%v", firstAny(params["imei"], req.Imei)), guid),
		"imsi":               firstNonEmpty(fmt.Sprintf("%v", firstAny(params["imsi"], req.Imsi)), ""),
		"manufacturer":       firstNonEmpty(fmt.Sprintf("%v", firstAny(params["manufacturer"], req.Manufacturer)), "Xiaomi"),
		"uuid":               firstNonEmpty(fmt.Sprintf("%v", firstAny(params["uuid"], req.Uuid)), guid),
		"accelerometer":      toBool(firstAny(params["accelerometer"], req.Accelerometer), false),
		"accelerometerValue": firstNonEmpty(fmt.Sprintf("%v", firstAny(params["accelerometerValue"], req.AccelerometerValue)), ""),
		"gravity":            toBool(firstAny(params["gravity"], req.Gravity), false),
		"gravityValue":       firstNonEmpty(fmt.Sprintf("%v", firstAny(params["gravityValue"], req.GravityValue)), ""),
		"gyroscope":          toBool(firstAny(params["gyroscope"], req.Gyroscope), false),
		"gyroscopeValue":     firstNonEmpty(fmt.Sprintf("%v", firstAny(params["gyroscopeValue"], req.GyroscopeValue)), ""),
		"light":              toBool(firstAny(params["light"], req.Light), false),
		"lightValue":         firstNonEmpty(fmt.Sprintf("%v", firstAny(params["lightValue"], req.LightValue)), ""),
		"magnetic":           toBool(firstAny(params["magnetic"], req.Magnetic), false),
		"magneticValue":      firstNonEmpty(fmt.Sprintf("%v", firstAny(params["magneticValue"], req.MagneticValue)), ""),
		"orientation":        toBool(firstAny(params["orientation"], req.Orientation), false),
		"orientationValue":   firstNonEmpty(fmt.Sprintf("%v", firstAny(params["orientationValue"], req.OrientationValue)), ""),
		"pressure":           toBool(firstAny(params["pressure"], req.Pressure), false),
		"pressureValue":      firstNonEmpty(fmt.Sprintf("%v", firstAny(params["pressureValue"], req.PressureValue)), ""),
		"step_counter":       toBool(firstAny(params["step_counter"], req.StepCounter), false),
		"step_counterValue":  firstNonEmpty(fmt.Sprintf("%v", firstAny(params["step_counterValue"], req.StepCounterValue)), ""),
		"temperature":        toBool(firstAny(params["temperature"], req.Temperature), false),
		"temperatureValue":   firstNonEmpty(fmt.Sprintf("%v", firstAny(params["temperatureValue"], req.TemperatureValue)), ""),
	}

	enc, err := playlistAesEncrypt(dataMap)
	if err != nil {
		return nil, err
	}

	pubKey := kugou.PublicRASKey
	if c.isLite {
		pubKey = kugou.PublicLiteRASKey
	}
	p, err := kugou.CryptoRSAEncryptPKCS1Hex(map[string]any{"aes": enc.Key, "uid": userid, "token": token}, pubKey)
	if err != nil {
		return nil, err
	}

	raw, err := c.core.CreateRequest(ctx, kugou.RequestConfig{
		Method:      "POST",
		BaseURL:     "https://userservice.kugou.com",
		URL:         "/risk/v2/r_register_dev",
		Data:        enc.CipherBase64,
		Params:      map[string]any{"part": 1, "platid": 1, "p": p},
		EncryptType: "android",
		Cookie:      cookies,
	})

	out := &Response{Status: raw.Status, RawBody: raw.Body, Headers: raw.Headers, Cookie: raw.Cookie}
	if len(raw.Body) > 0 {
		decoded, derr := playlistAesDecryptFromRaw(raw.Body, enc.Key)
		if derr == nil {
			out.Body = decoded
			if b, jerr := json.Marshal(decoded); jerr == nil {
				out.RawBody = b
			}
			if status, _ := decoded["status"].(float64); int(status) == 1 {
				if dm, ok := decoded["data"].(map[string]any); ok {
					if dfid := strings.TrimSpace(fmt.Sprintf("%v", dm["dfid"])); dfid != "" && dfid != "<nil>" {
						out.Cookie = append(out.Cookie, "dfid="+dfid)
					}
				}
			}
		}
	}
	if len(out.Cookie) > 0 {
		clean := dedupSetCookie(out.Cookie)
		out.Cookie = clean
		c.updateCookiePool(clean)
	}
	return (*RegisterDevResponse)(out), err
}

func (c *Client) UserVideoCollect(ctx context.Context, req UserVideoCollectRequest) (*UserVideoCollectResponse, error) {
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

	token := firstNonEmpty(fmt.Sprintf("%v", firstAny(params["token"], req.Token)), cookies["token"], "")
	userid := firstNonEmpty(fmt.Sprintf("%v", firstAny(params["userid"], req.Userid)), cookies["userid"], "0")
	dataMap := map[string]any{
		"userid":   userid,
		"token":    token,
		"page":     toInt(firstAny(params["page"], req.Page), 1),
		"pagesize": toInt(firstAny(params["pagesize"], req.Pagesize), 30),
	}

	resp, err := c.Call(ctx, RouteUserVideoCollect, Request{
		Method:      "POST",
		URL:         "/collectservice/v2/collect_list_mixvideo",
		Data:        dataMap,
		Params:      map[string]any{"plat": 1},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := UserVideoCollectResponse(*resp)
	return &out, nil
}

func (c *Client) UserVideoLove(ctx context.Context, req UserVideoLoveRequest) (*UserVideoLoveResponse, error) {
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

	userid := firstNonEmpty(fmt.Sprintf("%v", firstAny(params["userid"], req.Userid)), cookies["userid"], "0")
	pagesize := toInt(firstAny(params["pagesize"], req.Pagesize), 30)
	resp, err := c.Call(ctx, RouteUserVideoLove, Request{
		Method: "GET",
		URL:    "/m.comment.service/v1/get_user_like_video",
		Params: map[string]any{
			"kugouid":         userid,
			"pagesize":        pagesize,
			"load_video_info": 1,
			"p":               1,
			"plat":            1,
		},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := UserVideoLoveResponse(*resp)
	return &out, nil
}

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
		"ab_tag":          0,
		"ability":         511,
		"albumhide":       0,
		"apiver":          22,
		"area_code":       1,
		"clientver":       20125,
		"cursor":          0,
		"is_gpay":         0,
		"iscorrection":    1,
		"keyword":         keyword,
		"nocollect":       0,
		"osversion":       16.5,
		"platform":        "IOSFilter",
		"recver":          2,
		"req_ai":          1,
		"requestid":       util.MD5Hex("bdaa53d04e7475feb9024164a47032f9"+strconv.FormatInt(t, 10)) + "_0",
		"search_ability":  3,
		"sec_aggre":       1,
		"sec_aggre_bitmap": 0,
		"style_type":      3,
		"tag":             "em",
	}
	resp, err := c.Call(ctx, RouteSearchMixed, Request{
		Method: "GET",
		URL:    "/v3/search/mixed",
		Params: dataMap,
		Headers: map[string]string{
			"x-router":       "complexsearch.kugou.com",
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
		"kguid":      userid,
		"clienttime": dateTime,
		"mid":        cookies["KUGOU_API_MID"],
		"platform":   "android",
		"clientver":  clientver,
		"uid":        userid,
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

func (c *Client) CommentAlbum(ctx context.Context, req CommentAlbumRequest) (*CommentAlbumResponse, error) {
	resp, err := c.commentCommon(ctx, RouteCommentAlbum, map[string]any{
		"childrenid":       req.Id,
		"need_show_image":  1,
		"p":                toInt(req.Page, 1),
		"pagesize":         toInt(req.Pagesize, 30),
		"show_classify":    toInt(req.ShowClassify, 1),
		"show_hotword_list": toInt(req.ShowHotwordList, 1),
		"code":             "94f1792ced1df89aa68a7939eaf2efca",
	}, req.Cookie)
	if err != nil {
		return nil, err
	}
	out := CommentAlbumResponse(*resp)
	return &out, nil
}

func (c *Client) CommentFloor(ctx context.Context, req CommentFloorRequest) (*CommentFloorResponse, error) {
	resp, err := c.commentCommon(ctx, RouteCommentFloor, map[string]any{
		"childrenid":       req.SpecialId,
		"mixsongid":        req.Mixsongid,
		"need_show_image":  1,
		"p":                toInt(req.Page, 1),
		"pagesize":         toInt(req.Pagesize, 30),
		"show_classify":    toInt(req.ShowClassify, 1),
		"show_hotword_list": toInt(req.ShowHotwordList, 1),
		"code":             "fc4be23b4e972707f36b8a828a93ba8a",
		"tid":              req.Tid,
	}, req.Cookie)
	if err != nil {
		return nil, err
	}
	out := CommentFloorResponse(*resp)
	return &out, nil
}

func (c *Client) CommentMusic(ctx context.Context, req CommentMusicRequest) (*CommentMusicResponse, error) {
	resp, err := c.commentCommon(ctx, RouteCommentMusic, map[string]any{
		"mixsongid":        req.Mixsongid,
		"need_show_image":  1,
		"p":                toInt(req.Page, 1),
		"pagesize":         toInt(req.Pagesize, 30),
		"show_classify":    toInt(req.ShowClassify, 1),
		"show_hotword_list": toInt(req.ShowHotwordList, 1),
		"extdata":          "0",
		"code":             "fc4be23b4e972707f36b8a828a93ba8a",
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
		"childrenid":       req.Id,
		"need_show_image":  1,
		"p":                toInt(req.Page, 1),
		"pagesize":         toInt(req.Pagesize, 30),
		"show_classify":    toInt(req.ShowClassify, 1),
		"show_hotword_list": toInt(req.ShowHotwordList, 1),
		"code":             "ca53b96fe5a1d9c22d71c8f522ef7c4f",
		"content_type":     0,
		"tag":              5,
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
		"appid":                  appid,
		"clienttime":             dateTime,
		"mid":                    cookies["KUGOU_API_MID"],
		"action":                 firstNonEmpty(fmt.Sprintf("%v", req.Action), "play"),
		"recommend_source_locked": 0,
		"song_pool_id":           toInt(req.SongPoolId, 0),
		"callerid":               0,
		"m_type":                 1,
		"platform":               firstNonEmpty(fmt.Sprintf("%v", req.Platform), "ios"),
		"area_code":              1,
		"remain_songcnt":         toInt(req.RemainSongcnt, 0),
		"clientver":              clientver,
		"is_overplay":            ternaryInt(toBool(req.IsOverplay, false), 1, 0),
		"mode":                   firstNonEmpty(fmt.Sprintf("%v", req.Mode), "normal"),
		"fakem":                  "ca981cfc583a4c37f28d2d49000013c16a0a",
		"key":                    signParamsKey(strconv.FormatInt(dateTime, 10), c.isLite),
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
		"appid":       appid,
		"clientver":   clientver,
		"clienttime":  clienttime,
		"key":         signParamsKey(strconv.FormatInt(clienttime, 10), c.isLite),
		"userid":      firstNonEmpty(cookies["userid"], fmt.Sprintf("%v", req.Userid), "0"),
		"ugc":         1,
		"show_list":   1,
		"need_songs":  1,
		"data":        dataList,
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
		"withtag":            toInt(req.Withtag, 1),
		"withsong":           toInt(req.Withsong, 1),
		"sort":               toInt(req.Sort, 1),
		"ugc":                1,
		"is_selected":        0,
		"withrecommend":      1,
		"area_code":          1,
		"categoryid":         toInt(req.CategoryId, 0),
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
		"hash":         hash,
		"ssa_flag":     "is_fromtrack",
		"version":      "20102",
		"ssl":          0,
		"album_audio_id": toInt(firstAny(params["album_audio_id"], req.AlbumAudioId), 0),
		"pid":          20026,
		"audio_id":     toInt(firstAny(params["audio_id"], req.AudioId), 0),
		"kv_id":        2,
		"key":          signCloudKey(hash, pid),
		"bucket":       "musicclound",
		"name":         firstNonEmpty(fmt.Sprintf("%v", firstAny(params["name"], req.Name)), ""),
		"with_res_tag": 0,
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

func (c *Client) commentCommon(ctx context.Context, route string, params map[string]any, cookie map[string]string) (*Response, error) {
	resp, err := c.Call(ctx, route, Request{
		Method:      "POST",
		Params:      params,
		Cookie:      cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func audioRelatedSort(v any) int {
	s := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", v)))
	switch s {
	case "hot":
		return 2
	case "new":
		return 3
	default:
		return 1
	}
}

func md5SortedWithKey(m map[string]any, key string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+sigVal(m[k]))
	}
	return util.MD5Hex(key + strings.Join(parts, "") + key)
}

func md5AmpersandSign(m map[string]any) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+sigVal(m[k]))
	}
	sum := util.MD5Hex(strings.Join(pairs, "&") + "*s&iN#G70*")
	if len(sum) < 24 {
		return sum
	}
	return sum[8:24]
}

func sigVal(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

type playlistEncryptResult struct {
	Key          string
	CipherBase64 string
}

func playlistAesEncrypt(data any) (playlistEncryptResult, error) {
	plain, err := json.Marshal(data)
	if err != nil {
		return playlistEncryptResult{}, err
	}
	keySeed := strings.ToLower(util.RandomString(6))
	md5 := util.MD5Hex(keySeed)
	key := []byte(md5[:16])
	iv := []byte(md5[16:32])

	block, err := aes.NewCipher(key)
	if err != nil {
		return playlistEncryptResult{}, err
	}
	padded := pkcs7Pad(plain, aes.BlockSize)
	out := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out, padded)

	return playlistEncryptResult{Key: keySeed, CipherBase64: base64.StdEncoding.EncodeToString(out)}, nil
}

func playlistAesDecryptFromRaw(raw []byte, keySeed string) (map[string]any, error) {
	md5 := util.MD5Hex(keySeed)
	key := []byte(md5[:16])
	iv := []byte(md5[16:32])

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(raw))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(out, raw)
	plain, err := pkcs7Unpad(out, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(plain, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func toInt(v any, def int) int {
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "" || s == "<nil>" {
		return def
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return int(f)
	}
	return def
}

func toBool(v any, def bool) bool {
	s := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", v)))
	switch s {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	case "", "<nil>":
		return def
	default:
		return def
	}
}

func boolPtr(v bool) *bool { return &v }

func splitCSV(s string) []string {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil
	}
	parts := strings.Split(t, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func signCloudKey(hash, pid string) string {
	return util.MD5Hex("musicclound" + hash + pid + "ebd1ac3134c880bda6a2194537843caa0162e2e7")
}

func currentAppid(isLite bool) string {
	a, _ := config.PlatformConfig(isLite)
	return a
}

func currentClientVer(isLite bool) string {
	_, v := config.PlatformConfig(isLite)
	return v
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	pad := blockSize - len(data)%blockSize
	out := make([]byte, len(data)+pad)
	copy(out, data)
	for i := len(data); i < len(out); i++ {
		out[i] = byte(pad)
	}
	return out
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded data")
	}
	pad := int(data[len(data)-1])
	if pad <= 0 || pad > blockSize || pad > len(data) {
		return nil, fmt.Errorf("invalid pad size")
	}
	for i := len(data) - pad; i < len(data); i++ {
		if int(data[i]) != pad {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:len(data)-pad], nil
}
