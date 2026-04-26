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
	name := flag.String("name", "", "playlist name (optional). default: YYYY-MM-DD红心日推")
	private := flag.Bool("private", false, "create private playlist")
	debug := flag.Bool("debug", false, "print debug response")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	cfg := session.Load(session.DefaultPath())
	debugEnabled := cfg.Debug || *debug
	if !session.HasLoginCookie(cfg.Cookie) {
		fmt.Println("未检测到本地登录会话，请先登录后再创建歌单。")
		os.Exit(1)
	}

	client, err := sdk.New(sdk.WithCookie(cfg.Cookie))
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		os.Exit(1)
	}

	playlistName := *name
	if playlistName == "" {
		playlistName = time.Now().Format("2006-01-02") + "红心日推"
	}

	created, err := client.CreatePlaylist(ctx, playlistName, *private, cfg.Cookie)
	if err != nil {
		fmt.Printf("create playlist failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("playlist created: name=%s listid=%d\n", playlistName, created.ListID)
	if debugEnabled {
		printResp("create_playlist", created.Raw)
	}
}

func printResp(label string, resp *sdk.Response) {
	if resp == nil {
		fmt.Printf("[debug] %s: nil response\n", label)
		return
	}
	fmt.Printf("[debug] %s: http_status=%d\n", label, resp.Status)
	fmt.Printf("[debug] %s: raw=%s\n", label, string(resp.RawBody))
}
