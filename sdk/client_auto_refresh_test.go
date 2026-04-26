package sdk

import (
	"context"
	"fmt"
	"strings"
	"testing"

	corekugou "github.com/lfhy/kugou-music-api/core/kugou"
)

func TestCallAutoRefreshesExpiredLogin(t *testing.T) {
	client, err := New(WithCookie(map[string]string{"token": "old-token", "userid": "123"}))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	targetCalls := 0
	refreshCalls := 0
	validateCalls := 0
	client.requester = func(ctx context.Context, cfg corekugou.RequestConfig) (corekugou.Response, error) {
		switch {
		case strings.Contains(cfg.URL, "login_by_token"):
			refreshCalls++
			return stubCoreResponse(`{"status":1,"error_code":0,"data":{"token":"new-token","userid":123}}`, "token=new-token", "userid=123"), nil
		case cfg.URL == "/v3/get_my_info":
			validateCalls++
			if got := strings.TrimSpace(cfg.Cookie["token"]); got != "new-token" {
				t.Fatalf("validation cookie token = %q, want new-token", got)
			}
			return stubCoreResponse(`{"status":1,"error_code":0,"data":{"userid":123}}`), nil
		default:
			targetCalls++
			switch targetCalls {
			case 1:
				if got := fmt.Sprintf("%v", cfg.Params["token"]); got != "old-token" {
					t.Fatalf("first request token = %q, want old-token", got)
				}
				return stubCoreResponse(`{"status":0,"error_code":20018,"data":null}`), nil
			case 2:
				if got := fmt.Sprintf("%v", cfg.Params["token"]); got != "new-token" {
					t.Fatalf("retry request token = %q, want new-token", got)
				}
				if got := strings.TrimSpace(cfg.Cookie["token"]); got != "new-token" {
					t.Fatalf("retry cookie token = %q, want new-token", got)
				}
				return stubCoreResponse(`{"status":1,"error_code":0,"data":{"items":[]}}`), nil
			default:
				t.Fatalf("unexpected request #%d for %s", targetCalls, cfg.URL)
				return corekugou.Response{}, nil
			}
		}
	}

	resp, err := client.Call(context.Background(), RouteSearch, Request{
		Params: map[string]any{"keywords": "test", "token": "old-token", "userid": "123"},
	})
	if err != nil {
		t.Fatalf("Call() error = %v", err)
	}
	if targetCalls != 2 {
		t.Fatalf("targetCalls = %d, want 2", targetCalls)
	}
	if refreshCalls != 1 {
		t.Fatalf("refreshCalls = %d, want 1", refreshCalls)
	}
	if validateCalls != 1 {
		t.Fatalf("validateCalls = %d, want 1", validateCalls)
	}
	if !isBizSuccessCode(resp) {
		t.Fatalf("Call() returned non-success body: %s", string(resp.RawBody))
	}
	if got := strings.TrimSpace(client.Cookie()["token"]); got != "new-token" {
		t.Fatalf("cookie pool token = %q, want new-token", got)
	}
}

func TestCallReturnsErrorWhenAutoRefreshDisabled(t *testing.T) {
	client, err := New(
		WithCookie(map[string]string{"token": "old-token", "userid": "123"}),
		WithAutoRefresh(false),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	refreshCalls := 0
	client.requester = func(ctx context.Context, cfg corekugou.RequestConfig) (corekugou.Response, error) {
		if strings.Contains(cfg.URL, "login_by_token") {
			refreshCalls++
		}
		return stubCoreResponse(`{"status":0,"error_code":20018,"data":null}`), nil
	}

	resp, err := client.Call(context.Background(), RouteSearch, Request{
		Params: map[string]any{"keywords": "test", "token": "old-token", "userid": "123"},
	})
	if err == nil {
		t.Fatal("Call() error = nil, want auth error")
	}
	if !strings.Contains(err.Error(), "auto refresh disabled") {
		t.Fatalf("Call() error = %v, want auto refresh disabled", err)
	}
	if refreshCalls != 0 {
		t.Fatalf("refreshCalls = %d, want 0", refreshCalls)
	}
	if resp == nil || !isAuthExpiredResponse(resp) {
		t.Fatalf("Call() response = %#v, want auth-expired response", resp)
	}
}

func stubCoreResponse(body string, setCookies ...string) corekugou.Response {
	return corekugou.Response{
		Status:  200,
		Body:    []byte(body),
		Cookie:  setCookies,
		Headers: map[string]string{},
	}
}
