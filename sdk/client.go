package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/lfhy/kugou-music-api/core/kugou"
	"github.com/lfhy/kugou-music-api/core/util"
)

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
	core       *kugou.Client
	isLite     bool
	cookiePool map[string]string
	guid       string
	serverDev  string
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

func New(opts ...Option) (*Client, error) {
	platform := strings.ToLower(strings.TrimSpace(os.Getenv("platform")))
	c := &Client{
		isLite:     defaultLitePlatform(platform),
		cookiePool: map[string]string{},
		guid:       util.MD5Hex(util.RandomString(16)),
		serverDev:  strings.ToUpper(util.RandomString(10)),
	}
	for _, opt := range opts {
		opt(c)
	}

	c.core = kugou.NewClient(c.isLite)
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

func (c *Client) Call(ctx context.Context, route string, req Request) (*Response, error) {
	sp, ok := apiSpecMap[route]
	if !ok {
		return nil, errors.New("route not found: " + route)
	}
	return c.callSpec(ctx, sp, req)
}

func (c *Client) CallByIdentifier(ctx context.Context, identifier string, req Request) (*Response, error) {
	for route, sp := range apiSpecMap {
		if sp.Identifier == identifier {
			return c.callSpec(ctx, apiSpecMap[route], req)
		}
	}
	return nil, errors.New("identifier not found: " + identifier)
}

func (c *Client) callSpec(ctx context.Context, sp apiSpec, req Request) (*Response, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}

	cfg := kugou.RequestConfig{
		Method:             firstNonEmpty(req.Method, sp.Method),
		URL:                fallbackURL(firstNonEmpty(req.URL, sp.URL)),
		BaseURL:            firstNonEmpty(req.BaseURL, sp.BaseURL),
		Headers:            mergeHeaders(sp.Headers, req.Headers),
		EncryptType:        firstNonEmpty(req.EncryptType, sp.EncryptType),
		Cookie:             cookies,
		EncryptKey:         chooseBool(req.EncryptKey, sp.EncryptKey),
		ClearDefaultParams: chooseBool(req.ClearDefaultParams, sp.ClearDefaultParams),
		NotSignature:       chooseBool(req.NotSignature, sp.NotSignature),
	}

	args := map[string]any{}
	mergeAny(args, req.Params)
	mergeAny(args, req.Data)

	switch {
	case sp.UseParams && sp.UseData:
		cfg.Params = args
		if len(req.Data) > 0 {
			cfg.Data = req.Data
		} else {
			cfg.Data = args
		}
	case sp.UseData:
		if len(req.Data) > 0 {
			cfg.Data = req.Data
		} else {
			cfg.Data = args
		}
	default:
		cfg.Params = args
	}

	raw, err := c.core.CreateRequest(ctx, cfg)
	if len(raw.Cookie) > 0 {
		c.updateCookiePool(raw.Cookie)
	}

	out := &Response{Status: raw.Status, RawBody: raw.Body, Headers: raw.Headers, Cookie: raw.Cookie}
	var body map[string]any
	if json.Unmarshal(raw.Body, &body) == nil {
		out.Body = body
	}
	return out, err
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

func mergeAny(dst, src map[string]any) {
	for k, v := range src {
		dst[k] = v
	}
}

func mergeHeaders(base, extra map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(extra))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}

func firstNonEmpty(v ...string) string {
	for _, x := range v {
		if strings.TrimSpace(x) != "" {
			return x
		}
	}
	return ""
}

func fallbackURL(u string) string {
	if strings.Contains(u, "${") || strings.TrimSpace(u) == "" {
		return "/"
	}
	return u
}

func chooseBool(input *bool, def bool) bool {
	if input == nil {
		return def
	}
	return *input
}
