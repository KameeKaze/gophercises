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

var (
	UserPoints  = 0
	Filename    string
	TimeoutFlag int
)

func init() {
	flag.IntVar(&TimeoutFlag, "t", 0, "Set timeout - 0 means no timeout")
	flag.StringVar(&Filename, "f", "problems.csv", "input file")
	flag.Parse()
}

func main() {
	// open file
	file, err := os.Open(Filename)
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

		if TimeoutFlag != 0 {
			select {
			// time is up
			case <-time.After(time.Duration(TimeoutFlag) * time.Second):
				fmt.Println(Red + "Timed out, next!" + Reset)
			case <-timeout:
				UserPoints += checkAnswer(answer, v)
				// restet answer
				answer = ""
			}
		} else {
			<-timeout
			UserPoints += checkAnswer(answer, v)
			// restet answer
			answer = ""
		}
		fmt.Println()
	}
	fmt.Printf("You scored %d points out of %d!\n", UserPoints, len(data))

}

func checkAnswer(answer, solution string) int {
	if answer == solution {
		fmt.Println(Green + "Correct!" + Reset)
		return 1
	} else {
		fmt.Println(Red + "Incorrect!\n" + Reset + "The correct answer was " + solution)
		return 0
	}
}
