package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

var (
	DB Database
)

type Database struct {
	db     *bolt.DB
	tx     *bolt.Tx
	bucket *bolt.Bucket
}

func main() {
	err := ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer DB.db.Close()
	defer DB.tx.Rollback()

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
	if newUrl := DB.Get([]byte(url)); newUrl != "" {
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

func ConnectToDB() error {
	// connect to database
	var err error
	DB.db, err = bolt.Open("urls.db", 0600, nil)
	if err != nil {
		return err
	}
	// Start a writable transaction.
	DB.tx, err = DB.db.Begin(true)
	if err != nil {
		return err
	}
	// get bucket
	DB.bucket = DB.tx.Bucket([]byte("Redirects"))
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) Get(key []byte) string {
	url := db.bucket.Get(key)
	return string(url)
}
