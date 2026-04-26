package main

import (
	"context"
	"fmt"
	"math"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
	"github.com/lfhy/kugou-music-api/sdk"
)

type accountStatus struct {
	UserID         string
	Nickname       string
	SessionFile    string
	SessionUpdated string
	LastSignedDate string
	VIPExpireAt    string
	SignedToday    bool
}

func loadAccountStatus(ctx context.Context, sessionPath string) (accountStatus, session.Config, error) {
	sessionPath = filepath.Clean(sessionPath)
	sess := session.Load(sessionPath)
	if !session.HasLoginCookie(sess.Cookie) {
		return accountStatus{}, sess, fmt.Errorf("session missing token/userid: %s", sessionPath)
	}

	client, err := sdk.New(sdk.WithCookie(sess.Cookie))
	if err != nil {
		return accountStatus{}, sess, err
	}
	userResp, err := ensureUserDetail(ctx, client, sess.Cookie)
	if err != nil {
		return accountStatus{}, sess, err
	}
	vipResp, err := client.UserVipDetail(ctx, sdk.UserVipDetailRequest{})
	if err != nil {
		return accountStatus{}, sess, err
	}
	monthResp, err := client.YouthMonthVipRecord(ctx, sdk.YouthMonthVipRecordRequest{})
	if err != nil {
		return accountStatus{}, sess, err
	}

	sess.Cookie = client.Cookie()
	sess.LastUserID = normalizedUserID(firstString(
		pickMap(vipResp.Body, "data")["userid"],
		sess.Cookie["userid"],
		sess.LastUserID,
	))
	if err := session.Save(sessionPath, sess); err != nil {
		return accountStatus{}, sess, err
	}

	status := accountStatus{
		UserID:         firstString(pickMap(vipResp.Body, "data")["userid"], sess.LastUserID, sess.Cookie["userid"]),
		Nickname:       strings.TrimSpace(firstString(pickMap(userResp.Body, "data")["nickname"])),
		SessionFile:    sessionPath,
		SessionUpdated: sess.UpdatedAt,
		LastSignedDate: latestSignedDate(monthResp.Body),
		VIPExpireAt:    latestVIPEndTime(vipResp.Body),
	}
	status.UserID = normalizedUserID(status.UserID)
	status.SignedToday = status.LastSignedDate == chinaToday(time.Now())
	return status, sess, nil
}

func ensureUserDetail(ctx context.Context, client *sdk.Client, cookie map[string]string) (*sdk.UserDetailResponse, error) {
	resp, err := client.UserDetail(ctx, sdk.UserDetailRequest{Cookie: cookie})
	if err == nil && responseOK(resp.Body) {
		return resp, nil
	}
	_, refreshErr := client.LoginByToken(ctx, sdk.TokenLoginRequest{
		Token:  strings.TrimSpace(cookie["token"]),
		UserID: normalizedUserID(cookie["userid"]),
		Cookie: cookie,
	})
	if refreshErr != nil {
		if err != nil {
			return nil, fmt.Errorf("user detail failed: %v; token refresh failed: %v", err, refreshErr)
		}
		return nil, fmt.Errorf("token refresh failed: %v", refreshErr)
	}
	resp, err = client.UserDetail(ctx, sdk.UserDetailRequest{Cookie: client.Cookie()})
	if err != nil || !responseOK(resp.Body) {
		return nil, fmt.Errorf("user detail still failed after refresh: %v", err)
	}
	return resp, nil
}

func applyStatusSnapshot(account *accountConfig, status accountStatus, now time.Time) {
	account.UserID = status.UserID
	account.Nickname = status.Nickname
	account.SessionFile = status.SessionFile
	account.LastSignedDate = status.LastSignedDate
	account.VIPExpireAt = status.VIPExpireAt
	account.LastSyncedAt = now.Format(time.RFC3339)
	if status.SignedToday && account.LastCheckinStatus == "" {
		account.LastCheckinStatus = "already_signed"
	}
	if account.Enabled == false {
		account.Enabled = true
	}
}

func responseOK(body map[string]any) bool {
	return asInt(body["status"]) == 1 && asInt(body["error_code"]) == 0
}

func latestSignedDate(body map[string]any) string {
	data := pickMap(body, "data")
	list, _ := data["list"].([]any)
	days := make([]string, 0, len(list))
	for _, item := range list {
		m, _ := item.(map[string]any)
		if asInt(m["receive_vip"]) != 1 {
			continue
		}
		day := strings.TrimSpace(asString(m["day"]))
		if day != "" {
			days = append(days, day)
		}
	}
	if len(days) == 0 {
		return ""
	}
	sort.Strings(days)
	return days[len(days)-1]
}

func latestVIPEndTime(body map[string]any) string {
	data := pickMap(body, "data")
	list, _ := data["busi_vip"].([]any)
	latest := ""
	for _, item := range list {
		m, _ := item.(map[string]any)
		end := strings.TrimSpace(asString(m["vip_end_time"]))
		if end > latest {
			latest = end
		}
	}
	return latest
}

func displayLastCheckin(account accountConfig) string {
	if strings.TrimSpace(account.LastCheckinAt) != "" {
		return account.LastCheckinAt
	}
	if strings.TrimSpace(account.LastSignedDate) != "" {
		return account.LastSignedDate + " (接口仅返回日期)"
	}
	return "未知"
}

func chinaToday(now time.Time) string {
	return now.In(time.FixedZone("CST", 8*3600)).Format("2006-01-02")
}

func normalizedUserID(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if !strings.ContainsAny(v, ".eE") {
		return v
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return v
	}
	return strconv.FormatInt(int64(math.Round(f)), 10)
}

func pickMap(body map[string]any, key string) map[string]any {
	out, _ := body[key].(map[string]any)
	if out == nil {
		return map[string]any{}
	}
	return out
}

func firstString(values ...any) string {
	for _, value := range values {
		text := strings.TrimSpace(asString(value))
		if text != "" {
			return text
		}
	}
	return ""
}

func asString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	default:
		return fmt.Sprintf("%v", t)
	}
}

func asInt(v any) int {
	switch t := v.(type) {
	case nil:
		return 0
	case int:
		return t
	case int64:
		return int(t)
	case float64:
		return int(t)
	case string:
		i, _ := strconv.Atoi(strings.TrimSpace(t))
		return i
	default:
		i, _ := strconv.Atoi(strings.TrimSpace(fmt.Sprintf("%v", t)))
		return i
	}
}
