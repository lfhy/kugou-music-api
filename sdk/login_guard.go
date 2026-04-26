package sdk

import (
	"context"
	"fmt"
	"strings"
)

func (c *Client) ensureLoginValid(ctx context.Context, cookie map[string]string) (map[string]string, bool) {
	merged := c.Cookie()
	for k, v := range cookie {
		merged[k] = v
	}
	token := strings.TrimSpace(merged["token"])
	userid := strings.TrimSpace(merged["userid"])
	if token == "" || userid == "" || userid == "0" {
		return merged, false
	}

	resp, err := c.UserDetail(ctx, UserDetailRequest{Cookie: merged})
	if err == nil && isBizSuccessCode(resp) {
		return c.Cookie(), true
	}

	refreshResp, refreshErr := c.LoginByToken(ctx, TokenLoginRequest{
		Token:  token,
		UserID: userid,
		Cookie: merged,
	})
	if refreshErr != nil || !isBizSuccessCode(refreshResp) {
		return c.Cookie(), false
	}

	resp, err = c.UserDetail(ctx, UserDetailRequest{Cookie: c.Cookie()})
	if err != nil || !isBizSuccessCode(resp) {
		return c.Cookie(), false
	}
	return c.Cookie(), true
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

