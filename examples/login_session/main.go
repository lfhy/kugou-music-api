package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/sdk"
)

type SessionConfig struct {
	Username   string            `json:"username,omitempty"`
	Cookie     map[string]string `json:"cookie,omitempty"`
	UpdatedAt  string            `json:"updated_at,omitempty"`
	LastUserID string            `json:"last_userid,omitempty"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfgPath := defaultConfigPath()
	cfg := loadConfig(cfgPath)

	client, err := sdk.New(sdk.WithCookie(cfg.Cookie))
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		os.Exit(1)
	}

	if hasLoginCookie(cfg.Cookie) {
		ok, msg := fetchUserInfo(ctx, client)
		if ok {
			fmt.Println("已检测到本地登录会话，自动获取用户信息成功。")
			fmt.Println(msg)
			saveConfig(cfgPath, SessionConfig{
				Username:   cfg.Username,
				Cookie:     client.Cookie(),
				UpdatedAt:  time.Now().Format(time.RFC3339),
				LastUserID: client.Cookie()["userid"],
			})
			return
		}

		fmt.Println("检测到历史会话，但已失效，尝试 token 刷新登录...")
		_, _ = client.LoginByToken(ctx, sdk.TokenLoginRequest{
			Token:  cfg.Cookie["token"],
			UserID: cfg.Cookie["userid"],
			Cookie: cfg.Cookie,
		})
		ok, msg = fetchUserInfo(ctx, client)
		if ok {
			fmt.Println("token 刷新成功，已恢复登录态。")
			fmt.Println(msg)
			saveConfig(cfgPath, SessionConfig{
				Username:   cfg.Username,
				Cookie:     client.Cookie(),
				UpdatedAt:  time.Now().Format(time.RFC3339),
				LastUserID: client.Cookie()["userid"],
			})
			return
		}
	}

	fmt.Println("当前未登录，请先登录。")
	username, password, ok := promptLoginInput()
	if !ok {
		fmt.Println("已取消登录。")
		return
	}

	resp, err := client.LoginByPassword(ctx, sdk.PasswordLoginRequest{
		Username: username,
		Password: password,
		Cookie:   cfg.Cookie,
	})
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		printRespBrief(resp)
		os.Exit(1)
	}

	if !isBizSuccess(resp) {
		fmt.Println("登录请求已返回，但业务未成功。")
		printRespBrief(resp)
		os.Exit(1)
	}

	ok, msg := fetchUserInfo(ctx, client)
	if !ok {
		fmt.Println("登录成功，但获取用户信息失败。")
		fmt.Println(msg)
		os.Exit(1)
	}

	saveConfig(cfgPath, SessionConfig{
		Username:   username,
		Cookie:     client.Cookie(),
		UpdatedAt:  time.Now().Format(time.RFC3339),
		LastUserID: client.Cookie()["userid"],
	})
	fmt.Println("登录成功，已保存会话配置：", cfgPath)
	fmt.Println(msg)
}

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return ".kugou_music_api_session.json"
	}
	return filepath.Join(home, ".kugou_music_api_session.json")
}

func loadConfig(path string) SessionConfig {
	cfg := SessionConfig{Cookie: map[string]string{}}
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(b, &cfg)
	if cfg.Cookie == nil {
		cfg.Cookie = map[string]string{}
	}
	return cfg
}

func saveConfig(path string, cfg SessionConfig) {
	if cfg.Cookie == nil {
		cfg.Cookie = map[string]string{}
	}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	_ = os.WriteFile(path, b, 0o600)
}

func hasLoginCookie(cookie map[string]string) bool {
	if cookie == nil {
		return false
	}
	return strings.TrimSpace(cookie["token"]) != "" && strings.TrimSpace(cookie["userid"]) != ""
}

func promptLoginInput() (username, password string, ok bool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入账号(手机号/用户名，直接回车取消): ")
	u, _ := reader.ReadString('\n')
	u = strings.TrimSpace(u)
	if u == "" {
		return "", "", false
	}
	fmt.Print("请输入密码: ")
	p, _ := reader.ReadString('\n')
	p = strings.TrimSpace(p)
	if p == "" {
		return "", "", false
	}
	return u, p, true
}

func fetchUserInfo(ctx context.Context, client *sdk.Client) (bool, string) {
	resp, err := client.UserDetail(ctx, sdk.UserDetailRequest{Cookie: client.Cookie()})
	if err != nil {
		return false, fmt.Sprintf("UserDetail 请求失败: %v", err)
	}
	if !isBizSuccess(resp) {
		return false, briefBody(resp)
	}
	uid := client.Cookie()["userid"]
	if uid == "" {
		uid = pickString(resp.Body, "userid", "user_id", "uid")
	}
	nickname := pickString(resp.Body, "nickname", "name", "username", "user_name")
	if nickname == "" {
		nickname = "(未返回昵称字段)"
	}
	if uid == "" {
		uid = "(未返回用户ID字段)"
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
		if errorCode == "<nil>" || errorCode == "" || errorCode == "0" || errorCode == "0.0" {
			return true
		}
		// Some endpoints return status=1 with non-zero error_code in non-fatal cases.
		if strings.EqualFold(errorCode, "<nil>") {
			return true
		}
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

func printRespBrief(resp *sdk.Response) {
	fmt.Println(briefBody(resp))
}

func briefBody(resp *sdk.Response) string {
	if resp == nil || resp.Body == nil {
		if resp != nil {
			return string(resp.RawBody)
		}
		return "(empty response)"
	}
	status := fmt.Sprintf("%v", resp.Body["status"])
	errorCode := fmt.Sprintf("%v", resp.Body["error_code"])
	msg := fmt.Sprintf("%v", firstNonNil(resp.Body["error"], resp.Body["msg"], resp.Body["error_msg"]))
	return fmt.Sprintf("status=%s error_code=%s msg=%s", status, errorCode, msg)
}

func firstNonNil(values ...any) any {
	for _, v := range values {
		if v != nil {
			s := strings.TrimSpace(fmt.Sprintf("%v", v))
			if s != "" && s != "<nil>" {
				return v
			}
		}
	}
	return ""
}
