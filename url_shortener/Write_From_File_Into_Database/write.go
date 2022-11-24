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
	Filename    string
	pathsToUrls map[string]string
	DB          Database
)

type Database struct {
	db *bolt.DB
}

func init() {
	flag.StringVar(&Filename, "f", "PathsToUrls.yaml", "Accept YAML file as a flag")
	flag.Parse()
}

func ReadFile(filename string) []byte {
	// read file containing the urls
	file, err := os.ReadFile(Filename)
	if err != nil {
		return nil
	}
	return file
}

func main() {
	file := ReadFile(Filename)
	// check file extension and parse into map
	var err error
	if Filename[len(Filename)-4:] == "json" {
		pathsToUrls, err = ParseJSON(file)
	} else {
		pathsToUrls, err = ParseYAML(file)
	}
	// connect to database
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
func ParseYAML(data []byte) (yamlData map[string]string, err error) {
	// parse yaml data
	err = yaml.Unmarshal(data, &yamlData)
	// return yaml as map[string]string)
	return yamlData, err
}

// get urls from json
func ParseJSON(data []byte) (jsonData map[string]string, err error) {
	// parse json data
	err = json.Unmarshal(data, &jsonData)
	// return json as map[string]string)
	return jsonData, err
}

// put to database
func (db *Database) Put(key, value []byte) error {
	err := db.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Redirects"))
		return bucket.Put(key, value)
	})
	return err
}
