package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"

	"golang.org/x/net/html"
)

var (
	URL   string
	links []string
	wg    sync.WaitGroup
)

func init() {
	flag.StringVar(&URL, "u", "https://google.com", "URL")
	flag.Parse()
}

func main() {
	// start recursive link parsing
	wg.Add(1)
	go ParseLinks(URL)
	wg.Wait()

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
	fmt.Println("Sitemap generated to sitemap.xml")
}

func GetBody(url string) (string, error) {
	// make get request and get body
	resp, err := http.Get(url)
	if err != nil {
		return "", err

	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), err

}

func ParseLinks(url string) {
	// get url body
	body, err := GetBody(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	// parse body
	z := html.NewTokenizer(strings.NewReader(body))
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken: // end of the document
			wg.Done()
			return
		case html.StartTagToken: // found html tag
			t := z.Token()
			// find link tag
			if t.Data == "a" { // html link tag
				// find href
				for _, a := range t.Attr {
					if a.Key == "href" { // link
						// check if link already found
						if !contains(links, a.Val) {
							links = append(links, a.Val)
							// recursive parse on new link
							if strings.HasPrefix(a.Val, "/") {
								wg.Add(1)
								go ParseLinks(URL + a.Val)
							}
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
