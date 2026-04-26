package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
	"github.com/lfhy/kugou-music-api/sdk"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfgPath := session.DefaultPath()
	cfg := session.Load(cfgPath)
	if !session.HasLoginCookie(cfg.Cookie) {
		fmt.Println("未检测到本地登录会话，请先登录。")
		os.Exit(1)
	}

	client, err := sdk.New(sdk.WithCookie(cfg.Cookie))
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		os.Exit(1)
	}

	ok, info := fetchUserInfo(ctx, client)
	if ok {
		fmt.Println("token 有效。")
		fmt.Println(info)
		return
	}

	fmt.Printf("token 失效，尝试刷新: %s\n", info)
	refreshResp, refreshErr := client.LoginByToken(ctx, sdk.TokenLoginRequest{
		Token:  cfg.Cookie["token"],
		UserID: cfg.Cookie["userid"],
		Cookie: cfg.Cookie,
	})
	if refreshErr != nil || !isBizSuccess(refreshResp) {
		fmt.Println("刷新 token 失败。")
		if refreshResp != nil {
			fmt.Println(string(refreshResp.RawBody))
		}
		os.Exit(1)
	}

	ok, info = fetchUserInfo(ctx, client)
	if !ok {
		fmt.Printf("刷新后仍无法获取用户信息: %s\n", info)
		os.Exit(1)
	}

	cookie := client.Cookie()
	_ = session.Save(cfgPath, session.Config{
		Username:   cfg.Username,
		Cookie:     cookie,
		LastUserID: cookie["userid"],
	})
	fmt.Println("token 已刷新并保存。")
	fmt.Println(info)
}

func fetchUserInfo(ctx context.Context, client *sdk.Client) (bool, string) {
	resp, err := client.UserDetail(ctx, sdk.UserDetailRequest{Cookie: client.Cookie()})
	if err != nil {
		return false, err.Error()
	}
	if !isBizSuccess(resp) {
		return false, string(resp.RawBody)
	}
	nickname := pickString(resp.Body, "nickname", "name", "username", "user_name")
	if nickname == "" {
		nickname = "(未返回昵称字段)"
	}
	uid := strings.TrimSpace(client.Cookie()["userid"])
	if uid == "" {
		uid = pickString(resp.Body, "userid", "user_id", "uid")
	}
	if uid == "" {
		uid = "(未返回用户ID)"
	}
	return true, fmt.Sprintf("用户信息: nickname=%s, userid=%s", nickname, uid)
}

func isBizSuccess(resp *sdk.Response) bool {
	if resp == nil || resp.Body == nil {
		return false
	}
	status := fmt.Sprintf("%v", resp.Body["status"])
	errorCode := fmt.Sprintf("%v", resp.Body["error_code"])
	if status == "1" || status == "1.0" {
		return errorCode == "" || errorCode == "<nil>" || errorCode == "0" || errorCode == "0.0"
	}
	return false
}

func pickString(root map[string]any, keys ...string) string {
	for _, key := range keys {
		if v, ok := root[key]; ok {
			s := strings.TrimSpace(fmt.Sprintf("%v", v))
			if s != "" && s != "<nil>" {
				return s
			}
		}
	}
	for _, v := range root {
		if m, ok := v.(map[string]any); ok {
			if s := pickString(m, keys...); s != "" {
				return s
			}
		}
	}
	return ""
}

