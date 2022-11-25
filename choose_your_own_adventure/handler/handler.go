package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	utils "github.com/KameeKaze/gophercises/utils"
	"github.com/gorilla/mux"
)

func Routes() {
	r := mux.NewRouter()

	// define routes
	r.HandleFunc("/", home).Methods("GET")
	r.PathPrefix("/css").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("templates/css/")))).Methods("GET")
	r.HandleFunc("/{url}", storyHandler).Methods("GET")

	//start http server
	fmt.Println("Running on http://localhost" + ":8080")
	http.ListenAndServe(":8080", r)

}

func home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/intro", http.StatusPermanentRedirect)
}

func storyHandler(w http.ResponseWriter, r *http.Request) {
	url := mux.Vars(r)["url"]
	tmplt, err := template.ParseFiles("templates/chapter.html")
	if err != nil {
		log.Println(err)
	}
	event := utils.Story[url]
	tmplt.Execute(w, event)
}
