package sdk

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"encoding/base64"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	krcXORKey         = []byte{64, 71, 97, 119, 94, 50, 116, 71, 81, 54, 49, 45, 206, 210, 110, 105}
	krcLineTagPattern = regexp.MustCompile(`^\[(\d+),\d+\]`)
	krcWordTagPattern = regexp.MustCompile(`<\d+,\d+,\d+>`)
	lrcMetaTagPattern = regexp.MustCompile(`^\[[a-zA-Z]+:.*\]$`)
)

// DecodedContent returns decoded lyric text.
// For lrc/plain content it returns plain text, and for krc content it returns decoded krc text.
func (r *LyricResponse) DecodedContent() string {
	if r == nil || r.Body == nil {
		return ""
	}
	if cached := strings.TrimSpace(asString(r.Body["decodeContent"])); cached != "" {
		return cached
	}

	content := strings.TrimSpace(asString(r.Body["content"]))
	if content == "" {
		return ""
	}

	contentType := toInt(firstAny(r.Body["contenttype"], r.Body["contentType"]), 0)
	decoded := ""
	if contentType != 0 {
		decoded = decodeBase64LyricText(content)
	} else {
		decoded = decodeKRCLyricText(content)
		if decoded == "" {
			decoded = decodeBase64LyricText(content)
		}
	}
	if decoded != "" {
		r.Body["decodeContent"] = decoded
	}
	return decoded
}

// ToLrc converts lyric response content to common LRC text for external players/software.
func (r *LyricResponse) ToLrc() string {
	decoded := normalizeLyricNewlines(r.DecodedContent())
	if decoded == "" {
		return ""
	}
	if !containsKRCLine(decoded) {
		return decoded
	}

	lines := strings.Split(decoded, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		switch {
		case line == "":
			out = append(out, "")
		case lrcMetaTagPattern.MatchString(line):
			out = append(out, line)
		case krcLineTagPattern.MatchString(line):
			match := krcLineTagPattern.FindStringSubmatch(line)
			startMs, _ := strconv.Atoi(match[1])
			text := krcWordTagPattern.ReplaceAllString(line[len(match[0]):], "")
			out = append(out, formatLRCTimestamp(startMs)+text)
		default:
			out = append(out, krcWordTagPattern.ReplaceAllString(line, ""))
		}
	}

	return strings.TrimRight(strings.Join(out, "\n"), "\n")
}

func containsKRCLine(text string) bool {
	for _, line := range strings.Split(normalizeLyricNewlines(text), "\n") {
		if krcLineTagPattern.MatchString(line) {
			return true
		}
	}
	return false
}

func normalizeLyricNewlines(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	return text
}

func formatLRCTimestamp(ms int) string {
	if ms < 0 {
		ms = 0
	}
	totalSeconds := ms / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	centiseconds := (ms % 1000) / 10
	return "[" + pad2(minutes) + ":" + pad2(seconds) + "." + pad2(centiseconds) + "]"
}

func pad2(v int) string {
	if v < 10 {
		return "0" + strconv.Itoa(v)
	}
	return strconv.Itoa(v)
}

func decodeBase64LyricText(content string) string {
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(content))
	if err != nil {
		return ""
	}
	return string(raw)
}

func decodeKRCLyricText(content string) string {
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(content))
	if err != nil || len(raw) <= 4 {
		return ""
	}

	payload := append([]byte(nil), raw[4:]...)
	for i := range payload {
		payload[i] ^= krcXORKey[i%len(krcXORKey)]
	}

	if text := inflateLyricPayload(payload, true); text != "" {
		return text
	}
	return inflateLyricPayload(payload, false)
}

func inflateLyricPayload(payload []byte, useZlib bool) string {
	var reader io.ReadCloser
	if useZlib {
		zr, err := zlib.NewReader(bytes.NewReader(payload))
		if err != nil {
			return ""
		}
		reader = zr
	} else {
		reader = io.NopCloser(flate.NewReader(bytes.NewReader(payload)))
	}
	defer reader.Close()

	plain, err := io.ReadAll(reader)
	if err != nil {
		return ""
	}
	return string(plain)
}
