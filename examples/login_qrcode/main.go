package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
	"github.com/lfhy/kugou-music-api/sdk"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cfgPath := session.DefaultPath()
	cfg := session.Load(cfgPath)
	client, err := sdk.New(sdk.WithCookie(cfg.Cookie))
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		os.Exit(1)
	}

	keyResp, err := client.LoginQrKey(ctx, sdk.LoginQrKeyRequest{})
	if err != nil {
		fmt.Printf("获取二维码 key 失败: %v\n", err)
		os.Exit(1)
	}

	key := pickStringDeep(keyResp.Body, "qrcode", "key", "qrkey")
	if key == "" {
		fmt.Printf("未从响应中提取到二维码 key: %s\n", string(keyResp.RawBody))
		os.Exit(1)
	}

	qrURL := "https://h5.kugou.com/apps/loginQRCode/html/index.html?qrcode=" + url.QueryEscape(key)
	fmt.Println("请使用酷狗 App 扫码登录：")
	fmt.Println(qrURL)
	fmt.Println("轮询登录状态中（最多 120 秒）...")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("登录超时，请重试。")
			os.Exit(1)
		case <-ticker.C:
			checkResp, err := client.LoginQrCheck(context.Background(), sdk.LoginQrCheckRequest{Key: key})
			if err != nil {
				fmt.Printf("检查二维码状态失败: %v\n", err)
				continue
			}

			status := pickIntDeep(checkResp.Body, "status")
			switch status {
			case 0:
				fmt.Println("二维码已过期，请重新运行。")
				os.Exit(1)
			case 1:
				fmt.Println("等待扫码...")
			case 2:
				fmt.Println("已扫码，等待确认...")
			case 4:
				token := pickStringDeep(checkResp.Body, "token")
				userid := pickStringDeep(checkResp.Body, "userid", "user_id", "uid")
				if token != "" {
					client.SetCookie("token", token)
				}
				if userid != "" {
					client.SetCookie("userid", userid)
				}

				cookie := client.Cookie()
				_ = session.Save(cfgPath, session.Config{
					Username:   cfg.Username,
					Cookie:     cookie,
					LastUserID: cookie["userid"],
				})
				fmt.Println("扫码登录成功，会话已保存：", cfgPath)

				userResp, userErr := client.UserDetail(context.Background(), sdk.UserDetailRequest{})
				if userErr == nil {
					fmt.Printf("用户信息: %s\n", string(userResp.RawBody))
				}
				return
			default:
				fmt.Printf("未知二维码状态: %d, body=%s\n", status, string(checkResp.RawBody))
			}
		}
	}
}

func pickStringDeep(v any, keys ...string) string {
	m, ok := v.(map[string]any)
	if !ok || m == nil {
		return ""
	}
	for _, k := range keys {
		if x, ok := m[k]; ok {
			s := fmt.Sprintf("%v", x)
			if s != "" && s != "<nil>" {
				return s
			}
		}
	}
	for _, x := range m {
		if mm, ok := x.(map[string]any); ok {
			if s := pickStringDeep(mm, keys...); s != "" {
				return s
			}
		}
	}
	return ""
}

func pickIntDeep(v any, key string) int {
	m, ok := v.(map[string]any)
	if !ok || m == nil {
		return -1
	}
	if x, ok := m[key]; ok {
		var n int
		fmt.Sscanf(fmt.Sprintf("%v", x), "%d", &n)
		return n
	}
	for _, x := range m {
		if mm, ok := x.(map[string]any); ok {
			n := pickIntDeep(mm, key)
			if n != -1 {
				return n
			}
		}
	}
	return -1
}
