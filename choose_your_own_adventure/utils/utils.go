package utils

import (
	"encoding/json"
)

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

var (
	Story map[string]Chapter
)

// parse json into struct
func ParseJSON(data []byte) (story map[string]Chapter, err error) {
	// parse json data
	err = json.Unmarshal(data, &story)
	return story, err
}
