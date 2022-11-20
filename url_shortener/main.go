package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

var (
	pathsToUrls map[string]string
)

// parse yaml to urls
func init() {
	// open file containing the urls
	file, err := os.ReadFile("PathsToUrls.yaml")
	if err != nil {
		log.Fatal(err)
	}
	// parse file into yaml
	pathsToUrls = ParseYAML([]byte(file))
}

func main() {
	r := mux.NewRouter()

	// define routes
	r.HandleFunc("/{url}", UrlHandler).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	//start http server
	fmt.Println("Running on http://localhost" + ":8080")
	http.ListenAndServe(":8080", r)
}

func UrlHandler(w http.ResponseWriter, r *http.Request) {
	// get url parameter
	url := mux.Vars(r)["url"]
	// check if url exists, if it does redirect to the path
	if newUrl := pathsToUrls[url]; newUrl != "" {
		http.Redirect(w, r, newUrl, http.StatusPermanentRedirect)
	} else { // else 404 not found
		NotFoundHandler(w, r)
	}
}

// custom 404 handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("<h1>The requested URL was not found.</h1>\n"))
}

// get urls from yaml
func ParseYAML(data []byte) (yamlData map[string]string) {
	// parse yaml data
	err := yaml.Unmarshal(data, &yamlData)
	if err != nil {
		log.Fatal(err)
	}
	// return yaml as map[string]string)
	return yamlData
}
