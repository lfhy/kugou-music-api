package sdk

import (
	"context"
	"testing"

	corekugou "github.com/rushingrain/kugou-music-api/core/kugou"
)

func TestAlbumSongsMatchesJSRequestBody(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	client.requester = func(_ context.Context, cfg corekugou.RequestConfig) (corekugou.Response, error) {
		if cfg.URL != "/v1/album_audio/lite" || cfg.Method != "POST" {
			t.Fatalf("request = %s %s", cfg.Method, cfg.URL)
		}
		body, ok := cfg.Data.(map[string]any)
		if !ok {
			t.Fatalf("Data type = %T, want map[string]any", cfg.Data)
		}
		if body["album_id"] != "2997292" || body["page"] != 2 || body["pagesize"] != 30 {
			t.Fatalf("Data = %#v", body)
		}
		if body["is_buy"] != "" {
			t.Fatalf("is_buy = %#v, want empty string", body["is_buy"])
		}
		if _, ok := body["id"]; ok {
			t.Fatal("unexpected id field")
		}
		return stubCoreResponse(`{"status":1,"error_code":0,"data":{"songs":[]}}`), nil
	}

	_, err = client.AlbumSongs(context.Background(), AlbumSongsRequest{
		Id:       "2997292",
		Page:     2,
		Pagesize: 30,
	})
	if err != nil {
		t.Fatalf("AlbumSongs() error = %v", err)
	}
}
