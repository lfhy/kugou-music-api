package kugou

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rushingrain/kugou-music-api/core/config"
	"github.com/rushingrain/kugou-music-api/core/util"
)

type RequestConfig struct {
	Method             string
	URL                string
	BaseURL            string
	Params             map[string]any
	Data               any
	Headers            map[string]string
	EncryptType        string
	Cookie             map[string]string
	EncryptKey         bool
	ClearDefaultParams bool
	NotSignature       bool
	IP                 string
	Timeout            time.Duration
}

type Response struct {
	Status  int
	Body    []byte
	Cookie  []string
	Headers map[string]string
}

type Client struct {
	httpClient *http.Client
	isLite     bool
}

func NewClient(isLite bool) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 25 * time.Second},
		isLite:     isLite,
	}
}

func (c *Client) CreateRequest(ctx context.Context, cfg RequestConfig) (Response, error) {
	appid, clientver := config.PlatformConfig(c.isLite)

	if cfg.Cookie == nil {
		cfg.Cookie = map[string]string{}
	}

	dfid := getMapOr(cfg.Cookie, "dfid", "-")
	mid := getMapOr(cfg.Cookie, "KUGOU_API_MID", "")
	uuid := "-"
	token := getMapOr(cfg.Cookie, "token", "")
	userid := getMapOr(cfg.Cookie, "userid", "0")
	clienttime := strconv.FormatInt(time.Now().Unix(), 10)

	headers := map[string]string{
		"dfid":       dfid,
		"clienttime": clienttime,
		"mid":        mid,
		"kg-rc":      "1",
		"kg-thash":   "5d816a0",
		"kg-rec":     "1",
		"kg-rf":      "B9EDA08A64250DEFFBCADDEE00F8F25F",
	}
	if cfg.IP != "" {
		headers["X-Real-IP"] = cfg.IP
		headers["X-Forwarded-For"] = cfg.IP
	}

	defaultParams := map[string]any{
		"dfid":       dfid,
		"mid":        mid,
		"uuid":       uuid,
		"appid":      appid,
		"clientver":  clientver,
		"clienttime": clienttime,
	}
	if token != "" {
		defaultParams["token"] = token
	}
	if userid != "" && userid != "0" {
		defaultParams["userid"] = userid
	}

	params := map[string]any{}
	if !cfg.ClearDefaultParams {
		mergeAny(params, defaultParams)
	}
	mergeAny(params, cfg.Params)

	if cfg.EncryptKey {
		hash := fmt.Sprintf("%v", params["hash"])
		params["key"] = SignKey(hash, mid, userid, appid, c.isLite)
	}

	dataString := ""
	bodyReader := io.Reader(nil)
	if cfg.Data != nil {
		var raw []byte
		switch v := cfg.Data.(type) {
		case []byte:
			raw = v
		case string:
			raw = []byte(v)
		default:
			b, _ := json.Marshal(v)
			raw = b
		}
		dataString = string(raw)
		bodyReader = bytes.NewReader(raw)
	}

	if _, ok := params["signature"]; !ok && !cfg.NotSignature {
		switch strings.ToLower(cfg.EncryptType) {
		case "register":
			params["signature"] = SignatureRegisterParams(params)
		case "web":
			params["signature"] = SignatureWebParams(params)
		default:
			params["signature"] = SignatureAndroidParams(params, dataString, c.isLite)
		}
	}

	method := strings.ToUpper(strings.TrimSpace(cfg.Method))
	if method == "" {
		method = http.MethodGet
	}
	if cfg.URL == "" {
		cfg.URL = "/"
	}
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://gateway.kugou.com"
	}

	fullURL, err := buildURL(baseURL, cfg.URL, params)
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return Response{}, err
	}

	req.Header.Set("User-Agent", "Android15-1070-11083-46-0-DiscoveryDRADProtocol-wifi")
	if cfg.Data != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if cookieHeader := buildCookieHeader(cfg.Cookie); cookieHeader != "" {
		req.Header.Set("Cookie", cookieHeader)
	}

	mergeHeader(req.Header, cfg.Headers)
	mergeHeader(req.Header, headers)

	resp, err := c.requestHTTPClient(cfg.Timeout).Do(req)
	if err != nil {
		return Response{Status: 502, Body: []byte(`{"status":0,"msg":"request failed"}`)}, err
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)
	out := Response{Status: 200, Body: rawBody, Cookie: []string{}, Headers: map[string]string{}}

	if s := resp.Header.Get("ssa-code"); s != "" {
		out.Headers["ssa-code"] = s
	}

	for _, c := range resp.Header.Values("Set-Cookie") {
		out.Cookie = append(out.Cookie, util.ParseCookieString(c))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		out.Status = 502
		return out, fmt.Errorf("upstream status %d", resp.StatusCode)
	}

	return out, nil
}
