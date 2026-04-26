package sdk

import "testing"

func TestLoginQrKeyResponseQRCodeURL(t *testing.T) {
	resp := &LoginQrKeyResponse{
		Body: map[string]any{
			"status": 1,
			"data": map[string]any{
				"qrcode": "abc123",
			},
		},
	}

	if got := resp.QRCodeKey(); got != "abc123" {
		t.Fatalf("QRCodeKey() = %q, want %q", got, "abc123")
	}
	wantURL := "https://h5.kugou.com/apps/loginQRCode/html/index.html?qrcode=abc123"
	if got := resp.QRCodeURL(); got != wantURL {
		t.Fatalf("QRCodeURL() = %q, want %q", got, wantURL)
	}
}

func TestBuildLoginQrCreateResponse(t *testing.T) {
	resp, err := buildLoginQrCreateResponse(LoginQrCreateRequest{
		Key:   "abc123",
		Qrimg: true,
	})
	if err != nil {
		t.Fatalf("buildLoginQrCreateResponse() error = %v", err)
	}
	wantURL := "https://h5.kugou.com/apps/loginQRCode/html/index.html?qrcode=abc123"
	if got := resp.URL(); got != wantURL {
		t.Fatalf("URL() = %q, want %q", got, wantURL)
	}
	if got := resp.Base64(); len(got) == 0 {
		t.Fatal("Base64() returned empty string")
	}
}

func TestLoginQrCheckResponseStatusCodeUsesNestedData(t *testing.T) {
	resp := &LoginQrCheckResponse{
		Body: map[string]any{
			"status": 1,
			"data": map[string]any{
				"status": 4,
				"token":  "tok",
				"userid": "123",
			},
		},
	}

	if got := resp.StatusCode(); got != 4 {
		t.Fatalf("StatusCode() = %d, want %d", got, 4)
	}
	if got := resp.Token(); got != "tok" {
		t.Fatalf("Token() = %q, want %q", got, "tok")
	}
	if got := resp.UserID(); got != "123" {
		t.Fatalf("UserID() = %q, want %q", got, "123")
	}
}
