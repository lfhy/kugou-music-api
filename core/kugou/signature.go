package kugou

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
)

func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func SignatureWebParams(params map[string]any) string {
	key := "NVPh5oo715z5DIWAeQlhMDsWXXQV4hwt"
	keys := sortedKeys(params)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+toSigVal(params[k]))
	}
	body := strings.Join(parts, "")
	return md5Hex(key + body + key)
}

func SignatureAndroidParams(params map[string]any, data string, isLite bool) string {
	key := "OIlwieks28dk2k092lksi2UIkp"
	if isLite {
		key = "LnT6xpN3khm36zse0QzvmgTZ3waWdRSA"
	}

	keys := sortedKeys(params)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+toSigVal(params[k]))
	}
	body := strings.Join(parts, "")
	return md5Hex(key + body + data + key)
}

func SignatureRegisterParams(params map[string]any) string {
	vals := make([]string, 0, len(params))
	for _, v := range params {
		vals = append(vals, toSigVal(v))
	}
	sort.Strings(vals)
	return md5Hex("1014" + strings.Join(vals, "") + "1014")
}

func SignKey(hash, mid, userid, appid string, isLite bool) string {
	key := "57ae12eb6890223e355ccfcb74edf70d"
	if isLite {
		key = "185672dd44712f60bb1736df5a377e82"
	}
	if userid == "" {
		userid = "0"
	}
	return md5Hex(hash + key + appid + mid + userid)
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func toSigVal(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(b)
	}
}
