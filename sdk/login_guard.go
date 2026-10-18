package sdk

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	errAutoRefreshDisabled = errors.New("auto refresh disabled")
	errLoginExpired        = errors.New("login expired")
)

var authCookieKeys = []string{"token", "userid", "vip_token", "vip_type"}

func (c *Client) ensureLoginValid(ctx context.Context, cookie map[string]string) (map[string]string, bool) {
	merged := c.mergeRequestCookies(cookie, false)
	if err := requireLoginCookie(merged); err != nil {
		return merged, false
	}

	resp, err := c.userDetailOnce(ctx, UserDetailRequest{Cookie: merged})
	if err == nil && isBizSuccessCode(resp) {
		// 返回实际用于校验的 cookie（client 池 + 本次请求传入），而不是 client 池。
		// 与 JS 版一致：调用方只通过请求参数传登录态时，token / userid 不能在
		// 校验后被丢掉，否则后续请求会退化成未登录态。
		return merged, true
	}
	// Callers that disable automatic refresh must not trigger a refresh
	// indirectly through an authenticated helper method. Returning the
	// current cookie preserves the existing best-effort validation behavior;
	// the caller will receive the original request result or auth error.
	if !c.autoRefresh {
		return merged, false
	}
	if _, refreshErr := c.refreshLoginSession(ctx, merged); refreshErr != nil {
		return c.Cookie(), false
	}
	return c.Cookie(), true
}

func (c *Client) tryAutoRefresh(ctx context.Context, cookie map[string]string) error {
	if !c.autoRefresh {
		return fmt.Errorf("%w: disabled by client option", errAutoRefreshDisabled)
	}
	_, err := c.refreshLoginSession(ctx, cookie)
	return err
}

func (c *Client) refreshLoginSession(ctx context.Context, cookie map[string]string) (map[string]string, error) {
	merged := c.mergeRequestCookies(cookie, false)
	if err := requireLoginCookie(merged); err != nil {
		return merged, err
	}

	refreshResp, refreshErr := c.LoginByToken(ctx, TokenLoginRequest{
		Token:  strings.TrimSpace(merged["token"]),
		UserID: strings.TrimSpace(merged["userid"]),
		Cookie: merged,
	})
	if refreshErr != nil {
		return c.Cookie(), fmt.Errorf("%w: token refresh request failed: %w", errLoginExpired, refreshErr)
	}
	if !isBizSuccessCode(refreshResp) {
		return c.Cookie(), fmt.Errorf("%w: token refresh rejected", errLoginExpired)
	}

	validated, err := c.userDetailOnce(ctx, UserDetailRequest{Cookie: c.Cookie()})
	if err != nil || !isBizSuccessCode(validated) {
		return c.Cookie(), fmt.Errorf("%w: refreshed session validation failed", errLoginExpired)
	}
	return c.Cookie(), nil
}

func (c *Client) loginStateError(cookie map[string]string) error {
	if err := requireLoginCookie(cookie); err != nil {
		return err
	}
	if !c.autoRefresh {
		return fmt.Errorf("%w: disabled by client option", errAutoRefreshDisabled)
	}
	return fmt.Errorf("%w: token refresh failed or session cannot be recovered", errLoginExpired)
}

func requireLoginCookie(cookie map[string]string) error {
	token := strings.TrimSpace(cookie["token"])
	userid := strings.TrimSpace(cookie["userid"])
	if token == "" || userid == "" || userid == "0" {
		return fmt.Errorf("login required: missing token/userid")
	}
	return nil
}

func isBizSuccessCode(resp *Response) bool {
	if resp == nil || resp.Body == nil {
		return false
	}
	status := strings.TrimSpace(fmt.Sprintf("%v", resp.Body["status"]))
	errorCode := strings.TrimSpace(fmt.Sprintf("%v", resp.Body["error_code"]))
	if status != "1" && status != "1.0" {
		return false
	}
	return errorCode == "" || errorCode == "<nil>" || errorCode == "0" || errorCode == "0.0"
}

func isAuthExpiredResponse(resp *Response) bool {
	if resp == nil || resp.Body == nil {
		return false
	}
	errorCode := strings.TrimSpace(fmt.Sprintf("%v", resp.Body["error_code"]))
	return errorCode == "20018" || errorCode == "20018.0"
}

func isAuthCookieKey(key string) bool {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "token", "userid", "vip_token", "vip_type":
		return true
	default:
		return false
	}
}

func allowAutoRefreshRoute(route string) bool {
	switch route {
	case RouteLogin, RouteLoginCellphone, RouteLoginDevice, RouteLoginOpenplat,
		RouteLoginQrCheck, RouteLoginQrCreate, RouteLoginQrKey, RouteLoginToken,
		RouteLoginWxCheck, RouteLoginWxCreate:
		return false
	default:
		return true
	}
}

func wrapLoginResponseError(route string, resp *Response, cause error) error {
	summary := summarizeRawBody(nil)
	if resp != nil {
		summary = summarizeRawBody(resp.RawBody)
	}
	if summary == "" {
		return fmt.Errorf("%s: %w", route, cause)
	}
	return fmt.Errorf("%s: %w: %s", route, cause, summary)
}

func summarizeRawBody(raw []byte) string {
	text := strings.TrimSpace(string(raw))
	if text == "" {
		return ""
	}
	if len(text) > 180 {
		return text[:180] + "..."
	}
	return text
}
