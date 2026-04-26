package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
	"github.com/lfhy/kugou-music-api/sdk"
)

type smokeCase struct {
	Name       string
	NeedLogin  bool
	RiskyWrite bool
	Run        func(context.Context, *sdk.Client, map[string]string) (*sdk.Response, error)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cfg := session.Load(session.DefaultPath())
	client, err := sdk.New(sdk.WithCookie(cfg.Cookie))
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		os.Exit(1)
	}

	hasLogin := session.HasLoginCookie(client.Cookie())
	includeRisky := strings.TrimSpace(os.Getenv("SMOKE_INCLUDE_RISKY")) == "1"
	fmt.Printf("pending smoke start: login=%v include_risky=%v\n", hasLogin, includeRisky)

	cases := pendingCases()
	total := 0
	run := 0
	success := 0
	failed := 0
	skipped := 0

	for _, tc := range cases {
		total++
		if tc.RiskyWrite && !includeRisky {
			skipped++
			fmt.Printf("[SKIP] %s (risky write)\n", tc.Name)
			continue
		}
		if tc.NeedLogin && !hasLogin {
			skipped++
			fmt.Printf("[SKIP] %s (login required)\n", tc.Name)
			continue
		}

		run++
		resp, err := tc.Run(ctx, client, client.Cookie())
		if err != nil {
			failed++
			fmt.Printf("[FAIL] %s err=%v\n", tc.Name, err)
			continue
		}
		code := pick(resp, "error_code")
		st := pick(resp, "status")
		if st == "1" || st == "1.0" {
			success++
			fmt.Printf("[OK]   %s status=%s error_code=%s\n", tc.Name, st, code)
		} else {
			failed++
			fmt.Printf("[FAIL] %s status=%s error_code=%s\n", tc.Name, st, code)
		}
	}

	fmt.Println("---")
	fmt.Printf("summary: total=%d run=%d success=%d fail=%d skip=%d\n", total, run, success, failed, skipped)
	if failed > 0 {
		os.Exit(1)
	}
}

func pick(resp *sdk.Response, key string) string {
	if resp == nil || resp.Body == nil {
		return ""
	}
	v, ok := resp.Body[key]
	if !ok {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

func pendingCases() []smokeCase {
	return []smokeCase{
		{Name: "album_detail", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.AlbumDetail(ctx, sdk.AlbumDetailRequest{Cookie: cookie, Extra: map[string]any{"id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "album_shop", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.AlbumShop(ctx, sdk.AlbumShopRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "artist_detail", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.ArtistDetail(ctx, sdk.ArtistDetailRequest{Cookie: cookie, Extra: map[string]any{"id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "artist_follow_newsongs", NeedLogin: true, Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.ArtistFollowNewsongs(ctx, sdk.ArtistFollowNewsongsRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "artist_honour", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.ArtistHonour(ctx, sdk.ArtistHonourRequest{Cookie: cookie, Extra: map[string]any{"id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "everyday_friend", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.EverydayFriend(ctx, sdk.EverydayFriendRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "everyday_history", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.EverydayHistory(ctx, sdk.EverydayHistoryRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "everyday_style_recommend", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.EverydayStyleRecommend(ctx, sdk.EverydayStyleRecommendRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "favorite_count", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.FavoriteCount(ctx, sdk.FavoriteCountRequest{Cookie: cookie, Extra: map[string]any{"mixsongids": "1"}})
			return (*sdk.Response)(r), e
		}},
		{Name: "ip_zone", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.IpZone(ctx, sdk.IpZoneRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "kmr_audio_mv", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.KmrAudioMv(ctx, sdk.KmrAudioMvRequest{Cookie: cookie, Extra: map[string]any{"album_audio_id": "1"}})
			return (*sdk.Response)(r), e
		}},
		{Name: "krm_audio", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.KrmAudio(ctx, sdk.KrmAudioRequest{Cookie: cookie, Extra: map[string]any{"album_audio_id": "1"}})
			return (*sdk.Response)(r), e
		}},
		{Name: "longaudio_album_audios", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.LongaudioAlbumAudios(ctx, sdk.LongaudioAlbumAudiosRequest{Cookie: cookie, Extra: map[string]any{"album_id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "longaudio_album_detail", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.LongaudioAlbumDetail(ctx, sdk.LongaudioAlbumDetailRequest{Cookie: cookie, Extra: map[string]any{"album_id": "1"}})
			return (*sdk.Response)(r), e
		}},
		{Name: "longaudio_daily_recommend", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.LongaudioDailyRecommend(ctx, sdk.LongaudioDailyRecommendRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "longaudio_rank_recommend", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.LongaudioRankRecommend(ctx, sdk.LongaudioRankRecommendRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "longaudio_vip_recommend", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.LongaudioVipRecommend(ctx, sdk.LongaudioVipRecommendRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "longaudio_week_recommend", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.LongaudioWeekRecommend(ctx, sdk.LongaudioWeekRecommendRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "rank_info", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.RankInfo(ctx, sdk.RankInfoRequest{Cookie: cookie, Extra: map[string]any{"rankid": 666}})
			return (*sdk.Response)(r), e
		}},
		{Name: "rank_list", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.RankList(ctx, sdk.RankListRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "rank_top", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.RankTop(ctx, sdk.RankTopRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "rank_vol", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.RankVol(ctx, sdk.RankVolRequest{Cookie: cookie, Extra: map[string]any{"rankid": 666}})
			return (*sdk.Response)(r), e
		}},
		{Name: "scene_audio_list", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SceneAudioList(ctx, sdk.SceneAudioListRequest{Cookie: cookie, Extra: map[string]any{"id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "scene_collection_list", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SceneCollectionList(ctx, sdk.SceneCollectionListRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "scene_lists", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SceneLists(ctx, sdk.SceneListsRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "scene_lists_v2", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SceneListsV2(ctx, sdk.SceneListsV2Request{Cookie: cookie, Extra: map[string]any{"id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "scene_module", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SceneModule(ctx, sdk.SceneModuleRequest{Cookie: cookie, Extra: map[string]any{"id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "scene_module_info", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SceneModuleInfo(ctx, sdk.SceneModuleInfoRequest{Cookie: cookie, Extra: map[string]any{"id": 1, "module_id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "scene_music", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SceneMusic(ctx, sdk.SceneMusicRequest{Cookie: cookie, Extra: map[string]any{"id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "scene_video_list", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SceneVideoList(ctx, sdk.SceneVideoListRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "search_suggest", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SearchSuggest(ctx, sdk.SearchSuggestRequest{Cookie: cookie, Extra: map[string]any{"keywords": "周杰伦"}})
			return (*sdk.Response)(r), e
		}},
		{Name: "server_now", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.ServerNow(ctx, sdk.ServerNowRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "sheet_hot", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SheetHot(ctx, sdk.SheetHotRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "singer_list", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SingerList(ctx, sdk.SingerListRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "song_climax", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SongClimax(ctx, sdk.SongClimaxRequest{Cookie: cookie, Extra: map[string]any{"hash": "00000000000000000000000000000000"}})
			return (*sdk.Response)(r), e
		}},
		{Name: "song_ranking", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SongRanking(ctx, sdk.SongRankingRequest{Cookie: cookie, Extra: map[string]any{"album_audio_id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "song_ranking_filter", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.SongRankingFilter(ctx, sdk.SongRankingFilterRequest{Cookie: cookie, Extra: map[string]any{"album_audio_id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "user_vip_detail", NeedLogin: true, Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.UserVipDetail(ctx, sdk.UserVipDetailRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "youth_channel_all", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.YouthChannelAll(ctx, sdk.YouthChannelAllRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "youth_channel_amway", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.YouthChannelAmway(ctx, sdk.YouthChannelAmwayRequest{Cookie: cookie, Extra: map[string]any{"global_collection_id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "youth_channel_detail", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.YouthChannelDetail(ctx, sdk.YouthChannelDetailRequest{Cookie: cookie, Extra: map[string]any{"global_collection_id": "1"}})
			return (*sdk.Response)(r), e
		}},
		{Name: "youth_channel_similar", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.YouthChannelSimilar(ctx, sdk.YouthChannelSimilarRequest{Cookie: cookie, Extra: map[string]any{"channel_id": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "youth_channel_sub", NeedLogin: true, RiskyWrite: true, Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.YouthChannelSub(ctx, sdk.YouthChannelSubRequest{Cookie: cookie, Extra: map[string]any{"global_collection_id": 1, "t": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "youth_day_vip", NeedLogin: true, RiskyWrite: true, Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.YouthDayVip(ctx, sdk.YouthDayVipRequest{Cookie: cookie, Extra: map[string]any{"receive_day": 1}})
			return (*sdk.Response)(r), e
		}},
		{Name: "youth_dynamic", NeedLogin: true, Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.YouthDynamic(ctx, sdk.YouthDynamicRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "youth_dynamic_recent", NeedLogin: true, Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.YouthDynamicRecent(ctx, sdk.YouthDynamicRecentRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "youth_month_vip_record", NeedLogin: true, Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.YouthMonthVipRecord(ctx, sdk.YouthMonthVipRecordRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "yueku", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.Yueku(ctx, sdk.YuekuRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
		{Name: "yueku_fm", Run: func(ctx context.Context, c *sdk.Client, cookie map[string]string) (*sdk.Response, error) {
			r, e := c.YuekuFm(ctx, sdk.YuekuFmRequest{Cookie: cookie})
			return (*sdk.Response)(r), e
		}},
	}
}
