package main

import (
	"fmt"
	"os"
)

type Task struct {
	Name string
	Done bool
}

var (
	args = os.Args[1:]
)

func init() {
	if len(args) == 0 {
		fmt.Println("Too few arguements")
		os.Exit(1)
	}
}

func main() {
	switch args[0] {
	case "add":
		fmt.Println("Add a task:")
	case "do":
		fmt.Println("Enter number of task that you finnsihed:")
	default:
		fmt.Printf("Unrecognised option %s.\n", args[0])
	}
}
