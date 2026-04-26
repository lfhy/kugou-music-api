package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type subsonicRoot struct {
	Resp struct {
		Status    string `json:"status"`
		Version   string `json:"version"`
		Playlists struct {
			Playlist []playlist `json:"playlist"`
		} `json:"playlists"`
		Playlist playlistDetail `json:"playlist"`
		Error    struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	} `json:"subsonic-response"`
}

type playlist struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SongCount int64  `json:"songCount"`
}

type playlistDetail struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	SongCount int64   `json:"songCount"`
	Entry     []entry `json:"entry"`
}

type entry struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

func main() {
	base := strings.TrimRight(getenv("SUBSONIC_BASE", "http://127.0.0.1:8089"), "/")
	user := getenv("SUBSONIC_USER", "admin")
	pass := getenv("SUBSONIC_PASSWORD", "admin")
	clientName := getenv("SUBSONIC_CLIENT", "subsonic-playlists-example")
	salt := "c19b2d"
	token := md5Hex(pass + salt)

	playlistsURL := fmt.Sprintf("%s/rest/getPlaylists.view?u=%s&t=%s&s=%s&v=1.16.1&c=%s&f=json",
		base, url.QueryEscape(user), token, salt, url.QueryEscape(clientName))

	root, err := getJSON(playlistsURL)
	if err != nil {
		panic(err)
	}
	if root.Resp.Status != "ok" {
		panic(fmt.Sprintf("subsonic failed: code=%d message=%s", root.Resp.Error.Code, root.Resp.Error.Message))
	}

	fmt.Printf("Playlists: %d\n", len(root.Resp.Playlists.Playlist))
	for i, p := range root.Resp.Playlists.Playlist {
		fmt.Printf("%d. %s (%s) songs=%d\n", i+1, p.Name, p.ID, p.SongCount)
		if i >= 4 {
			break
		}
	}
	if len(root.Resp.Playlists.Playlist) == 0 {
		return
	}

	first := root.Resp.Playlists.Playlist[0]
	detailURL := fmt.Sprintf("%s/rest/getPlaylist.view?u=%s&t=%s&s=%s&v=1.16.1&c=%s&f=json&id=%s",
		base, url.QueryEscape(user), token, salt, url.QueryEscape(clientName), url.QueryEscape(first.ID))
	detailRoot, err := getJSON(detailURL)
	if err != nil {
		panic(err)
	}
	if detailRoot.Resp.Status != "ok" {
		panic(fmt.Sprintf("getPlaylist failed: code=%d message=%s", detailRoot.Resp.Error.Code, detailRoot.Resp.Error.Message))
	}

	d := detailRoot.Resp.Playlist
	fmt.Printf("\nPlaylist Detail: %s (%s) songs=%d entries=%d\n", d.Name, d.ID, d.SongCount, len(d.Entry))
	for i, e := range d.Entry {
		fmt.Printf("- %s / %s / %s\n", e.Title, e.Artist, e.Album)
		if i >= 9 {
			break
		}
	}
}

func getJSON(rawURL string) (*subsonicRoot, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var out subsonicRoot
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode failed: %w, body=%s", err, string(body))
	}
	return &out, nil
}

func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func getenv(k, dv string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return dv
}
