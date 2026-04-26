package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
)

type autoConfig struct {
	Accounts  []accountConfig `json:"accounts"`
	UpdatedAt string          `json:"updated_at,omitempty"`
}

type accountConfig struct {
	UserID            string `json:"user_id"`
	Nickname          string `json:"nickname,omitempty"`
	SessionFile       string `json:"session_file"`
	AddedAt           string `json:"added_at,omitempty"`
	LastSyncedAt      string `json:"last_synced_at,omitempty"`
	LastSignedDate    string `json:"last_signed_date,omitempty"`
	LastCheckinAt     string `json:"last_checkin_at,omitempty"`
	LastCheckinStatus string `json:"last_checkin_status,omitempty"`
	LastError         string `json:"last_error,omitempty"`
	VIPExpireAt       string `json:"vip_expire_at,omitempty"`
	Enabled           bool   `json:"enabled"`
}

func loadAutoConfig(path string) (autoConfig, error) {
	cfg := autoConfig{Accounts: []accountConfig{}}
	body, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if len(strings.TrimSpace(string(body))) == 0 {
		return cfg, nil
	}
	if err := json.Unmarshal(body, &cfg); err != nil {
		return cfg, err
	}
	if cfg.Accounts == nil {
		cfg.Accounts = []accountConfig{}
	}
	for i := range cfg.Accounts {
		if !cfg.Accounts[i].Enabled {
			cfg.Accounts[i].Enabled = true
		}
	}
	return cfg, nil
}

func saveAutoConfig(path string, cfg autoConfig) error {
	cfg.UpdatedAt = time.Now().Format(time.RFC3339)
	if cfg.Accounts == nil {
		cfg.Accounts = []accountConfig{}
	}
	if err := ensureParentDir(path); err != nil {
		return err
	}
	body, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o644)
}

func upsertAccount(cfg *autoConfig, account accountConfig) bool {
	for i := range cfg.Accounts {
		if cfg.Accounts[i].UserID == account.UserID {
			cfg.Accounts[i] = mergeAccount(cfg.Accounts[i], account)
			return false
		}
	}
	cfg.Accounts = append(cfg.Accounts, account)
	return true
}

func mergeAccount(oldAccount, newAccount accountConfig) accountConfig {
	if newAccount.AddedAt == "" {
		newAccount.AddedAt = oldAccount.AddedAt
	}
	if newAccount.LastCheckinAt == "" {
		newAccount.LastCheckinAt = oldAccount.LastCheckinAt
	}
	if newAccount.LastCheckinStatus == "" {
		newAccount.LastCheckinStatus = oldAccount.LastCheckinStatus
	}
	if newAccount.LastError == "" {
		newAccount.LastError = oldAccount.LastError
	}
	if newAccount.Enabled == oldAccount.Enabled {
		return newAccount
	}
	if oldAccount.Enabled && !newAccount.Enabled {
		return newAccount
	}
	newAccount.Enabled = true
	return newAccount
}

func removeAccount(cfg *autoConfig, userID string) bool {
	for i := range cfg.Accounts {
		if cfg.Accounts[i].UserID != userID {
			continue
		}
		cfg.Accounts = append(cfg.Accounts[:i], cfg.Accounts[i+1:]...)
		return true
	}
	return false
}

func findAccount(cfg autoConfig, userID string) (accountConfig, bool) {
	for _, account := range cfg.Accounts {
		if account.UserID == userID {
			return account, true
		}
	}
	return accountConfig{}, false
}

// migrateDefaultSession auto-registers the shared SDK session when no account exists yet.
func migrateDefaultSession(ctx context.Context, cfg *autoConfig, sourceSessionPath string) (bool, error) {
	if len(cfg.Accounts) > 0 {
		return false, nil
	}
	resolved := strings.TrimSpace(sourceSessionPath)
	if resolved == "" {
		resolved = session.DefaultPath()
	}
	resolved = filepath.Clean(resolved)
	sess := session.Load(resolved)
	if !session.HasLoginCookie(sess.Cookie) {
		return false, nil
	}
	status, _, err := loadAccountStatus(ctx, resolved)
	if err != nil {
		return false, err
	}
	account := newAccountConfig(status, resolved, time.Now())
	return upsertAccount(cfg, account), nil
}

func newAccountConfig(status accountStatus, sessionPath string, now time.Time) accountConfig {
	return accountConfig{
		UserID:         status.UserID,
		Nickname:       status.Nickname,
		SessionFile:    filepath.Clean(sessionPath),
		AddedAt:        now.Format(time.RFC3339),
		LastSyncedAt:   now.Format(time.RFC3339),
		LastSignedDate: status.LastSignedDate,
		VIPExpireAt:    status.VIPExpireAt,
		Enabled:        true,
	}
}
