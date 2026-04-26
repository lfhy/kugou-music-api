package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
	"github.com/lfhy/kugou-music-api/sdk"
)

func main() {
	mode := flag.String("mode", "heart", "radio mode: heart | new | niche")
	name := flag.String("name", "", "playlist name (optional). default: YYYY-MM-DD<mode>日推")
	limit := flag.Int("limit", 50, "tracks to save (recommended 30-50)")
	private := flag.Bool("private", false, "create private playlist")
	debug := flag.Bool("debug", false, "print debug response")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()

	cfg := session.Load(session.DefaultPath())
	debugEnabled := cfg.Debug || *debug
	if session.HasLoginCookie(cfg.Cookie) == false {
		fmt.Println("未检测到本地登录会话，请先登录（examples/login_session 或 login_cellphone/login_qrcode）。")
		os.Exit(1)
	}

	client, err := sdk.New(sdk.WithCookie(cfg.Cookie))
	if err == nil {
		// continue
	} else {
		fmt.Printf("init sdk failed: %v\n", err)
		os.Exit(1)
	}

	target := normalizeLimit(*limit)
	radioMode := sdk.PersonalRadioMode(strings.TrimSpace(*mode))
	if radioMode == "" {
		radioMode = sdk.PersonalRadioHeart
	}

	plName := *name
	if plName == "" {
		label := map[string]string{"heart": "红心", "new": "新歌", "niche": "小众"}[string(radioMode)]
		if label == "" {
			label = string(radioMode)
		}
		plName = time.Now().Format("2006-01-02") + label + "日推"
	}

	created, err := client.CreatePlaylist(ctx, plName, *private, cfg.Cookie)
	if err == nil {
		// continue
	} else {
		fmt.Printf("create playlist failed: %v\n", err)
		os.Exit(1)
	}
	if debugEnabled {
		printResp("create_playlist", created.Raw)
	}
	fmt.Printf("playlist created: name=%s listid=%d\n", plName, created.ListID)

	seen := map[string]struct{}{}
	totalAdded := 0
	page := 1
	round := 0
	staleRounds := 0
	maxRounds := target*4 + 10
	source := ""

	for totalAdded < target && round < maxRounds {
		round++
		wantThisRound := minInt(5, target-totalAdded)
		req := sdk.PersonalRadioRequest{
			Mode:     radioMode,
			PageSize: wantThisRound,
			Cookie:   cfg.Cookie,
		}
		if radioMode == sdk.PersonalRadioNew {
			req.Page = page
			page++
		} else {
			req.Page = 1
		}

		radio, getErr := client.GetPersonalRadio(ctx, req)
		if getErr == nil {
			// continue
		} else {
			fmt.Printf("get radio round=%d failed: %v\n", round, getErr)
			os.Exit(1)
		}
		if source == "" {
			source = radio.Source
			fmt.Printf("radio: mode=%s source=%s target=%d\n", radio.Mode, radio.Source, target)
		}
		if debugEnabled {
			printResp(fmt.Sprintf("radio_round_%d", round), radio.Raw)
		}

		batch := pickUniqueTracks(radio.Tracks, seen, wantThisRound)
		if len(batch) == 0 {
			staleRounds++
			if staleRounds >= 8 {
				fmt.Printf("no new tracks after %d rounds, stop early\n", staleRounds)
				break
			}
			continue
		}
		staleRounds = 0

		added, addErr := client.AddTracksToPlaylist(ctx, created.ListID, batch, cfg.Cookie)
		if debugEnabled {
			printResp(fmt.Sprintf("add_round_%d", round), nilIfNilAdded(added))
		}
		if addErr == nil {
			// continue
		} else {
			fmt.Printf("add tracks round=%d failed: %v\n", round, addErr)
			os.Exit(1)
		}
		if added.Raw == nil || added.Raw.Body == nil {
			fmt.Printf("round=%d pulled=%d unique=%d added=%d\n", round, len(radio.Tracks), len(batch), added.Added)
		} else {
			fmt.Printf("round=%d pulled=%d unique=%d added=%d biz_status=%v error_code=%v\n",
				round, len(radio.Tracks), len(batch), added.Added, added.Raw.Body["status"], added.Raw.Body["error_code"])
			if isBizSuccess(added.Raw) == false {
				fmt.Println("add tracks business failed")
				os.Exit(1)
			}
		}

		totalAdded += added.Added
	}

	fmt.Printf("tracks added total: %d/%d\n", totalAdded, target)
	if totalAdded < target {
		fmt.Println("warning: did not reach target limit, you can rerun to补齐更多歌曲")
	}
}

func pickUniqueTracks(tracks []sdk.RadioTrack, seen map[string]struct{}, limit int) []sdk.RadioTrack {
	out := make([]sdk.RadioTrack, 0, minInt(len(tracks), limit))
	for _, t := range tracks {
		h := strings.ToLower(strings.TrimSpace(t.Hash))
		if h == "" {
			continue
		}
		_, exists := seen[h]
		if exists {
			continue
		}
		seen[h] = struct{}{}
		out = append(out, t)
		if len(out) >= limit {
			break
		}
	}
	return out
}

func normalizeLimit(v int) int {
	if v <= 0 {
		return 30
	}
	if v > 50 {
		return 50
	}
	return v
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func isBizSuccess(resp *sdk.Response) bool {
	if resp == nil || resp.Body == nil {
		return false
	}
	status := strings.TrimSpace(fmt.Sprintf("%v", resp.Body["status"]))
	errorCode := strings.TrimSpace(fmt.Sprintf("%v", resp.Body["error_code"]))
	if status == "1" || status == "1.0" {
		if errorCode == "" || errorCode == "<nil>" || errorCode == "0" || errorCode == "0.0" {
			return true
		}
	}
	return false
}

func printResp(label string, resp *sdk.Response) {
	if resp == nil {
		fmt.Printf("[debug] %s: nil response\n", label)
		return
	}
	fmt.Printf("[debug] %s: http_status=%d\n", label, resp.Status)
	fmt.Printf("[debug] %s: raw=%s\n", label, string(resp.RawBody))
}

func nilIfNilAdded(r *sdk.PlaylistAddTracksResult) *sdk.Response {
	if r == nil {
		return nil
	}
	return r.Raw
}
