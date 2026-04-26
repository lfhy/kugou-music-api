package sdk

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
)

func TestLyricLiveToLRC(t *testing.T) {
	if strings.TrimSpace(os.Getenv("KUGOU_LIVE_TESTS")) != "1" {
		t.Skip("set KUGOU_LIVE_TESTS=1 to enable live lyric test")
	}

	cfg := session.Load(session.DefaultPath())
	client, err := New(WithCookie(cfg.Cookie))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	id := strings.TrimSpace(os.Getenv("KUGOU_LYRIC_ID"))
	accesskey := strings.TrimSpace(os.Getenv("KUGOU_LYRIC_ACCESSKEY"))
	if id == "" || accesskey == "" {
		t.Skip("missing KUGOU_LYRIC_ID or KUGOU_LYRIC_ACCESSKEY; current upstream search_lyric is not stable enough to use as a required pre-step")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	resp, err := client.Lyric(ctx, LyricRequest{
		Id:        id,
		Accesskey: accesskey,
		Decode:    true,
		Cookie:    client.Cookie(),
	})
	if err != nil {
		t.Fatalf("Lyric() error = %v", err)
	}
	if resp == nil {
		t.Fatal("Lyric() returned nil response")
	}
	if got := strings.TrimSpace(resp.ToLrc()); got == "" {
		t.Fatalf("ToLrc() returned empty content, raw=%s", string(resp.RawBody))
	}
}
