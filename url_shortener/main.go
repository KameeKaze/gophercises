package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

var (
	pathsToUrls map[string]string
	DB          Database
)

type Database struct {
	db     *bolt.DB
	tx     *bolt.Tx
	bucket *bolt.Bucket
}

func init() {
	// YAML file as a flag
	var filename string
	flag.StringVar(&filename, "f", "PathsToUrls.yaml", "Accept YAML file as a flag")
	flag.Parse()

	// open file containing the urls
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	// check file extension and parse into map
	if filename[len(filename)-4:] == "json" {
		pathsToUrls = ParseJSON([]byte(file))
	} else {
		pathsToUrls = ParseYAML([]byte(file))
	}
}

func main() {
	// connect to database
	var err error
	DB.db, err = bolt.Open("urls.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.db.Close()
	// Start a writable transaction.
	DB.tx, err = DB.db.Begin(true)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.tx.Rollback()
	// create bucket
	DB.bucket, err = DB.tx.CreateBucket([]byte("Redirects"))
	if err != nil {
		log.Fatal(err)
	}
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

// get urls from json
func ParseJSON(data []byte) (jsonData map[string]string) {
	// parse json data
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatal(err)
	}
	// return json as map[string]string)
	return jsonData
}
