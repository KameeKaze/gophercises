package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

var (
	args       = os.Args[1:]
	db         *bolt.DB
	bucketName = "tasks"
)

func main() {
	// no arguements
	if len(args) == 0 {
		fmt.Println("Too few arguements")
		os.Exit(1)
	}
	err := DBInit()
	if err != nil {
		log.Fatal(err)
	}

	switch args[0] {
	case "add": // add task to the list
		// task needs a name
		if len(args) == 1 {
			fmt.Println("Enter a name!")
			os.Exit(1)
		}
		// add task into database
		task := Task{}
		task.Name = strings.Join(args[1:], " ")
		Update(task)
		fmt.Printf("Added \"%s\" to your task list.\n", task.Name)
	case "do":
		// need task number
		if len(args) == 1 {
			fmt.Println("Enter a number!")
			os.Exit(1)
		}
		// convert to int
		number, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Enter a number!")
			os.Exit(1)
		}
		fmt.Println(number)
		// TODO
	case "list": // print all undone tasks
		tasks, err := ViewAll()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You have the following tasks:")
		for i := range tasks {
			if !tasks[i].Done {
				fmt.Printf("%d. %s\n", tasks[i].ID, tasks[i].Name)
			}
		}

	// invalid option
	default:
		fmt.Printf("Unrecognised option %s.\n", args[0])
	}
}

// connect to database
func DBInit() error {
	var err error
	db, err = bolt.Open("tasks.db", 0600, nil)
	if err != nil {
		return err
	}
	// create database bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	return err
}

// convert int to byte to use as database key
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// add task to the database
func Update(task Task) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		// create task id
		id, _ := b.NextSequence()
		task.ID = int(id)

		// convert task to json
		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}
		// store id and json in database
		err = b.Put(itob(task.ID), buf)
		return err
	})
}

// return all elements from database
func ViewAll() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()
		// iterate over all elements in database
		for k, v := c.First(); k != nil; k, v = c.Next() {
			// convert json back to task
			var task *Task
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}
			tasks = append(tasks, *task)
		}
		return nil
	})
	// return all tasks
	return tasks, err
}

// get element from database from given key
// convert int to byte with itob()
func View(index []byte) ([]byte, error) {
	var v []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte((bucketName)))
		v = b.Get((index))
		return nil
	})
	return v, err
}
