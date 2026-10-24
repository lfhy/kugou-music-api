package sdk

import (
	"context"
	"fmt"
)

// PlaylistDetail mirrors the JS client: ids are converted to the request body's
// data array, with auth values sourced from explicit fields or cookies.
func (c *Client) PlaylistDetail(ctx context.Context, req PlaylistDetailRequest) (*PlaylistDetailResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}

	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}

	data := make([]map[string]any, 0)
	ids := firstNonEmpty(fmt.Sprintf("%v", firstAny(params["ids"], nil)), "")
	for _, id := range splitCSV(ids) {
		data = append(data, map[string]any{"global_collection_id": id})
	}

	resp, err := c.Call(ctx, RoutePlaylistDetail, Request{
		Data: map[string]any{
			"data":   data,
			"userid": firstAny(params["userid"], firstAny(cookies["userid"], 0)),
			"token":  firstAny(params["token"], firstAny(cookies["token"], "")),
		},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := PlaylistDetailResponse(*resp)
	return &out, nil
}
