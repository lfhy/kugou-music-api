package main

import (
	"fmt"
	"os"
	"strings"
)

func writeFixStatus(models []apiModel) error {
	var b strings.Builder
	b.WriteString("# API Fix Status\n\n")
	b.WriteString("说明：`已校对修复` 表示已进入 Go 兼容适配链路并完成校对；`待实测` 表示手写封装已落地，待人工联调验收。\n\n")
	b.WriteString("| Identifier | Route | Status | Mode | Note |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, m := range models {
		status, mode, note := fixStatusForModel(m)
		b.WriteString(fmt.Sprintf("| `%s` | `%s` | `%s` | `%s` | %s |\n", m.Identifier, m.Route, status, mode, note))
	}
	return os.WriteFile("sdk/API_FIX_STATUS.md", []byte(b.String()), 0o644)
}

func fixStatusForModel(m apiModel) (status, mode, note string) {
	manual := map[string]bool{
		"captcha_sent":              true,
		"daily_recommend":           true,
		"song_url":                  true,
		"song_url_new":              true,
		"login":                     true,
		"login_cellphone":           true,
		"login_token":               true,
		"login_qr_key":              true,
		"login_qr_create":           true,
		"login_qr_check":            true,
		"login_openplat":            true,
		"login_wx_create":           true,
		"login_wx_check":            true,
		"login_device":              true,
		"everyday_recommend":        true,
		"recommend_songs":           true,
		"user_detail":               true,
		"audio_related":             true,
		"audio_accompany_matching":  true,
		"audio_ktv_total":           true,
		"brush":                     true,
		"register_dev":              true,
		"user_video_collect":        true,
		"user_video_love":           true,
		"search_mixed":              true,
		"fm_class":                  true,
		"fm_image":                  true,
		"fm_recommend":              true,
		"fm_songs":                  true,
		"ai_recommend":              true,
		"album":                     true,
		"artist_audios":             true,
		"artist_follow":             true,
		"artist_unfollow":           true,
		"audio":                     true,
		"comment_album":             true,
		"comment_floor":             true,
		"comment_music":             true,
		"comment_music_hotword":     true,
		"comment_playlist":          true,
		"personal_fm":               true,
		"playlist_del":              true,
		"playlist_add":              true,
		"playlist_tracks_add":       true,
		"playlist_tracks_del":       true,
		"playlist_similar":          true,
		"top_card":                  true,
		"top_playlist":              true,
		"user_cloud":                true,
		"user_cloud_url":            true,
		"user_follow":               true,
		"user_listen":               true,
		"video_detail":              true,
		"video_privilege":           true,
		"album_songs":               true,
		"artist_albums":             true,
		"artist_lists":              true,
		"artist_videos":             true,
		"comment_music_classify":    true,
		"lastest_songs_listen":      true,
		"lyric":                     true,
		"playhistory_upload":        true,
		"playlist_track_all":        true,
		"playlist_track_all_new":    true,
		"privilege_lite":            true,
		"rank_audio":                true,
		"search_complex":            true,
		"search_default":            true,
		"search_lyric":              true,
		"sheet_collection":          true,
		"sheet_collection_detail":   true,
		"sheet_detail":              true,
		"sheet_list":                true,
		"theme_music":               true,
		"theme_music_detail":        true,
		"theme_playlist":            true,
		"theme_playlist_track":      true,
		"top_card_youth":            true,
		"top_ip":                    true,
		"top_song":                  true,
		"user_history":              true,
		"video_url":                 true,
		"youth_channel_song":        true,
		"youth_channel_song_detail": true,
		"youth_day_vip_upgrade":     true,
		"youth_listen_song":         true,
		"youth_union_vip":           true,
		"youth_user_song":           true,
		"youth_vip":                 true,
		"yueku_banner":              true,
	}
	if manual[m.Identifier] {
		return "已校对修复", "manual", "手写封装优先"
	}
	manualPending := map[string]bool{
		"album_detail":              true,
		"album_shop":                true,
		"artist_detail":             true,
		"artist_follow_newsongs":    true,
		"artist_honour":             true,
		"everyday_friend":           true,
		"everyday_history":          true,
		"everyday_style_recommend":  true,
		"favorite_count":            true,
		"ip_zone":                   true,
		"kmr_audio_mv":              true,
		"krm_audio":                 true,
		"longaudio_album_audios":    true,
		"longaudio_album_detail":    true,
		"longaudio_daily_recommend": true,
		"longaudio_rank_recommend":  true,
		"longaudio_vip_recommend":   true,
		"longaudio_week_recommend":  true,
		"rank_info":                 true,
		"rank_list":                 true,
		"rank_top":                  true,
		"rank_vol":                  true,
		"scene_audio_list":          true,
		"scene_collection_list":     true,
		"scene_lists":               true,
		"scene_lists_v2":            true,
		"scene_module":              true,
		"scene_module_info":         true,
		"scene_music":               true,
		"scene_video_list":          true,
		"search_suggest":            true,
		"server_now":                true,
		"sheet_hot":                 true,
		"singer_list":               true,
		"song_climax":               true,
		"song_ranking":              true,
		"song_ranking_filter":       true,
		"user_vip_detail":           true,
		"youth_channel_all":         true,
		"youth_channel_amway":       true,
		"youth_channel_detail":      true,
		"youth_channel_similar":     true,
		"youth_channel_sub":         true,
		"youth_day_vip":             true,
		"youth_dynamic":             true,
		"youth_dynamic_recent":      true,
		"youth_month_vip_record":    true,
		"yueku":                     true,
		"yueku_fm":                  true,
	}
	if manualPending[m.Identifier] {
		return "待实测", "manual-pending", "手写封装已完成，待联调校验"
	}
	rules := len(m.Compat.DataMapRules) + len(m.Compat.ParamsMapRules) + len(m.Compat.CookieRules)
	if rules > 0 && len(m.Compat.UnsupportedExprs) == 0 {
		return "已校对修复", "auto-compat", fmt.Sprintf("自动规则 %d 条", rules)
	}
	if rules > 0 && len(m.Compat.UnsupportedExprs) > 0 {
		return "部分修复", "auto-compat", fmt.Sprintf("自动规则 %d 条, 未覆盖表达式 %d 条", rules, len(m.Compat.UnsupportedExprs))
	}
	return "待校对", "none", "-"
}
