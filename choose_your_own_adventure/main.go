package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

func main() {
	// read story file
	file, err := os.ReadFile("gopher.json")
	if err != nil {
		log.Fatal(err)
	}
	Story = ParseJSON([]byte(file))

	r := mux.NewRouter()

	// define routes
	r.HandleFunc("/", Home).Methods("GET")
	r.PathPrefix("/css").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("templates/css/")))).Methods("GET")
	r.HandleFunc("/{url}", StoryHandler).Methods("GET")

	//start http server
	fmt.Println("Running on http://localhost" + ":8080")
	http.ListenAndServe(":8080", r)
}

func Home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/intro", http.StatusPermanentRedirect)

}

func StoryHandler(w http.ResponseWriter, r *http.Request) {
	url := mux.Vars(r)["url"]
	tmplt, err := template.ParseFiles("templates/chapter.html")
	if err != nil {
		log.Fatal(err)
	}
	event := Story[url]
	tmplt.Execute(w, event)
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
