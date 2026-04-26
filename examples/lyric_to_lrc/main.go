package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	kg "github.com/lfhy/kugou-music-api"
	"github.com/lfhy/kugou-music-api/examples/shared/session"
)

func main() {
	var (
		id         = flag.String("id", strings.TrimSpace(os.Getenv("KUGOU_LYRIC_ID")), "lyric id")
		accesskey  = flag.String("accesskey", strings.TrimSpace(os.Getenv("KUGOU_LYRIC_ACCESSKEY")), "lyric accesskey")
		format     = flag.String("fmt", firstNonEmpty(strings.TrimSpace(os.Getenv("KUGOU_LYRIC_FMT")), "krc"), "lyric format: krc or lrc")
		keyword    = flag.String("keyword", strings.TrimSpace(os.Getenv("KUGOU_LYRIC_KEYWORD")), "optional keyword for search_lyric fallback")
		hash       = flag.String("hash", strings.TrimSpace(os.Getenv("KUGOU_LYRIC_HASH")), "optional hash for search_lyric fallback")
		albumAudio = flag.Int("album-audio-id", envInt("KUGOU_LYRIC_ALBUM_AUDIO_ID"), "optional album_audio_id for search_lyric fallback")
	)
	flag.Parse()

	cfg := session.Load(session.DefaultPath())
	client, err := kg.New(kg.WithCookie(cfg.Cookie))
	if err != nil {
		fmt.Printf("init client failed: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(*id) == "" || strings.TrimSpace(*accesskey) == "" {
		if err := fillLyricCredentialFromSearch(client, keyword, hash, albumAudio, id, accesskey); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	resp, err := client.Lyric(ctx, kg.LyricRequest{
		Id:        strings.TrimSpace(*id),
		Accesskey: strings.TrimSpace(*accesskey),
		Fmt:       strings.TrimSpace(*format),
		Decode:    true,
		Cookie:    client.Cookie(),
	})
	if err != nil {
		fmt.Printf("fetch lyric failed: %v\n", err)
		if resp != nil && len(resp.RawBody) > 0 {
			fmt.Println(string(resp.RawBody))
		}
		os.Exit(1)
	}

	lrc := resp.ToLrc()
	if strings.TrimSpace(lrc) == "" {
		fmt.Println("lyric fetched but ToLrc() returned empty content")
		if len(resp.RawBody) > 0 {
			fmt.Println(string(resp.RawBody))
		}
		os.Exit(1)
	}

	fmt.Print(lrc)
	if !strings.HasSuffix(lrc, "\n") {
		fmt.Println()
	}
}

func fillLyricCredentialFromSearch(client *kg.Client, keyword, hash *string, albumAudio *int, id, accesskey *string) error {
	if strings.TrimSpace(*keyword) == "" && strings.TrimSpace(*hash) == "" {
		return fmt.Errorf("missing lyric id/accesskey; provide -id and -accesskey, or use -keyword/-hash to try search_lyric fallback")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	resp, err := client.SearchLyric(ctx, kg.SearchLyricRequest{
		Keywords:     strings.TrimSpace(*keyword),
		Hash:         strings.TrimSpace(*hash),
		AlbumAudioId: *albumAudio,
		Cookie:       client.Cookie(),
	})
	if err != nil {
		return fmt.Errorf("search_lyric failed: %v", err)
	}
	if resp == nil || len(resp.RawBody) == 0 || resp.Body == nil {
		return fmt.Errorf("search_lyric returned empty body; current upstream may reject this route, try passing -id and -accesskey directly")
	}

	candidates, _ := resp.Body["candidates"].([]any)
	if len(candidates) == 0 {
		return fmt.Errorf("search_lyric returned no candidates; try passing -id and -accesskey directly")
	}
	first, _ := candidates[0].(map[string]any)
	*id = strings.TrimSpace(fmt.Sprintf("%v", first["id"]))
	*accesskey = strings.TrimSpace(fmt.Sprintf("%v", first["accesskey"]))
	if *id == "" || *accesskey == "" || *id == "<nil>" || *accesskey == "<nil>" {
		return fmt.Errorf("search_lyric did not return usable id/accesskey; try passing them directly")
	}
	return nil
}

func firstNonEmpty(v ...string) string {
	for _, s := range v {
		if strings.TrimSpace(s) != "" {
			return s
		}
	}
	return ""
}

func envInt(key string) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return 0
	}
	var out int
	fmt.Sscanf(raw, "%d", &out)
	return out
}
