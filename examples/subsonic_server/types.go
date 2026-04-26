package main

import (
	"encoding/xml"
	"net/http"
	"sync"

	"github.com/lfhy/kugou-music-api/sdk"
)

// Subsonic XML models and in-memory server state are defined here.
const (
	subsonicVersion = "1.16.1"
	xmlNS           = "http://subsonic.org/restapi"
)

type server struct {
	kg       *sdk.Client
	user     string
	pass     string
	addr     string
	dataFile string
	plMode   string
	httpCli  *http.Client

	mu      sync.RWMutex
	songMap map[string]songMeta
	cover   map[string]string
	plMap   map[string]playlistMeta
}

type songMeta struct {
	SongID       string
	AlbumAudioID int64
	AlbumID      int64
	ArtistID     int64
	Hash         string
	Title        string
	Album        string
	Artist       string
	DurationSec  int64
	BitRate      int64
	Size         int64
	Suffix       string
	CoverURL     string
	CoverID      string
	Track        int64
	Year         int64
}

type subError struct {
	Code    int    `xml:"code,attr" json:"code"`
	Message string `xml:"message,attr" json:"message"`
}

type xmlEnvelope struct {
	XMLName xml.Name  `xml:"subsonic-response"`
	Xmlns   string    `xml:"xmlns,attr"`
	Status  string    `xml:"status,attr"`
	Version string    `xml:"version,attr"`
	Error   *subError `xml:"error,omitempty"`
	Payload any       `xml:",any,omitempty"`
}

type license struct {
	XMLName xml.Name `xml:"license"`
	Valid   bool     `xml:"valid,attr" json:"valid"`
	Email   string   `xml:"email,attr,omitempty" json:"email,omitempty"`
}

type musicFolders struct {
	XMLName xml.Name      `xml:"musicFolders"`
	Folders []musicFolder `xml:"musicFolder" json:"musicFolder"`
}

type musicFolder struct {
	ID   int64  `xml:"id,attr" json:"id"`
	Name string `xml:"name,attr" json:"name"`
}

type indexes struct {
	XMLName   xml.Name   `xml:"indexes"`
	LastMod   int64      `xml:"lastModified,attr" json:"lastModified"`
	Ignored   string     `xml:"ignoredArticles,attr,omitempty" json:"ignoredArticles,omitempty"`
	IndexList []idxGroup `xml:"index" json:"index"`
}

type idxGroup struct {
	Name    string      `xml:"name,attr" json:"name"`
	Artists []subArtist `xml:"artist" json:"artist"`
}

type subArtist struct {
	ID         string `xml:"id,attr" json:"id"`
	Name       string `xml:"name,attr" json:"name"`
	AlbumCount int64  `xml:"albumCount,attr,omitempty" json:"albumCount,omitempty"`
	CoverArt   string `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty"`
	Starred    string `xml:"starred,attr,omitempty" json:"starred,omitempty"`
}

type albumList2 struct {
	XMLName xml.Name   `xml:"albumList2"`
	Albums  []subAlbum `xml:"album" json:"album"`
}

type subAlbum struct {
	ID       string    `xml:"id,attr" json:"id"`
	Name     string    `xml:"name,attr" json:"name"`
	Artist   string    `xml:"artist,attr,omitempty" json:"artist,omitempty"`
	ArtistID string    `xml:"artistId,attr,omitempty" json:"artistId,omitempty"`
	CoverArt string    `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty"`
	SongCnt  int64     `xml:"songCount,attr,omitempty" json:"songCount,omitempty"`
	Duration int64     `xml:"duration,attr,omitempty" json:"duration,omitempty"`
	Created  string    `xml:"created,attr,omitempty" json:"created,omitempty"`
	Starred  string    `xml:"starred,attr,omitempty" json:"starred,omitempty"`
	Songs    []subSong `xml:"song,omitempty" json:"song,omitempty"`
}

type subSong struct {
	ID          string `xml:"id,attr" json:"id"`
	Parent      string `xml:"parent,attr,omitempty" json:"parent,omitempty"`
	Title       string `xml:"title,attr" json:"title"`
	Album       string `xml:"album,attr,omitempty" json:"album,omitempty"`
	Artist      string `xml:"artist,attr,omitempty" json:"artist,omitempty"`
	IsDir       bool   `xml:"isDir,attr" json:"isDir"`
	CoverArt    string `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty"`
	Duration    int64  `xml:"duration,attr,omitempty" json:"duration,omitempty"`
	BitRate     int64  `xml:"bitRate,attr,omitempty" json:"bitRate,omitempty"`
	Track       int64  `xml:"track,attr,omitempty" json:"track,omitempty"`
	Year        int64  `xml:"year,attr,omitempty" json:"year,omitempty"`
	Size        int64  `xml:"size,attr,omitempty" json:"size,omitempty"`
	Suffix      string `xml:"suffix,attr,omitempty" json:"suffix,omitempty"`
	ContentType string `xml:"contentType,attr,omitempty" json:"contentType,omitempty"`
	AlbumID     string `xml:"albumId,attr,omitempty" json:"albumId,omitempty"`
	ArtistID    string `xml:"artistId,attr,omitempty" json:"artistId,omitempty"`
	Starred     string `xml:"starred,attr,omitempty" json:"starred,omitempty"`
	Type        string `xml:"type,attr,omitempty" json:"type,omitempty"`
}

type persistedState struct {
	Version   int                     `json:"version"`
	Songs     map[string]songMeta     `json:"songs"`
	Covers    map[string]string       `json:"covers"`
	Playlists map[string]playlistMeta `json:"playlists"`
}

type playlistMeta struct {
	ID                 string   `json:"id"`
	Source             string   `json:"source,omitempty"`
	ListID             int64    `json:"list_id,omitempty"`
	Name               string   `json:"name"`
	Comment            string   `json:"comment,omitempty"`
	Owner              string   `json:"owner,omitempty"`
	Public             bool     `json:"public"`
	SongCount          int64    `json:"song_count,omitempty"`
	Duration           int64    `json:"duration,omitempty"`
	Created            string   `json:"created,omitempty"`
	Changed            string   `json:"changed,omitempty"`
	CoverArt           string   `json:"cover_art,omitempty"`
	SongIDs            []string `json:"song_ids,omitempty"`
	GlobalCollectionID string   `json:"global_collection_id,omitempty"`
}
