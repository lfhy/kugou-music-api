package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lfhy/kugou-music-api/examples/shared/session"
	"github.com/lfhy/kugou-music-api/sdk"
)

// main wires the KuGou SDK, session loading, and HTTP server startup together.
func main() {
	kg, err := sdk.New()
	if err != nil {
		log.Fatalf("init kugou sdk failed: %v", err)
	}
	sessionCfg := session.Load(session.DefaultPath())
	if session.HasLoginCookie(sessionCfg.Cookie) {
		for k, v := range sessionCfg.Cookie {
			kg.SetCookie(k, v)
		}
		log.Printf("loaded kugou session: userid=%s", sessionCfg.Cookie["userid"])
	} else {
		log.Printf("no valid kugou session found at %s", session.DefaultPath())
	}

	s := &server{
		kg:       kg,
		user:     getenv("SUBSONIC_USER", "admin"),
		pass:     getenv("SUBSONIC_PASSWORD", "admin"),
		addr:     getenv("SUBSONIC_ADDR", ":8089"),
		dataFile: getenv("SUBSONIC_DATA_FILE", "./examples/subsonic_server/data/state.json"),
		plMode:   strings.ToLower(strings.TrimSpace(getenv("SUBSONIC_PLAYLIST_MODE", "user"))),
		httpCli:  &http.Client{Timeout: 60 * time.Second},
		songMap:  map[string]songMeta{},
		cover:    map[string]string{},
		plMap:    map[string]playlistMeta{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/rest/", s.handleRest)
	mux.HandleFunc("/rest", s.handleRest)
	mux.HandleFunc("/", s.handleNotFound)

	log.Printf("subsonic server listening on %s", s.addr)
	log.Printf("default user=%s password=%s", s.user, s.pass)
	log.Printf("state file: %s", s.dataFile)
	log.Fatal(http.ListenAndServe(s.addr, mux))
}
