package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
	"github.com/lfhy/kugou-music-api/sdk"
)

func runStartCommand(ctx context.Context, paths appPaths, interval time.Duration) error {
	if err := ensureDir(paths.BaseDir); err != nil {
		return err
	}
	if err := ensureDir(paths.SessionsDir); err != nil {
		return err
	}
	if err := acquireDaemonLock(paths.PID); err != nil {
		return err
	}
	defer releaseDaemonLock(paths.PID)

	logger := newLogger(paths.Log, strings.TrimSpace(os.Getenv("KUGOU_AUTO_CHECKIN_DAEMON")) == "")
	logger.Printf("auto_daily_checkin started")

	cfg, err := loadAutoConfig(paths.Config)
	if err != nil {
		return err
	}
	migrated, err := migrateDefaultSession(ctx, &cfg, paths.SourceSession)
	if err != nil {
		return err
	}
	if migrated {
		logger.Printf("migrated default session into auto config")
	}
	if err := saveAutoConfig(paths.Config, cfg); err != nil {
		return err
	}
	if len(cfg.Accounts) == 0 {
		return fmt.Errorf("no accounts configured; run add-account first")
	}

	if err := syncAndCheckin(ctx, &cfg, logger); err != nil {
		logger.Printf("initial sync failed: %v", err)
	}
	if err := saveAutoConfig(paths.Config, cfg); err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Printf("auto_daily_checkin stopped")
			return nil
		case <-ticker.C:
			cfg, err = loadAutoConfig(paths.Config)
			if err != nil {
				logger.Printf("load config failed: %v", err)
				continue
			}
			if len(cfg.Accounts) == 0 {
				continue
			}
			if err := syncAndCheckin(ctx, &cfg, logger); err != nil {
				logger.Printf("sync loop failed: %v", err)
			}
			if err := saveAutoConfig(paths.Config, cfg); err != nil {
				logger.Printf("save config failed: %v", err)
			}
		}
	}
}

func syncAndCheckin(ctx context.Context, cfg *autoConfig, logger *appLogger) error {
	now := time.Now()
	var firstErr error
	for i := range cfg.Accounts {
		account := &cfg.Accounts[i]
		if !account.Enabled {
			continue
		}
		status, _, err := loadAccountStatus(ctx, account.SessionFile)
		if err != nil {
			account.LastError = err.Error()
			account.LastSyncedAt = now.Format(time.RFC3339)
			if firstErr == nil {
				firstErr = err
			}
			logger.Printf("status fetch failed for %s: %v", fallbackText(account.UserID, account.SessionFile), err)
			continue
		}
		applyStatusSnapshot(account, status, now)
		account.LastError = ""
		if status.SignedToday {
			wasSigned := account.LastCheckinStatus == "already_signed" && account.LastSignedDate == status.LastSignedDate
			account.LastCheckinStatus = "already_signed"
			if !wasSigned {
				logger.Printf("%s already signed today", fallbackText(account.Nickname, account.UserID))
			}
			continue
		}
		if now.In(time.FixedZone("CST", 8*3600)).Hour() < 1 {
			if account.LastCheckinStatus != "pending" {
				logger.Printf("%s pending until 01:00", fallbackText(account.Nickname, account.UserID))
			}
			account.LastCheckinStatus = "pending"
			continue
		}
		if err := performDailyCheckin(ctx, account, logger); err != nil {
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func performDailyCheckin(ctx context.Context, account *accountConfig, logger *appLogger) error {
	client, sess, err := newCheckinClient(ctx, account.SessionFile)
	if err != nil {
		account.LastError = err.Error()
		account.LastCheckinStatus = "error"
		return err
	}
	label := fallbackText(account.Nickname, account.UserID)

	listenResp, err := client.YouthListenSong(ctx, sdk.YouthListenSongRequest{Cookie: sess.Cookie})
	if err != nil {
		account.LastError = err.Error()
		account.LastCheckinStatus = "error"
		logger.Printf("%s listen_song failed: %v", label, err)
		return err
	}
	listenCode := asInt(listenResp.Body["error_code"])
	if responseOK(listenResp.Body) {
		logger.Printf("%s listen_song success", label)
	} else if listenCode == 130012 {
		logger.Printf("%s already claimed via listen_song today", label)
		return finalizeCheckinStatus(ctx, account, "already_signed", "", logger)
	} else {
		err = fmt.Errorf("listen_song failed: %s", strings.TrimSpace(asString(listenResp.RawBody)))
		account.LastError = err.Error()
		account.LastCheckinStatus = "error"
		logger.Printf("%s listen_song failed: %v", label, err)
		return err
	}

	for attempt := 1; attempt <= 8; attempt++ {
		vipResp, vipErr := client.YouthVip(ctx, sdk.YouthVipRequest{Cookie: client.Cookie()})
		if vipErr != nil {
			account.LastError = vipErr.Error()
			account.LastCheckinStatus = "error"
			logger.Printf("%s youth_vip failed at attempt %d: %v", label, attempt, vipErr)
			return vipErr
		}
		vipCode := asInt(vipResp.Body["error_code"])
		switch {
		case responseOK(vipResp.Body):
			logger.Printf("%s youth_vip success attempt %d/8", label, attempt)
		case vipCode == 30002:
			logger.Printf("%s youth_vip exhausted for today", label)
			return finalizeCheckinStatus(ctx, account, "success", "", logger)
		default:
			err = fmt.Errorf("youth_vip failed: %s", strings.TrimSpace(asString(vipResp.RawBody)))
			account.LastError = err.Error()
			account.LastCheckinStatus = "error"
			logger.Printf("%s youth_vip failed at attempt %d: %v", label, attempt, err)
			return err
		}
		if attempt < 8 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(vipAttemptInterval()):
			}
		}
	}
	return finalizeCheckinStatus(ctx, account, "success", "", logger)
}

func newCheckinClient(ctx context.Context, sessionPath string) (*sdk.Client, session.Config, error) {
	status, sess, err := loadAccountStatus(ctx, sessionPath)
	if err != nil {
		return nil, sess, err
	}
	client, err := sdk.New(sdk.WithCookie(sess.Cookie))
	if err != nil {
		return nil, sess, err
	}
	if status.UserID != "" {
		sess.LastUserID = status.UserID
	}
	return client, sess, nil
}

func finalizeCheckinStatus(ctx context.Context, account *accountConfig, result, errText string, logger *appLogger) error {
	status, _, err := loadAccountStatus(ctx, account.SessionFile)
	if err != nil {
		account.LastError = err.Error()
		account.LastCheckinStatus = "error"
		return err
	}
	applyStatusSnapshot(account, status, time.Now())
	account.LastCheckinStatus = result
	account.LastCheckinAt = time.Now().Format(time.RFC3339)
	account.LastError = strings.TrimSpace(errText)
	logger.Printf("%s final status: signed_today=%t, vip_end=%s", fallbackText(account.Nickname, account.UserID), status.SignedToday, fallbackText(status.VIPExpireAt, "unknown"))
	return nil
}

func vipAttemptInterval() time.Duration {
	if raw := strings.TrimSpace(strings.ToLower(os.Getenv("KUGOU_AUTO_CHECKIN_VIP_INTERVAL"))); raw != "" {
		if d, err := time.ParseDuration(raw); err == nil && d > 0 {
			return d
		}
	}
	return 30 * time.Second
}
