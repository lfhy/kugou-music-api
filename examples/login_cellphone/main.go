package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
	"github.com/lfhy/kugou-music-api/sdk"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cfgPath := session.DefaultPath()
	cfg := session.Load(cfgPath)

	client, err := sdk.New(sdk.WithCookie(cfg.Cookie))
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入手机号: ")
	mobile, _ := reader.ReadString('\n')
	mobile = strings.TrimSpace(mobile)
	if mobile == "" {
		fmt.Println("手机号不能为空")
		return
	}

	captchaResp, err := client.SendCaptcha(ctx, sdk.SendCaptchaRequest{
		Mobile: mobile,
		Cookie: cfg.Cookie,
	})
	if err != nil {
		fmt.Printf("发送验证码失败: %v\n", err)
		os.Exit(1)
	}
	if captchaResp == nil || fmt.Sprintf("%v", captchaResp.Body["status"]) != "1" {
		fmt.Printf("发送验证码失败，响应: %s\n", string(captchaResp.RawBody))
		os.Exit(1)
	}
	fmt.Println("验证码已发送，请查收短信。")

	fmt.Print("请输入短信验证码: ")
	code, _ := reader.ReadString('\n')
	code = strings.TrimSpace(code)
	if code == "" {
		fmt.Println("验证码不能为空")
		return
	}

	resp, err := client.LoginByCellphone(ctx, sdk.CellphoneLoginRequest{
		Mobile: mobile,
		Code:   code,
		Cookie: cfg.Cookie,
	})
	if err != nil {
		fmt.Printf("验证码登录失败: %v\n", err)
		if resp != nil {
			fmt.Println(string(resp.RawBody))
		}
		os.Exit(1)
	}
	if resp == nil || fmt.Sprintf("%v", resp.Body["status"]) != "1" {
		fmt.Printf("验证码登录失败，响应: %s\n", string(resp.RawBody))
		os.Exit(1)
	}
	cookie := client.Cookie()
	if strings.TrimSpace(cookie["token"]) == "" || strings.TrimSpace(cookie["userid"]) == "" || strings.TrimSpace(cookie["userid"]) == "0" {
		fmt.Printf("登录返回成功，但未拿到有效会话(token/userid)，响应: %s\n", string(resp.RawBody))
		os.Exit(1)
	}

	userResp, userErr := client.UserDetail(ctx, sdk.UserDetailRequest{Cookie: client.Cookie()})
	if userErr != nil {
		fmt.Printf("登录成功，但获取用户信息失败: %v\n", userErr)
		os.Exit(1)
	}
	if userResp == nil || fmt.Sprintf("%v", userResp.Body["status"]) != "1" {
		c := client.Cookie()
		fmt.Printf("登录成功，但用户信息接口失败: %s\n", string(userResp.RawBody))
		fmt.Printf("当前会话: userid=%s token_prefix=%s\n", c["userid"], tokenPrefix(c["token"]))
		os.Exit(1)
	}
	fmt.Printf("登录成功，用户信息: %s\n", string(userResp.RawBody))

	_ = session.Save(cfgPath, session.Config{
		Username:   cfg.Username,
		Cookie:     cookie,
		LastUserID: cookie["userid"],
	})
	fmt.Println("会话已保存：", cfgPath)
}

func tokenPrefix(token string) string {
	token = strings.TrimSpace(token)
	if len(token) <= 8 {
		return token
	}
	return token[:8] + "..."
}
