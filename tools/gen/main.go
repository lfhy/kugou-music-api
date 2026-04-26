package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/lfhy/kugou-music-api/core/parser"
)

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
