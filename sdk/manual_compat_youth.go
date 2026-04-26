package sdk

import (
	"context"
	"time"
)

// Youth and yueku compat endpoints stay split out to keep each file small.
func (c *Client) YouthChannelSong(ctx context.Context, req YouthChannelSongRequest) (*YouthChannelSongResponse, error) {
	resp, err := compatCall(ctx, c, RouteYouthChannelSong, "youth_channel_song", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := YouthChannelSongResponse(*resp)
	return &out, nil
}

func (c *Client) YouthChannelSongDetail(ctx context.Context, req YouthChannelSongDetailRequest) (*YouthChannelSongDetailResponse, error) {
	resp, err := compatCall(ctx, c, RouteYouthChannelSongDetail, "youth_channel_song_detail", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := YouthChannelSongDetailResponse(*resp)
	return &out, nil
}

func (c *Client) YouthDayVipUpgrade(ctx context.Context, req YouthDayVipUpgradeRequest) (*YouthDayVipUpgradeResponse, error) {
	resp, err := compatCall(ctx, c, RouteYouthDayVipUpgrade, "youth_day_vip_upgrade", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := YouthDayVipUpgradeResponse(*resp)
	return &out, nil
}

func (c *Client) YouthListenSong(ctx context.Context, req YouthListenSongRequest) (*YouthListenSongResponse, error) {
	resp, err := compatCall(ctx, c, RouteYouthListenSong, "youth_listen_song", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := YouthListenSongResponse(*resp)
	return &out, nil
}

func (c *Client) YouthUnionVip(ctx context.Context, req YouthUnionVipRequest) (*YouthUnionVipResponse, error) {
	resp, err := compatCall(ctx, c, RouteYouthUnionVip, "youth_union_vip", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := YouthUnionVipResponse(*resp)
	return &out, nil
}

func (c *Client) YouthUserSong(ctx context.Context, req YouthUserSongRequest) (*YouthUserSongResponse, error) {
	resp, err := compatCall(ctx, c, RouteYouthUserSong, "youth_user_song", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := YouthUserSongResponse(*resp)
	return &out, nil
}

func (c *Client) YouthVip(ctx context.Context, req YouthVipRequest) (*YouthVipResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, requireLoginCookie(cookies)
	}

	nowMS := time.Now().UnixMilli()
	resp, err := c.Call(ctx, RouteYouthVip, Request{
		Method: "POST",
		URL:    "/youth/v1/ad/play_report",
		Data: map[string]any{
			"ad_id":      12307537187,
			"play_end":   nowMS,
			"play_start": nowMS - 30000,
		},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := YouthVipResponse(*resp)
	return &out, nil
}

func (c *Client) YuekuBanner(ctx context.Context, req YuekuBannerRequest) (*YuekuBannerResponse, error) {
	resp, err := compatCall(ctx, c, RouteYuekuBanner, "yueku_banner", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := YuekuBannerResponse(*resp)
	return &out, nil
}
