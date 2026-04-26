package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
)

type appPaths struct {
	BaseDir       string
	Config        string
	SourceSession string
	Log           string
	PID           string
	SessionsDir   string
}

func defaultAppPaths() appPaths {
	base := defaultBaseDir()
	return appPaths{
		BaseDir:       base,
		Config:        filepath.Join(base, "config.json"),
		SourceSession: session.DefaultPath(),
		Log:           filepath.Join(base, "auto_daily_checkin.log"),
		PID:           filepath.Join(base, "auto_daily_checkin.pid"),
		SessionsDir:   filepath.Join(base, "sessions"),
	}
}

func defaultBaseDir() string {
	if v := strings.TrimSpace(os.Getenv("KUGOU_AUTO_CHECKIN_DIR")); v != "" {
		return filepath.Clean(v)
	}
	cfgDir, err := os.UserConfigDir()
	if err == nil && strings.TrimSpace(cfgDir) != "" {
		return filepath.Join(cfgDir, "kugou-music-api", "auto_daily_checkin")
	}
	return ".kugou-auto-daily-checkin"
}

func ensureParentDir(path string) error {
	dir := filepath.Dir(path)
	if strings.TrimSpace(dir) == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func ensureDir(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	return os.MkdirAll(path, 0o755)
}

func managedSessionPath(dir string) string {
	stamp := time.Now().Format("20060102-150405")
	return filepath.Join(dir, fmt.Sprintf("session-%s.json", stamp))
}

func findRepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	cur := wd
	for {
		if _, err := os.Stat(filepath.Join(cur, "go.mod")); err == nil {
			return cur, nil
		}
		next := filepath.Dir(cur)
		if next == cur {
			return "", fmt.Errorf("repo root not found from %s", wd)
		}
		cur = next
	}
}
