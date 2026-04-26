package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lfhy/kugou-music-api/core/kugou"
)

// Device registration and user video mutation endpoints share bespoke signatures.
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
