package sdk

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type PersonalRadioMode string

const (
	PersonalRadioHeart PersonalRadioMode = "heart" // 红心/猜你喜欢（私人FM normal）
	PersonalRadioNew   PersonalRadioMode = "new"   // 新歌速递（top_song）
	PersonalRadioNiche PersonalRadioMode = "niche" // 小众（私人FM small）
)

type PersonalRadioRequest struct {
	Mode PersonalRadioMode

	// For fm:
	SongPoolID int

	// For list-like modes:
	Page     int
	PageSize int

	Cookie map[string]string
}

type RadioTrack struct {
	Name         string
	Singer       string
	Hash         string
	AlbumAudioID int
	SongID       int
	URL          string
}

type PersonalRadioResponse struct {
	Mode   PersonalRadioMode
	Source string
	Raw    *Response
	Tracks []RadioTrack
}

func (c *Client) GetPersonalRadio(ctx context.Context, req PersonalRadioRequest) (*PersonalRadioResponse, error) {
	mode := PersonalRadioHeart
	if strings.TrimSpace(string(req.Mode)) != "" {
		mode = req.Mode
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 30
	}

	switch mode {
	case PersonalRadioHeart:
		raw, tracks, err := c.fetchPersonalFmBatch(ctx, req.Cookie, "normal", req.SongPoolID, pageSize)
		if err != nil {
			return nil, err
		}
		return &PersonalRadioResponse{Mode: mode, Source: "personal_fm", Raw: (*Response)(raw), Tracks: tracks}, nil

	case PersonalRadioNiche:
		raw, tracks, err := c.fetchPersonalFmBatch(ctx, req.Cookie, "small", req.SongPoolID, pageSize)
		if err != nil {
			return nil, err
		}
		return &PersonalRadioResponse{Mode: mode, Source: "personal_fm", Raw: (*Response)(raw), Tracks: tracks}, nil

	case PersonalRadioNew:
		raw, err := c.TopSong(ctx, TopSongRequest{
			Page:     page,
			Pagesize: pageSize,
			Type:     21608,
			Cookie:   req.Cookie,
		})
		if err != nil {
			return nil, err
		}
		tracks := extractRadioTracks(raw.Body)
		return &PersonalRadioResponse{Mode: mode, Source: "top_song", Raw: (*Response)(raw), Tracks: tracks}, nil

	default:
		return nil, fmt.Errorf("unknown personal radio mode: %q", mode)
	}
}

func (c *Client) fetchPersonalFmBatch(ctx context.Context, cookie map[string]string, mode string, songPoolID, want int) (*PersonalFmResponse, []RadioTrack, error) {
	if want <= 0 {
		want = 30
	}
	if want > 50 {
		want = 50
	}

	all := make([]RadioTrack, 0, want)
	seen := map[string]struct{}{}
	var lastRaw *PersonalFmResponse

	// Personal FM commonly returns a small chunk per call (often 5 tracks).
	// Keep polling until target count is reached, while tolerating duplicate/empty rounds.
	maxAttempts := want * 3
	if maxAttempts < 12 {
		maxAttempts = 12
	}
	if maxAttempts > 80 {
		maxAttempts = 80
	}

	staleRounds := 0
	lastHash := ""
	lastSongID := 0
	remain := 0
	for i := 0; i < maxAttempts && len(all) < want; i++ {
		raw, err := c.PersonalFm(ctx, PersonalFmRequest{
			Mode:          mode,
			Action:        "play",
			Platform:      "ios",
			SongPoolId:    songPoolID,
			RemainSongcnt: remain,
			Hash:          lastHash,
			Songid:        lastSongID,
			Cookie:        cookie,
		})
		if err != nil {
			if lastRaw != nil && len(all) > 0 {
				staleRounds++
				if staleRounds >= 3 {
					break
				}
				if !sleepWithContext(ctx, 150*time.Millisecond) {
					break
				}
				continue
			}
			return nil, nil, err
		}
		lastRaw = raw
		chunk := extractRadioTracks(raw.Body)
		if len(chunk) == 0 {
			staleRounds++
			if staleRounds >= 3 {
				break
			}
			if !sleepWithContext(ctx, 150*time.Millisecond) {
				break
			}
			continue
		}

		addedThisRound := 0
		for _, t := range chunk {
			h := strings.ToLower(strings.TrimSpace(t.Hash))
			if h == "" {
				continue
			}
			if _, dup := seen[h]; dup {
				continue
			}
			seen[h] = struct{}{}
			all = append(all, t)
			addedThisRound++
			if len(all) >= want {
				break
			}
		}

		last := chunk[len(chunk)-1]
		lastHash = strings.TrimSpace(last.Hash)
		lastSongID = last.SongID
		remain = maxInt(want-len(all), 0)

		if addedThisRound == 0 {
			staleRounds++
			if staleRounds >= 4 {
				break
			}
		} else {
			staleRounds = 0
		}
		if len(all) >= want {
			break
		}
		if !sleepWithContext(ctx, 150*time.Millisecond) {
			break
		}
	}

	if len(all) > want {
		all = all[:want]
	}
	return lastRaw, all, nil
}

func sleepWithContext(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		return true
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func extractRadioTracks(root map[string]any) []RadioTrack {
	items := findFirstTrackList(root)
	out := make([]RadioTrack, 0, len(items))
	for _, it := range items {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		t := RadioTrack{
			Name:         pickAnyString(m, "songname", "song_name", "name", "title", "audio_name"),
			Singer:       pickAnyString(m, "singername", "singer_name", "author_name", "singer", "authors_name"),
			Hash:         strings.ToLower(strings.TrimSpace(pickAnyString(m, "hash", "audio_hash"))),
			AlbumAudioID: pickAnyInt(m, "album_audio_id", "albumAudioId", "album_audioid", "mixsongid", "mixsong_id"),
			SongID:       pickAnyInt(m, "songid", "song_id", "audio_id"),
		}
		out = append(out, t)
	}
	return out
}

func findFirstTrackList(root map[string]any) []any {
	if root == nil {
		return nil
	}
	// Prefer common containers.
	candidates := []any{root["data"], root["songs"], root["songlist"], root["list"], root["info"], root["items"]}
	for _, c := range candidates {
		if arr, ok := c.([]any); ok && looksLikeTrackList(arr) {
			return arr
		}
	}
	// Walk shallowly.
	for _, v := range root {
		switch t := v.(type) {
		case []any:
			if looksLikeTrackList(t) {
				return t
			}
		case map[string]any:
			if arr := findFirstTrackList(t); len(arr) > 0 {
				return arr
			}
		}
	}
	return nil
}

func looksLikeTrackList(arr []any) bool {
	if len(arr) == 0 {
		return false
	}
	m, ok := arr[0].(map[string]any)
	if !ok {
		return false
	}
	// Heuristic: any of these indicates it's a song-like item.
	for _, k := range []string{"hash", "songname", "song_name", "album_audio_id", "audio_id", "songid"} {
		if _, ok := m[k]; ok {
			return true
		}
	}
	return false
}

func pickAnyString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			s := strings.TrimSpace(fmt.Sprintf("%v", v))
			if s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

func pickAnyInt(m map[string]any, keys ...string) int {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			return toInt(v, 0)
		}
	}
	return 0
}
