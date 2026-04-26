package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/lfhy/kugou-music-api/sdk"
)

type songItem struct {
	Name         string
	Hash         string
	AlbumID      int
	AlbumAudioID int
}

func main() {
	client, err := sdk.New()
	if err != nil {
		log.Fatalf("init sdk failed: %v", err)
	}

	cookie := map[string]string{}
	if token := strings.TrimSpace(os.Getenv("KUGOU_TOKEN")); token != "" {
		cookie["token"] = token
	}
	if uid := strings.TrimSpace(os.Getenv("KUGOU_USERID")); uid != "" {
		cookie["userid"] = uid
	}

	resp, err := client.GetDailyRecommendGuest(context.Background(), cookie)
	if err != nil {
		log.Fatalf("daily recommend failed: %v", err)
	}

	fmt.Printf("status: %d\n", resp.Status)
	if resp.Body == nil {
		fmt.Printf("raw body: %s\n", string(resp.RawBody))
		return
	}

	playlistName := extractPlaylistName(resp.Body)
	if playlistName != "" {
		fmt.Printf("playlist: %s\n", playlistName)
	} else {
		fmt.Println("playlist: (unknown)")
	}

	songs := extractSongs(resp.Body)
	fmt.Printf("songs: %d\n", len(songs))

	// Daily recommend is usually 30-50 tracks. Print more by default.
	limit := len(songs)
	if limit > 30 {
		limit = 30
	}
	for i := 0; i < limit; i++ {
		s := songs[i]
		fmt.Printf("%d. %s\n", i+1, s.Name)
		if s.Hash == "" {
			fmt.Println("   url: (no hash)")
			continue
		}

		u := ""
		u, _ = client.GetSongPlayURL(context.Background(), sdk.SongPlayURLRequest{
			Hash:         s.Hash,
			AlbumID:      s.AlbumID,
			AlbumAudioID: s.AlbumAudioID,
			FreePart:     true,
			Cookie:       cookie,
		})
		if u == "" && s.Hash != "" {
			// JS module/song_url.js lower-cases hash before request.
			u, _ = client.GetSongPlayURL(context.Background(), sdk.SongPlayURLRequest{
				Hash:         strings.ToLower(s.Hash),
				AlbumID:      s.AlbumID,
				AlbumAudioID: s.AlbumAudioID,
				FreePart:     true,
				Cookie:       cookie,
			})
		}
		if u == "" {
			fmt.Println("   url: (unavailable)")
		} else {
			fmt.Printf("   url: %s\n", u)
		}
	}
}

func extractPlaylistName(root map[string]any) string {
	candidates := []string{"playlist_name", "list_name", "name", "title", "theme_name"}
	return findFirstStringByKeys(root, candidates)
}

func extractSongs(root map[string]any) []songItem {
	arr := findLikelySongArray(root)
	out := make([]songItem, 0, len(arr))
	for _, it := range arr {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		name := firstString(m, "songname", "filename", "name", "audio_name", "song_name")
		hash := firstString(m, "hash", "audio_hash", "file_hash", "hash_128", "hash_320", "hash_flac")
		albumID := firstInt(m, "album_id", "albumid")
		albumAudioID := firstInt(m, "album_audio_id", "albumaudioid")
		if name == "" && hash == "" {
			continue
		}
		out = append(out, songItem{
			Name:         name,
			Hash:         strings.ToLower(hash),
			AlbumID:      albumID,
			AlbumAudioID: albumAudioID,
		})
	}
	return out
}

func extractSongURL(root map[string]any) string {
	return findFirstURL(root)
}

func findLikelySongArray(v any) []any {
	switch t := v.(type) {
	case map[string]any:
		for _, key := range []string{"song_list", "songlist", "songs", "list", "info"} {
			if arr, ok := t[key].([]any); ok && seemsSongs(arr) {
				return arr
			}
		}
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if arr := findLikelySongArray(t[k]); len(arr) > 0 {
				return arr
			}
		}
	case []any:
		if seemsSongs(t) {
			return t
		}
		for _, item := range t {
			if arr := findLikelySongArray(item); len(arr) > 0 {
				return arr
			}
		}
	}
	return nil
}

func seemsSongs(arr []any) bool {
	if len(arr) == 0 {
		return false
	}
	probe := 0
	if len(arr) < 5 {
		probe = len(arr)
	} else {
		probe = 5
	}
	for i := 0; i < probe; i++ {
		m, ok := arr[i].(map[string]any)
		if !ok {
			continue
		}
		if firstString(m, "hash", "audio_hash", "filename", "songname", "name") != "" {
			return true
		}
	}
	return false
}

func findFirstURL(v any) string {
	switch t := v.(type) {
	case map[string]any:
		for _, k := range []string{"url", "play_url", "backup_url", "song_url", "sq_url", "hq_url"} {
			if s := asString(t[k]); strings.HasPrefix(s, "http") {
				return s
			}
		}
		for _, k := range []string{"urls", "backup_urls", "url_list", "sq_urls", "hq_urls"} {
			if arr, ok := t[k].([]any); ok {
				for _, x := range arr {
					s := asString(x)
					if strings.HasPrefix(s, "http") {
						return s
					}
				}
			}
		}
		for _, v2 := range t {
			if s := findFirstURL(v2); s != "" {
				return s
			}
		}
	case []any:
		for _, x := range t {
			if s := findFirstURL(x); s != "" {
				return s
			}
		}
	}
	return ""
}

func findFirstStringByKeys(v any, keys []string) string {
	switch t := v.(type) {
	case map[string]any:
		for _, k := range keys {
			if s := asString(t[k]); s != "" {
				return s
			}
		}
		for _, v2 := range t {
			if s := findFirstStringByKeys(v2, keys); s != "" {
				return s
			}
		}
	case []any:
		for _, x := range t {
			if s := findFirstStringByKeys(x, keys); s != "" {
				return s
			}
		}
	}
	return ""
}

func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if s := asString(m[k]); s != "" {
			return s
		}
	}
	return ""
}

func asString(v any) string {
	if v == nil {
		return ""
	}
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "" || s == "<nil>" || s == "0" {
		return ""
	}
	return s
}

func firstInt(m map[string]any, keys ...string) int {
	for _, k := range keys {
		if n := asInt(m[k]); n > 0 {
			return n
		}
	}
	return 0
}

func asInt(v any) int {
	if v == nil {
		return 0
	}
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "" || s == "<nil>" {
		return 0
	}
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return 0
	}
	return n
}
