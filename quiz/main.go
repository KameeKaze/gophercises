package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
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

	// iterate over questions
	userPoints := 0
	for k, v := range data {
		// print question
		fmt.Println(k)
		// read answer user input
		var answer string
		fmt.Println("Your answer:")
		fmt.Scanln(&answer)

		// check answer
		if answer == v {
			userPoints++
			fmt.Println(Green + "Correct!" + Reset)
		} else {
			fmt.Println(Red + "Incorrect!\n" + Reset + "The correct answer was " + v)
		}
		fmt.Println()
	}

	// print final score
	fmt.Printf("You scored %d points\n", userPoints)
}
