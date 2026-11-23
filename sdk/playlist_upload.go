package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rushingrain/kugou-music-api/core/kugou"
	"github.com/rushingrain/kugou-music-api/core/util"
)

// PlaylistUpload sends a local image or raw bytes to the playlist cover service.
// A non-nil Data slice takes precedence over File, including an empty slice.
func (c *Client) PlaylistUpload(ctx context.Context, req PlaylistUploadRequest) (*PlaylistUploadResponse, error) {
	image := req.Data
	if image == nil {
		var err error
		image, err = os.ReadFile(req.File)
		if err != nil {
			out := playlistUploadResult(0, map[string]any{"status": 0, "msg": "文件不存在"})
			return out, fmt.Errorf("playlist_upload: %w", err)
		}
	}
	md5 := req.Md5
	if md5 == "" {
		md5 = util.MD5Hex(time.Now().Format("20060102") + "hewry678WEK23D")
	}
	imageType, extension := req.Type, req.ExtendName
	if imageType == "" {
		imageType = "custom"
	}
	if extension == "" {
		extension = ".jpg"
	}
	raw, err := c.doRequest(ctx, kugou.RequestConfig{
		Method: "POST", BaseURL: "http://imgphp.kugou.com", URL: "/imageupload/stream.php",
		Params: map[string]any{
			"type": imageType, "extendName": extension, "md5": md5, "jsonResponse": 1,
		},
		Data: image, Headers: map[string]string{"Content-Type": "application/octet-stream"},
		ClearDefaultParams: true, NotSignature: true, Timeout: 30 * time.Second,
	})
	if err != nil {
		return playlistUploadResult(502, map[string]any{"status": 0, "msg": err.Error()}), err
	}
	var body map[string]any
	if err := json.Unmarshal(raw.Body, &body); err != nil {
		return playlistUploadResult(502, map[string]any{"status": 0, "msg": err.Error()}), err
	}
	if playlistUploadTruthy(body["IsSuccess"]) {
		result := map[string]any{"status": 1}
		if fileName, ok := body["FileName"]; ok {
			result["FileName"] = fileName
		}
		return playlistUploadResult(200, result), nil
	}
	message := body["Message"]
	if !playlistUploadTruthy(message) {
		message = "上传失败"
	}
	return playlistUploadResult(502, map[string]any{"status": 0, "msg": message}), nil
}

func playlistUploadResult(status int, body map[string]any) *PlaylistUploadResponse {
	raw, _ := json.Marshal(body)
	return &PlaylistUploadResponse{Status: status, Body: body, RawBody: raw}
}

// JSON truthiness differs from Go zero values: for example, "0" is truthy in JS.
func playlistUploadTruthy(v any) bool {
	switch value := v.(type) {
	case nil:
		return false
	case bool:
		return value
	case float64:
		return value != 0
	case string:
		return value != ""
	default:
		return true
	}
}
