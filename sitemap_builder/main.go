package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

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
	// parse links into xml file
	link, err := CreateXML(links)
	if err != nil {
		log.Fatal(err)
	}
	// write xml into file
	err = WriteIntoFile("sitemap.xml", link)
	if err != nil {
		log.Fatal(err)
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

func WriteIntoFile(filename string, data []byte) error {
	err := os.WriteFile(filename, data, 0644)
	return err
}

func CreateXML(links []string) ([]byte, error) {
	t, err := template.ParseFiles("templates/sitemap.xml")
	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, links)
	if err != nil {
		return nil, err
	}
	return tpl.Bytes(), err
}
