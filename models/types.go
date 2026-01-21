// Package models has the data structures all for the r-a-d.io API.
package models

// Song a singular track from the API.
type Song struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Requestable bool   `json:"requestable"`
}

// SearchResponse is the JSON envelope for searches.
type SearchResponse struct {
	Total int    `json:"total"`
	Data  []Song `json:"data"`
}

// RadioStatus holds the high level state of the stream.
type RadioStatus struct {
	Main MainStream `json:"main"`
}

// MainStream real time stats, current song playing and show queue.
type MainStream struct {
	NowPlaying string  `json:"np"`
	Listeners  int     `json:"listeners"`
	DJ         DJ      `json:"dj"`
	Queue      []Track `json:"queue"`
}

// DJ shows the current dj...
// If not mistaken the clanker dj is Hanyuu-sama.
// When Hanyuu-sama is dj requesting is allowed.
// But when an actual human dj is on, then requests are usually disabled.
type DJ struct {
	ID   int    `json:"id"`
	Name string `json:"djname"`
}

// Track defines a queued item in the stream.
type Track struct {
	MetaData  string `json:"meta"`
	Type      int    `json:"type"` // Detect if real human requested the song or clanker slop
	Timestamp int64  `json:"timestamp"`
}
