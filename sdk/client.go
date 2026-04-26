package sdk

import (
	"context"
	"os"
	"strings"

	"github.com/lfhy/kugou-music-api/core/kugou"
	"github.com/lfhy/kugou-music-api/core/util"
)

type requestExecutor func(context.Context, kugou.RequestConfig) (kugou.Response, error)

type Option func(*Client)

type Request struct {
	Params map[string]any
	Data   map[string]any
	Cookie map[string]string

	Method      string
	URL         string
	BaseURL     string
	Headers     map[string]string
	EncryptType string

	EncryptKey         *bool
	ClearDefaultParams *bool
	NotSignature       *bool
}

type Response struct {
	Status  int
	RawBody []byte
	Body    map[string]any
	Headers map[string]string
	Cookie  []string
}

type Client struct {
	core        *kugou.Client
	requester   requestExecutor
	isLite      bool
	autoRefresh bool
	cookiePool  map[string]string
	guid        string
	serverDev   string
}

func WithLite(v bool) Option {
	return func(c *Client) { c.isLite = v }
}

func WithCookie(cookie map[string]string) Option {
	return func(c *Client) {
		for k, v := range cookie {
			c.cookiePool[k] = v
		}
	}
}

// WithAutoRefresh controls whether the SDK should try token refresh on 20018-style auth failures.
func WithAutoRefresh(enabled bool) Option {
	return func(c *Client) { c.autoRefresh = enabled }
}

func New(opts ...Option) (*Client, error) {
	platform := strings.ToLower(strings.TrimSpace(os.Getenv("platform")))
	c := &Client{
		isLite:      defaultLitePlatform(platform),
		autoRefresh: true,
		cookiePool:  map[string]string{},
		guid:        util.MD5Hex(util.RandomString(16)),
		serverDev:   strings.ToUpper(util.RandomString(10)),
	}
	for _, opt := range opts {
		opt(c)
	}

	c.core = kugou.NewClient(c.isLite)
	c.requester = c.core.CreateRequest
	c.injectPlatformCookies()
	return c, nil
}

func (c *Client) Endpoints() []APIInfo {
	out := make([]APIInfo, len(APIList))
	copy(out, APIList)
	return out
}

func (c *Client) RouteByIdentifier(identifier string) (string, bool) {
	for _, x := range APIList {
		if x.Identifier == identifier {
			return x.Route, true
		}
	}
	return "", false
}

func (c *Client) SetCookie(k, v string) { c.cookiePool[k] = v }

func (c *Client) Cookie() map[string]string {
	out := make(map[string]string, len(c.cookiePool))
	for k, v := range c.cookiePool {
		out[k] = v
	}
	return out
}

func (c *Client) doRequest(ctx context.Context, cfg kugou.RequestConfig) (kugou.Response, error) {
	if c.requester != nil {
		return c.requester(ctx, cfg)
	}
	return c.core.CreateRequest(ctx, cfg)
}

func (c *Client) updateCookiePool(setCookies []string) {
	for _, sc := range setCookies {
		seg := strings.Split(strings.TrimSpace(sc), ";")
		if len(seg) == 0 {
			continue
		}
		kv := strings.SplitN(seg[0], "=", 2)
		if len(kv) != 2 {
			continue
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		if k != "" {
			c.cookiePool[k] = v
		}
	}
}

func (c *Client) injectPlatformCookies() {
	guid := firstNonEmpty(strings.TrimSpace(os.Getenv("KUGOU_API_GUID")), c.guid)
	mid := util.CalculateMid(guid)
	dev := strings.ToUpper(firstNonEmpty(strings.TrimSpace(os.Getenv("KUGOU_API_DEV")), c.serverDev))
	mac := strings.ToUpper(firstNonEmpty(strings.TrimSpace(os.Getenv("KUGOU_API_MAC")), "02:00:00:00:00:00"))
	platform := strings.TrimSpace(os.Getenv("platform"))
	if platform == "" {
		platform = "lite"
	}

	c.cookiePool["KUGOU_API_PLATFORM"] = platform
	c.cookiePool["KUGOU_API_MID"] = mid
	c.cookiePool["KUGOU_API_GUID"] = guid
	c.cookiePool["KUGOU_API_DEV"] = dev
	c.cookiePool["KUGOU_API_MAC"] = mac
}

func defaultLitePlatform(platform string) bool {
	switch strings.ToLower(strings.TrimSpace(platform)) {
	case "", "lite":
		return true
	case "normal", "official", "standard", "std":
		return false
	default:
		return true
	}
}
