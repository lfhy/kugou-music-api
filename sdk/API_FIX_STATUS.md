# API Fix Status

说明：`已校对修复` 表示已进入 Go 兼容适配链路并完成校对；`待实测` 表示手写封装已落地，待人工联调验收。

| Identifier | Route | Status | Mode | Note |
| --- | --- | --- | --- | --- |
| `ai_recommend` | `/ai/recommend` | `已校对修复` | `manual` | 手写封装优先 |
| `album` | `/album` | `已校对修复` | `manual` | 手写封装优先 |
| `album_detail` | `/album/detail` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `album_shop` | `/album/shop` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `album_songs` | `/album/songs` | `已校对修复` | `manual` | 手写封装优先 |
| `artist_albums` | `/artist/albums` | `已校对修复` | `manual` | 手写封装优先 |
| `artist_audios` | `/artist/audios` | `已校对修复` | `manual` | 手写封装优先 |
| `artist_detail` | `/artist/detail` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `artist_follow` | `/artist/follow` | `已校对修复` | `manual` | 手写封装优先 |
| `artist_follow_newsongs` | `/artist/follow/newsongs` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `artist_honour` | `/artist/honour` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `artist_lists` | `/artist/lists` | `已校对修复` | `manual` | 手写封装优先 |
| `artist_unfollow` | `/artist/unfollow` | `已校对修复` | `manual` | 手写封装优先 |
| `artist_videos` | `/artist/videos` | `已校对修复` | `manual` | 手写封装优先 |
| `audio` | `/audio` | `已校对修复` | `manual` | 手写封装优先 |
| `audio_accompany_matching` | `/audio/accompany/matching` | `已校对修复` | `manual` | 手写封装优先 |
| `audio_ktv_total` | `/audio/ktv/total` | `已校对修复` | `manual` | 手写封装优先 |
| `audio_related` | `/audio/related` | `已校对修复` | `manual` | 手写封装优先 |
| `brush` | `/brush` | `已校对修复` | `manual` | 手写封装优先 |
| `captcha_sent` | `/captcha/sent` | `已校对修复` | `manual` | 手写封装优先 |
| `comment_album` | `/comment/album` | `已校对修复` | `manual` | 手写封装优先 |
| `comment_count` | `/comment/count` | `已校对修复` | `auto-compat` | 自动规则 2 条 |
| `comment_floor` | `/comment/floor` | `已校对修复` | `manual` | 手写封装优先 |
| `comment_music` | `/comment/music` | `已校对修复` | `manual` | 手写封装优先 |
| `comment_music_classify` | `/comment/music/classify` | `已校对修复` | `manual` | 手写封装优先 |
| `comment_music_hotword` | `/comment/music/hotword` | `已校对修复` | `manual` | 手写封装优先 |
| `comment_playlist` | `/comment/playlist` | `已校对修复` | `manual` | 手写封装优先 |
| `everyday_friend` | `/everyday/friend` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `everyday_history` | `/everyday/history` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `everyday_recommend` | `/everyday/recommend` | `已校对修复` | `manual` | 手写封装优先 |
| `everyday_style_recommend` | `/everyday/style/recommend` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `favorite_count` | `/favorite/count` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `fm_class` | `/fm/class` | `已校对修复` | `manual` | 手写封装优先 |
| `fm_image` | `/fm/image` | `已校对修复` | `manual` | 手写封装优先 |
| `fm_recommend` | `/fm/recommend` | `已校对修复` | `manual` | 手写封装优先 |
| `fm_songs` | `/fm/songs` | `已校对修复` | `manual` | 手写封装优先 |
| `images` | `/images` | `已校对修复` | `auto-compat` | 自动规则 5 条 |
| `images_audio` | `/images/audio` | `已校对修复` | `auto-compat` | 自动规则 4 条 |
| `ip` | `/ip` | `已校对修复` | `auto-compat` | 自动规则 6 条 |
| `ip_dateil` | `/ip/dateil` | `已校对修复` | `auto-compat` | 自动规则 1 条 |
| `ip_playlist` | `/ip/playlist` | `已校对修复` | `auto-compat` | 自动规则 3 条 |
| `ip_zone` | `/ip/zone` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `ip_zone_home` | `/ip/zone/home` | `已校对修复` | `auto-compat` | 自动规则 2 条 |
| `kmr_audio_mv` | `/kmr/audio/mv` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `krm_audio` | `/krm/audio` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `lastest_songs_listen` | `/lastest/songs/listen` | `已校对修复` | `manual` | 手写封装优先 |
| `login` | `/login` | `已校对修复` | `manual` | 手写封装优先 |
| `login_cellphone` | `/login/cellphone` | `已校对修复` | `manual` | 手写封装优先 |
| `login_device` | `/login/device` | `已校对修复` | `manual` | 手写封装优先 |
| `login_openplat` | `/login/openplat` | `已校对修复` | `manual` | 手写封装优先 |
| `login_qr_check` | `/login/qr/check` | `已校对修复` | `manual` | 手写封装优先 |
| `login_qr_create` | `/login/qr/create` | `已校对修复` | `manual` | 手写封装优先 |
| `login_qr_key` | `/login/qr/key` | `已校对修复` | `manual` | 手写封装优先 |
| `login_token` | `/login/token` | `已校对修复` | `manual` | 手写封装优先 |
| `login_wx_check` | `/login/wx/check` | `已校对修复` | `manual` | 手写封装优先 |
| `login_wx_create` | `/login/wx/create` | `已校对修复` | `manual` | 手写封装优先 |
| `longaudio_album_audios` | `/longaudio/album/audios` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `longaudio_album_detail` | `/longaudio/album/detail` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `longaudio_daily_recommend` | `/longaudio/daily/recommend` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `longaudio_rank_recommend` | `/longaudio/rank/recommend` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `longaudio_vip_recommend` | `/longaudio/vip/recommend` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `longaudio_week_recommend` | `/longaudio/week/recommend` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `lyric` | `/lyric` | `已校对修复` | `manual` | 手写封装优先 |
| `pc_diantai` | `/pc/diantai` | `已校对修复` | `auto-compat` | 自动规则 2 条 |
| `personal_fm` | `/personal/fm` | `已校对修复` | `manual` | 手写封装优先 |
| `playhistory_upload` | `/playhistory/upload` | `已校对修复` | `manual` | 手写封装优先 |
| `playlist_add` | `/playlist/add` | `已校对修复` | `manual` | 手写封装优先 |
| `playlist_del` | `/playlist/del` | `已校对修复` | `manual` | 手写封装优先 |
| `playlist_detail` | `/playlist/detail` | `已校对修复` | `auto-compat` | 自动规则 2 条 |
| `playlist_effect` | `/playlist/effect` | `已校对修复` | `auto-compat` | 自动规则 2 条 |
| `playlist_similar` | `/playlist/similar` | `已校对修复` | `manual` | 手写封装优先 |
| `playlist_tags` | `/playlist/tags` | `已校对修复` | `auto-compat` | 自动规则 3 条 |
| `playlist_track_all` | `/playlist/track/all` | `已校对修复` | `manual` | 手写封装优先 |
| `playlist_track_all_new` | `/playlist/track/all/new` | `已校对修复` | `manual` | 手写封装优先 |
| `playlist_tracks_add` | `/playlist/tracks/add` | `已校对修复` | `manual` | 手写封装优先 |
| `playlist_tracks_del` | `/playlist/tracks/del` | `已校对修复` | `manual` | 手写封装优先 |
| `privilege_lite` | `/privilege/lite` | `已校对修复` | `manual` | 手写封装优先 |
| `rank_audio` | `/rank/audio` | `已校对修复` | `manual` | 手写封装优先 |
| `rank_info` | `/rank/info` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `rank_list` | `/rank/list` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `rank_top` | `/rank/top` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `rank_vol` | `/rank/vol` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `recommend_songs` | `/recommend/songs` | `已校对修复` | `manual` | 手写封装优先 |
| `register_dev` | `/register/dev` | `已校对修复` | `manual` | 手写封装优先 |
| `scene_audio_list` | `/scene/audio/list` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `scene_collection_list` | `/scene/collection/list` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `scene_lists` | `/scene/lists` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `scene_lists_v2` | `/scene/lists/v2` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `scene_module` | `/scene/module` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `scene_module_info` | `/scene/module/info` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `scene_music` | `/scene/music` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `scene_video_list` | `/scene/video/list` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `search` | `/search` | `已校对修复` | `auto-compat` | 自动规则 7 条 |
| `search_complex` | `/search/complex` | `已校对修复` | `manual` | 手写封装优先 |
| `search_default` | `/search/default` | `已校对修复` | `manual` | 手写封装优先 |
| `search_hot` | `/search/hot` | `已校对修复` | `auto-compat` | 自动规则 2 条 |
| `search_lyric` | `/search/lyric` | `已校对修复` | `manual` | 手写封装优先 |
| `search_mixed` | `/search/mixed` | `已校对修复` | `manual` | 手写封装优先 |
| `search_suggest` | `/search/suggest` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `server_now` | `/server/now` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `sheet_collection` | `/sheet/collection` | `已校对修复` | `manual` | 手写封装优先 |
| `sheet_collection_detail` | `/sheet/collection/detail` | `已校对修复` | `manual` | 手写封装优先 |
| `sheet_detail` | `/sheet/detail` | `已校对修复` | `manual` | 手写封装优先 |
| `sheet_hot` | `/sheet/hot` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `sheet_list` | `/sheet/list` | `已校对修复` | `manual` | 手写封装优先 |
| `singer_list` | `/singer/list` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `song_climax` | `/song/climax` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `song_ranking` | `/song/ranking` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `song_ranking_filter` | `/song/ranking/filter` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `song_url` | `/song/url` | `已校对修复` | `manual` | 手写封装优先 |
| `song_url_new` | `/song/url/new` | `已校对修复` | `manual` | 手写封装优先 |
| `theme_music` | `/theme/music` | `已校对修复` | `manual` | 手写封装优先 |
| `theme_music_detail` | `/theme/music/detail` | `已校对修复` | `manual` | 手写封装优先 |
| `theme_playlist` | `/theme/playlist` | `已校对修复` | `manual` | 手写封装优先 |
| `theme_playlist_track` | `/theme/playlist/track` | `已校对修复` | `manual` | 手写封装优先 |
| `top_album` | `/top/album` | `已校对修复` | `auto-compat` | 自动规则 4 条 |
| `top_card` | `/top/card` | `已校对修复` | `manual` | 手写封装优先 |
| `top_card_youth` | `/top/card/youth` | `已校对修复` | `manual` | 手写封装优先 |
| `top_ip` | `/top/ip` | `已校对修复` | `manual` | 手写封装优先 |
| `top_playlist` | `/top/playlist` | `已校对修复` | `manual` | 手写封装优先 |
| `top_song` | `/top/song` | `已校对修复` | `manual` | 手写封装优先 |
| `user_cloud` | `/user/cloud` | `已校对修复` | `manual` | 手写封装优先 |
| `user_cloud_url` | `/user/cloud/url` | `已校对修复` | `manual` | 手写封装优先 |
| `user_detail` | `/user/detail` | `已校对修复` | `manual` | 手写封装优先 |
| `user_follow` | `/user/follow` | `已校对修复` | `manual` | 手写封装优先 |
| `user_history` | `/user/history` | `已校对修复` | `manual` | 手写封装优先 |
| `user_listen` | `/user/listen` | `已校对修复` | `manual` | 手写封装优先 |
| `user_playlist` | `/user/playlist` | `已校对修复` | `auto-compat` | 自动规则 4 条 |
| `user_video_collect` | `/user/video/collect` | `已校对修复` | `manual` | 手写封装优先 |
| `user_video_love` | `/user/video/love` | `已校对修复` | `manual` | 手写封装优先 |
| `user_vip_detail` | `/user/vip/detail` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `video_detail` | `/video/detail` | `已校对修复` | `manual` | 手写封装优先 |
| `video_privilege` | `/video/privilege` | `已校对修复` | `manual` | 手写封装优先 |
| `video_url` | `/video/url` | `已校对修复` | `manual` | 手写封装优先 |
| `youth_channel_all` | `/youth/channel/all` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `youth_channel_amway` | `/youth/channel/amway` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `youth_channel_detail` | `/youth/channel/detail` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `youth_channel_similar` | `/youth/channel/similar` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `youth_channel_song` | `/youth/channel/song` | `已校对修复` | `manual` | 手写封装优先 |
| `youth_channel_song_detail` | `/youth/channel/song/detail` | `已校对修复` | `manual` | 手写封装优先 |
| `youth_channel_sub` | `/youth/channel/sub` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `youth_day_vip` | `/youth/day/vip` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `youth_day_vip_upgrade` | `/youth/day/vip/upgrade` | `已校对修复` | `manual` | 手写封装优先 |
| `youth_dynamic` | `/youth/dynamic` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `youth_dynamic_recent` | `/youth/dynamic/recent` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `youth_listen_song` | `/youth/listen/song` | `已校对修复` | `manual` | 手写封装优先 |
| `youth_month_vip_record` | `/youth/month/vip/record` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `youth_union_vip` | `/youth/union/vip` | `已校对修复` | `manual` | 手写封装优先 |
| `youth_user_song` | `/youth/user/song` | `已校对修复` | `manual` | 手写封装优先 |
| `youth_vip` | `/youth/vip` | `已校对修复` | `manual` | 手写封装优先 |
| `yueku` | `/yueku` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
| `yueku_banner` | `/yueku/banner` | `已校对修复` | `manual` | 手写封装优先 |
| `yueku_fm` | `/yueku/fm` | `待实测` | `manual-pending` | 手写封装已完成，待联调校验 |
