package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	// open file containing the urls
	file, err := ioutil.ReadFile("ex1.html")
	if err != nil {
		log.Fatal(err)
	}
	z := html.NewTokenizer(strings.NewReader(string(file)))

	links := []Link{}
	defer fmt.Println(&links)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			return
		case html.StartTagToken:
			t := z.Token()
			// find link tag
			if t.Data == "a" {
				link := Link{}
				// find href
				for _, a := range t.Attr {
					if a.Key == "href" {
						link.Href = a.Val
						break
					}
				}
				// find text
				z.Next()
				// trim whitespace characters
				text := strings.TrimSpace(string(z.Text()))
				link.Text = string(text)

				links = append(links, link)

			}
		}
	}

}
