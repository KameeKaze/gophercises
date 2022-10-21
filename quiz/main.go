package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// parse cli flag arguements
	filename := flag.String("file", "problems.csv", "input file")
	flag.Parse()

	// open file
	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// read file into data
	var data = make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// split line by comma
		line := strings.Split(scanner.Text(), ",")
		// check if line has one comma
		if len(line) != 2 {
			panic("Invalid file content")
		}
		// add to data
		data[line[0]] = line[1]
	}

	// print data
	for k, v := range data {
		fmt.Println(k, v)
	}
}
