package main

import (
	"context"
	"fmt"
	"os"
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
		fmt.Println("未检测到本地登录会话，请先运行登录示例。")
		os.Exit(1)
	}

	client, err := sdk.New(sdk.WithCookie(cfg.Cookie))
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		os.Exit(1)
	}

	resp, err := client.LoginByToken(ctx, sdk.TokenLoginRequest{
		Token:  cfg.Cookie["token"],
		UserID: cfg.Cookie["userid"],
		Cookie: cfg.Cookie,
	})
	if err != nil {
		fmt.Printf("token 续期失败: %v\n", err)
		if resp != nil {
			fmt.Println(string(resp.RawBody))
		}
		os.Exit(1)
	}

	cookie := client.Cookie()
	_ = session.Save(cfgPath, session.Config{
		Username:   cfg.Username,
		Cookie:     cookie,
		LastUserID: cookie["userid"],
	})
	fmt.Println("token 续期成功，配置已更新：", cfgPath)

	userResp, userErr := client.UserDetail(ctx, sdk.UserDetailRequest{})
	if userErr != nil {
		fmt.Printf("获取用户信息失败: %v\n", userErr)
		os.Exit(1)
	}
	fmt.Printf("用户信息: %s\n", string(userResp.RawBody))
}
