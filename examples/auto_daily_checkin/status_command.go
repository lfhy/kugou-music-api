package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func runStatusCommand(ctx context.Context, paths appPaths) error {
	if err := ensureDir(paths.BaseDir); err != nil {
		return err
	}
	cfg, err := loadAutoConfig(paths.Config)
	if err != nil {
		return err
	}
	migrated, err := migrateDefaultSession(ctx, &cfg, paths.SourceSession)
	if err != nil {
		return err
	}
	if migrated {
		if err := saveAutoConfig(paths.Config, cfg); err != nil {
			return err
		}
	}

	state, stale := daemonState(paths.PID)
	if state.PID > 0 && !stale {
		fmt.Printf("daemon: running (pid=%d, started_at=%s)\n", state.PID, state.StartedAt)
	} else if state.PID > 0 {
		fmt.Printf("daemon: stale pid file (pid=%d)\n", state.PID)
	} else {
		fmt.Println("daemon: stopped")
	}

	if len(cfg.Accounts) == 0 {
		fmt.Println("未配置自动签到账号，请先运行 add-account。")
		return nil
	}

	now := time.Now()
	for i := range cfg.Accounts {
		status, _, err := loadAccountStatus(ctx, cfg.Accounts[i].SessionFile)
		if err != nil {
			cfg.Accounts[i].LastError = err.Error()
			cfg.Accounts[i].LastSyncedAt = now.Format(time.RFC3339)
			fmt.Printf("\n账号 %d\n", i+1)
			fmt.Printf("  user_id: %s\n", cfg.Accounts[i].UserID)
			fmt.Printf("  error: %s\n", err)
			continue
		}
		applyStatusSnapshot(&cfg.Accounts[i], status, now)
		cfg.Accounts[i].LastError = ""

		fmt.Printf("\n账号 %d\n", i+1)
		fmt.Printf("  名称: %s\n", fallbackText(cfg.Accounts[i].Nickname, "未知"))
		fmt.Printf("  UserID: %s\n", fallbackText(cfg.Accounts[i].UserID, "未知"))
		fmt.Printf("  Session: %s\n", cfg.Accounts[i].SessionFile)
		fmt.Printf("  最后一次签到时间: %s\n", displayLastCheckin(cfg.Accounts[i]))
		fmt.Printf("  会员过期时间: %s\n", fallbackText(cfg.Accounts[i].VIPExpireAt, "未知"))
		fmt.Printf("  今日是否已签到: %s\n", yesNo(status.SignedToday))
		fmt.Printf("  最近同步时间: %s\n", fallbackText(cfg.Accounts[i].LastSyncedAt, "未知"))
	}

	return saveAutoConfig(paths.Config, cfg)
}

func fallbackText(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func yesNo(v bool) string {
	if v {
		return "是"
	}
	return "否"
}
