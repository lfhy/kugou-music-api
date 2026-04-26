package sdk

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// sleepWithContext lets batch fetching stop promptly on cancellation.
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
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		out = append(out, RadioTrack{
			Name:         pickAnyString(m, "songname", "name", "SongName", "audio_name"),
			Singer:       pickAnyString(m, "singername", "author_name", "SingerName", "author"),
			Hash:         pickAnyString(m, "hash", "Hash"),
			AlbumAudioID: pickAnyInt(m, "album_audio_id", "AlbumAudioID", "audioid", "MixSongID"),
			SongID:       pickAnyInt(m, "songid", "SongID", "audioid"),
			URL:          pickAnyString(m, "url", "play_url", "play_backup_url"),
		})
	}
	return out
}

func findFirstTrackList(root map[string]any) []any {
	if root == nil {
		return nil
	}
	candidates := []string{
		"data", "list", "lists", "songlist", "songs", "items",
		"recommend_list", "album_audio_list", "song_list",
	}
	for _, key := range candidates {
		if arr, ok := root[key].([]any); ok && looksLikeTrackList(arr) {
			return arr
		}
		if sub, ok := root[key].(map[string]any); ok {
			if arr := findFirstTrackList(sub); len(arr) > 0 {
				return arr
			}
		}
	}
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
	limit := len(arr)
	if limit > 3 {
		limit = 3
	}
	for i := 0; i < limit; i++ {
		m, ok := arr[i].(map[string]any)
		if !ok {
			continue
		}
		if pickAnyString(m, "hash", "songname", "name", "audio_name") != "" {
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
			if n := toInt(v, 0); n != 0 {
				return n
			}
		}
	}
	return 0
}
