package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
	"github.com/lfhy/kugou-music-api/sdk"
)

func main() {
	mode := flag.String("mode", "heart", "radio mode: heart | new | niche")
	limit := flag.Int("limit", 30, "print top N tracks")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	cfg := session.Load(session.DefaultPath())
	client, err := sdk.New(sdk.WithCookie(cfg.Cookie))
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		os.Exit(1)
	}

	resp, err := client.GetPersonalRadio(ctx, sdk.PersonalRadioRequest{
		Mode:   sdk.PersonalRadioMode(*mode),
		Cookie: cfg.Cookie,
	})
	if err != nil {
		fmt.Printf("get personal radio failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("mode: %s\n", resp.Mode)
	fmt.Printf("source: %s\n", resp.Source)
	if resp.Raw != nil {
		fmt.Printf("status: %d\n", resp.Raw.Status)
		if resp.Raw.Body != nil {
			fmt.Printf("biz_status: %v error_code: %v\n", resp.Raw.Body["status"], resp.Raw.Body["error_code"])
		}
	}
	fmt.Printf("tracks: %d\n", len(resp.Tracks))

	n := *limit
	if n <= 0 || n > len(resp.Tracks) {
		n = len(resp.Tracks)
	}

	for i := 0; i < n; i++ {
		t := resp.Tracks[i]
		fmt.Printf("%d. %s - %s\n", i+1, fallback(t.Name, "(unknown)"), fallback(t.Singer, "(unknown)"))

		url := "(unavailable)"
		if t.Hash != "" {
			playURL, uerr := client.GetSongPlayURL(ctx, sdk.SongPlayURLRequest{
				Hash:         t.Hash,
				AlbumAudioID: t.AlbumAudioID,
				Cookie:       cfg.Cookie,
			})
			if uerr == nil && playURL != "" {
				url = playURL
			}
		}
		fmt.Printf("   hash: %s\n", fallback(t.Hash, "-"))
		fmt.Printf("   url: %s\n", url)
	}
}

func fallback(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
