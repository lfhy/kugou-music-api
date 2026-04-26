package sdk

import (
	"context"
	"strings"
)

// normalizeOptionalLogin keeps guest-capable endpoints working when saved login cookies expire.
func normalizeOptionalLogin(ctx context.Context, c *Client, cookie map[string]string) map[string]string {
	token := strings.TrimSpace(cookie["token"])
	userid := strings.TrimSpace(cookie["userid"])
	if token == "" || userid == "" || userid == "0" {
		return cookie
	}
	if merged, ok := c.ensureLoginValid(ctx, cookie); ok {
		return merged
	}
	delete(cookie, "token")
	delete(cookie, "userid")
	delete(cookie, "vip_token")
	delete(cookie, "vip_type")
	return cookie
}
