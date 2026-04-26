package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Username    string            `json:"username,omitempty"`
	Cookie      map[string]string `json:"cookie,omitempty"`
	Debug       bool              `json:"debug,omitempty"`
	DownloadDir string            `json:"download_dir,omitempty"`
	UpdatedAt   string            `json:"updated_at,omitempty"`
	LastUserID  string            `json:"last_userid,omitempty"`
}

func DefaultPath() string {
	if p := strings.TrimSpace(os.Getenv("KUGOU_SESSION_FILE")); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return ".kugou_music_api_session.json"
	}
	return filepath.Join(home, ".kugou_music_api_session.json")
}

func Load(path string) Config {
	cfg := Config{Cookie: map[string]string{}}
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

func Save(path string, cfg Config) error {
	if cfg.Cookie == nil {
		cfg.Cookie = map[string]string{}
	}
	if cfg.UpdatedAt == "" {
		cfg.UpdatedAt = time.Now().Format(time.RFC3339)
	}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(path, b, 0o600)
}

func HasLoginCookie(cookie map[string]string) bool {
	if cookie == nil {
		return false
	}
	return strings.TrimSpace(cookie["token"]) != "" && strings.TrimSpace(cookie["userid"]) != ""
}
