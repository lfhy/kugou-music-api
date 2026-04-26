package sdk

import "context"

// Collection and theme endpoints share the generic compat request path.
func (c *Client) SheetCollection(ctx context.Context, req SheetCollectionRequest) (*SheetCollectionResponse, error) {
	resp, err := compatCall(ctx, c, RouteSheetCollection, "sheet_collection", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := SheetCollectionResponse(*resp)
	return &out, nil
}

func (c *Client) SheetCollectionDetail(ctx context.Context, req SheetCollectionDetailRequest) (*SheetCollectionDetailResponse, error) {
	resp, err := compatCall(ctx, c, RouteSheetCollectionDetail, "sheet_collection_detail", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := SheetCollectionDetailResponse(*resp)
	return &out, nil
}

func (c *Client) SheetDetail(ctx context.Context, req SheetDetailRequest) (*SheetDetailResponse, error) {
	resp, err := compatCall(ctx, c, RouteSheetDetail, "sheet_detail", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := SheetDetailResponse(*resp)
	return &out, nil
}

func (c *Client) SheetList(ctx context.Context, req SheetListRequest) (*SheetListResponse, error) {
	resp, err := compatCall(ctx, c, RouteSheetList, "sheet_list", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := SheetListResponse(*resp)
	return &out, nil
}

func (c *Client) ThemeMusic(ctx context.Context, req ThemeMusicRequest) (*ThemeMusicResponse, error) {
	resp, err := compatCall(ctx, c, RouteThemeMusic, "theme_music", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := ThemeMusicResponse(*resp)
	return &out, nil
}

func (c *Client) ThemeMusicDetail(ctx context.Context, req ThemeMusicDetailRequest) (*ThemeMusicDetailResponse, error) {
	resp, err := compatCall(ctx, c, RouteThemeMusicDetail, "theme_music_detail", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := ThemeMusicDetailResponse(*resp)
	return &out, nil
}

func (c *Client) ThemePlaylist(ctx context.Context, req ThemePlaylistRequest) (*ThemePlaylistResponse, error) {
	resp, err := compatCall(ctx, c, RouteThemePlaylist, "theme_playlist", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := ThemePlaylistResponse(*resp)
	return &out, nil
}

func (c *Client) ThemePlaylistTrack(ctx context.Context, req ThemePlaylistTrackRequest) (*ThemePlaylistTrackResponse, error) {
	resp, err := compatCall(ctx, c, RouteThemePlaylistTrack, "theme_playlist_track", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := ThemePlaylistTrackResponse(*resp)
	return &out, nil
}

func (c *Client) TopCardYouth(ctx context.Context, req TopCardYouthRequest) (*TopCardYouthResponse, error) {
	resp, err := compatCall(ctx, c, RouteTopCardYouth, "top_card_youth", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := TopCardYouthResponse(*resp)
	return &out, nil
}

func (c *Client) TopIp(ctx context.Context, req TopIpRequest) (*TopIpResponse, error) {
	resp, err := compatCall(ctx, c, RouteTopIp, "top_ip", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := TopIpResponse(*resp)
	return &out, nil
}

func (c *Client) TopSong(ctx context.Context, req TopSongRequest) (*TopSongResponse, error) {
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
	if merged, ok := c.ensureLoginValid(ctx, cookies); ok {
		cookies = merged
	} else {
		delete(cookies, "token")
		delete(cookies, "userid")
	}

	rankID := toInt(firstAny(params["type"], params["rank_id"]), 21608)
	if rankID == 0 {
		rankID = 21608
	}
	userid := firstNonEmpty(asIntString(firstAny(params["userid"], nil)), cookies["userid"], "0")

	resp, err := c.Call(ctx, RouteTopSong, Request{
		Method: "POST",
		URL:    "/musicadservice/container/v1/newsong_publish",
		Data: map[string]any{
			"rank_id":  rankID,
			"userid":   userid,
			"page":     toInt(firstAny(params["page"], nil), 1),
			"pagesize": toInt(firstAny(params["pagesize"], nil), 30),
			"tags":     []any{},
		},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := TopSongResponse(*resp)
	return &out, nil
}

func (c *Client) UserHistory(ctx context.Context, req UserHistoryRequest) (*UserHistoryResponse, error) {
	resp, err := compatCall(ctx, c, RouteUserHistory, "user_history", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := UserHistoryResponse(*resp)
	return &out, nil
}

func (c *Client) VideoUrl(ctx context.Context, req VideoUrlRequest) (*VideoUrlResponse, error) {
	resp, err := compatCall(ctx, c, RouteVideoUrl, "video_url", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := VideoUrlResponse(*resp)
	return &out, nil
}
