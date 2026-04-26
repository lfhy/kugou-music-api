package sdk

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/core/config"
	"github.com/lfhy/kugou-music-api/core/util"
)

type SongPlayURLRequest struct {
	Hash         string
	AlbumID      int
	AlbumAudioID int
	Quality      any
	FreePart     bool
	Cookie       map[string]string
}

type SongPlayURLResult struct {
	Primary string
	URLs    []string
	Source  string
}

type SongPlayURLOption func(*songPlayURLConfig)

type songPlayURLConfig struct {
	fallbackToNew bool
	collectAll    bool
	dfid          string
}

func defaultSongPlayURLConfig() songPlayURLConfig {
	return songPlayURLConfig{
		fallbackToNew: true,
		collectAll:    false,
		dfid:          "",
	}
}

// WithSongURLFallback controls whether fallback to song_url_new is enabled.
func WithSongURLFallback(enabled bool) SongPlayURLOption {
	return func(c *songPlayURLConfig) { c.fallbackToNew = enabled }
}

// WithSongURLAll controls whether all discovered URLs are returned.
func WithSongURLAll(enabled bool) SongPlayURLOption {
	return func(c *songPlayURLConfig) { c.collectAll = enabled }
}

// WithSongURLDFID forces a specific dfid for URL resolving.
func WithSongURLDFID(dfid string) SongPlayURLOption {
	return func(c *songPlayURLConfig) { c.dfid = strings.TrimSpace(dfid) }
}

// GetSongPlayURL resolves a playable URL with module-compatible defaults.
// It keeps backward compatibility and returns the first available URL.
func (c *Client) GetSongPlayURL(ctx context.Context, req SongPlayURLRequest, opts ...SongPlayURLOption) (string, error) {
	result, err := c.ResolveSongPlayURL(ctx, req, opts...)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", nil
	}
	return result.Primary, nil
}

// ResolveSongPlayURL returns detailed URL resolving result.
func (c *Client) ResolveSongPlayURL(ctx context.Context, req SongPlayURLRequest, opts ...SongPlayURLOption) (*SongPlayURLResult, error) {
	hash := strings.ToLower(strings.TrimSpace(req.Hash))
	if hash == "" {
		return nil, fmt.Errorf("empty hash")
	}

	cfg := defaultSongPlayURLConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	if cfg.dfid != "" {
		cookies["dfid"] = cfg.dfid
	}
	if strings.TrimSpace(cookies["dfid"]) == "" {
		cookies["dfid"] = "-"
	}

	quality := req.Quality
	if quality == nil || strings.TrimSpace(fmt.Sprintf("%v", quality)) == "" {
		quality = 128
	}
	if q := strings.TrimSpace(fmt.Sprintf("%v", quality)); q != "" {
		switch q {
		case "piano", "acappella", "subwoofer", "ancient", "dj", "surnay":
			quality = "magic_" + q
		}
	}

	pageID := 151369488
	ppageID := "463467626,350369493,788954147"
	if c.isLite {
		pageID = 967177915
		ppageID = "356753938,823673182,967485191"
	}

	resp, err := c.Call(ctx, RouteSongUrl, Request{
		Params: map[string]any{
			"album_id":       req.AlbumID,
			"area_code":      1,
			"hash":           hash,
			"ssa_flag":       "is_fromtrack",
			"version":        11436,
			"page_id":        pageID,
			"quality":        quality,
			"album_audio_id": req.AlbumAudioID,
			"behavior":       "play",
			"pid":            ternaryInt(c.isLite, 411, 2),
			"cmd":            26,
			"pidversion":     3001,
			"IsFreePart":     ternaryInt(req.FreePart, 1, 0),
			"ppage_id":       ppageID,
			"cdnBackup":      1,
			"kcard":          0,
			"module":         "",
		},
		Cookie: cookies,
	})
	if err == nil && resp != nil && resp.Body != nil {
		urls := extractSongURLs(resp.Body)
		if len(urls) > 0 {
			out := &SongPlayURLResult{Primary: urls[0], Source: "song_url"}
			if cfg.collectAll {
				out.URLs = urls
			}
			return out, nil
		}
	}

	if !cfg.fallbackToNew {
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	userid := fmt.Sprintf("%v", firstAny(cookies["userid"], "0"))
	vipType := fmt.Sprintf("%v", firstAny(cookies["vip_type"], "0"))
	token := fmt.Sprintf("%v", firstAny(cookies["token"], ""))
	vipToken := fmt.Sprintf("%v", firstAny(cookies["vip_token"], ""))
	mid := fmt.Sprintf("%v", firstAny(cookies["KUGOU_API_MID"], ""))
	appid, _ := config.PlatformConfig(c.isLite)
	clientTimeMs := time.Now().UnixMilli()

	respNew, errNew := c.Call(ctx, RouteSongUrlNew, Request{
		Data: map[string]any{
			"area_code": "1",
			"behavior":  "play",
			"qualities": []any{"128", "320", "flac", "high", "multitrack", "viper_atmos", "viper_tape", "viper_clear"},
			"resource": map[string]any{
				"album_audio_id":  req.AlbumAudioID,
				"collect_list_id": "3",
				"collect_time":    clientTimeMs,
				"hash":            hash,
				"id":              0,
				"page_id":         1,
				"type":            "audio",
			},
			"token": token,
			"tracker_param": map[string]any{
				"all_m":         1,
				"auth":          "",
				"is_free_part":  ternaryInt(req.FreePart, 1, 0),
				"key":           util.MD5Hex(hash + "185672dd44712f60bb1736df5a377e82" + appid + mid + userid),
				"module_id":     0,
				"need_climax":   1,
				"need_xcdn":     1,
				"open_time":     "",
				"pid":           "411",
				"pidversion":    "3001",
				"priv_vip_type": "6",
				"viptoken":      vipToken,
			},
			"userid": userid,
			"vip":    vipType,
		},
		Cookie: cookies,
	})
	if errNew != nil {
		if err != nil {
			return nil, err
		}
		return nil, errNew
	}
	if respNew != nil && respNew.Body != nil {
		urls := extractSongURLs(respNew.Body)
		if len(urls) > 0 {
			out := &SongPlayURLResult{Primary: urls[0], Source: "song_url_new"}
			if cfg.collectAll {
				out.URLs = urls
			}
			return out, nil
		}
	}

	return nil, nil
}

func ternaryInt(ok bool, a, b int) int {
	if ok {
		return a
	}
	return b
}

func extractSongURLs(root map[string]any) []string {
	ordered := []string{"url", "play_url", "backup_url", "song_url", "sq_url", "hq_url", "backupUrl", "sqUrl", "hqUrl", "urls", "backup_urls", "url_list", "sq_urls", "hq_urls"}
	seen := map[string]struct{}{}
	out := make([]string, 0, 4)

	if root != nil {
		for _, k := range ordered {
			if v, ok := root[k]; ok {
				for _, u := range urlsFromAny(v) {
					if _, dup := seen[u]; dup {
						continue
					}
					seen[u] = struct{}{}
					out = append(out, u)
				}
			}
		}
	}

	var walk func(any)
	walk = func(v any) {
		switch t := v.(type) {
		case map[string]any:
			for _, vv := range t {
				walk(vv)
			}
		case []any:
			for _, x := range t {
				walk(x)
			}
		default:
			if s := asStringAny(t); strings.HasPrefix(s, "http") {
				if _, dup := seen[s]; !dup {
					seen[s] = struct{}{}
					out = append(out, s)
				}
			}
		}
	}
	walk(root)

	return out
}

func urlsFromAny(v any) []string {
	out := []string{}
	if s := asStringAny(v); strings.HasPrefix(s, "http") {
		return append(out, s)
	}
	if arr, ok := v.([]any); ok {
		for _, x := range arr {
			if s := asStringAny(x); strings.HasPrefix(s, "http") {
				out = append(out, s)
			}
		}
	}
	return out
}

func asStringAny(v any) string {
	if v == nil {
		return ""
	}
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "" || s == "<nil>" || s == "0" {
		return ""
	}
	return s
}
