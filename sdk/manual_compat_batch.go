package sdk

import (
	"context"
	"fmt"

	"github.com/lfhy/kugou-music-api/core/config"
)

func compatRequest(req any, identifier string, cookie map[string]string, extra map[string]any) (map[string]any, map[string]string) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	if compat, ok := buildCompatParams(identifier, params, cookie); ok {
		params = compat
	}
	for k, v := range extra {
		params[k] = v
	}
	return params, applyCompatCookie(identifier, cookie)
}

func compatCall(ctx context.Context, c *Client, route, identifier string, req any, cookie map[string]string, extra map[string]any) (*Response, error) {
	params, mergedCookie := compatRequest(req, identifier, cookie, extra)
	resp, err := c.Call(ctx, route, Request{Params: params, Cookie: mergedCookie})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func compatCallRequireLogin(ctx context.Context, c *Client, route, identifier string, req any, cookie map[string]string, extra map[string]any) (*Response, error) {
	cookies := c.Cookie()
	for k, v := range cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, requireLoginCookie(cookies)
	}
	params, mergedCookie := compatRequest(req, identifier, cookies, extra)
	resp, err := c.Call(ctx, route, Request{Params: params, Cookie: mergedCookie})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) AlbumSongs(ctx context.Context, req AlbumSongsRequest) (*AlbumSongsResponse, error) {
	resp, err := compatCall(ctx, c, RouteAlbumSongs, "album_songs", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := AlbumSongsResponse(*resp)
	return &out, nil
}

func (c *Client) ArtistAlbums(ctx context.Context, req ArtistAlbumsRequest) (*ArtistAlbumsResponse, error) {
	resp, err := compatCall(ctx, c, RouteArtistAlbums, "artist_albums", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := ArtistAlbumsResponse(*resp)
	return &out, nil
}

func (c *Client) ArtistLists(ctx context.Context, req ArtistListsRequest) (*ArtistListsResponse, error) {
	resp, err := compatCall(ctx, c, RouteArtistLists, "artist_lists", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := ArtistListsResponse(*resp)
	return &out, nil
}

func (c *Client) ArtistVideos(ctx context.Context, req ArtistVideosRequest) (*ArtistVideosResponse, error) {
	resp, err := compatCall(ctx, c, RouteArtistVideos, "artist_videos", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := ArtistVideosResponse(*resp)
	return &out, nil
}

func (c *Client) CommentMusicClassify(ctx context.Context, req CommentMusicClassifyRequest) (*CommentMusicClassifyResponse, error) {
	resp, err := compatCall(ctx, c, RouteCommentMusicClassify, "comment_music_classify", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := CommentMusicClassifyResponse(*resp)
	return &out, nil
}

func (c *Client) LastestSongsListen(ctx context.Context, req LastestSongsListenRequest) (*LastestSongsListenResponse, error) {
	resp, err := compatCallRequireLogin(ctx, c, RouteLastestSongsListen, "lastest_songs_listen", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := LastestSongsListenResponse(*resp)
	return &out, nil
}

func (c *Client) Lyric(ctx context.Context, req LyricRequest) (*LyricResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteLyric, Request{
		Method:  "GET",
		BaseURL: "https://lyrics.kugou.com",
		URL:     "/download",
		Params: map[string]any{
			"ver":       1,
			"client":    firstNonEmpty(fmt.Sprintf("%v", firstAny(params["client"], req.Client)), "android"),
			"id":        firstAny(params["id"], req.Id),
			"accesskey": firstAny(params["accesskey"], req.Accesskey),
			"fmt":       firstNonEmpty(fmt.Sprintf("%v", firstAny(params["fmt"], req.Fmt)), "krc"),
			"charset":   "utf8",
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := LyricResponse(*resp)
	if toBool(firstAny(params["decode"], req.Decode), false) {
		_ = out.DecodedContent()
	}
	return &out, nil
}

func (c *Client) PlayhistoryUpload(ctx context.Context, req PlayhistoryUploadRequest) (*PlayhistoryUploadResponse, error) {
	resp, err := compatCallRequireLogin(ctx, c, RoutePlayhistoryUpload, "playhistory_upload", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := PlayhistoryUploadResponse(*resp)
	return &out, nil
}

func (c *Client) PlaylistTrackAll(ctx context.Context, req PlaylistTrackAllRequest) (*PlaylistTrackAllResponse, error) {
	resp, err := compatCall(ctx, c, RoutePlaylistTrackAll, "playlist_track_all", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := PlaylistTrackAllResponse(*resp)
	return &out, nil
}

func (c *Client) PlaylistTrackAllNew(ctx context.Context, req PlaylistTrackAllNewRequest) (*PlaylistTrackAllNewResponse, error) {
	resp, err := compatCall(ctx, c, RoutePlaylistTrackAllNew, "playlist_track_all_new", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := PlaylistTrackAllNewResponse(*resp)
	return &out, nil
}

func (c *Client) PrivilegeLite(ctx context.Context, req PrivilegeLiteRequest) (*PrivilegeLiteResponse, error) {
	resp, err := compatCall(ctx, c, RoutePrivilegeLite, "privilege_lite", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := PrivilegeLiteResponse(*resp)
	return &out, nil
}

func (c *Client) RankAudio(ctx context.Context, req RankAudioRequest) (*RankAudioResponse, error) {
	resp, err := compatCall(ctx, c, RouteRankAudio, "rank_audio", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := RankAudioResponse(*resp)
	return &out, nil
}

func (c *Client) SearchComplex(ctx context.Context, req SearchComplexRequest) (*SearchComplexResponse, error) {
	resp, err := compatCall(ctx, c, RouteSearchComplex, "search_complex", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := SearchComplexResponse(*resp)
	return &out, nil
}

func (c *Client) SearchDefault(ctx context.Context, req SearchDefaultRequest) (*SearchDefaultResponse, error) {
	resp, err := compatCall(ctx, c, RouteSearchDefault, "search_default", req, req.Cookie, req.Extra)
	if err != nil {
		return nil, err
	}
	out := SearchDefaultResponse(*resp)
	return &out, nil
}

func (c *Client) SearchLyric(ctx context.Context, req SearchLyricRequest) (*SearchLyricResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	appid, clientver := config.PlatformConfig(c.isLite)
	clear := true
	noSign := true
	resp, err := c.Call(ctx, RouteSearchLyric, Request{
		Method:  "GET",
		BaseURL: "https://lyrics.kugou.com",
		URL:     "/v1/search",
		Params: map[string]any{
			"album_audio_id": toInt(firstAny(params["album_audio_id"], req.AlbumAudioId), 0),
			"appid":          appid,
			"clientver":      clientver,
			"duration":       0,
			"hash":           firstNonEmpty(fmt.Sprintf("%v", firstAny(params["hash"], req.Hash)), ""),
			"keyword":        firstNonEmpty(fmt.Sprintf("%v", firstAny(params["keywords"], req.Keywords)), ""),
			"lrctxt":         1,
			"man":            firstNonEmpty(fmt.Sprintf("%v", firstAny(params["man"], req.Man)), "no"),
		},
		Cookie:             req.Cookie,
		EncryptType:        "android",
		ClearDefaultParams: &clear,
		NotSignature:       &noSign,
	})
	if err != nil {
		return nil, err
	}
	out := SearchLyricResponse(*resp)
	return &out, nil
}

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
	// Optional login: if present and valid, use it; otherwise keep guest.
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
	resp, err := compatCallRequireLogin(ctx, c, RouteUserHistory, "user_history", req, req.Cookie, req.Extra)
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
	resp, err := compatCall(ctx, c, RouteYouthVip, "youth_vip", req, req.Cookie, req.Extra)
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
