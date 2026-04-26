# 酷狗音乐Go SDK

中文 | [English](README_EN.md)

> 注意：本项目为使用 Codex 基于源项目 [MakcRe/KuGouMusicApi](https://github.com/MakcRe/KuGouMusicApi)（JavaScript 版本）迁移和更新的 Go 版本。由于迁移实现与接口行为可能存在偏差，当前版本可能存在问题，请谨慎使用并自行评估风险。

## 快速开始

安装：

```bash
go get github.com/lfhy/kugou-music-api
```

推荐直接从根包导入：

```go
import kg "github.com/lfhy/kugou-music-api"
```

最小示例：

```go
package main

import (
	"context"
	"fmt"
	"log"

	kg "github.com/lfhy/kugou-music-api"
)

func main() {
	client, err := kg.New()
	if err != nil {
		log.Fatalf("init client failed: %v", err)
	}

	resp, err := client.Search(context.Background(), kg.SearchRequest{
		Keywords: "周杰伦",
		Page:     1,
		Pagesize: 10,
	})
	if err != nil {
		log.Fatalf("search failed: %v", err)
	}

	fmt.Printf("status: %d\n", resp.Status)
	fmt.Printf("body: %s\n", string(resp.RawBody))
}
```

兼容说明：
- 推荐新项目直接使用根包 `github.com/lfhy/kugou-music-api`
- 旧代码仍可继续使用子包 `github.com/lfhy/kugou-music-api/sdk`

接口文档入口：
- 技术清单（路由/方法/模型）：`sdk/API_CATALOG.md`
- 中文说明清单（自动提取注释）：`sdk/API_CATALOG_ZH.md`
- 兼容性审计清单（逐接口风险）：`sdk/API_COMPAT_AUDIT.md`
- 修复状态清单（逐接口校对进度）：`sdk/API_FIX_STATUS.md`
- 签名/加密检查清单（逐接口）：`sdk/API_SIGN_CHECK.md`

说明：
- 中文说明来自原项目 `module/*.js` 注释自动提取；没有注释的会显示“(暂无注释)”。
- 调用函数清单仍在本 README 底部（`(*Client).Xxx(...)`）。

---

# API Catalog

| Identifier | Route | Method | Request Model | Response Model |
| --- | --- | --- | --- | --- |
| `ai_recommend` | `/ai/recommend` | `POST` | `AiRecommendRequest` | `AiRecommendResponse` |
| `album` | `/album` | `POST` | `AlbumRequest` | `AlbumResponse` |
| `album_detail` | `/album/detail` | `POST` | `AlbumDetailRequest` | `AlbumDetailResponse` |
| `album_shop` | `/album/shop` | `GET` | `AlbumShopRequest` | `AlbumShopResponse` |
| `album_songs` | `/album/songs` | `POST` | `AlbumSongsRequest` | `AlbumSongsResponse` |
| `artist_albums` | `/artist/albums` | `POST` | `ArtistAlbumsRequest` | `ArtistAlbumsResponse` |
| `artist_audios` | `/artist/audios` | `POST` | `ArtistAudiosRequest` | `ArtistAudiosResponse` |
| `artist_detail` | `/artist/detail` | `POST` | `ArtistDetailRequest` | `ArtistDetailResponse` |
| `artist_follow` | `/artist/follow` | `POST` | `ArtistFollowRequest` | `ArtistFollowResponse` |
| `artist_follow_newsongs` | `/artist/follow/newsongs` | `POST` | `ArtistFollowNewsongsRequest` | `ArtistFollowNewsongsResponse` |
| `artist_honour` | `/artist/honour` | `POST` | `ArtistHonourRequest` | `ArtistHonourResponse` |
| `artist_lists` | `/artist/lists` | `GET` | `ArtistListsRequest` | `ArtistListsResponse` |
| `artist_unfollow` | `/artist/unfollow` | `POST` | `ArtistUnfollowRequest` | `ArtistUnfollowResponse` |
| `artist_videos` | `/artist/videos` | `GET` | `ArtistVideosRequest` | `ArtistVideosResponse` |
| `audio` | `/audio` | `POST` | `AudioRequest` | `AudioResponse` |
| `audio_accompany_matching` | `/audio/accompany/matching` | `GET` | `AudioAccompanyMatchingRequest` | `AudioAccompanyMatchingResponse` |
| `audio_ktv_total` | `/audio/ktv/total` | `GET` | `AudioKtvTotalRequest` | `AudioKtvTotalResponse` |
| `audio_related` | `/audio/related` | `GET` | `AudioRelatedRequest` | `AudioRelatedResponse` |
| `brush` | `/brush` | `POST` | `BrushRequest` | `BrushResponse` |
| `captcha_sent` | `/captcha/sent` | `POST` | `CaptchaSentRequest` | `CaptchaSentResponse` |
| `comment_album` | `/comment/album` | `POST` | `CommentAlbumRequest` | `CommentAlbumResponse` |
| `comment_count` | `/comment/count` | `GET` | `CommentCountRequest` | `CommentCountResponse` |
| `comment_floor` | `/comment/floor` | `POST` | `CommentFloorRequest` | `CommentFloorResponse` |
| `comment_music` | `/comment/music` | `POST` | `CommentMusicRequest` | `CommentMusicResponse` |
| `comment_music_classify` | `/comment/music/classify` | `POST` | `CommentMusicClassifyRequest` | `CommentMusicClassifyResponse` |
| `comment_music_hotword` | `/comment/music/hotword` | `POST` | `CommentMusicHotwordRequest` | `CommentMusicHotwordResponse` |
| `comment_playlist` | `/comment/playlist` | `POST` | `CommentPlaylistRequest` | `CommentPlaylistResponse` |
| `everyday_friend` | `/everyday/friend` | `POST` | `EverydayFriendRequest` | `EverydayFriendResponse` |
| `everyday_history` | `/everyday/history` | `POST` | `EverydayHistoryRequest` | `EverydayHistoryResponse` |
| `everyday_recommend` | `/everyday/recommend` | `POST` | `EverydayRecommendRequest` | `EverydayRecommendResponse` |
| `everyday_style_recommend` | `/everyday/style/recommend` | `POST` | `EverydayStyleRecommendRequest` | `EverydayStyleRecommendResponse` |
| `favorite_count` | `/favorite/count` | `GET` | `FavoriteCountRequest` | `FavoriteCountResponse` |
| `fm_class` | `/fm/class` | `POST` | `FmClassRequest` | `FmClassResponse` |
| `fm_image` | `/fm/image` | `POST` | `FmImageRequest` | `FmImageResponse` |
| `fm_recommend` | `/fm/recommend` | `POST` | `FmRecommendRequest` | `FmRecommendResponse` |
| `fm_songs` | `/fm/songs` | `POST` | `FmSongsRequest` | `FmSongsResponse` |
| `images` | `/images` | `GET` | `ImagesRequest` | `ImagesResponse` |
| `images_audio` | `/images/audio` | `GET` | `ImagesAudioRequest` | `ImagesAudioResponse` |
| `ip` | `/ip` | `POST` | `IpRequest` | `IpResponse` |
| `ip_dateil` | `/ip/dateil` | `POST` | `IpDateilRequest` | `IpDateilResponse` |
| `ip_playlist` | `/ip/playlist` | `POST` | `IpPlaylistRequest` | `IpPlaylistResponse` |
| `ip_zone` | `/ip/zone` | `GET` | `IpZoneRequest` | `IpZoneResponse` |
| `ip_zone_home` | `/ip/zone/home` | `GET` | `IpZoneHomeRequest` | `IpZoneHomeResponse` |
| `kmr_audio_mv` | `/kmr/audio/mv` | `POST` | `KmrAudioMvRequest` | `KmrAudioMvResponse` |
| `krm_audio` | `/krm/audio` | `POST` | `KrmAudioRequest` | `KrmAudioResponse` |
| `lastest_songs_listen` | `/lastest/songs/listen` | `POST` | `LastestSongsListenRequest` | `LastestSongsListenResponse` |
| `login` | `/login` | `POST` | `LoginRequest` | `LoginResponse` |
| `login_cellphone` | `/login/cellphone` | `POST` | `LoginCellphoneRequest` | `LoginCellphoneResponse` |
| `login_device` | `/login/device` | `POST` | `LoginDeviceRequest` | `LoginDeviceResponse` |
| `login_openplat` | `/login/openplat` | `POST` | `LoginOpenplatRequest` | `LoginOpenplatResponse` |
| `login_qr_check` | `/login/qr/check` | `GET` | `LoginQrCheckRequest` | `LoginQrCheckResponse` |
| `login_qr_create` | `/login/qr/create` | `GET` | `LoginQrCreateRequest` | `LoginQrCreateResponse` |
| `login_qr_key` | `/login/qr/key` | `GET` | `LoginQrKeyRequest` | `LoginQrKeyResponse` |
| `login_token` | `/login/token` | `POST` | `LoginTokenRequest` | `LoginTokenResponse` |
| `login_wx_check` | `/login/wx/check` | `GET` | `LoginWxCheckRequest` | `LoginWxCheckResponse` |
| `login_wx_create` | `/login/wx/create` | `GET` | `LoginWxCreateRequest` | `LoginWxCreateResponse` |
| `longaudio_album_audios` | `/longaudio/album/audios` | `POST` | `LongaudioAlbumAudiosRequest` | `LongaudioAlbumAudiosResponse` |
| `longaudio_album_detail` | `/longaudio/album/detail` | `POST` | `LongaudioAlbumDetailRequest` | `LongaudioAlbumDetailResponse` |
| `longaudio_daily_recommend` | `/longaudio/daily/recommend` | `POST` | `LongaudioDailyRecommendRequest` | `LongaudioDailyRecommendResponse` |
| `longaudio_rank_recommend` | `/longaudio/rank/recommend` | `GET` | `LongaudioRankRecommendRequest` | `LongaudioRankRecommendResponse` |
| `longaudio_vip_recommend` | `/longaudio/vip/recommend` | `POST` | `LongaudioVipRecommendRequest` | `LongaudioVipRecommendResponse` |
| `longaudio_week_recommend` | `/longaudio/week/recommend` | `POST` | `LongaudioWeekRecommendRequest` | `LongaudioWeekRecommendResponse` |
| `lyric` | `/lyric` | `GET` | `LyricRequest` | `LyricResponse` |
| `pc_diantai` | `/pc/diantai` | `POST` | `PcDiantaiRequest` | `PcDiantaiResponse` |
| `personal_fm` | `/personal/fm` | `POST` | `PersonalFmRequest` | `PersonalFmResponse` |
| `playhistory_upload` | `/playhistory/upload` | `POST` | `PlayhistoryUploadRequest` | `PlayhistoryUploadResponse` |
| `playlist_add` | `/playlist/add` | `POST` | `PlaylistAddRequest` | `PlaylistAddResponse` |
| `playlist_del` | `/playlist/del` | `POST` | `PlaylistDelRequest` | `PlaylistDelResponse` |
| `playlist_detail` | `/playlist/detail` | `POST` | `PlaylistDetailRequest` | `PlaylistDetailResponse` |
| `playlist_effect` | `/playlist/effect` | `POST` | `PlaylistEffectRequest` | `PlaylistEffectResponse` |
| `playlist_similar` | `/playlist/similar` | `POST` | `PlaylistSimilarRequest` | `PlaylistSimilarResponse` |
| `playlist_tags` | `/playlist/tags` | `POST` | `PlaylistTagsRequest` | `PlaylistTagsResponse` |
| `playlist_track_all` | `/playlist/track/all` | `GET` | `PlaylistTrackAllRequest` | `PlaylistTrackAllResponse` |
| `playlist_track_all_new` | `/playlist/track/all/new` | `POST` | `PlaylistTrackAllNewRequest` | `PlaylistTrackAllNewResponse` |
| `playlist_tracks_add` | `/playlist/tracks/add` | `POST` | `PlaylistTracksAddRequest` | `PlaylistTracksAddResponse` |
| `playlist_tracks_del` | `/playlist/tracks/del` | `POST` | `PlaylistTracksDelRequest` | `PlaylistTracksDelResponse` |
| `privilege_lite` | `/privilege/lite` | `POST` | `PrivilegeLiteRequest` | `PrivilegeLiteResponse` |
| `rank_audio` | `/rank/audio` | `POST` | `RankAudioRequest` | `RankAudioResponse` |
| `rank_info` | `/rank/info` | `GET` | `RankInfoRequest` | `RankInfoResponse` |
| `rank_list` | `/rank/list` | `GET` | `RankListRequest` | `RankListResponse` |
| `rank_top` | `/rank/top` | `GET` | `RankTopRequest` | `RankTopResponse` |
| `rank_vol` | `/rank/vol` | `GET` | `RankVolRequest` | `RankVolResponse` |
| `recommend_songs` | `/recommend/songs` | `POST` | `RecommendSongsRequest` | `RecommendSongsResponse` |
| `register_dev` | `/register/dev` | `POST` | `RegisterDevRequest` | `RegisterDevResponse` |
| `scene_audio_list` | `/scene/audio/list` | `POST` | `SceneAudioListRequest` | `SceneAudioListResponse` |
| `scene_collection_list` | `/scene/collection/list` | `POST` | `SceneCollectionListRequest` | `SceneCollectionListResponse` |
| `scene_lists` | `/scene/lists` | `GET` | `SceneListsRequest` | `SceneListsResponse` |
| `scene_lists_v2` | `/scene/lists/v2` | `POST` | `SceneListsV2Request` | `SceneListsV2Response` |
| `scene_module` | `/scene/module` | `POST` | `SceneModuleRequest` | `SceneModuleResponse` |
| `scene_module_info` | `/scene/module/info` | `GET` | `SceneModuleInfoRequest` | `SceneModuleInfoResponse` |
| `scene_music` | `/scene/music` | `POST` | `SceneMusicRequest` | `SceneMusicResponse` |
| `scene_video_list` | `/scene/video/list` | `POST` | `SceneVideoListRequest` | `SceneVideoListResponse` |
| `search` | `/search` | `GET` | `SearchRequest` | `SearchResponse` |
| `search_complex` | `/search/complex` | `GET` | `SearchComplexRequest` | `SearchComplexResponse` |
| `search_default` | `/search/default` | `POST` | `SearchDefaultRequest` | `SearchDefaultResponse` |
| `search_hot` | `/search/hot` | `GET` | `SearchHotRequest` | `SearchHotResponse` |
| `search_lyric` | `/search/lyric` | `GET` | `SearchLyricRequest` | `SearchLyricResponse` |
| `search_mixed` | `/search/mixed` | `GET` | `SearchMixedRequest` | `SearchMixedResponse` |
| `search_suggest` | `/search/suggest` | `GET` | `SearchSuggestRequest` | `SearchSuggestResponse` |
| `server_now` | `/server/now` | `POST` | `ServerNowRequest` | `ServerNowResponse` |
| `sheet_collection` | `/sheet/collection` | `GET` | `SheetCollectionRequest` | `SheetCollectionResponse` |
| `sheet_collection_detail` | `/sheet/collection/detail` | `GET` | `SheetCollectionDetailRequest` | `SheetCollectionDetailResponse` |
| `sheet_detail` | `/sheet/detail` | `GET` | `SheetDetailRequest` | `SheetDetailResponse` |
| `sheet_hot` | `/sheet/hot` | `GET` | `SheetHotRequest` | `SheetHotResponse` |
| `sheet_list` | `/sheet/list` | `GET` | `SheetListRequest` | `SheetListResponse` |
| `singer_list` | `/singer/list` | `GET` | `SingerListRequest` | `SingerListResponse` |
| `song_climax` | `/song/climax` | `GET` | `SongClimaxRequest` | `SongClimaxResponse` |
| `song_ranking` | `/song/ranking` | `GET` | `SongRankingRequest` | `SongRankingResponse` |
| `song_ranking_filter` | `/song/ranking/filter` | `GET` | `SongRankingFilterRequest` | `SongRankingFilterResponse` |
| `song_url` | `/song/url` | `GET` | `SongUrlRequest` | `SongUrlResponse` |
| `song_url_new` | `/song/url/new` | `POST` | `SongUrlNewRequest` | `SongUrlNewResponse` |
| `theme_music` | `/theme/music` | `POST` | `ThemeMusicRequest` | `ThemeMusicResponse` |
| `theme_music_detail` | `/theme/music/detail` | `POST` | `ThemeMusicDetailRequest` | `ThemeMusicDetailResponse` |
| `theme_playlist` | `/theme/playlist` | `POST` | `ThemePlaylistRequest` | `ThemePlaylistResponse` |
| `theme_playlist_track` | `/theme/playlist/track` | `POST` | `ThemePlaylistTrackRequest` | `ThemePlaylistTrackResponse` |
| `top_album` | `/top/album` | `POST` | `TopAlbumRequest` | `TopAlbumResponse` |
| `top_card` | `/top/card` | `POST` | `TopCardRequest` | `TopCardResponse` |
| `top_card_youth` | `/top/card/youth` | `POST` | `TopCardYouthRequest` | `TopCardYouthResponse` |
| `top_ip` | `/top/ip` | `POST` | `TopIpRequest` | `TopIpResponse` |
| `top_playlist` | `/top/playlist` | `POST` | `TopPlaylistRequest` | `TopPlaylistResponse` |
| `top_song` | `/top/song` | `POST` | `TopSongRequest` | `TopSongResponse` |
| `user_cloud` | `/user/cloud` | `POST` | `UserCloudRequest` | `UserCloudResponse` |
| `user_cloud_url` | `/user/cloud/url` | `GET` | `UserCloudUrlRequest` | `UserCloudUrlResponse` |
| `user_detail` | `/user/detail` | `POST` | `UserDetailRequest` | `UserDetailResponse` |
| `user_follow` | `/user/follow` | `POST` | `UserFollowRequest` | `UserFollowResponse` |
| `user_history` | `/user/history` | `POST` | `UserHistoryRequest` | `UserHistoryResponse` |
| `user_listen` | `/user/listen` | `POST` | `UserListenRequest` | `UserListenResponse` |
| `user_playlist` | `/user/playlist` | `POST` | `UserPlaylistRequest` | `UserPlaylistResponse` |
| `user_video_collect` | `/user/video/collect` | `POST` | `UserVideoCollectRequest` | `UserVideoCollectResponse` |
| `user_video_love` | `/user/video/love` | `GET` | `UserVideoLoveRequest` | `UserVideoLoveResponse` |
| `user_vip_detail` | `/user/vip/detail` | `GET` | `UserVipDetailRequest` | `UserVipDetailResponse` |
| `video_detail` | `/video/detail` | `POST` | `VideoDetailRequest` | `VideoDetailResponse` |
| `video_privilege` | `/video/privilege` | `POST` | `VideoPrivilegeRequest` | `VideoPrivilegeResponse` |
| `video_url` | `/video/url` | `GET` | `VideoUrlRequest` | `VideoUrlResponse` |
| `youth_channel_all` | `/youth/channel/all` | `GET` | `YouthChannelAllRequest` | `YouthChannelAllResponse` |
| `youth_channel_amway` | `/youth/channel/amway` | `GET` | `YouthChannelAmwayRequest` | `YouthChannelAmwayResponse` |
| `youth_channel_detail` | `/youth/channel/detail` | `POST` | `YouthChannelDetailRequest` | `YouthChannelDetailResponse` |
| `youth_channel_similar` | `/youth/channel/similar` | `POST` | `YouthChannelSimilarRequest` | `YouthChannelSimilarResponse` |
| `youth_channel_song` | `/youth/channel/song` | `GET` | `YouthChannelSongRequest` | `YouthChannelSongResponse` |
| `youth_channel_song_detail` | `/youth/channel/song/detail` | `GET` | `YouthChannelSongDetailRequest` | `YouthChannelSongDetailResponse` |
| `youth_channel_sub` | `/youth/channel/sub` | `GET` | `YouthChannelSubRequest` | `YouthChannelSubResponse` |
| `youth_day_vip` | `/youth/day/vip` | `POST` | `YouthDayVipRequest` | `YouthDayVipResponse` |
| `youth_day_vip_upgrade` | `/youth/day/vip/upgrade` | `POST` | `YouthDayVipUpgradeRequest` | `YouthDayVipUpgradeResponse` |
| `youth_dynamic` | `/youth/dynamic` | `GET` | `YouthDynamicRequest` | `YouthDynamicResponse` |
| `youth_dynamic_recent` | `/youth/dynamic/recent` | `GET` | `YouthDynamicRecentRequest` | `YouthDynamicRecentResponse` |
| `youth_listen_song` | `/youth/listen/song` | `POST` | `YouthListenSongRequest` | `YouthListenSongResponse` |
| `youth_month_vip_record` | `/youth/month/vip/record` | `GET` | `YouthMonthVipRecordRequest` | `YouthMonthVipRecordResponse` |
| `youth_union_vip` | `/youth/union/vip` | `GET` | `YouthUnionVipRequest` | `YouthUnionVipResponse` |
| `youth_user_song` | `/youth/user/song` | `GET` | `YouthUserSongRequest` | `YouthUserSongResponse` |
| `youth_vip` | `/youth/vip` | `POST` | `YouthVipRequest` | `YouthVipResponse` |
| `yueku` | `/yueku` | `GET` | `YuekuRequest` | `YuekuResponse` |
| `yueku_banner` | `/yueku/banner` | `POST` | `YuekuBannerRequest` | `YuekuBannerResponse` |
| `yueku_fm` | `/yueku/fm` | `GET` | `YuekuFmRequest` | `YuekuFmResponse` |

## 函数用处（Client 方法）

| 函数 | 用处 | 对应接口 |
| --- | --- | --- |
| `(*Client).Call(ctx, route, req)` | 通用调用任意路由接口 | 通用 |
| `(*Client).CallByIdentifier(ctx, identifier, req)` | 按接口标识符调用 | 通用 |
| `(*Client).Endpoints()` | 获取全部可用接口清单 | 通用 |
| `(*Client).RouteByIdentifier(identifier)` | 根据标识符查询路由 | 通用 |
| `(*Client).SetCookie(key, value)` | 设置会话 Cookie | 通用 |
| `(*Client).Cookie()` | 获取当前 Cookie 池 | 通用 |
| `(*Client).LoginByPassword(...)` | 账号密码登录（含加密与 token 解包） | 登录封装 |
| `(*Client).LoginByCellphone(...)` | 手机验证码登录（含加密与 token 解包） | 登录封装 |
| `(*Client).LoginByToken(...)` | token 刷新登录 | 登录封装 |
| `(*Client).GetDailyRecommendGuest(...)` | 获取每日推荐（自动 fallback） | 每日推荐封装 |
| `(*Client).GetSongPlayURL(...)` | 获取歌曲可播放地址（按原 JS 逻辑） | 播放地址封装 |
| `(*Client).ResolveSongPlayURL(...)` | 获取歌曲地址明细（支持 option） | 播放地址封装 |
| `sdk.WithSongURLFallback(bool)` | 控制是否回退 song_url_new | 播放地址 Option |
| `sdk.WithSongURLAll(bool)` | 返回全部可用地址（而非仅首个） | 播放地址 Option |
| `sdk.WithSongURLDFID(string)` | 指定 dfid 参与地址解析 | 播放地址 Option |
| `(*Client).AiRecommend(...)` | (暂无注释) | `/ai/recommend` |
| `(*Client).Album(...)` | (暂无注释) | `/album` |
| `(*Client).AlbumDetail(...)` | 专辑详情 | `/album/detail` |
| `(*Client).AlbumShop(...)` | 唱片店 | `/album/shop` |
| `(*Client).AlbumSongs(...)` | 专辑音乐列表 | `/album/songs` |
| `(*Client).ArtistAlbums(...)` | 获取歌手专辑 | `/artist/albums` |
| `(*Client).ArtistAudios(...)` | 获取歌手单曲 | `/artist/audios` |
| `(*Client).ArtistDetail(...)` | 歌手详情 | `/artist/detail` |
| `(*Client).ArtistFollow(...)` | 关注歌手 | `/artist/follow` |
| `(*Client).ArtistFollowNewsongs(...)` | 获取关注歌手新歌 | `/artist/follow/newsongs` |
| `(*Client).ArtistHonour(...)` | 歌手荣誉详情 | `/artist/honour` |
| `(*Client).ArtistLists(...)` | 歌手列表 | `/artist/lists` |
| `(*Client).ArtistUnfollow(...)` | 取消关注歌手 | `/artist/unfollow` |
| `(*Client).ArtistVideos(...)` | 获取歌手mv | `/artist/videos` |
| `(*Client).Audio(...)` | (暂无注释) | `/audio` |
| `(*Client).AudioAccompanyMatching(...)` | (暂无注释) | `/audio/accompany/matching` |
| `(*Client).AudioKtvTotal(...)` | (暂无注释) | `/audio/ktv/total` |
| `(*Client).AudioRelated(...)` | https://listkmrp3cdnretry.kugou.com/v3/album_audio/related | `/audio/related` |
| `(*Client).Brush(...)` | (暂无注释) | `/brush` |
| `(*Client).CaptchaSent(...)` | 手机验证码发送 | `/captcha/sent` |
| `(*Client).CommentAlbum(...)` | 歌曲评论 | `/comment/album` |
| `(*Client).CommentCount(...)` | 歌曲评论数 | `/comment/count` |
| `(*Client).CommentFloor(...)` | 歌曲评论 | `/comment/floor` |
| `(*Client).CommentMusic(...)` | 歌曲评论 | `/comment/music` |
| `(*Client).CommentMusicClassify(...)` | 歌曲评论-根据分类获取评论 | `/comment/music/classify` |
| `(*Client).CommentMusicHotword(...)` | 歌曲评论-更具热词获取评论 | `/comment/music/hotword` |
| `(*Client).CommentPlaylist(...)` | 歌曲评论 | `/comment/playlist` |
| `(*Client).EverydayFriend(...)` | (暂无注释) | `/everyday/friend` |
| `(*Client).EverydayHistory(...)` | mode list ,song | `/everyday/history` |
| `(*Client).EverydayRecommend(...)` | (暂无注释) | `/everyday/recommend` |
| `(*Client).EverydayStyleRecommend(...)` | (暂无注释) | `/everyday/style/recommend` |
| `(*Client).FavoriteCount(...)` | (暂无注释) | `/favorite/count` |
| `(*Client).FmClass(...)` | (暂无注释) | `/fm/class` |
| `(*Client).FmImage(...)` | (暂无注释) | `/fm/image` |
| `(*Client).FmRecommend(...)` | (暂无注释) | `/fm/recommend` |
| `(*Client).FmSongs(...)` | fmType 生成 | `/fm/songs` |
| `(*Client).Images(...)` | (暂无注释) | `/images` |
| `(*Client).ImagesAudio(...)` | (暂无注释) | `/images/audio` |
| `(*Client).Ip(...)` | 根据 ip id 获取相对应的 歌曲/专辑/视频/歌手 | `/ip` |
| `(*Client).IpDateil(...)` | 获取ip详情 | `/ip/dateil` |
| `(*Client).IpPlaylist(...)` | 根据 ip 获取相对于歌单 | `/ip/playlist` |
| `(*Client).IpZone(...)` | ip 专区 | `/ip/zone` |
| `(*Client).IpZoneHome(...)` | 获取今日推荐信息，有可能为空 | `/ip/zone/home` |
| `(*Client).KmrAudioMv(...)` | 根据 album_audio_id/MixSongID 获取歌曲 相对应的 mv | `/kmr/audio/mv` |
| `(*Client).KrmAudio(...)` | 根据 album_audio_id/MixSongID 获取歌曲 相对应的 歌手/专辑/歌曲信息 | `/krm/audio` |
| `(*Client).LastestSongsListen(...)` | 获取继续播放信息 | `/lastest/songs/listen` |
| `(*Client).Login(...)` | (暂无注释) | `/login` |
| `(*Client).LoginCellphone(...)` | 手机登录 | `/login/cellphone` |
| `(*Client).LoginDevice(...)` | (暂无注释) | `/login/device` |
| `(*Client).LoginOpenplat(...)` | 开放平台登录 | `/login/openplat` |
| `(*Client).LoginQrCheck(...)` | 酷狗二维码状态检测 | `/login/qr/check` |
| `(*Client).LoginQrCreate(...)` | 酷狗二维码生成 | `/login/qr/create` |
| `(*Client).LoginQrKey(...)` | 二维码 key 生成接口 | `/login/qr/key` |
| `(*Client).LoginToken(...)` | 刷新登录 | `/login/token` |
| `(*Client).LoginWxCheck(...)` | (暂无注释) | `/login/wx/check` |
| `(*Client).LoginWxCreate(...)` | (暂无注释) | `/login/wx/create` |
| `(*Client).LongaudioAlbumAudios(...)` | (暂无注释) | `/longaudio/album/audios` |
| `(*Client).LongaudioAlbumDetail(...)` | (暂无注释) | `/longaudio/album/detail` |
| `(*Client).LongaudioDailyRecommend(...)` | (暂无注释) | `/longaudio/daily/recommend` |
| `(*Client).LongaudioRankRecommend(...)` | (暂无注释) | `/longaudio/rank/recommend` |
| `(*Client).LongaudioVipRecommend(...)` | (暂无注释) | `/longaudio/vip/recommend` |
| `(*Client).LongaudioWeekRecommend(...)` | (暂无注释) | `/longaudio/week/recommend` |
| `(*Client).Lyric(...)` | 歌词获取 | `/lyric` |
| `(*Client).PcDiantai(...)` | 电台 banner | `/pc/diantai` |
| `(*Client).PersonalFm(...)` | (暂无注释) | `/personal/fm` |
| `(*Client).PlayhistoryUpload(...)` | 提交听歌历史 | `/playhistory/upload` |
| `(*Client).PlaylistAdd(...)` | 收藏歌单 | `/playlist/add` |
| `(*Client).PlaylistDel(...)` | 取消收藏歌单 | `/playlist/del` |
| `(*Client).PlaylistDetail(...)` | 获取歌单详情 | `/playlist/detail` |
| `(*Client).PlaylistEffect(...)` | 获取音效歌单 | `/playlist/effect` |
| `(*Client).PlaylistSimilar(...)` | (暂无注释) | `/playlist/similar` |
| `(*Client).PlaylistTags(...)` | 获取歌单分类 | `/playlist/tags` |
| `(*Client).PlaylistTrackAll(...)` | 获取歌单所有歌曲 | `/playlist/track/all` |
| `(*Client).PlaylistTrackAllNew(...)` | 获取歌单所有歌曲 | `/playlist/track/all/new` |
| `(*Client).PlaylistTracksAdd(...)` | 对歌单添加歌曲 | `/playlist/tracks/add` |
| `(*Client).PlaylistTracksDel(...)` | 对歌单删除歌曲 | `/playlist/tracks/del` |
| `(*Client).PrivilegeLite(...)` | 获取歌曲信息 | `/privilege/lite` |
| `(*Client).RankAudio(...)` | 获取排行榜音乐列表 | `/rank/audio` |
| `(*Client).RankInfo(...)` | 获取排行榜详情 | `/rank/info` |
| `(*Client).RankList(...)` | 获取排行榜列表 | `/rank/list` |
| `(*Client).RankTop(...)` | 获取排行榜推荐列表 | `/rank/top` |
| `(*Client).RankVol(...)` | 获取排行榜往期列表 | `/rank/vol` |
| `(*Client).RecommendSongs(...)` | 每日推荐歌曲 | `/recommend/songs` |
| `(*Client).RegisterDev(...)` | 可用内存，单位是字节 | `/register/dev` |
| `(*Client).SceneAudioList(...)` | (暂无注释) | `/scene/audio/list` |
| `(*Client).SceneCollectionList(...)` | (暂无注释) | `/scene/collection/list` |
| `(*Client).SceneLists(...)` | (暂无注释) | `/scene/lists` |
| `(*Client).SceneListsV2(...)` | (暂无注释) | `/scene/lists/v2` |
| `(*Client).SceneModule(...)` | (暂无注释) | `/scene/module` |
| `(*Client).SceneModuleInfo(...)` | (暂无注释) | `/scene/module/info` |
| `(*Client).SceneMusic(...)` | (暂无注释) | `/scene/music` |
| `(*Client).SceneVideoList(...)` | (暂无注释) | `/scene/video/list` |
| `(*Client).Search(...)` | 搜索 | `/search` |
| `(*Client).SearchComplex(...)` | 综合搜索 | `/search/complex` |
| `(*Client).SearchDefault(...)` | (暂无注释) | `/search/default` |
| `(*Client).SearchHot(...)` | 热搜 | `/search/hot` |
| `(*Client).SearchLyric(...)` | 歌词搜索 | `/search/lyric` |
| `(*Client).SearchMixed(...)` | 综合搜索 | `/search/mixed` |
| `(*Client).SearchSuggest(...)` | (暂无注释) | `/search/suggest` |
| `(*Client).ServerNow(...)` | 获取服务器时间 | `/server/now` |
| `(*Client).SheetCollection(...)` | 乐谱详情 | `/sheet/collection` |
| `(*Client).SheetCollectionDetail(...)` | 乐谱合集详情 | `/sheet/collection/detail` |
| `(*Client).SheetDetail(...)` | 乐谱详情 | `/sheet/detail` |
| `(*Client).SheetHot(...)` | 推荐乐谱 | `/sheet/hot` |
| `(*Client).SheetList(...)` | 乐谱列表 // 0全部，1、钢琴，2、吉他，3、鼓谱，98：简谱，99:其他 | `/sheet/list` |
| `(*Client).SingerList(...)` | 获取歌手列表 | `/singer/list` |
| `(*Client).SongClimax(...)` | 获取音频高潮部分 | `/song/climax` |
| `(*Client).SongRanking(...)` | 歌曲成绩单 | `/song/ranking` |
| `(*Client).SongRankingFilter(...)` | 歌曲成绩单 | `/song/ranking/filter` |
| `(*Client).SongUrl(...)` | 获取音乐urls | `/song/url` |
| `(*Client).SongUrlNew(...)` | const quality = ['piano', 'acappella', 'subwoofer', 'ancient', 'dj', 'surnay'].includes(params.quality) | `/song/url/new` |
| `(*Client).ThemeMusic(...)` | 获取主题音乐 | `/theme/music` |
| `(*Client).ThemeMusicDetail(...)` | 获取主题音乐详情 | `/theme/music/detail` |
| `(*Client).ThemePlaylist(...)` | 主题歌单 | `/theme/playlist` |
| `(*Client).ThemePlaylistTrack(...)` | 获取主题歌单说有歌曲 | `/theme/playlist/track` |
| `(*Client).TopAlbum(...)` | 推荐专辑 | `/top/album` |
| `(*Client).TopCard(...)` | 热门好歌精选 | `/top/card` |
| `(*Client).TopCardYouth(...)` | 热门好歌精选 | `/top/card/youth` |
| `(*Client).TopIp(...)` | (暂无注释) | `/top/ip` |
| `(*Client).TopPlaylist(...)` | 歌单 | `/top/playlist` |
| `(*Client).TopSong(...)` | (暂无注释) | `/top/song` |
| `(*Client).UserCloud(...)` | (暂无注释) | `/user/cloud` |
| `(*Client).UserCloudUrl(...)` | 获取云盘音乐url | `/user/cloud/url` |
| `(*Client).UserDetail(...)` | (暂无注释) | `/user/detail` |
| `(*Client).UserFollow(...)` | (暂无注释) | `/user/follow` |
| `(*Client).UserHistory(...)` | 获取用户听歌排行 | `/user/history` |
| `(*Client).UserListen(...)` | (暂无注释) | `/user/listen` |
| `(*Client).UserPlaylist(...)` | 获取用户歌单 | `/user/playlist` |
| `(*Client).UserVideoCollect(...)` | (暂无注释) | `/user/video/collect` |
| `(*Client).UserVideoLove(...)` | (暂无注释) | `/user/video/love` |
| `(*Client).UserVipDetail(...)` | (暂无注释) | `/user/vip/detail` |
| `(*Client).VideoDetail(...)` | 获取视频详情 | `/video/detail` |
| `(*Client).VideoPrivilege(...)` | 获取视频特权 | `/video/privilege` |
| `(*Client).VideoUrl(...)` | 获取视频urls | `/video/url` |
| `(*Client).YouthChannelAll(...)` | (暂无注释) | `/youth/channel/all` |
| `(*Client).YouthChannelAmway(...)` | (暂无注释) | `/youth/channel/amway` |
| `(*Client).YouthChannelDetail(...)` | (暂无注释) | `/youth/channel/detail` |
| `(*Client).YouthChannelSimilar(...)` | (暂无注释) | `/youth/channel/similar` |
| `(*Client).YouthChannelSong(...)` | (暂无注释) | `/youth/channel/song` |
| `(*Client).YouthChannelSongDetail(...)` | (暂无注释) | `/youth/channel/song/detail` |
| `(*Client).YouthChannelSub(...)` | (暂无注释) | `/youth/channel/sub` |
| `(*Client).YouthDayVip(...)` | 领取vip(领取一天) 需要登录 | `/youth/day/vip` |
| `(*Client).YouthDayVipUpgrade(...)` | 升级vip | `/youth/day/vip/upgrade` |
| `(*Client).YouthDynamic(...)` | (暂无注释) | `/youth/dynamic` |
| `(*Client).YouthDynamicRecent(...)` | (暂无注释) | `/youth/dynamic/recent` |
| `(*Client).YouthListenSong(...)` | 听歌领取vip 需要登录 | `/youth/listen/song` |
| `(*Client).YouthMonthVipRecord(...)` | (暂无注释) | `/youth/month/vip/record` |
| `(*Client).YouthUnionVip(...)` | 领取vip 需要登录 | `/youth/union/vip` |
| `(*Client).YouthUserSong(...)` | (暂无注释) | `/youth/user/song` |
| `(*Client).YouthVip(...)` | 领取vip 需要登录 | `/youth/vip` |
| `(*Client).Yueku(...)` | 获取安卓乐库相关内容 | `/yueku` |
| `(*Client).YuekuBanner(...)` | 获取乐库下的 banner | `/yueku/banner` |
| `(*Client).YuekuFm(...)` | 获取乐库下的 fm | `/yueku/fm` |

## 登录会话示例

可复用示例：
- 路径：`examples/login_session/main.go`
- 行为：
  - 首次运行：提示输入账号和密码登录
  - 登录成功后：将会话保存到 `~/.kugou_music_api_session.json`
  - 再次运行：自动复用会话并获取用户信息
  - 会话过期：自动尝试 token 刷新；失败则重新提示登录

运行：

```bash
go run ./examples/login_session
```

## 更多示例

- 手机验证码登录并保存会话：`go run ./examples/login_cellphone`
- 二维码登录并保存会话：`go run ./examples/login_qrcode`
- 使用本地会话进行 token 续期：`go run ./examples/token_refresh`
- 获取当前用户信息并自动校验/刷新 token：`go run ./examples/user_info_refresh`
- 批量联调待实测接口（默认跳过风险写操作）：`go run ./examples/pending_smoke`
- 个性电台三模式示例（红心/新歌/小众）：`go run ./examples/personal_radio -mode=heart`
- 创建歌单（支持 debug=true 打印响应）：`go run ./examples/create_playlist -name="2026-03-22红心日推"`
- 给歌单添加歌曲（支持 source=radio|daily，支持自动下载）：`go run ./examples/add_playlist_tracks -listid=123456 -source=daily -limit=50 -download=1`

会话文件默认：`~/.kugou_music_api_session.json`，可通过环境变量 `KUGOU_SESSION_FILE` 覆盖。

短信验证码接口建议使用 `(*Client).SendCaptcha(...)`，该方法已按原 JS 逻辑固定 `businessid=5`、`plat=3`。
