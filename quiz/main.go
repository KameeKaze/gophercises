package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
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
	Filename    string
	TimeoutFlag int
)

func init() {
	flag.IntVar(&TimeoutFlag, "t", 0, "Set timeout - 0 means no timeout")
	flag.StringVar(&Filename, "f", "problems.csv", "input file")
	flag.Parse()
}

type Question struct {
	Question string
	Answer   string
}

func main() {
	quiz, err := createQuiz()
	if err != nil {
		log.Fatal(err)
	}

	// channel for timeout
	timeout := make(chan bool)
	// always scan for answers
	var answer string
	go scanAnswer(&answer, timeout)

	// iterate over questions
	var userPoints int
	for i := range quiz {
		// print question
		fmt.Println("Question:", quiz[i].Question)
		fmt.Println("Your answer:")

		if TimeoutFlag != 0 {
			select {
			// time is up
			case <-time.After(time.Duration(TimeoutFlag) * time.Second):
				fmt.Println(Red + "Timed out, next!" + Reset)
			case <-timeout:
				userPoints += checkAnswer(answer, quiz[i].Answer)
				// restet answer
				answer = ""
			}
		} else {
			<-timeout
			userPoints += checkAnswer(answer, quiz[i].Answer)
			// restet answer
			answer = ""
		}
		fmt.Println()
	}
	fmt.Printf("You scored %d points out of %d!\n", userPoints, len(quiz))
}

func scanAnswer(answer *string, timeout chan<- bool) {
	for {
		fmt.Scanln(answer)
		timeout <- true
	}
}

func readFile(filename string) ([]byte, error) {
	// read file and return
	return os.ReadFile(filename)
}

func createQuiz() (quiz []*Question, err error) {
	// read file
	file, err := readFile(Filename)
	if err != nil {
		return
	}
	// split file by new line
	for _, line := range strings.Split(string(file), "\n") {
		// split line by comma
		line := strings.Split(string(line), ",")
		// check if question and answer seperated by a comme
		if len(line) != 2 {
			err = errors.New("invalid file content")
			return
		}
		// append to quiz
		quiz = append(quiz, &Question{
			line[0],
			line[1],
		})
	}
	return quiz, nil
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
