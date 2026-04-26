package sdk

import (
	"context"
	"fmt"
	"strings"
)

// Album, artist, social, and zone endpoints with hand-tuned request shaping.
func (c *Client) AlbumDetail(ctx context.Context, req AlbumDetailRequest) (*AlbumDetailResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteAlbumDetail, Request{
		Method: "POST",
		URL:    "/kmr/v2/albums",
		Data: map[string]any{
			"data": []map[string]any{{
				"album_id": firstAny(params["id"], params["album_id"]),
			}},
			"is_buy": toInt(firstAny(params["is_buy"], nil), 0),
			"fields": "album_id,album_name,publish_date,sizable_cover,intro,language,is_publish,heat,type,quality,authors,exclusive,author_name,trans_param",
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "openapi.kugou.com", "kg-tid": "255"},
	})
	if err != nil {
		return nil, err
	}
	out := AlbumDetailResponse(*resp)
	return &out, nil
}

func (c *Client) AlbumShop(ctx context.Context, req AlbumShopRequest) (*AlbumShopResponse, error) {
	resp, err := c.Call(ctx, RouteAlbumShop, Request{
		Method:      "GET",
		URL:         "/zhuanjidata/v3/album_shop_v2/get_classify_data",
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := AlbumShopResponse(*resp)
	return &out, nil
}

func (c *Client) ArtistDetail(ctx context.Context, req ArtistDetailRequest) (*ArtistDetailResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteArtistDetail, Request{
		Method: "POST",
		URL:    "/kmr/v3/author",
		Data: map[string]any{
			"author_id": firstAny(params["id"], params["author_id"]),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "openapi.kugou.com", "kg-tid": "36"},
	})
	if err != nil {
		return nil, err
	}
	out := ArtistDetailResponse(*resp)
	return &out, nil
}

func (c *Client) ArtistFollowNewsongs(ctx context.Context, req ArtistFollowNewsongsRequest) (*ArtistFollowNewsongsResponse, error) {
	cookies := c.Cookie()
	for k, v := range req.Cookie {
		cookies[k] = v
	}
	var ok bool
	cookies, ok = c.ensureLoginValid(ctx, cookies)
	if !ok {
		return nil, c.loginStateError(cookies)
	}
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	lastAlbumID := toInt(firstAny(params["last_album_id"], nil), 0)
	resp, err := c.Call(ctx, RouteArtistFollowNewsongs, Request{
		Method: "POST",
		URL:    "/feed/v1/follow/newsong_album_list",
		Params: map[string]any{
			"last_album_id": lastAlbumID,
			"page_size":     toInt(firstAny(params["pagesize"], params["page_size"]), 30),
			"opt_sort":      ternaryInt(toInt(firstAny(params["opt_sort"], nil), 1) == 2, 2, 1),
		},
		Data: map[string]any{
			"last_album_id": lastAlbumID,
		},
		Cookie:      cookies,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := ArtistFollowNewsongsResponse(*resp)
	return &out, nil
}

func (c *Client) ArtistHonour(ctx context.Context, req ArtistHonourRequest) (*ArtistHonourResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteArtistHonour, Request{
		Method:  "POST",
		BaseURL: "http://h5activity.kugou.com",
		URL:     "/v1/query_singer_honour_detail",
		Params: map[string]any{
			"singer_id": firstAny(params["id"], params["singer_id"]),
			"pagesize":  toInt(firstAny(params["pagesize"], nil), 30),
			"page":      toInt(firstAny(params["page"], nil), 1),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := ArtistHonourResponse(*resp)
	return &out, nil
}

func (c *Client) EverydayFriend(ctx context.Context, req EverydayFriendRequest) (*EverydayFriendResponse, error) {
	resp, err := c.Call(ctx, RouteEverydayFriend, Request{
		Method:  "POST",
		BaseURL: "https://acsing.service.kugou.com",
		URL:     "/sing7/relation/json/v3/friend_rec_by_using_song_list",
		Data: map[string]any{
			"list": []any{
				map[string]any{
					"user_id": 853927886,
					"mixsong_ids": []any{
						290083753, 251724346, 571554587, 250126644, 208831644, 40328518, 250504076, 581706850, 318347675, 585258401,
						288481998, 407414475, 28239430, 280584633, 291957521, 64556644, 243149863, 488725103, 32114153, 39951172,
						29019580, 40397606, 327507651, 32029382, 32218359, 340353127, 276448762, 177071956, 100031397, 249251602,
					},
				},
			},
		},
		Params:      map[string]any{"channel": 130, "isteen": 0, "platform": 2, "usemkv": 1},
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"pid": "126556797"},
	})
	if err != nil {
		return nil, err
	}
	out := EverydayFriendResponse(*resp)
	return &out, nil
}

func (c *Client) EverydayHistory(ctx context.Context, req EverydayHistoryRequest) (*EverydayHistoryResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	paramsMap := map[string]any{
		"mode":     firstNonEmpty(fmt.Sprintf("%v", firstAny(params["mode"], nil)), "list"),
		"platform": firstNonEmpty(fmt.Sprintf("%v", firstAny(params["platform"], nil)), "ios"),
	}
	if v := strings.TrimSpace(fmt.Sprintf("%v", firstAny(params["history_name"], nil))); v != "" && v != "<nil>" {
		paramsMap["history_name"] = v
	}
	if v := strings.TrimSpace(fmt.Sprintf("%v", firstAny(params["date"], nil))); v != "" && v != "<nil>" {
		paramsMap["date"] = v
	}
	resp, err := c.Call(ctx, RouteEverydayHistory, Request{
		Method:      "POST",
		URL:         "/everyday/api/v1/get_history",
		Params:      paramsMap,
		Cookie:      req.Cookie,
		EncryptType: "android",
		Headers:     map[string]string{"x-router": "everydayrec.service.kugou.com"},
	})
	if err != nil {
		return nil, err
	}
	out := EverydayHistoryResponse(*resp)
	return &out, nil
}

func (c *Client) EverydayStyleRecommend(ctx context.Context, req EverydayStyleRecommendRequest) (*EverydayStyleRecommendResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteEverydayStyleRecommend, Request{
		Method: "POST",
		URL:    "/everydayrec.service/everyday_style_recommend",
		Params: map[string]any{
			"tagids": firstNonEmpty(fmt.Sprintf("%v", firstAny(params["tagids"], nil)), ""),
		},
		Data:        map[string]any{},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := EverydayStyleRecommendResponse(*resp)
	return &out, nil
}

func (c *Client) FavoriteCount(ctx context.Context, req FavoriteCountRequest) (*FavoriteCountResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	for k, v := range req.Extra {
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteFavoriteCount, Request{
		Method: "GET",
		URL:    "/count/v1/audio/mget_collect",
		Params: map[string]any{
			"mixsongids": firstAny(params["mixsongids"], nil),
		},
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := FavoriteCountResponse(*resp)
	return &out, nil
}
