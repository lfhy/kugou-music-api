package sdk

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// IpZone normalizes the embedded special_link payload into an integer ip_id.
func (c *Client) IpZone(ctx context.Context, req IpZoneRequest) (*IpZoneResponse, error) {
	resp, err := c.Call(ctx, RouteIpZone, Request{
		Method:      "GET",
		URL:         "/v1/zone/index",
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "yuekucategory.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		if status, _ := resp.Body["status"].(float64); int(status) == 1 {
			if dm, ok := resp.Body["data"].(map[string]any); ok {
				if list, ok := dm["list"].([]any); ok {
					for i, item := range list {
						im, ok := item.(map[string]any)
						if !ok {
							continue
						}
						rawLink := strings.TrimSpace(fmt.Sprintf("%v", im["special_link"]))
						if rawLink == "" || rawLink == "<nil>" {
							continue
						}
						linkQS, err := url.ParseQuery(rawLink)
						if err != nil {
							continue
						}
						pathRaw := strings.TrimSpace(linkQS.Get("path"))
						if pathRaw == "" {
							continue
						}
						pathQS, err := url.ParseQuery(pathRaw)
						if err != nil {
							continue
						}
						ipID := strings.TrimSpace(pathQS.Get("ip_id"))
						if ipID == "" {
							continue
						}
						if n, err := strconv.Atoi(ipID); err == nil {
							im["ip_id"] = n
						}
						list[i] = im
					}
					dm["list"] = list
					resp.Body["data"] = dm
				}
			}
		}
	}
	out := IpZoneResponse(*resp)
	return &out, nil
}
