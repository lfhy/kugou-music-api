# API Sign/Encrypt Check

说明：按原 JS 模块静态扫描是否存在自定义签名/加密字段（非基础请求头签名）。

| Identifier | Route | Need | Type | Status | Mode | Note |
| --- | --- | --- | --- | --- | --- | --- |
| `ai_recommend` | `/ai/recommend` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `album` | `/album` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `artist_audios` | `/artist/audios` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `artist_follow` | `/artist/follow` | `YES` | `custom-sign,custom-aes,signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `artist_unfollow` | `/artist/unfollow` | `YES` | `custom-sign,custom-aes,signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `audio` | `/audio` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `audio_accompany_matching` | `/audio/accompany/matching` | `YES` | `custom-hash` | `已校对修复` | `manual` | 手写封装优先 |
| `audio_ktv_total` | `/audio/ktv/total` | `YES` | `custom-hash` | `已校对修复` | `manual` | 手写封装优先 |
| `audio_related` | `/audio/related` | `YES` | `custom-hash` | `已校对修复` | `manual` | 手写封装优先 |
| `brush` | `/brush` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `comment_album` | `/comment/album` | `YES` | `signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `comment_floor` | `/comment/floor` | `YES` | `signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `comment_music` | `/comment/music` | `YES` | `signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `comment_music_hotword` | `/comment/music/hotword` | `YES` | `signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `comment_playlist` | `/comment/playlist` | `YES` | `signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `fm_class` | `/fm/class` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `fm_image` | `/fm/image` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `fm_recommend` | `/fm/recommend` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `fm_songs` | `/fm/songs` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `images` | `/images` | `YES` | `custom-sign` | `已校对修复` | `auto-compat` | 自动规则 5 条 |
| `images_audio` | `/images/audio` | `YES` | `custom-sign` | `已校对修复` | `auto-compat` | 自动规则 4 条 |
| `login` | `/login` | `YES` | `custom-sign,custom-aes,signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `login_cellphone` | `/login/cellphone` | `YES` | `custom-sign,custom-aes,signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `login_device` | `/login/device` | `YES` | `custom-sign,custom-aes,signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `login_openplat` | `/login/openplat` | `YES` | `custom-sign,custom-aes` | `已校对修复` | `manual` | 手写封装优先 |
| `login_token` | `/login/token` | `YES` | `custom-sign,custom-aes` | `已校对修复` | `manual` | 手写封装优先 |
| `login_wx_create` | `/login/wx/create` | `YES` | `custom-hash` | `已校对修复` | `manual` | 手写封装优先 |
| `personal_fm` | `/personal/fm` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `playlist_del` | `/playlist/del` | `YES` | `custom-sign,custom-aes` | `已校对修复` | `manual` | 手写封装优先 |
| `playlist_similar` | `/playlist/similar` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `register_dev` | `/register/dev` | `YES` | `custom-sign,custom-aes` | `已校对修复` | `manual` | 手写封装优先 |
| `search_mixed` | `/search/mixed` | `YES` | `custom-hash` | `已校对修复` | `manual` | 手写封装优先 |
| `song_url_new` | `/song/url/new` | `YES` | `custom-hash` | `已校对修复` | `manual` | 手写封装优先 |
| `top_card` | `/top/card` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `top_playlist` | `/top/playlist` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `user_cloud` | `/user/cloud` | `YES` | `custom-sign,custom-aes` | `已校对修复` | `manual` | 手写封装优先 |
| `user_cloud_url` | `/user/cloud/url` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `user_detail` | `/user/detail` | `YES` | `custom-sign,custom-aes,signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `user_follow` | `/user/follow` | `YES` | `custom-sign,signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `user_listen` | `/user/listen` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `user_video_collect` | `/user/video/collect` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |
| `user_video_love` | `/user/video/love` | `YES` | `custom-sign,signed-field` | `已校对修复` | `manual` | 手写封装优先 |
| `video_detail` | `/video/detail` | `YES` | `custom-sign,custom-hash` | `已校对修复` | `manual` | 手写封装优先 |
| `video_privilege` | `/video/privilege` | `YES` | `custom-sign` | `已校对修复` | `manual` | 手写封装优先 |

总计接口: 153，需自定义签名/加密: 44，其中 已校对修复: 44，部分修复: 0，待校对: 0
