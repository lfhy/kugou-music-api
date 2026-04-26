package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runAddAccountCommand(ctx context.Context, paths appPaths, method, targetSession string) error {
	if err := ensureDir(paths.BaseDir); err != nil {
		return err
	}
	if err := ensureDir(paths.SessionsDir); err != nil {
		return err
	}
	if strings.TrimSpace(method) == "" {
		var err error
		method, err = chooseLoginMethod()
		if err != nil {
			return err
		}
	}
	if strings.TrimSpace(targetSession) == "" {
		targetSession = managedSessionPath(paths.SessionsDir)
	}
	cfg, err := loadAutoConfig(paths.Config)
	if err != nil {
		return err
	}
	if err := runLoginExample(ctx, method, targetSession); err != nil {
		return err
	}
	status, _, err := loadAccountStatus(ctx, targetSession)
	if err != nil {
		return err
	}
	account := newAccountConfig(status, targetSession, time.Now())
	added := upsertAccount(&cfg, account)
	if err := saveAutoConfig(paths.Config, cfg); err != nil {
		return err
	}
	if added {
		fmt.Printf("已添加账号: %s (%s)\n", fallbackText(account.Nickname, "未知"), account.UserID)
	} else {
		fmt.Printf("已更新账号: %s (%s)\n", fallbackText(account.Nickname, "未知"), account.UserID)
	}
	fmt.Printf("session: %s\n", targetSession)
	return nil
}

func runDeleteAccountCommand(paths appPaths, userID string) error {
	cfg, err := loadAutoConfig(paths.Config)
	if err != nil {
		return err
	}
	if len(cfg.Accounts) == 0 {
		return fmt.Errorf("no accounts configured")
	}
	if strings.TrimSpace(userID) == "" {
		var err error
		userID, err = chooseDeleteAccount(cfg)
		if err != nil {
			return err
		}
	}
	if !removeAccount(&cfg, normalizedUserID(userID)) {
		return fmt.Errorf("account not found: %s", userID)
	}
	if err := saveAutoConfig(paths.Config, cfg); err != nil {
		return err
	}
	fmt.Printf("已移除账号: %s\n", normalizedUserID(userID))
	fmt.Println("注意：仅从自动签到配置中移除，不会删除原 session 文件。")
	return nil
}

func runLoginExample(ctx context.Context, method, sessionPath string) error {
	root, err := findRepoRoot()
	if err != nil {
		return err
	}
	pkg := "./examples/login_qrcode"
	if strings.EqualFold(strings.TrimSpace(method), "cellphone") || strings.EqualFold(strings.TrimSpace(method), "phone") {
		pkg = "./examples/login_cellphone"
	}
	cmd := exec.CommandContext(ctx, "go", "run", pkg)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "KUGOU_SESSION_FILE="+sessionPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func chooseLoginMethod() (string, error) {
	fmt.Println("请选择登录方式：")
	fmt.Println("1. 二维码登录")
	fmt.Println("2. 手机验证码登录")
	for {
		choice, err := promptLine("请输入 1 或 2: ")
		if err != nil {
			return "", err
		}
		switch strings.TrimSpace(choice) {
		case "1", "qr":
			return "qr", nil
		case "2", "cellphone", "phone":
			return "cellphone", nil
		default:
			fmt.Println("输入无效，请重新输入。")
		}
	}
}

func chooseDeleteAccount(cfg autoConfig) (string, error) {
	fmt.Println("当前账号：")
	for i, account := range cfg.Accounts {
		fmt.Printf("%d. %s (%s)\n", i+1, fallbackText(account.Nickname, "未知"), account.UserID)
	}
	for {
		choice, err := promptLine("请输入要删除的序号: ")
		if err != nil {
			return "", err
		}
		idx := asInt(choice) - 1
		if idx >= 0 && idx < len(cfg.Accounts) {
			return cfg.Accounts[idx].UserID, nil
		}
		fmt.Println("输入无效，请重新输入。")
	}
}

func promptLine(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
