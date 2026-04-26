package main

import (
	"sort"
	"strings"
)

// Field extraction helpers keep the generator's inferred request models consistent.
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
