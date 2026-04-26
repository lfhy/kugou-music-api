package sdk

import (
	"context"
	"fmt"
	"strings"
)

type PlaylistCreateResult struct {
	ListID int
	Raw    *Response
}

type PlaylistAddTracksResult struct {
	Added int
	Raw   *Response
}

// CreatePlaylist creates a new playlist for the current logged-in user.
// It uses PlaylistAdd(...) to stay aligned with upstream JS request shaping.
func (c *Client) CreatePlaylist(ctx context.Context, name string, isPrivate bool, cookie map[string]string) (*PlaylistCreateResult, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("empty playlist name")
	}

	cookies := c.Cookie()
	for k, v := range cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, requireLoginCookie(cookies)
	}

	respRaw, err := c.PlaylistAdd(ctx, PlaylistAddRequest{
		Name:   name,
		Type:   0,
		Source: 1,
		IsPri:  ternaryInt(isPrivate, 1, 0),
		Cookie: cookies,
	})
	if err != nil {
		return nil, err
	}
	resp := (*Response)(respRaw)

	listID := 0
	if resp.Body != nil {
		listID = findFirstInt(resp.Body, "listid", "list_id", "listId", "id")
	}
	if listID == 0 {
		return nil, fmt.Errorf("create playlist: cannot parse listid, raw=%s", string(resp.RawBody))
	}
	return &PlaylistCreateResult{ListID: listID, Raw: resp}, nil
}

// AddTracksToPlaylist adds tracks to the given playlist id.
// It skips tracks with empty hash and deduplicates by hash.
// It uses PlaylistTracksAdd(...) to stay aligned with upstream JS request shaping.
func (c *Client) AddTracksToPlaylist(ctx context.Context, listID int, tracks []RadioTrack, cookie map[string]string) (*PlaylistAddTracksResult, error) {
	if listID <= 0 {
		return nil, fmt.Errorf("invalid listID")
	}

	cookies := c.Cookie()
	for k, v := range cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, requireLoginCookie(cookies)
	}

	seen := map[string]struct{}{}
	parts := make([]string, 0, len(tracks))
	for _, t := range tracks {
		h := strings.TrimSpace(t.Hash)
		if h == "" {
			continue
		}
		if _, dup := seen[h]; dup {
			continue
		}
		seen[h] = struct{}{}
		name := strings.TrimSpace(t.Name)
		mixsongid := ternaryInt(t.AlbumAudioID > 0, t.AlbumAudioID, 0)
		// JS format: name|hash|album_id|mixsongid
		parts = append(parts, fmt.Sprintf("%s|%s|%d|%d", name, h, 0, mixsongid))
	}
	if len(parts) == 0 {
		return &PlaylistAddTracksResult{Added: 0, Raw: &Response{Status: 0, Body: map[string]any{"status": 1}}}, nil
	}

	respRaw, err := c.PlaylistTracksAdd(ctx, PlaylistTracksAddRequest{
		Listid: listID,
		Data:   strings.Join(parts, ","),
		Cookie: cookies,
	})
	if err != nil {
		return nil, err
	}
	resp := (*Response)(respRaw)
	return &PlaylistAddTracksResult{Added: len(parts), Raw: resp}, nil
}

func findFirstInt(root map[string]any, keys ...string) int {
	if root == nil {
		return 0
	}
	for _, k := range keys {
		if v, ok := root[k]; ok {
			if n := toInt(v, 0); n != 0 {
				return n
			}
		}
	}
	for _, v := range root {
		switch t := v.(type) {
		case map[string]any:
			if n := findFirstInt(t, keys...); n != 0 {
				return n
			}
		case []any:
			for _, x := range t {
				if m, ok := x.(map[string]any); ok {
					if n := findFirstInt(m, keys...); n != 0 {
						return n
					}
				}
			}
		}
	}
	return 0
}
