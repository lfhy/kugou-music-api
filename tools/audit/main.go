package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lfhy/kugou-music-api/core/parser"
)

type item struct {
	Identifier string
	Route      string
	Risk       string
	Reasons    []string
}

func main() {
	moduleDir := os.Getenv("KUGOU_MODULE_DIR")
	if strings.TrimSpace(moduleDir) == "" {
		moduleDir = "../module"
	}

	specs, err := parser.LoadModuleSpecs(moduleDir)
	if err != nil {
		panic(err)
	}

	items := make([]item, 0, len(specs))
	for _, sp := range specs {
		p := filepath.Join(moduleDir, sp.Identifier+".js")
		b, err := os.ReadFile(p)
		if err != nil {
			panic(err)
		}
		s := string(b)
		reasons := analyze(s, sp)
		risk := "LOW"
		if len(reasons) > 0 {
			risk = "HIGH"
		}
		items = append(items, item{Identifier: sp.Identifier, Route: sp.Route, Risk: risk, Reasons: reasons})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Route < items[j].Route })

	out := "sdk/API_COMPAT_AUDIT.md"
	f, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	fmt.Fprintln(w, "# API Compatibility Audit")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "说明：该报告逐个模块检查 Go 自动生成调用与原 JS 逻辑的一致性风险。")
	fmt.Fprintln(w, "`HIGH` 表示该接口很可能需要手写适配层（默认值/参数重组/加解密/后处理）。")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "| Identifier | Route | Risk | Reasons |")
	fmt.Fprintln(w, "| --- | --- | --- | --- |")

	high := 0
	for _, it := range items {
		reason := "-"
		if len(it.Reasons) > 0 {
			reason = strings.Join(it.Reasons, "; ")
		}
		fmt.Fprintf(w, "| `%s` | `%s` | `%s` | %s |\n", it.Identifier, it.Route, it.Risk, reason)
		if it.Risk == "HIGH" {
			high++
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "总计: %d, HIGH: %d, LOW: %d\n", len(items), high, len(items)-high)
}

func analyze(src string, sp parser.ModuleSpec) []string {
	reasons := []string{}
	add := func(s string) {
		for _, x := range reasons {
			if x == s {
				return
			}
		}
		reasons = append(reasons, s)
	}

	if strings.Contains(src, "Date.now(") || strings.Contains(src, "new Date(") || strings.Contains(src, "Math.floor(") {
		add("dynamic timestamp")
	}
	if strings.Contains(src, "randomString(") || strings.Contains(src, "randomNumber(") {
		add("random/default device fields")
	}
	if strings.Contains(src, "crypto") || strings.Contains(src, "sign") || strings.Contains(src, "rsa") || strings.Contains(src, "AES") || strings.Contains(src, "playlistAes") {
		add("custom crypto/sign logic")
	}
	if strings.Contains(src, "new Promise") || strings.Contains(src, ".then(") || strings.Contains(src, ".catch(") {
		add("custom response post-processing")
	}
	if strings.Contains(src, "decodeLyrics(") || strings.Contains(src, "Buffer.from(") {
		add("custom body decoding")
	}
	if strings.Contains(src, "||") && strings.Contains(src, "params") {
		add("default fallback params in module")
	}
	if strings.Contains(src, "const dataMap") || strings.Contains(src, "const paramsMap") || strings.Contains(src, "let dataMap") {
		add("manual request map assembly")
	}
	if strings.Contains(src, "cookie:") && strings.Contains(src, "params?.cookie") == false {
		add("custom cookie mapping")
	}
	if strings.Contains(sp.UpstreamURL, "${") {
		add("dynamic template URL")
	}

	return reasons
}
