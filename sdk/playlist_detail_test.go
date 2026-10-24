package sdk

import (
	"context"
	"testing"

	corekugou "github.com/rushingrain/kugou-music-api/core/kugou"
)

func TestPlaylistDetailMatchesJSRequestBody(t *testing.T) {
	client, err := New(WithCookie(map[string]string{
		"token":  "token-value",
		"userid": "123456",
	}))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	client.requester = func(_ context.Context, cfg corekugou.RequestConfig) (corekugou.Response, error) {
		if cfg.URL != "/v3/get_list_info" {
			t.Fatalf("URL = %q, want /v3/get_list_info", cfg.URL)
		}
		if cfg.Method != "POST" {
			t.Fatalf("Method = %q, want POST", cfg.Method)
		}

		body, ok := cfg.Data.(map[string]any)
		if !ok {
			t.Fatalf("Data type = %T, want map[string]any", cfg.Data)
		}
		data, ok := body["data"].([]map[string]any)
		if !ok || len(data) != 2 {
			t.Fatalf("data = %#v, want two playlist IDs", body["data"])
		}
		if data[0]["global_collection_id"] != "collection-1" || data[1]["global_collection_id"] != "collection-2" {
			t.Fatalf("data = %#v", data)
		}
		if body["userid"] != "123456" || body["token"] != "token-value" {
			t.Fatalf("auth fields = userid:%v token:%v", body["userid"], body["token"])
		}
		if _, ok := cfg.Params["ids"]; ok {
			t.Fatal("unexpected ids query parameter")
		}

		return stubCoreResponse(`{"status":1,"error_code":0,"data":{}}`), nil
	}

	_, err = client.PlaylistDetail(context.Background(), PlaylistDetailRequest{
		Ids: "collection-1,collection-2",
	})
	if err != nil {
		t.Fatalf("PlaylistDetail() error = %v", err)
	}
}
