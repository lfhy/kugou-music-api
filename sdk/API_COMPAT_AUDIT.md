# API Compatibility Audit

说明：该报告逐个模块检查 Go 自动生成调用与原 JS 逻辑的一致性风险。
`HIGH` 表示该接口很可能需要手写适配层（默认值/参数重组/加解密/后处理）。

| Identifier | Route | Risk | Reasons |
| --- | --- | --- | --- |
| `ai_recommend` | `/ai/recommend` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `album` | `/album` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `album_detail` | `/album/detail` | `HIGH` | default fallback params in module |
| `album_shop` | `/album/shop` | `HIGH` | default fallback params in module |
| `album_songs` | `/album/songs` | `HIGH` | default fallback params in module; manual request map assembly |
| `artist_albums` | `/artist/albums` | `HIGH` | default fallback params in module; manual request map assembly |
| `artist_audios` | `/artist/audios` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `artist_detail` | `/artist/detail` | `HIGH` | default fallback params in module |
| `artist_follow` | `/artist/follow` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `artist_follow_newsongs` | `/artist/follow/newsongs` | `HIGH` | default fallback params in module; manual request map assembly |
| `artist_honour` | `/artist/honour` | `HIGH` | default fallback params in module |
| `artist_lists` | `/artist/lists` | `HIGH` | default fallback params in module; manual request map assembly |
| `artist_unfollow` | `/artist/unfollow` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `artist_videos` | `/artist/videos` | `HIGH` | default fallback params in module; manual request map assembly |
| `audio` | `/audio` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `audio_accompany_matching` | `/audio/accompany/matching` | `HIGH` | custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `audio_ktv_total` | `/audio/ktv/total` | `HIGH` | custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `audio_related` | `/audio/related` | `HIGH` | custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `brush` | `/brush` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `captcha_sent` | `/captcha/sent` | `HIGH` | manual request map assembly |
| `comment_album` | `/comment/album` | `HIGH` | default fallback params in module; manual request map assembly |
| `comment_count` | `/comment/count` | `HIGH` | default fallback params in module; manual request map assembly |
| `comment_floor` | `/comment/floor` | `HIGH` | default fallback params in module; manual request map assembly |
| `comment_music` | `/comment/music` | `HIGH` | default fallback params in module; manual request map assembly |
| `comment_music_classify` | `/comment/music/classify` | `HIGH` | default fallback params in module; manual request map assembly |
| `comment_music_hotword` | `/comment/music/hotword` | `HIGH` | default fallback params in module; manual request map assembly |
| `comment_playlist` | `/comment/playlist` | `HIGH` | default fallback params in module; manual request map assembly |
| `everyday_friend` | `/everyday/friend` | `HIGH` | default fallback params in module |
| `everyday_history` | `/everyday/history` | `HIGH` | default fallback params in module; manual request map assembly |
| `everyday_recommend` | `/everyday/recommend` | `HIGH` | default fallback params in module |
| `everyday_style_recommend` | `/everyday/style/recommend` | `HIGH` | default fallback params in module; manual request map assembly |
| `favorite_count` | `/favorite/count` | `HIGH` | default fallback params in module |
| `fm_class` | `/fm/class` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `fm_image` | `/fm/image` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `fm_recommend` | `/fm/recommend` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `fm_songs` | `/fm/songs` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `images` | `/images` | `HIGH` | custom crypto/sign logic; default fallback params in module; manual request map assembly; dynamic template URL |
| `images_audio` | `/images/audio` | `HIGH` | custom crypto/sign logic; default fallback params in module; manual request map assembly; dynamic template URL |
| `ip` | `/ip` | `HIGH` | default fallback params in module; manual request map assembly; dynamic template URL |
| `ip_dateil` | `/ip/dateil` | `HIGH` | default fallback params in module; manual request map assembly |
| `ip_playlist` | `/ip/playlist` | `HIGH` | default fallback params in module; manual request map assembly |
| `ip_zone` | `/ip/zone` | `HIGH` | custom response post-processing; default fallback params in module |
| `ip_zone_home` | `/ip/zone/home` | `HIGH` | default fallback params in module; manual request map assembly |
| `kmr_audio_mv` | `/kmr/audio/mv` | `HIGH` | default fallback params in module; manual request map assembly |
| `krm_audio` | `/krm/audio` | `HIGH` | default fallback params in module; manual request map assembly |
| `lastest_songs_listen` | `/lastest/songs/listen` | `HIGH` | default fallback params in module; manual request map assembly |
| `login` | `/login` | `HIGH` | dynamic timestamp; custom crypto/sign logic; custom response post-processing; default fallback params in module; manual request map assembly |
| `login_cellphone` | `/login/cellphone` | `HIGH` | dynamic timestamp; random/default device fields; custom crypto/sign logic; custom response post-processing; default fallback params in module; manual request map assembly |
| `login_device` | `/login/device` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `login_openplat` | `/login/openplat` | `HIGH` | dynamic timestamp; custom crypto/sign logic; custom response post-processing; default fallback params in module; manual request map assembly |
| `login_qr_check` | `/login/qr/check` | `HIGH` | custom response post-processing; default fallback params in module |
| `login_qr_create` | `/login/qr/create` | `HIGH` | custom response post-processing |
| `login_qr_key` | `/login/qr/key` | `HIGH` | default fallback params in module |
| `login_token` | `/login/token` | `HIGH` | dynamic timestamp; custom crypto/sign logic; custom response post-processing; default fallback params in module; manual request map assembly |
| `login_wx_check` | `/login/wx/check` | `HIGH` | custom response post-processing; default fallback params in module; custom cookie mapping; dynamic template URL |
| `login_wx_create` | `/login/wx/create` | `HIGH` | dynamic timestamp; random/default device fields; custom crypto/sign logic; custom response post-processing; custom cookie mapping |
| `longaudio_album_audios` | `/longaudio/album/audios` | `HIGH` | default fallback params in module |
| `longaudio_album_detail` | `/longaudio/album/detail` | `HIGH` | default fallback params in module |
| `longaudio_daily_recommend` | `/longaudio/daily/recommend` | `HIGH` | default fallback params in module |
| `longaudio_rank_recommend` | `/longaudio/rank/recommend` | `HIGH` | default fallback params in module |
| `longaudio_vip_recommend` | `/longaudio/vip/recommend` | `HIGH` | default fallback params in module |
| `longaudio_week_recommend` | `/longaudio/week/recommend` | `HIGH` | default fallback params in module |
| `lyric` | `/lyric` | `HIGH` | custom response post-processing; custom body decoding; default fallback params in module; manual request map assembly |
| `pc_diantai` | `/pc/diantai` | `HIGH` | default fallback params in module; manual request map assembly |
| `personal_fm` | `/personal/fm` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `playhistory_upload` | `/playhistory/upload` | `HIGH` | dynamic timestamp; default fallback params in module; manual request map assembly |
| `playlist_add` | `/playlist/add` | `HIGH` | dynamic timestamp; default fallback params in module; manual request map assembly |
| `playlist_del` | `/playlist/del` | `HIGH` | dynamic timestamp; custom crypto/sign logic; custom response post-processing; default fallback params in module; manual request map assembly |
| `playlist_detail` | `/playlist/detail` | `HIGH` | default fallback params in module; manual request map assembly |
| `playlist_effect` | `/playlist/effect` | `HIGH` | default fallback params in module; manual request map assembly |
| `playlist_similar` | `/playlist/similar` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `playlist_tags` | `/playlist/tags` | `HIGH` | default fallback params in module; manual request map assembly |
| `playlist_track_all` | `/playlist/track/all` | `HIGH` | default fallback params in module; manual request map assembly |
| `playlist_track_all_new` | `/playlist/track/all/new` | `HIGH` | default fallback params in module; manual request map assembly |
| `playlist_tracks_add` | `/playlist/tracks/add` | `HIGH` | dynamic timestamp; default fallback params in module; manual request map assembly |
| `playlist_tracks_del` | `/playlist/tracks/del` | `HIGH` | default fallback params in module; manual request map assembly |
| `privilege_lite` | `/privilege/lite` | `HIGH` | default fallback params in module; manual request map assembly |
| `rank_audio` | `/rank/audio` | `HIGH` | default fallback params in module; manual request map assembly |
| `rank_info` | `/rank/info` | `HIGH` | default fallback params in module |
| `rank_list` | `/rank/list` | `HIGH` | default fallback params in module |
| `rank_top` | `/rank/top` | `HIGH` | default fallback params in module |
| `rank_vol` | `/rank/vol` | `HIGH` | default fallback params in module |
| `recommend_songs` | `/recommend/songs` | `HIGH` | default fallback params in module; manual request map assembly |
| `register_dev` | `/register/dev` | `HIGH` | custom crypto/sign logic; custom response post-processing; default fallback params in module; manual request map assembly |
| `scene_audio_list` | `/scene/audio/list` | `HIGH` | default fallback params in module; manual request map assembly |
| `scene_collection_list` | `/scene/collection/list` | `HIGH` | default fallback params in module; manual request map assembly |
| `scene_lists` | `/scene/lists` | `HIGH` | default fallback params in module |
| `scene_lists_v2` | `/scene/lists/v2` | `HIGH` | default fallback params in module |
| `scene_module` | `/scene/module` | `HIGH` | default fallback params in module |
| `scene_module_info` | `/scene/module/info` | `HIGH` | default fallback params in module |
| `scene_music` | `/scene/music` | `HIGH` | default fallback params in module |
| `scene_video_list` | `/scene/video/list` | `HIGH` | default fallback params in module; manual request map assembly |
| `search` | `/search` | `HIGH` | default fallback params in module; manual request map assembly |
| `search_complex` | `/search/complex` | `HIGH` | default fallback params in module; manual request map assembly |
| `search_default` | `/search/default` | `HIGH` | default fallback params in module; manual request map assembly |
| `search_hot` | `/search/hot` | `HIGH` | default fallback params in module; manual request map assembly |
| `search_lyric` | `/search/lyric` | `HIGH` | default fallback params in module; manual request map assembly |
| `search_mixed` | `/search/mixed` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `search_suggest` | `/search/suggest` | `HIGH` | default fallback params in module |
| `server_now` | `/server/now` | `HIGH` | default fallback params in module |
| `sheet_collection` | `/sheet/collection` | `HIGH` | default fallback params in module; manual request map assembly |
| `sheet_collection_detail` | `/sheet/collection/detail` | `HIGH` | default fallback params in module; manual request map assembly |
| `sheet_detail` | `/sheet/detail` | `HIGH` | default fallback params in module; manual request map assembly |
| `sheet_hot` | `/sheet/hot` | `HIGH` | default fallback params in module; manual request map assembly |
| `sheet_list` | `/sheet/list` | `HIGH` | default fallback params in module; manual request map assembly |
| `singer_list` | `/singer/list` | `HIGH` | default fallback params in module |
| `song_climax` | `/song/climax` | `HIGH` | default fallback params in module |
| `song_ranking` | `/song/ranking` | `HIGH` | default fallback params in module |
| `song_ranking_filter` | `/song/ranking/filter` | `HIGH` | default fallback params in module |
| `song_url` | `/song/url` | `HIGH` | random/default device fields; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `song_url_new` | `/song/url/new` | `HIGH` | dynamic timestamp; random/default device fields; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `theme_music` | `/theme/music` | `HIGH` | dynamic timestamp; default fallback params in module; manual request map assembly |
| `theme_music_detail` | `/theme/music/detail` | `HIGH` | dynamic timestamp; default fallback params in module; manual request map assembly |
| `theme_playlist` | `/theme/playlist` | `HIGH` | dynamic timestamp; default fallback params in module; manual request map assembly |
| `theme_playlist_track` | `/theme/playlist/track` | `HIGH` | dynamic timestamp; default fallback params in module; manual request map assembly |
| `top_album` | `/top/album` | `HIGH` | default fallback params in module; manual request map assembly |
| `top_card` | `/top/card` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `top_card_youth` | `/top/card/youth` | `HIGH` | default fallback params in module; manual request map assembly |
| `top_ip` | `/top/ip` | `HIGH` | custom response post-processing; default fallback params in module; manual request map assembly |
| `top_playlist` | `/top/playlist` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `top_song` | `/top/song` | `HIGH` | default fallback params in module; manual request map assembly |
| `user_cloud` | `/user/cloud` | `HIGH` | dynamic timestamp; custom crypto/sign logic; custom response post-processing; custom body decoding; default fallback params in module; manual request map assembly |
| `user_cloud_url` | `/user/cloud/url` | `HIGH` | custom crypto/sign logic; manual request map assembly |
| `user_detail` | `/user/detail` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `user_follow` | `/user/follow` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `user_history` | `/user/history` | `HIGH` | default fallback params in module; manual request map assembly |
| `user_listen` | `/user/listen` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `user_playlist` | `/user/playlist` | `HIGH` | default fallback params in module; manual request map assembly |
| `user_video_collect` | `/user/video/collect` | `HIGH` | custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `user_video_love` | `/user/video/love` | `HIGH` | custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `user_vip_detail` | `/user/vip/detail` | `HIGH` | default fallback params in module |
| `video_detail` | `/video/detail` | `HIGH` | dynamic timestamp; custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `video_privilege` | `/video/privilege` | `HIGH` | custom crypto/sign logic; default fallback params in module; manual request map assembly |
| `video_url` | `/video/url` | `HIGH` | manual request map assembly |
| `youth_channel_all` | `/youth/channel/all` | `HIGH` | default fallback params in module |
| `youth_channel_amway` | `/youth/channel/amway` | `LOW` | - |
| `youth_channel_detail` | `/youth/channel/detail` | `HIGH` | default fallback params in module |
| `youth_channel_similar` | `/youth/channel/similar` | `HIGH` | default fallback params in module |
| `youth_channel_song` | `/youth/channel/song` | `HIGH` | default fallback params in module; manual request map assembly |
| `youth_channel_song_detail` | `/youth/channel/song/detail` | `HIGH` | manual request map assembly |
| `youth_channel_sub` | `/youth/channel/sub` | `HIGH` | dynamic template URL |
| `youth_day_vip` | `/youth/day/vip` | `LOW` | - |
| `youth_day_vip_upgrade` | `/youth/day/vip/upgrade` | `HIGH` | default fallback params in module; manual request map assembly |
| `youth_dynamic` | `/youth/dynamic` | `LOW` | - |
| `youth_dynamic_recent` | `/youth/dynamic/recent` | `LOW` | - |
| `youth_listen_song` | `/youth/listen/song` | `HIGH` | default fallback params in module; manual request map assembly |
| `youth_month_vip_record` | `/youth/month/vip/record` | `LOW` | - |
| `youth_union_vip` | `/youth/union/vip` | `HIGH` | manual request map assembly |
| `youth_user_song` | `/youth/user/song` | `HIGH` | default fallback params in module; manual request map assembly |
| `youth_vip` | `/youth/vip` | `HIGH` | dynamic timestamp; manual request map assembly |
| `yueku` | `/yueku` | `HIGH` | default fallback params in module |
| `yueku_banner` | `/yueku/banner` | `HIGH` | default fallback params in module; manual request map assembly |
| `yueku_fm` | `/yueku/fm` | `HIGH` | default fallback params in module |

总计: 153, HIGH: 148, LOW: 5
