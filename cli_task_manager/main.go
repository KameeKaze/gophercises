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

const (
	usage = `task is a CLI for managing your TODOs.

	Usage:
		task [command]
	
	Available Commands:
		add         Add a new task to your TODO list
		do          Mark a task on your TODO list as complete
		list        List all of your incomplete tasks
		rm          Remove a task from your TODO list
		completed   List all of your completed tasks`
)

func main() {
	// no arguements
	if len(args) == 0 {
		fmt.Println(usage)
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
			return
		}
		// add task into database
		task := Task{}
		task.Name = strings.Join(args[1:], " ")
		err = Add(task)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Added \"%s\" to your task list.\n", task.Name)
	case "do": // change task status to done
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
		// get task from database
		task, err := Get(itob(number))
		if err != nil {
			fmt.Println("No task found with this number.")
		}
		// check if task was already done
		if task.Done {
			fmt.Println("This task was already done.")
			return
		}
		// change task status
		task.Done = true
		// update task in database
		err = Update(task)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("You have completed the \"%s\" task.\n", task.Name)
	case "rm": // print all undone tasks
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
		// get task from database
		task, err := Get(itob(number))
		if err != nil {
			fmt.Println("No task found with this number.")
		}
		err = Delete(itob(task.ID))
		if err != nil {
			log.Fatal(err)
		}
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
	case "completed": // print all finished tasks
		tasks, err := ViewAll()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You have the following tasks done:")
		for i := range tasks {
			if tasks[i].Done {
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
func Add(task Task) error {
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

// add task to the database
func Update(task Task) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		// convert task to json
		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}
		// update id and json in database
		return b.Put(itob(task.ID), buf)
	})
}

// add task to the database
func Delete(index []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		// update id and json in database
		return b.Delete(index)
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
func Get(index []byte) (Task, error) {
	var task Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte((bucketName)))
		v := b.Get((index))
		err := json.Unmarshal(v, &task)
		if err != nil {
			return err
		}
		return nil
	})
	return task, err
}
