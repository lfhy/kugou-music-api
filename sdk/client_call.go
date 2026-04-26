package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/lfhy/kugou-music-api/core/kugou"
)

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
	return c.callSpecWithRetry(ctx, sp, req, false)
}

func (c *Client) callSpecWithRetry(ctx context.Context, sp apiSpec, req Request, retried bool) (*Response, error) {
	cookies := c.mergeRequestCookies(req.Cookie, retried)
	resp, err := c.executeJSONRequest(ctx, c.buildRequestConfig(sp, req, cookies))
	if !isAuthExpiredResponse(resp) || !allowAutoRefreshRoute(sp.Route) {
		return resp, err
	}
	if retried {
		return resp, wrapLoginResponseError(sp.Route, resp, c.loginStateError(cookies))
	}
	if refreshErr := c.tryAutoRefresh(ctx, cookies); refreshErr != nil {
		return resp, wrapLoginResponseError(sp.Route, resp, refreshErr)
	}
	return c.callSpecWithRetry(ctx, sp, c.prepareRetryRequest(req, c.Cookie()), true)
}

func (c *Client) buildRequestConfig(sp apiSpec, req Request, cookies map[string]string) kugou.RequestConfig {
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
	return cfg
}

func (c *Client) executeJSONRequest(ctx context.Context, cfg kugou.RequestConfig) (*Response, error) {
	raw, err := c.doRequest(ctx, cfg)
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

func (c *Client) mergeRequestCookies(extra map[string]string, preferClientAuth bool) map[string]string {
	out := c.Cookie()
	for k, v := range extra {
		if preferClientAuth && isAuthCookieKey(k) && strings.TrimSpace(out[k]) != "" {
			continue
		}
		out[k] = v
	}
	return out
}

func (c *Client) prepareRetryRequest(req Request, refreshed map[string]string) Request {
	retry := req
	retry.Cookie = mergeRetryCookies(req.Cookie, refreshed)
	retry.Params = cloneAnyMap(req.Params)
	retry.Data = cloneAnyMap(req.Data)
	applyRetryAuthFields(retry.Params, refreshed)
	applyRetryAuthFields(retry.Data, refreshed)
	return retry
}

func mergeRetryCookies(requestCookie, refreshed map[string]string) map[string]string {
	out := cloneStringMap(requestCookie)
	for _, key := range authCookieKeys {
		if val := strings.TrimSpace(refreshed[key]); val != "" {
			out[key] = val
		}
	}
	return out
}

func applyRetryAuthFields(dst map[string]any, refreshed map[string]string) {
	if len(dst) == 0 {
		return
	}
	for _, key := range authCookieKeys {
		if _, ok := dst[key]; ok {
			dst[key] = refreshed[key]
		}
	}
}

func cloneAnyMap(src map[string]any) map[string]any {
	if len(src) == 0 {
		return nil
	}
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}

func cloneStringMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return map[string]string{}
	}
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
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
