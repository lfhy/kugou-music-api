// Package kugou exposes the SDK from the module root so callers can import
// github.com/lfhy/kugou-music-api directly.
package kugou

import "github.com/lfhy/kugou-music-api/sdk"

func New(opts ...Option) (*Client, error) {
	return sdk.New(opts...)
}

func NewClient(opts ...Option) (*Client, error) {
	return sdk.New(opts...)
}

func WithLite(v bool) Option {
	return sdk.WithLite(v)
}

func WithCookie(cookie map[string]string) Option {
	return sdk.WithCookie(cookie)
}

func WithSongURLFallback(enabled bool) SongPlayURLOption {
	return sdk.WithSongURLFallback(enabled)
}

func WithSongURLAll(enabled bool) SongPlayURLOption {
	return sdk.WithSongURLAll(enabled)
}

func WithSongURLDFID(dfid string) SongPlayURLOption {
	return sdk.WithSongURLDFID(dfid)
}

func Endpoints() []APIInfo {
	out := make([]APIInfo, len(sdk.APIList))
	copy(out, sdk.APIList)
	return out
}

func RouteByIdentifier(identifier string) (string, bool) {
	for _, x := range sdk.APIList {
		if x.Identifier == identifier {
			return x.Route, true
		}
	}
	return "", false
}
