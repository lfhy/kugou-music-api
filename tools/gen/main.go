package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/lfhy/kugou-music-api/core/parser"
)

type fieldType int

const (
	tAny fieldType = iota
	tString
	tInt
	tBool
)

type apiModel struct {
	Identifier string
	Route      string
	Method     string
	Name       string
	ReqName    string
	RespName   string
	Spec       parser.ModuleSpec
	Compat     parser.CompatSpec
	Fields     map[string]fieldType
}

var manualSkipIdentifiers = map[string]bool{
	"user_detail":               true,
	"audio_related":             true,
	"audio_accompany_matching":  true,
	"audio_ktv_total":           true,
	"brush":                     true,
	"register_dev":              true,
	"user_video_collect":        true,
	"user_video_love":           true,
	"search_mixed":              true,
	"everyday_recommend":        true,
	"recommend_songs":           true,
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

func shouldSkipManual(identifier string) bool {
	id := strings.ToLower(strings.TrimSpace(identifier))
	if manualSkipIdentifiers[id] {
		return true
	}
	id2 := strings.ReplaceAll(id, "_", "")
	switch id2 {
	case "userdetail", "audiorelated", "audioaccompanymatching", "audioktvtotal", "brush", "registerdev", "uservideocollect", "uservideolove", "searchmixed", "everydayrecommend", "recommendsongs", "fmclass", "fmimage", "fmrecommend", "fmsongs", "airecommend", "album", "artistaudios", "artistfollow", "artistunfollow", "audio", "commentalbum", "commentfloor", "commentmusic", "commentmusichotword", "commentplaylist", "personalfm", "playlistdel", "playlistadd", "playlisttracksadd", "playlisttracksdel", "playlistsimilar", "topcard", "topplaylist", "usercloud", "usercloudurl", "userfollow", "userlisten", "videodetail", "videoprivilege", "albumsongs", "artistalbums", "artistlists", "artistvideos", "commentmusicclassify", "lastestsongslisten", "lyric", "playhistoryupload", "playlisttrackall", "playlisttrackallnew", "privilegelite", "rankaudio", "searchcomplex", "searchdefault", "searchlyric", "sheetcollection", "sheetcollectiondetail", "sheetdetail", "sheetlist", "thememusic", "thememusicdetail", "themeplaylist", "themeplaylisttrack", "topcardyouth", "topip", "topsong", "userhistory", "videourl", "youthchannelsong", "youthchannelsongdetail", "youthdayvipupgrade", "youthlistensong", "youthunionvip", "youthusersong", "youthvip", "yuekubanner":
		return true
	default:
		return false
	}
}

var (
	reParamDotQ      = regexp.MustCompile(`params\?\.([a-zA-Z_][a-zA-Z0-9_]*)`)
	reParamDot       = regexp.MustCompile(`params\.([a-zA-Z_][a-zA-Z0-9_]*)`)
	reParamBracketQ  = regexp.MustCompile(`params\?\[['\"]([a-zA-Z_][a-zA-Z0-9_]*)['\"]\]`)
	reParamBracket   = regexp.MustCompile(`params\[['\"]([a-zA-Z_][a-zA-Z0-9_]*)['\"]\]`)
	reIntDefaultQ    = regexp.MustCompile(`params\?\.([a-zA-Z_][a-zA-Z0-9_]*)\s*\|\|\s*[0-9]+`)
	reIntNumberWrapQ = regexp.MustCompile(`Number\(params\?\.([a-zA-Z_][a-zA-Z0-9_]*)`)
	reIntNumberWrap  = regexp.MustCompile(`Number\(params\.([a-zA-Z_][a-zA-Z0-9_]*)`)
	reStringDefaultQ = regexp.MustCompile(`params\?\.([a-zA-Z_][a-zA-Z0-9_]*)\s*\|\|\s*['\"]`)
	reBoolTernaryQ   = regexp.MustCompile(`params\?\.([a-zA-Z_][a-zA-Z0-9_]*)\s*\?\s*1\s*:\s*0`)
)

func main() {
	moduleDir := os.Getenv("KUGOU_MODULE_DIR")
	if strings.TrimSpace(moduleDir) == "" {
		moduleDir = "../module"
	}

	specs, err := parser.LoadModuleSpecs(moduleDir)
	if err != nil {
		panic(err)
	}

	models := make([]apiModel, 0, len(specs))
	for _, sp := range specs {
		srcPath := filepath.Join(moduleDir, sp.Identifier+".js")
		src, _ := os.ReadFile(srcPath)
		m := apiModel{
			Identifier: sp.Identifier,
			Route:      sp.Route,
			Method:     strings.ToUpper(sp.Method),
			Name:       toExportName(sp.Identifier),
			Spec:       sp,
			Compat:     parser.ExtractCompatSpec(sp.Identifier, string(src)),
			Fields:     extractFields(string(src)),
		}
		m.ReqName = m.Name + "Request"
		m.RespName = m.Name + "Response"
		models = append(models, m)
	}

	sort.Slice(models, func(i, j int) bool { return models[i].Route < models[j].Route })

	if err := writeGeneratedGo(models); err != nil {
		panic(err)
	}
	if err := writeCatalog(models); err != nil {
		panic(err)
	}
	if err := writeCompatGenerated(models); err != nil {
		panic(err)
	}
	if err := writeFixStatus(models); err != nil {
		panic(err)
	}
	fmt.Printf("generated %d apis\n", len(models))
}

func writeGeneratedGo(models []apiModel) error {
	var b bytes.Buffer
	b.WriteString("// Code generated by tools/gen; DO NOT EDIT.\n")
	b.WriteString("package sdk\n\n")
	b.WriteString("import (\n\t\"context\"\n)\n\n")

	b.WriteString("type APIInfo struct {\n")
	b.WriteString("\tIdentifier string\n\tRoute string\n\tMethod string\n\tRequestModel string\n\tResponseModel string\n}\n\n")

	b.WriteString("type apiSpec struct {\n")
	b.WriteString("\tIdentifier string\n\tRoute string\n\tMethod string\n\tURL string\n\tBaseURL string\n\tEncryptType string\n\tEncryptKey bool\n\tClearDefaultParams bool\n\tNotSignature bool\n\tUseParams bool\n\tUseData bool\n\tHeaders map[string]string\n}\n\n")

	for _, m := range models {
		b.WriteString(fmt.Sprintf("const Route%s = %q\n", m.Name, m.Route))
	}
	b.WriteString("\n")

	b.WriteString("var APIList = []APIInfo{\n")
	for _, m := range models {
		b.WriteString(fmt.Sprintf("\t{Identifier: %q, Route: Route%s, Method: %q, RequestModel: %q, ResponseModel: %q},\n", m.Identifier, m.Name, m.Method, m.ReqName, m.RespName))
	}
	b.WriteString("}\n\n")

	b.WriteString("var apiSpecMap = map[string]apiSpec{\n")
	for _, m := range models {
		b.WriteString(fmt.Sprintf("\tRoute%s: {Identifier: %q, Route: Route%s, Method: %q, URL: %q, BaseURL: %q, EncryptType: %q, EncryptKey: %t, ClearDefaultParams: %t, NotSignature: %t, UseParams: %t, UseData: %t, Headers: map[string]string{",
			m.Name, m.Identifier, m.Name, m.Method, m.Spec.UpstreamURL, m.Spec.BaseURL, m.Spec.EncryptType, m.Spec.EncryptKey, m.Spec.ClearDefaultParams, m.Spec.NotSignature, m.Spec.UseParams, m.Spec.UseData))
		hKeys := make([]string, 0, len(m.Spec.Headers))
		for k := range m.Spec.Headers {
			hKeys = append(hKeys, k)
		}
		sort.Strings(hKeys)
		for _, k := range hKeys {
			b.WriteString(fmt.Sprintf("%q:%q,", k, m.Spec.Headers[k]))
		}
		b.WriteString("}},\n")
	}
	b.WriteString("}\n\n")

	for _, m := range models {
		fields := sortedFieldKeys(m.Fields)
		b.WriteString(fmt.Sprintf("type %s struct {\n", m.ReqName))
		for _, k := range fields {
			if k == "cookie" {
				continue
			}
			ft := goType(m.Fields[k])
			b.WriteString(fmt.Sprintf("\t%s %s `json:\"%s,omitempty\"`\n", toExportName(k), ft, k))
		}
		b.WriteString("\tCookie map[string]string `json:\"-\"`\n")
		b.WriteString("\tExtra map[string]any `json:\"-\"`\n")
		b.WriteString("}\n\n")

		b.WriteString(fmt.Sprintf("type %s = Response\n\n", m.RespName))

		if shouldSkipManual(m.Identifier) || hasManualMethod(m.Identifier) ||
			m.Name == "FmClass" || m.Name == "FmImage" || m.Name == "FmRecommend" || m.Name == "FmSongs" ||
			m.Name == "AiRecommend" || m.Name == "Album" || m.Name == "ArtistAudios" || m.Name == "ArtistFollow" || m.Name == "ArtistUnfollow" ||
			m.Name == "Audio" || m.Name == "CommentAlbum" || m.Name == "CommentFloor" || m.Name == "CommentMusic" || m.Name == "CommentMusicHotword" || m.Name == "CommentPlaylist" ||
			m.Name == "PersonalFm" || m.Name == "PlaylistDel" || m.Name == "PlaylistSimilar" || m.Name == "TopCard" || m.Name == "TopPlaylist" ||
			m.Name == "UserCloud" || m.Name == "UserCloudUrl" || m.Name == "UserFollow" || m.Name == "UserListen" || m.Name == "VideoDetail" || m.Name == "VideoPrivilege" {
			continue
		}

		b.WriteString(fmt.Sprintf("func (c *Client) %s(ctx context.Context, req %s) (*%s, error) {\n", m.Name, m.ReqName, m.RespName))
		b.WriteString("\tparams := structToMap(req)\n")
		b.WriteString("\tdelete(params, \"Cookie\")\n\tdelete(params, \"Extra\")\n")
		b.WriteString(fmt.Sprintf("\tif compat, ok := buildCompatParams(%q, params, req.Cookie); ok { params = compat }\n", m.Identifier))
		b.WriteString("\tfor k, v := range req.Extra { params[k] = v }\n")
		b.WriteString(fmt.Sprintf("\tcookie := applyCompatCookie(%q, req.Cookie)\n", m.Identifier))
		b.WriteString(fmt.Sprintf("\tresp, err := c.Call(ctx, Route%s, Request{Params: params, Cookie: cookie})\n", m.Name))
		b.WriteString("\tif err != nil { return nil, err }\n")
		b.WriteString(fmt.Sprintf("\tout := %s(*resp)\n", m.RespName))
		b.WriteString("\treturn &out, nil\n")
		b.WriteString("}\n\n")
	}

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		_ = os.WriteFile("sdk/generated_apis.go.unformatted", b.Bytes(), 0o644)
		return err
	}
	return os.WriteFile("sdk/generated_apis.go", formatted, 0o644)
}

func hasManualMethod(identifier string) bool {
	id := strings.ToLower(strings.TrimSpace(identifier))
	return manualSkipIdentifiers[id]
}

func writeCatalog(models []apiModel) error {
	var b strings.Builder
	b.WriteString("# API Catalog\n\n")
	b.WriteString("| Identifier | Route | Method | Request Model | Response Model |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, m := range models {
		b.WriteString(fmt.Sprintf("| `%s` | `%s` | `%s` | `%s` | `%s` |\n", m.Identifier, m.Route, m.Method, m.ReqName, m.RespName))
	}
	return os.WriteFile("sdk/API_CATALOG.md", []byte(b.String()), 0o644)
}

func writeCompatGenerated(models []apiModel) error {
	var b bytes.Buffer
	b.WriteString("// Code generated by tools/gen; DO NOT EDIT.\n")
	b.WriteString("package sdk\n\n")
	b.WriteString("import (\n\t\"fmt\"\n\t\"strconv\"\n)\n\n")

	b.WriteString("func buildCompatParams(identifier string, in map[string]any, cookie map[string]string) (map[string]any, bool) {\n")
	b.WriteString("\tswitch identifier {\n")
	for _, m := range models {
		rules := m.Compat.DataMapRules
		if len(rules) == 0 {
			rules = m.Compat.ParamsMapRules
		}
		if len(rules) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("\tcase %q:\n", m.Identifier))
		b.WriteString("\t\tout := map[string]any{}\n")
		for _, r := range rules {
			expr := toGoCompatExpr(r.Expr, "in", "cookie")
			if expr == "" {
				continue
			}
			b.WriteString(fmt.Sprintf("\t\tout[%q] = %s\n", r.Field, expr))
		}
		b.WriteString("\t\tfor k, v := range in { if _, ok := out[k]; !ok { out[k] = v } }\n")
		b.WriteString("\t\treturn out, true\n")
	}
	b.WriteString("\tdefault:\n\t\treturn in, false\n\t}\n}\n\n")

	b.WriteString("func applyCompatCookie(identifier string, cookie map[string]string) map[string]string {\n")
	b.WriteString("\tout := map[string]string{}\n")
	b.WriteString("\tfor k, v := range cookie { out[k] = v }\n")
	b.WriteString("\tswitch identifier {\n")
	for _, m := range models {
		if len(m.Compat.CookieRules) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("\tcase %q:\n", m.Identifier))
		for _, r := range m.Compat.CookieRules {
			b.WriteString(fmt.Sprintf("\t\tif v := cookie[%q]; v != \"\" { out[%q] = v }\n", r.Expr.CookieKey, r.Field))
		}
	}
	b.WriteString("\t}\n\treturn out\n}\n\n")

	b.WriteString("func compatFirstAny(v ...any) any {\n")
	b.WriteString("\tfor _, x := range v {\n")
	b.WriteString("\t\tswitch vv := x.(type) {\n")
	b.WriteString("\t\tcase nil:\n")
	b.WriteString("\t\t\tcontinue\n")
	b.WriteString("\t\tcase string:\n")
	b.WriteString("\t\t\tif vv != \"\" { return vv }\n")
	b.WriteString("\t\tdefault:\n")
	b.WriteString("\t\t\treturn x\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn nil\n")
	b.WriteString("}\n\n")

	b.WriteString("func compatFirstAnyString(v any, def string) string {\n")
	b.WriteString("\tif v == nil { return def }\n")
	b.WriteString("\ts := fmt.Sprintf(\"%v\", v)\n")
	b.WriteString("\tif s == \"\" || s == \"<nil>\" { return def }\n")
	b.WriteString("\treturn s\n")
	b.WriteString("}\n\n")

	b.WriteString("func compatFirstAnyInt(v any, def int) int {\n")
	b.WriteString("\tif v == nil { return def }\n")
	b.WriteString("\ts := fmt.Sprintf(\"%v\", v)\n")
	b.WriteString("\tif s == \"\" || s == \"<nil>\" { return def }\n")
	b.WriteString("\tn, err := strconv.Atoi(s)\n")
	b.WriteString("\tif err != nil || n == 0 { return def }\n")
	b.WriteString("\treturn n\n")
	b.WriteString("}\n\n")

	b.WriteString("func compatDefaultLiteral(kind string, s string, i int, b bool) any {\n")
	b.WriteString("\tswitch kind {\n")
	b.WriteString("\tcase \"string\":\n\t\treturn s\n")
	b.WriteString("\tcase \"int\":\n\t\treturn i\n")
	b.WriteString("\tcase \"bool\":\n\t\treturn b\n")
	b.WriteString("\tdefault:\n\t\treturn nil\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n")

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		_ = os.WriteFile("sdk/generated_compat.go.unformatted", b.Bytes(), 0o644)
		return err
	}
	return os.WriteFile("sdk/generated_compat.go", formatted, 0o644)
}

func toGoCompatExpr(expr parser.CompatExpr, inVar, cookieVar string) string {
	switch expr.Kind {
	case parser.ExprLiteralString:
		return fmt.Sprintf("%q", expr.DefaultStr)
	case parser.ExprLiteralInt:
		return fmt.Sprintf("%d", expr.DefaultInt)
	case parser.ExprLiteralBool:
		if expr.DefaultBool {
			return "true"
		}
		return "false"
	case parser.ExprParam:
		return fmt.Sprintf("%s[%q]", inVar, expr.ParamKey)
	case parser.ExprCookie:
		return fmt.Sprintf("%s[%q]", cookieVar, expr.CookieKey)
	case parser.ExprParamOrString:
		return fmt.Sprintf("compatFirstAnyString(%s[%q], %q)", inVar, expr.ParamKey, expr.DefaultStr)
	case parser.ExprParamOrInt:
		return fmt.Sprintf("compatFirstAnyInt(%s[%q], %d)", inVar, expr.ParamKey, expr.DefaultInt)
	case parser.ExprParamOrCookie:
		return fmt.Sprintf("compatFirstAny(%s[%q], %s[%q], compatDefaultLiteral(%q, %q, %d, %t))", inVar, expr.ParamKey, cookieVar, expr.CookieKey, expr.DefaultType, expr.DefaultStr, expr.DefaultInt, expr.DefaultBool)
	case parser.ExprTemplateParam:
		return fmt.Sprintf("fmt.Sprintf(\"%%v\", %s[%q])", inVar, expr.ParamKey)
	case parser.ExprNumberOrInt:
		return fmt.Sprintf("compatFirstAnyInt(%s[%q], %d)", inVar, expr.ParamKey, expr.DefaultInt)
	default:
		return ""
	}
}

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

func extractFields(src string) map[string]fieldType {
	fields := map[string]fieldType{}
	add := func(name string, t fieldType) {
		if strings.TrimSpace(name) == "" {
			return
		}
		if name == "cookie" || name == "body" {
			return
		}
		if cur, ok := fields[name]; !ok || t > cur {
			fields[name] = t
		}
	}
	for _, m := range reParamDotQ.FindAllStringSubmatch(src, -1) {
		add(m[1], tAny)
	}
	for _, m := range reParamDot.FindAllStringSubmatch(src, -1) {
		add(m[1], tAny)
	}
	for _, m := range reParamBracketQ.FindAllStringSubmatch(src, -1) {
		add(m[1], tAny)
	}
	for _, m := range reParamBracket.FindAllStringSubmatch(src, -1) {
		add(m[1], tAny)
	}
	for _, m := range reIntDefaultQ.FindAllStringSubmatch(src, -1) {
		add(m[1], tInt)
	}
	for _, m := range reIntNumberWrapQ.FindAllStringSubmatch(src, -1) {
		add(m[1], tInt)
	}
	for _, m := range reIntNumberWrap.FindAllStringSubmatch(src, -1) {
		add(m[1], tInt)
	}
	for _, m := range reStringDefaultQ.FindAllStringSubmatch(src, -1) {
		add(m[1], tString)
	}
	for _, m := range reBoolTernaryQ.FindAllStringSubmatch(src, -1) {
		add(m[1], tBool)
	}
	return fields
}

func sortedFieldKeys(m map[string]fieldType) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func goType(t fieldType) string {
	switch t {
	case tString:
		return "string"
	case tInt:
		return "int"
	case tBool:
		return "bool"
	default:
		return "any"
	}
}

func toExportName(s string) string {
	if s == "" {
		return "X"
	}
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' ' || r == '/'
	})
	if len(parts) == 0 {
		parts = []string{s}
	}
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	out := strings.Join(parts, "")
	if out == "" {
		out = "X"
	}
	if out[0] >= '0' && out[0] <= '9' {
		out = "API" + out
	}
	return out
}
