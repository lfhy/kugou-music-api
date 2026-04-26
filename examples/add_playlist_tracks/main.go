package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
	"github.com/lfhy/kugou-music-api/sdk"
)

func main() {
	listID := flag.Int("listid", 0, "target playlist listid")
	source := flag.String("source", "radio", "tracks source: radio | daily")
	mode := flag.String("mode", "heart", "radio mode when source=radio: heart | new | niche")
	limit := flag.Int("limit", 50, "tracks to add (daily/radio)")
	download := flag.Bool("download", false, "download resolved urls")
	dir := flag.String("dir", "", "download dir (optional). default from session.download_dir or ./downloads")
	debug := flag.Bool("debug", false, "print debug response")
	flag.Parse()

	if *listID <= 0 {
		fmt.Println("missing -listid")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cfg := session.Load(session.DefaultPath())
	debugEnabled := cfg.Debug || *debug
	if !session.HasLoginCookie(cfg.Cookie) {
		fmt.Println("未检测到本地登录会话，请先登录后再添加歌曲。")
		os.Exit(1)
	}

	client, err := sdk.New(sdk.WithCookie(cfg.Cookie))
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		os.Exit(1)
	}

	tracks, err := loadTracks(ctx, client, cfg.Cookie, *source, *mode, *limit, debugEnabled)
	if err != nil {
		fmt.Printf("load tracks failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("tracks loaded: %d\n", len(tracks))

	added, err := client.AddTracksToPlaylist(ctx, *listID, tracks, cfg.Cookie)
	if err != nil {
		fmt.Printf("add tracks failed: %v\n", err)
		if debugEnabled {
			printResp("add_tracks", added.Raw)
		}
		os.Exit(1)
	}
	fmt.Printf("tracks added: %d\n", added.Added)
	if debugEnabled {
		printResp("add_tracks", added.Raw)
	}

	if *download {
		d := *dir
		if d == "" {
			d = strings.TrimSpace(cfg.DownloadDir)
		}
		if d == "" {
			d = "./downloads"
		}
		if err := os.MkdirAll(d, 0o755); err != nil {
			fmt.Printf("mkdir failed: %v\n", err)
			os.Exit(1)
		}
		downloadTracks(ctx, client, cfg.Cookie, d, tracks, debugEnabled)
	}
}

func loadTracks(ctx context.Context, client *sdk.Client, cookie map[string]string, source, mode string, limit int, debug bool) ([]sdk.RadioTrack, error) {
	source = strings.ToLower(strings.TrimSpace(source))
	if limit <= 0 {
		limit = 50
	}
	if limit > 50 {
		limit = 50
	}

	switch source {
	case "daily":
		resp, err := client.GetDailyRecommendGuest(ctx, cookie)
		if err != nil {
			return nil, err
		}
		if debug {
			printResp("daily_recommend", resp)
		}
		tracks := extractTracks(resp.Body)
		if len(tracks) > limit {
			tracks = tracks[:limit]
		}
		return tracks, nil

	case "radio":
		resp, err := client.GetPersonalRadio(ctx, sdk.PersonalRadioRequest{Mode: sdk.PersonalRadioMode(mode), PageSize: limit, Cookie: cookie})
		if err != nil {
			return nil, err
		}
		if debug {
			printResp("personal_radio", resp.Raw)
		}
		tracks := resp.Tracks
		if len(tracks) > limit {
			tracks = tracks[:limit]
		}
		return tracks, nil
	default:
		return nil, fmt.Errorf("unknown source: %s", source)
	}
}

func extractTracks(root map[string]any) []sdk.RadioTrack {
	// Reuse sdk internal heuristics by going through PersonalRadioResponse shape.
	// Minimal extraction here: rely on hash/songname keys.
	if root == nil {
		return nil
	}

	var out []sdk.RadioTrack
	var walk func(any)
	walk = func(v any) {
		switch t := v.(type) {
		case []any:
			for _, x := range t {
				walk(x)
			}
		case map[string]any:
			name := pickAnyString(t, "songname", "song_name", "name", "title", "audio_name")
			hash := strings.ToLower(strings.TrimSpace(pickAnyString(t, "hash", "audio_hash")))
			albumAudioID := pickAnyInt(t, "album_audio_id", "mixsongid", "mixsong_id")
			if hash != "" && name != "" {
				out = append(out, sdk.RadioTrack{Name: name, Hash: hash, AlbumAudioID: albumAudioID})
			}
			for _, vv := range t {
				walk(vv)
			}
		}
	}
	walk(root)

	seen := map[string]struct{}{}
	dedup := make([]sdk.RadioTrack, 0, len(out))
	for _, t := range out {
		if t.Hash == "" {
			continue
		}
		if _, ok := seen[t.Hash]; ok {
			continue
		}
		seen[t.Hash] = struct{}{}
		dedup = append(dedup, t)
	}
	return dedup
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
			s := strings.TrimSpace(fmt.Sprintf("%v", v))
			if s == "" || s == "<nil>" {
				continue
			}
			var n int
			_, _ = fmt.Sscanf(s, "%d", &n)
			if n > 0 {
				return n
			}
		}
	}
	return 0
}

func printResp(label string, resp *sdk.Response) {
	if resp == nil {
		fmt.Printf("[debug] %s: nil response\n", label)
		return
	}
	fmt.Printf("[debug] %s: http_status=%d\n", label, resp.Status)
	fmt.Printf("[debug] %s: raw=%s\n", label, string(resp.RawBody))
}

func downloadTracks(ctx context.Context, client *sdk.Client, cookie map[string]string, dir string, tracks []sdk.RadioTrack, debug bool) {
	fmt.Printf("download dir: %s\n", dir)
	for i, t := range tracks {
		if t.Hash == "" {
			continue
		}
		url, err := client.GetSongPlayURL(ctx, sdk.SongPlayURLRequest{Hash: t.Hash, AlbumAudioID: t.AlbumAudioID, Cookie: cookie})
		if err != nil || url == "" {
			if debug {
				fmt.Printf("[debug] url unavailable: %d %s err=%v\n", i+1, t.Hash, err)
			}
			continue
		}
		fn := safeFileName(fmt.Sprintf("%02d_%s_%s.mp3", i+1, t.Name, t.Hash[:8]))
		path := filepath.Join(dir, fn)
		if _, err := os.Stat(path); err == nil {
			continue
		}
		if err := downloadOne(ctx, url, path); err != nil {
			fmt.Printf("download failed: %s err=%v\n", fn, err)
		} else {
			fmt.Printf("downloaded: %s\n", fn)
		}
	}
}

func downloadOne(ctx context.Context, url, path string) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("http %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func safeFileName(s string) string {
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, "?", "")
	s = strings.ReplaceAll(s, "*", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "<", "")
	s = strings.ReplaceAll(s, ">", "")
	s = strings.ReplaceAll(s, "|", "")
	return strings.TrimSpace(s)
}
