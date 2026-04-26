package parser

import (
	"regexp"
	"strconv"
	"strings"
)

type CompatExprKind string

const (
	ExprLiteralString CompatExprKind = "literal_string"
	ExprLiteralInt    CompatExprKind = "literal_int"
	ExprLiteralBool   CompatExprKind = "literal_bool"
	ExprParam         CompatExprKind = "param"
	ExprCookie        CompatExprKind = "cookie"
	ExprParamOrString CompatExprKind = "param_or_string"
	ExprParamOrInt    CompatExprKind = "param_or_int"
	ExprParamOrCookie CompatExprKind = "param_or_cookie_or_default"
	ExprTemplateParam CompatExprKind = "template_param"
	ExprNumberOrInt   CompatExprKind = "number_param_or_int"
)

type CompatExpr struct {
	Kind         CompatExprKind
	ParamKey     string
	CookieKey    string
	DefaultType  string
	DefaultStr   string
	DefaultInt   int
	DefaultBool  bool
	Stringify    bool
	OriginalExpr string
}

type CompatRule struct {
	Field string
	Expr  CompatExpr
}

type CompatSpec struct {
	Identifier       string
	DataMapRules     []CompatRule
	ParamsMapRules   []CompatRule
	CookieRules      []CompatRule
	UnsupportedExprs []string
}

var (
	reMapBlock = regexp.MustCompile(`(?s)const\s+(dataMap|paramsMap)\s*=\s*\{(.*?)\n\s*\};`)
	reKVLine   = regexp.MustCompile(`^\s*['"]?([a-zA-Z_][a-zA-Z0-9_]*)['"]?\s*:\s*(.+?)\s*,?\s*$`)

	reLitString      = regexp.MustCompile("^['\"`](.*)['\"`]$")
	reLitInt         = regexp.MustCompile(`^([0-9]+)$`)
	reLitBool        = regexp.MustCompile(`^(true|false)$`)
	reParam          = regexp.MustCompile(`^params\?\.([a-zA-Z_][a-zA-Z0-9_]*)$`)
	reCookie         = regexp.MustCompile(`^params\?\.cookie\?\.([a-zA-Z_][a-zA-Z0-9_]*)$`)
	reParamOrString  = regexp.MustCompile("^params\\?\\.([a-zA-Z_][a-zA-Z0-9_]*)\\s*\\|\\|\\s*['\"`](.*)['\"`]$")
	reParamOrInt     = regexp.MustCompile(`^params\?\.([a-zA-Z_][a-zA-Z0-9_]*)\s*\|\|\s*([0-9]+)$`)
	reParamOrCookie  = regexp.MustCompile("^params\\?\\.([a-zA-Z_][a-zA-Z0-9_]*)\\s*\\|\\|\\s*params\\?\\.cookie\\?\\.([a-zA-Z_][a-zA-Z0-9_]*)\\s*\\|\\|\\s*(['\"`].*['\"`]|[0-9]+|true|false)$")
	reTemplateParam  = regexp.MustCompile("^`\\$\\{params\\?\\.([a-zA-Z_][a-zA-Z0-9_]*)\\}`$")
	reNumberParamInt = regexp.MustCompile(`^Number\(params\?\.([a-zA-Z_][a-zA-Z0-9_]*)\)\s*\|\|\s*([0-9]+)$`)
)

func ExtractCompatSpec(identifier, src string) CompatSpec {
	out := CompatSpec{Identifier: identifier}

	for _, m := range reMapBlock.FindAllStringSubmatch(src, -1) {
		if len(m) != 3 {
			continue
		}
		name := m[1]
		body := m[2]
		rules, unsupported := parseObjectRules(body)
		if name == "dataMap" {
			out.DataMapRules = rules
		} else {
			out.ParamsMapRules = rules
		}
		out.UnsupportedExprs = append(out.UnsupportedExprs, unsupported...)
	}

	out.CookieRules = parseCookieRules(src)
	return out
}

func parseObjectRules(body string) ([]CompatRule, []string) {
	rules := make([]CompatRule, 0)
	unsupported := make([]string, 0)
	for _, line := range strings.Split(body, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "//") {
			continue
		}
		m := reKVLine.FindStringSubmatch(t)
		if len(m) != 3 {
			continue
		}
		field := strings.TrimSpace(m[1])
		expr := strings.TrimSpace(strings.TrimSuffix(m[2], ","))
		ce, ok := parseExpr(expr)
		if !ok {
			unsupported = append(unsupported, field+": "+expr)
			continue
		}
		rules = append(rules, CompatRule{Field: field, Expr: ce})
	}
	return rules, unsupported
}

func parseExpr(expr string) (CompatExpr, bool) {
	e := CompatExpr{OriginalExpr: expr}
	if m := reLitString.FindStringSubmatch(expr); len(m) == 2 {
		e.Kind = ExprLiteralString
		e.DefaultType = "string"
		e.DefaultStr = m[1]
		return e, true
	}
	if m := reLitInt.FindStringSubmatch(expr); len(m) == 2 {
		n, _ := strconv.Atoi(m[1])
		e.Kind = ExprLiteralInt
		e.DefaultType = "int"
		e.DefaultInt = n
		return e, true
	}
	if m := reLitBool.FindStringSubmatch(expr); len(m) == 2 {
		e.Kind = ExprLiteralBool
		e.DefaultType = "bool"
		e.DefaultBool = m[1] == "true"
		return e, true
	}
	if m := reTemplateParam.FindStringSubmatch(expr); len(m) == 2 {
		e.Kind = ExprTemplateParam
		e.ParamKey = m[1]
		e.Stringify = true
		return e, true
	}
	if m := reParamOrString.FindStringSubmatch(expr); len(m) == 3 {
		e.Kind = ExprParamOrString
		e.ParamKey = m[1]
		e.DefaultType = "string"
		e.DefaultStr = m[2]
		return e, true
	}
	if m := reParamOrInt.FindStringSubmatch(expr); len(m) == 3 {
		n, _ := strconv.Atoi(m[2])
		e.Kind = ExprParamOrInt
		e.ParamKey = m[1]
		e.DefaultType = "int"
		e.DefaultInt = n
		return e, true
	}
	if m := reParamOrCookie.FindStringSubmatch(expr); len(m) == 4 {
		e.Kind = ExprParamOrCookie
		e.ParamKey = m[1]
		e.CookieKey = m[2]
		d := strings.TrimSpace(m[3])
		if mm := reLitString.FindStringSubmatch(d); len(mm) == 2 {
			e.DefaultType = "string"
			e.DefaultStr = mm[1]
			return e, true
		}
		if mm := reLitInt.FindStringSubmatch(d); len(mm) == 2 {
			n, _ := strconv.Atoi(mm[1])
			e.DefaultType = "int"
			e.DefaultInt = n
			return e, true
		}
		if mm := reLitBool.FindStringSubmatch(d); len(mm) == 2 {
			e.DefaultType = "bool"
			e.DefaultBool = mm[1] == "true"
			return e, true
		}
		return CompatExpr{}, false
	}
	if m := reNumberParamInt.FindStringSubmatch(expr); len(m) == 3 {
		n, _ := strconv.Atoi(m[2])
		e.Kind = ExprNumberOrInt
		e.ParamKey = m[1]
		e.DefaultType = "int"
		e.DefaultInt = n
		return e, true
	}
	if m := reCookie.FindStringSubmatch(expr); len(m) == 2 {
		e.Kind = ExprCookie
		e.CookieKey = m[1]
		return e, true
	}
	if m := reParam.FindStringSubmatch(expr); len(m) == 2 {
		e.Kind = ExprParam
		e.ParamKey = m[1]
		return e, true
	}
	return CompatExpr{}, false
}

func parseCookieRules(src string) []CompatRule {
	cookieRules := make([]CompatRule, 0)
	reCookieBlock := regexp.MustCompile(`(?s)cookie\s*:\s*\{(.*?)\}`)
	reCookieKV := regexp.MustCompile(`['"]?([a-zA-Z_][a-zA-Z0-9_]*)['"]?\s*:\s*params\?\.cookie\?\.([a-zA-Z_][a-zA-Z0-9_]*)`)
	m := reCookieBlock.FindStringSubmatch(src)
	if len(m) != 2 {
		return cookieRules
	}
	for _, kv := range reCookieKV.FindAllStringSubmatch(m[1], -1) {
		if len(kv) == 3 {
			cookieRules = append(cookieRules, CompatRule{Field: kv[1], Expr: CompatExpr{Kind: ExprCookie, CookieKey: kv[2], OriginalExpr: "params?.cookie?." + kv[2]}})
		}
	}
	return cookieRules
}
