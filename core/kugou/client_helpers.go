package kugou

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Transport helpers handle URL encoding, headers, and per-request timeouts.
func (c *Client) requestHTTPClient(timeout time.Duration) *http.Client {
	if timeout <= 0 {
		return c.httpClient
	}
	client := *c.httpClient
	client.Timeout = timeout
	return &client
}

func buildURL(baseURL, endpoint string, params map[string]any) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	ep, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	u = u.ResolveReference(ep)

	q := u.Query()
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		q.Set(k, stringify(params[k]))
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func mergeAny(dst, src map[string]any) {
	for k, v := range src {
		dst[k] = v
	}
}

func mergeHeader(h http.Header, m map[string]string) {
	for k, v := range m {
		h.Set(k, v)
	}
}

func buildCookieHeader(cookie map[string]string) string {
	if len(cookie) == 0 {
		return ""
	}
	keys := make([]string, 0, len(cookie))
	for k, v := range cookie {
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+cookie[k])
	}
	return strings.Join(parts, "; ")
}

func stringify(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(t, 10)
	case json.Number:
		return t.String()
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

func getMapOr(m map[string]string, k, d string) string {
	if v, ok := m[k]; ok && v != "" {
		return v
	}
	return d
}
