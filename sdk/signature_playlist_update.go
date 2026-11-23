package sdk

import (
	"context"
	"math"
	"strconv"
	"strings"
)

// PlaylistPic and PlaylistUpdate modify cloud playlists using the cloud cipher.
// PlaylistPic accepts comma-separated listid|type|pic entries, as in the JS API.
func (c *Client) PlaylistPic(ctx context.Context, req PlaylistPicRequest) (*PlaylistPicResponse, error) {
	entries := strings.Split(req.Data, ",")
	resource := make([]map[string]any, 0, len(entries))
	for _, entry := range entries {
		parts := strings.Split(entry, "|")
		typeValue, pic := "0", "undefined"
		if len(parts) > 1 {
			typeValue = parts[1]
		}
		if len(parts) > 2 {
			pic = parts[2]
		}
		if !strings.HasPrefix(pic, "stdmusic/") {
			pic = "stdmusic/" + pic[strings.LastIndex(pic, "/")+1:]
		}
		resource = append(resource, map[string]any{
			"listid": playlistPicNumber(parts[0]),
			"type":   playlistPicNumber(typeValue),
			"pic":    pic,
		})
	}
	return c.callPlaylistCloud(ctx, "/v1/modify_list_pic", map[string]any{
		"total_ver": req.TotalVer,
		"data":      resource,
	}, req.Cookie)
}

// PlaylistUpdate omits an empty name or pic, while tags and intro always apply.
// Pic is custom/<FileName>, where FileName comes from PlaylistUpload.
func (c *Client) PlaylistUpdate(ctx context.Context, req PlaylistUpdateRequest) (*PlaylistUpdateResponse, error) {
	data := map[string]any{
		"total_ver": req.TotalVer,
		"listid":    req.Listid,
		"type":      req.Type,
		"sort":      req.Sort,
		"tags":      req.Tags,
		"intro":     req.Intro,
	}
	if req.Name != "" {
		data["name"] = req.Name
	}
	if req.Pic != "" {
		data["pic"] = req.Pic
	}
	return c.callPlaylistCloud(ctx, "/v1/modify_list", data, req.Cookie)
}

// JS Number accepts empty strings and radix prefixes. Non-finite results become
// JSON null when JS stringifies the encrypted payload, rather than becoming 0.
func playlistPicNumber(value string) any {
	value = strings.TrimSpace(value)
	if strings.Contains(value, "_") {
		return nil
	}
	if value == "" {
		return 0
	}
	if len(value) > 2 && value[0] == '0' && strings.ContainsAny(value[1:2], "xXbBoO") {
		n, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return nil
		}
		return float64(n)
	}
	n, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
		return nil
	}
	return n
}
