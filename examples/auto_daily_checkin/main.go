package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	command, args := normalizeCommand(os.Args[1:])
	defaults := defaultAppPaths()

	var err error
	switch command {
	case "start":
		fs := flag.NewFlagSet("start", flag.ExitOnError)
		configPath := fs.String("config", defaults.Config, "auto config path")
		sourceSession := fs.String("source-session-file", defaults.SourceSession, "shared SDK session path")
		logPath := fs.String("log", defaults.Log, "log file path")
		pidPath := fs.String("pid", defaults.PID, "pid file path")
		interval := fs.Duration("interval", time.Minute, "loop interval")
		_ = fs.Parse(args)
		err = runStartCommand(ctx, appPaths{
			BaseDir:       defaults.BaseDir,
			Config:        *configPath,
			SourceSession: *sourceSession,
			Log:           *logPath,
			PID:           *pidPath,
			SessionsDir:   defaults.SessionsDir,
		}, *interval)
	case "daemon":
		fs := flag.NewFlagSet("daemon", flag.ExitOnError)
		configPath := fs.String("config", defaults.Config, "auto config path")
		sourceSession := fs.String("source-session-file", defaults.SourceSession, "shared SDK session path")
		logPath := fs.String("log", defaults.Log, "log file path")
		pidPath := fs.String("pid", defaults.PID, "pid file path")
		interval := fs.Duration("interval", time.Minute, "loop interval")
		_ = fs.Parse(args)
		err = runDaemonCommand(appPaths{
			BaseDir:       defaults.BaseDir,
			Config:        *configPath,
			SourceSession: *sourceSession,
			Log:           *logPath,
			PID:           *pidPath,
			SessionsDir:   defaults.SessionsDir,
		}, *interval)
	case "stop":
		fs := flag.NewFlagSet("stop", flag.ExitOnError)
		pidPath := fs.String("pid", defaults.PID, "pid file path")
		_ = fs.Parse(args)
		err = stopDaemon(*pidPath)
	case "status":
		fs := flag.NewFlagSet("status", flag.ExitOnError)
		configPath := fs.String("config", defaults.Config, "auto config path")
		sourceSession := fs.String("source-session-file", defaults.SourceSession, "shared SDK session path")
		pidPath := fs.String("pid", defaults.PID, "pid file path")
		_ = fs.Parse(args)
		err = runStatusCommand(ctx, appPaths{
			BaseDir:       defaults.BaseDir,
			Config:        *configPath,
			SourceSession: *sourceSession,
			PID:           *pidPath,
			Log:           defaults.Log,
			SessionsDir:   defaults.SessionsDir,
		})
	case "add-account":
		fs := flag.NewFlagSet("add-account", flag.ExitOnError)
		configPath := fs.String("config", defaults.Config, "auto config path")
		method := fs.String("method", "", "qr or cellphone")
		targetSession := fs.String("session-file", "", "target SDK session path")
		_ = fs.Parse(args)
		err = runAddAccountCommand(ctx, appPaths{
			BaseDir:       defaults.BaseDir,
			Config:        *configPath,
			SourceSession: defaults.SourceSession,
			Log:           defaults.Log,
			PID:           defaults.PID,
			SessionsDir:   defaults.SessionsDir,
		}, *method, *targetSession)
	case "del-account":
		fs := flag.NewFlagSet("del-account", flag.ExitOnError)
		configPath := fs.String("config", defaults.Config, "auto config path")
		userID := fs.String("user", "", "user id to remove")
		_ = fs.Parse(args)
		err = runDeleteAccountCommand(appPaths{Config: *configPath}, *userID)
	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func normalizeCommand(args []string) (string, []string) {
	if len(args) == 0 {
		return "start", nil
	}
	cmd := strings.TrimSpace(args[0])
	switch cmd {
	case "start", "daemon", "stop", "status", "add-account", "del-account":
		return cmd, args[1:]
	default:
		if strings.HasPrefix(cmd, "-") {
			return "start", args
		}
		return cmd, args[1:]
	}
}

func printUsage() {
	fmt.Println("用法: go run ./examples/auto_daily_checkin [command] [flags]")
	fmt.Println("commands:")
	fmt.Println("  start         启动自动签到主循环，默认命令")
	fmt.Println("  daemon        后台启动 start")
	fmt.Println("  stop          停止当前后台实例")
	fmt.Println("  status        查看账号自动签到状态")
	fmt.Println("  add-account   交互式登录并添加账号")
	fmt.Println("  del-account   从自动签到配置中移除账号")
}
