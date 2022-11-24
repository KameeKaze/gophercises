package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

var (
	url   string
	links []string
)

func init() {
	flag.StringVar(&url, "u", "https://calhoun.io", "URL")
	flag.Parse()
}

func main() {
	// make get request and get body
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// parse links in body
	done := make(chan (bool))
	go ParseLinks(string(body), done)
	<-done
	// print links
	for x := range links {
		fmt.Println(links[x])
	}
}

func ParseLinks(body string, done chan<- bool) {
	z := html.NewTokenizer(strings.NewReader(string(body)))
	for {
		tt := z.Next()
		switch tt {
		// end of the document
		case html.ErrorToken:
			done <- true
			return
		case html.StartTagToken:
			t := z.Token()
			// find link tag
			if t.Data == "a" {
				// find href
				for _, a := range t.Attr {
					if a.Key == "href" {
						if !contains(links, a.Val) {
							links = append(links, a.Val)
						}
						break
					}
				}
			}
		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
