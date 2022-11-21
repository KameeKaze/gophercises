package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func main() {
	// read story file
	file, err := os.ReadFile("gopher.json")
	if err != nil {
		log.Fatal(err)
	}
	story := ParseJSON([]byte(file))
	fmt.Println(story)
}

// parse json into struct
func ParseJSON(data []byte) (story map[string]Chapter) {
	// parse json data
	err := json.Unmarshal(data, &story)
	if err != nil {
		log.Fatal(err)
	}
	return story
}
