package sdk

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"testing"
)

func TestLyricResponseDecodedContentLRC(t *testing.T) {
	plain := "[ti:Test]\n[00:01.23]hello world"
	resp := &LyricResponse{
		Body: map[string]any{
			"content":     base64.StdEncoding.EncodeToString([]byte(plain)),
			"contenttype": 1,
		},
	}

	got := resp.DecodedContent()
	if got != plain {
		t.Fatalf("DecodedContent() = %q, want %q", got, plain)
	}
	if gotLRC := resp.ToLrc(); gotLRC != plain {
		t.Fatalf("ToLrc() = %q, want %q", gotLRC, plain)
	}
	if cached := asString(resp.Body["decodeContent"]); cached != plain {
		t.Fatalf("decodeContent cache = %q, want %q", cached, plain)
	}
}

func TestLyricResponseToLrcFromKRC(t *testing.T) {
	krcPlain := "[ti:Test Song]\n[ar:Test Singer]\n[1230,900]<0,450,0>你<450,450,0>好\n[3450,1200]<0,600,0>世<600,600,0>界"
	resp := &LyricResponse{
		Body: map[string]any{
			"content":     encodeKRCTestPayload(t, krcPlain),
			"contenttype": 0,
		},
	}

	wantDecoded := krcPlain
	if got := resp.DecodedContent(); got != wantDecoded {
		t.Fatalf("DecodedContent() = %q, want %q", got, wantDecoded)
	}

	wantLRC := "[ti:Test Song]\n[ar:Test Singer]\n[00:01.23]你好\n[00:03.45]世界"
	if got := resp.ToLrc(); got != wantLRC {
		t.Fatalf("ToLrc() = %q, want %q", got, wantLRC)
	}
}

func encodeKRCTestPayload(t *testing.T, plain string) string {
	t.Helper()

	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write([]byte(plain)); err != nil {
		t.Fatalf("zlib write failed: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zlib close failed: %v", err)
	}

	payload := buf.Bytes()
	for i := range payload {
		payload[i] ^= krcXORKey[i%len(krcXORKey)]
	}

	raw := append([]byte("krc1"), payload...)
	return base64.StdEncoding.EncodeToString(raw)
}
