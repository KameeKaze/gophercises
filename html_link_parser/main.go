package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	// open file containing the urls
	file, err := os.ReadFile("ex1.html")
	if err != nil {
		log.Fatal(err)
	}
	z := html.NewTokenizer(strings.NewReader(string(file)))

	links := []Link{}
	done := make(chan (bool))
	go func() {
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
	}()
	// print out links
	<-done
	for i := range links {
		fmt.Printf("Href: %s\n", links[i].Href)
		fmt.Printf("Text: %s\n\n", links[i].Text)
	}

}
