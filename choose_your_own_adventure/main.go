package main

import (
	"log"
	"os"

	handler "github.com/KameeKaze/gophercises/handler"
	utils "github.com/KameeKaze/gophercises/utils"
)

func main() {
	// read story file
	file, err := os.ReadFile("source/gopher.json")
	if err != nil {
		log.Fatal(err)
	}
	utils.Story, err = utils.ParseJSON(file)
	if err != nil {
		log.Fatal(err)
	}
	// start
	handler.Routes()
}
