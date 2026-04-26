package parser

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type ModuleSpec struct {
	Identifier         string
	Route              string
	UpstreamURL        string
	Method             string
	BaseURL            string
	EncryptType        string
	EncryptKey         bool
	ClearDefaultParams bool
	NotSignature       bool
	Headers            map[string]string
	UseParams          bool
	UseData            bool
}

var (
	reURL         = regexp.MustCompile(`(?s)\burl\s*:\s*["'` + "`" + `](.+?)["'` + "`" + `]`)
	reMethod      = regexp.MustCompile(`(?s)\bmethod\s*:\s*["'` + "`" + `](.+?)["'` + "`" + `]`)
	reBaseURL     = regexp.MustCompile(`(?s)\bbaseURL\s*:\s*["'` + "`" + `](.+?)["'` + "`" + `]`)
	reEncryptType = regexp.MustCompile(`(?s)\bencryptType\s*:\s*["'` + "`" + `](.+?)["'` + "`" + `]`)
	reHeadersBlk  = regexp.MustCompile(`(?s)\bheaders\s*:\s*\{(.*?)\}`)
	reHeaderKV    = regexp.MustCompile(`["']([^"']+)["']\s*:\s*["']([^"']*)["']`)
)

func LoadModuleSpecs(moduleDir string) ([]ModuleSpec, error) {
	entries, err := os.ReadDir(moduleDir)
	if err != nil {
		return nil, err
	}

	specs := make([]ModuleSpec, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".js") || strings.HasPrefix(name, "_") {
			continue
		}

		p := filepath.Join(moduleDir, name)
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		src := string(data)

		id := strings.TrimSuffix(name, ".js")
		spec := ModuleSpec{
			Identifier:  id,
			Route:       "/" + strings.ReplaceAll(id, "_", "/"),
			UpstreamURL: "/",
			Method:      "GET",
			EncryptType: "android",
			Headers:     map[string]string{},
			UseParams:   strings.Contains(src, "params:"),
			UseData:     strings.Contains(src, "data:"),
		}

		if m := reURL.FindStringSubmatch(src); len(m) == 2 {
			spec.UpstreamURL = strings.TrimSpace(m[1])
		}
		if strings.Contains(spec.UpstreamURL, "${") && id == "search" {
			spec.UpstreamURL = "/v3/search/song"
		}
		if m := reMethod.FindStringSubmatch(src); len(m) == 2 {
			spec.Method = strings.ToUpper(strings.TrimSpace(m[1]))
		}
		if m := reBaseURL.FindStringSubmatch(src); len(m) == 2 {
			spec.BaseURL = strings.TrimSpace(m[1])
		}
		if m := reEncryptType.FindStringSubmatch(src); len(m) == 2 {
			spec.EncryptType = strings.TrimSpace(m[1])
		}

		spec.EncryptKey = strings.Contains(src, "encryptKey: true")
		spec.ClearDefaultParams = strings.Contains(src, "clearDefaultParams: true")
		// Keep behavior aligned with upstream JS createRequest:
		// it only checks `notSignature`, not `notSign`.
		spec.NotSignature = strings.Contains(src, "notSignature: true")

		if m := reHeadersBlk.FindStringSubmatch(src); len(m) == 2 {
			for _, kv := range reHeaderKV.FindAllStringSubmatch(m[1], -1) {
				if len(kv) == 3 {
					spec.Headers[kv[1]] = kv[2]
				}
			}
		}

		specs = append(specs, spec)
	}

	sort.Slice(specs, func(i, j int) bool {
		return specs[i].Route < specs[j].Route
	})
	return specs, nil
}
