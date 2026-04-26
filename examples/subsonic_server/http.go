package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// HTTP entrypoints, auth, and access logging stay close together.
func (s *server) handleRest(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	name := endpointName(r.URL.Path)
	log.Printf("rest request: method=%s path=%s endpoint=%s query=%s", r.Method, r.URL.Path, name, r.URL.RawQuery)
	if ok, code, msg := s.auth(r); !ok {
		s.writeError(w, r, code, msg)
		return
	}

	switch name {
	case "ping":
		s.writeOK(w, r, nil, map[string]any{})
	case "getLicense":
		x := license{Valid: true, Email: "kugou@local"}
		s.writeOK(w, r, x, map[string]any{"license": map[string]any{"valid": true, "email": x.Email}})
	case "getMusicFolders":
		payload := musicFolders{Folders: []musicFolder{{ID: 1, Name: "KuGou"}}}
		s.writeOK(w, r, payload, map[string]any{"musicFolders": map[string]any{"musicFolder": []map[string]any{{"id": 1, "name": "KuGou"}}}})
	default:
		s.writeError(w, r, 70, fmt.Sprintf("Not implemented: %s", name))
	}
}

func (s *server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func endpointName(path string) string {
	name := strings.TrimPrefix(path, "/rest/")
	name = strings.TrimSuffix(name, ".view")
	return strings.TrimSpace(name)
}

func (s *server) auth(r *http.Request) (bool, int, string) {
	u := r.Form.Get("u")
	t := r.Form.Get("t")
	salt := r.Form.Get("s")
	v := r.Form.Get("v")
	c := r.Form.Get("c")
	if u == "" || t == "" || salt == "" || v == "" || c == "" {
		return false, 10, "Missing required auth parameters"
	}
	if u != s.user {
		return false, 40, "Wrong username or password"
	}
	if !strings.EqualFold(md5Hex(s.pass+salt), t) {
		return false, 40, "Wrong username or password"
	}
	return true, 0, ""
}

func (s *server) writeOK(w http.ResponseWriter, r *http.Request, xmlPayload any, jsonFields map[string]any) {
	if strings.EqualFold(r.Form.Get("f"), "json") {
		payload := map[string]any{
			"subsonic-response": map[string]any{
				"status":  "ok",
				"version": subsonicVersion,
			},
		}
		for k, v := range jsonFields {
			payload["subsonic-response"].(map[string]any)[k] = v
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(payload)
		return
	}
	env := xmlEnvelope{Xmlns: xmlNS, Status: "ok", Version: subsonicVersion, Payload: xmlPayload}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(xml.Header))
	_ = xml.NewEncoder(w).Encode(env)
}

func (s *server) writeError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	if strings.EqualFold(r.Form.Get("f"), "json") {
		payload := map[string]any{
			"subsonic-response": map[string]any{
				"status":  "failed",
				"version": subsonicVersion,
				"error":   map[string]any{"code": code, "message": msg},
			},
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(payload)
		return
	}
	env := xmlEnvelope{Xmlns: xmlNS, Status: "failed", Version: subsonicVersion, Error: &subError{Code: code, Message: msg}}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(xml.Header))
	_ = xml.NewEncoder(w).Encode(env)
}
