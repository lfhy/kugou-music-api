package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/core/kugou"
)

// PlaylistAdd: create playlist (type=0).
func (c *Client) PlaylistAdd(ctx context.Context, req PlaylistAddRequest) (*PlaylistAddResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if ok == false {
		return nil, requireLoginCookie(cookies)
	}

	name := strings.TrimSpace(fmt.Sprintf("%v", firstAny(req.Name, "")))
	if name == "" || name == "<nil>" {
		return nil, fmt.Errorf("playlist_add: empty name")
	}

	userid := firstNonEmpty(strings.TrimSpace(cookies["userid"]), strings.TrimSpace(fmt.Sprintf("%v", req.Userid)), "0")
	token := firstNonEmpty(strings.TrimSpace(cookies["token"]), strings.TrimSpace(fmt.Sprintf("%v", req.Token)), "")
	if userid == "0" || token == "" {
		return nil, requireLoginCookie(cookies)
	}

	clienttime := time.Now().Unix()
	tp := toInt(firstAny(req.Type, 0), 0)
	source := toInt(firstAny(req.Source, 1), 1)
	if toInt(firstAny(req.Source, 1), 1) == 0 {
		source = 0
	}

	data := map[string]any{
		"userid":             userid,
		"token":              token,
		"total_ver":          0,
		"name":               name,
		"type":               tp,
		"source":             source,
		"is_pri":             0,
		"list_create_userid": req.ListCreateUserid,
		"list_create_listid": req.ListCreateListid,
		"list_create_gid":    firstNonEmpty(fmt.Sprintf("%v", req.ListCreateGid), ""),
		"from_shupinmv":      0,
	}
	if tp == 0 {
		data["is_pri"] = toInt(firstAny(req.IsPri, 0), 0)
	}

	params := map[string]any{}
	if tp == 0 {
		params = map[string]any{
			"last_time": clienttime,
			"last_area": "gztx",
			"userid":    userid,
			"token":     token,
		}
	}

	resp, err := c.callManual(ctx, kugou.RequestConfig{
		Method:      "POST",
		URL:         "/cloudlist.service/v5/add_list",
		Params:      params,
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := PlaylistAddResponse(*resp)
	return &out, nil
}

// PlaylistTracksAdd: add songs into playlist.
func (c *Client) PlaylistTracksAdd(ctx context.Context, req PlaylistTracksAddRequest) (*PlaylistTracksAddResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if ok == false {
		return nil, requireLoginCookie(cookies)
	}

	userid := firstNonEmpty(strings.TrimSpace(cookies["userid"]), strings.TrimSpace(fmt.Sprintf("%v", req.Userid)), "0")
	token := firstNonEmpty(strings.TrimSpace(cookies["token"]), strings.TrimSpace(fmt.Sprintf("%v", req.Token)), "")
	if userid == "0" || token == "" {
		return nil, requireLoginCookie(cookies)
	}

	listid := toInt(firstAny(req.Listid, 0), 0)
	if listid == 0 {
		return nil, fmt.Errorf("playlist_tracks_add: empty listid")
	}

	rawData := strings.TrimSpace(fmt.Sprintf("%v", firstAny(req.Data, "")))
	if rawData == "" || rawData == "<nil>" {
		return nil, fmt.Errorf("playlist_tracks_add: empty data")
	}

	resource := make([]map[string]any, 0, 16)
	for _, s := range splitCSV(rawData) {
		seg := strings.Split(s, "|")
		name := ""
		hash := ""
		albumID := 0
		mixsongid := 0
		if len(seg) > 0 {
			name = seg[0]
		}
		if len(seg) > 1 {
			hash = seg[1]
		}
		if len(seg) > 2 {
			albumID = toInt(seg[2], 0)
		}
		if len(seg) > 3 {
			mixsongid = toInt(seg[3], 0)
		}
		resource = append(resource, map[string]any{
			"number":    1,
			"name":      name,
			"hash":      strings.TrimSpace(hash),
			"size":      0,
			"sort":      0,
			"timelen":   0,
			"bitrate":   0,
			"album_id":  albumID,
			"mixsongid": mixsongid,
		})
	}

	clienttime := time.Now().Unix()
	params := map[string]any{
		"last_time": clienttime,
		"last_area": "gztx",
		"userid":    userid,
		"token":     token,
	}
	data := map[string]any{
		"userid":      userid,
		"token":       token,
		"listid":      listid,
		"list_ver":    0,
		"type":        0,
		"slow_upload": 1,
		"scene":       "false;null",
		"data":        resource,
	}

	resp, err := c.callManual(ctx, kugou.RequestConfig{
		Method:      "POST",
		URL:         "/cloudlist.service/v6/add_song",
		Params:      params,
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := PlaylistTracksAddResponse(*resp)
	return &out, nil
}

// PlaylistTracksDel: delete songs by fileids.
func (c *Client) PlaylistTracksDel(ctx context.Context, req PlaylistTracksDelRequest) (*PlaylistTracksDelResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if ok == false {
		return nil, requireLoginCookie(cookies)
	}

	userid := firstNonEmpty(strings.TrimSpace(cookies["userid"]), strings.TrimSpace(fmt.Sprintf("%v", req.Userid)), "0")
	token := firstNonEmpty(strings.TrimSpace(cookies["token"]), strings.TrimSpace(fmt.Sprintf("%v", req.Token)), "")
	if userid == "0" || token == "" {
		return nil, requireLoginCookie(cookies)
	}

	listid := toInt(firstAny(req.Listid, 0), 0)
	if listid == 0 {
		return nil, fmt.Errorf("playlist_tracks_del: empty listid")
	}
	fileidsRaw := strings.TrimSpace(fmt.Sprintf("%v", firstAny(req.Fileids, "")))
	if fileidsRaw == "" || fileidsRaw == "<nil>" {
		return nil, fmt.Errorf("playlist_tracks_del: empty fileids")
	}

	resource := make([]map[string]any, 0, 16)
	for _, s := range splitCSV(fileidsRaw) {
		n := toInt(s, 0)
		if n <= 0 {
			continue
		}
		resource = append(resource, map[string]any{"fileid": n})
	}

	data := map[string]any{
		"listid":   listid,
		"userid":   userid,
		"data":     resource,
		"type":     0,
		"token":    token,
		"list_ver": 0,
	}
	resp, err := c.Call(ctx, RoutePlaylistTracksDel, Request{
		Method:      "POST",
		URL:         "/v4/delete_songs",
		Data:        data,
		Cookie:      cookies,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "cloudlist.service.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := PlaylistTracksDelResponse(*resp)
	return &out, nil
}

func (c *Client) callManual(ctx context.Context, cfg kugou.RequestConfig) (*Response, error) {
	raw, err := c.core.CreateRequest(ctx, cfg)
	if len(raw.Cookie) > 0 {
		c.updateCookiePool(raw.Cookie)
	}

	out := &Response{Status: raw.Status, RawBody: raw.Body, Headers: raw.Headers, Cookie: raw.Cookie}
	var body map[string]any
	if json.Unmarshal(raw.Body, &body) == nil {
		out.Body = body
	}
	return out, err
}
