package sdk

import (
	"context"
	"testing"

	corekugou "github.com/rushingrain/kugou-music-api/core/kugou"
)

// TestEnsureLoginValidKeepsRequestCookie 校验登录守卫不会丢弃只通过请求参数传入的登录态。
//
// 与 JS 版一致：调用方既可以把登录态放进 client（WithCookie / SetCookie），
// 也可以只放在单次请求的 cookie 里；后者在校验通过后仍必须保留，
// 否则后续接口会以未登录态发出（例如每日推荐返回 200103）。
func TestEnsureLoginValidKeepsRequestCookie(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	client.requester = func(_ context.Context, cfg corekugou.RequestConfig) (corekugou.Response, error) {
		if cfg.URL != "/v3/get_my_info" {
			t.Fatalf("unexpected request url = %q", cfg.URL)
		}
		if got := cfg.Cookie["token"]; got != "request-token" {
			t.Fatalf("validation cookie token = %q, want request-token", got)
		}
		return stubCoreResponse(`{"status":1,"error_code":0,"data":{"userid":123}}`), nil
	}

	merged, ok := client.ensureLoginValid(context.Background(), map[string]string{
		"token":  "request-token",
		"userid": "123",
	})
	if !ok {
		t.Fatal("ensureLoginValid() = false, want true")
	}
	if merged["token"] != "request-token" || merged["userid"] != "123" {
		t.Fatalf("merged cookie = %v, want 保留请求传入的 token / userid", merged)
	}
}
