package main

import "testing"

func TestNormalizedUserID(t *testing.T) {
	if got := normalizedUserID("5.1815528e+07"); got != "51815528" {
		t.Fatalf("normalizedUserID() = %q", got)
	}
}

func TestLatestSignedDate(t *testing.T) {
	body := map[string]any{
		"data": map[string]any{
			"list": []any{
				map[string]any{"day": "2026-04-25", "receive_vip": 1},
				map[string]any{"day": "2026-04-26", "receive_vip": 1},
			},
		},
	}
	if got := latestSignedDate(body); got != "2026-04-26" {
		t.Fatalf("latestSignedDate() = %q", got)
	}
}

func TestLatestVIPEndTime(t *testing.T) {
	body := map[string]any{
		"data": map[string]any{
			"busi_vip": []any{
				map[string]any{"vip_end_time": "2026-12-17 22:15:41"},
				map[string]any{"vip_end_time": "2026-12-18 22:15:41"},
			},
		},
	}
	if got := latestVIPEndTime(body); got != "2026-12-18 22:15:41" {
		t.Fatalf("latestVIPEndTime() = %q", got)
	}
}
