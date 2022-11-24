package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v3"
)

var (
	pathsToUrls map[string]string
	DB          Database
)

type Database struct {
	db *bolt.DB
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

	// create bucket
	DB.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("Redirects"))
		if err != nil {
			return err
		}
		return nil
	})

	// read into database from file
	for k, v := range pathsToUrls {
		err = DB.Put([]byte(k), []byte(v))
		if err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
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

// put to database
func (db *Database) Put(key, value []byte) error {
	err := db.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Redirects"))
		return bucket.Put(key, value)
	})
	return err
}
