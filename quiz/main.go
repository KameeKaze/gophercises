package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
)

func main() {
	// parse cli flag arguements
	filename := flag.String("f", "problems.csv", "input file")
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

	// print final score at the end
	userPoints := 0
	defer fmt.Printf("You scored %d points out of %d!\n", userPoints, len(data))
	var answer string
	// channel for timeout
	timeout := make(chan int)
	// always scan for answers
	go func() {
		for {
			fmt.Println("Your answer:")
			fmt.Scanln(&answer)
			timeout <- 1
		}
	}()
	// iterate over questions
	for k, v := range data {
		// print question
		fmt.Println("Question:", k)

		select {
		// time is up
		case <-time.After(3 * time.Second):
			fmt.Println(Red + "Timed out, next!" + Reset)
		// check user answer
		case <-timeout:
			if answer == v {
				userPoints++
				fmt.Println(Green + "Correct!" + Reset)
			} else {
				fmt.Println(Red + "Incorrect!\n" + Reset + "The correct answer was " + v)
			}
			// restet answer
			answer = ""
		}
		fmt.Println()
	}

}
