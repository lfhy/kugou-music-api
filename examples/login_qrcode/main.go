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

	key := keyResp.QRCodeKey()
	if key == "" {
		fmt.Printf("未从响应中提取到二维码 key: %s\n", string(keyResp.RawBody))
		os.Exit(1)
	}

	qrCreateResp, err := client.LoginQrCreate(ctx, sdk.LoginQrCreateRequest{Key: key, Qrimg: true})
	if err != nil {
		fmt.Printf("生成二维码信息失败: %v\n", err)
		os.Exit(1)
	}
	qrURL := qrCreateResp.URL()
	if qrURL == "" {
		qrURL = keyResp.QRCodeURL()
	}
	qrPNGPath, saveErr := saveQRCodePNG(key, qrCreateResp.Base64())
	if saveErr != nil {
		fmt.Printf("保存二维码 PNG 失败: %v\n", saveErr)
	}
	fmt.Println("请使用酷狗 App 扫码登录：")
	renderTerminalQRCode(qrURL)
	if qrPNGPath != "" {
		fmt.Println("本地二维码 PNG：")
		fmt.Println(qrPNGPath)
	}
	fmt.Println("二维码地址：")
	fmt.Println(qrURL)
	fmt.Println("轮询登录状态中（最多 120 秒）...")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	lastStatus := -999

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

			status := checkResp.StatusCode()
			if status == lastStatus && (status == 1 || status == 2) {
				continue
			}
			lastStatus = status
			switch status {
			case 0:
				fmt.Println("二维码已过期，请重新运行。")
				os.Exit(1)
			case 1:
				fmt.Println("等待扫码...")
			case 2:
				fmt.Println("已扫码，等待确认...")
			case 4:
				token := checkResp.Token()
				userid := checkResp.UserID()
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
