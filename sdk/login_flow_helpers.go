package sdk

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/lfhy/kugou-music-api/core/config"
)

// dedupSetCookie keeps only the last value for each cookie key.
func dedupSetCookie(input []string) []string {
	m := map[string]string{}
	for _, x := range input {
		x = strings.TrimSpace(x)
		if x == "" {
			continue
		}
		seg := strings.SplitN(x, ";", 2)[0]
		kv := strings.SplitN(seg, "=", 2)
		if len(kv) != 2 {
			continue
		}
		m[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	out := make([]string, 0, len(m))
	for k, v := range m {
		out = append(out, k+"="+v)
	}
	return out
}

func firstAny(v any, def any) any {
	if v == nil {
		return def
	}
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "" || s == "<nil>" {
		return def
	}
	return v
}

func asString(v any) string {
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "<nil>" {
		return ""
	}
	return s
}

func asIntString(v any) string {
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "" || s == "<nil>" {
		return "0"
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return strconv.FormatInt(i, 10)
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return strconv.FormatInt(int64(f), 10)
	}
	return "0"
}

func ternaryString(ok bool, a, b string) string {
	if ok {
		return a
	}
	return b
}

func maskMobile(mobile string) string {
	mobile = strings.TrimSpace(mobile)
	if len(mobile) < 7 {
		return mobile
	}
	return mobile[:3] + "****" + mobile[len(mobile)-4:]
}

func signParamsKey(data string, isLite bool) string {
	appid, clientver := config.PlatformConfig(isLite)
	secret := "OIlwieks28dk2k092lksi2UIkp"
	if isLite {
		secret = "LnT6xpN3khm36zse0QzvmgTZ3waWdRSA"
	}
	sum := md5.Sum([]byte(appid + secret + clientver + data))
	return hex.EncodeToString(sum[:])
}
