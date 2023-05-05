package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/shaurya947/gophercises-task/store"
	"github.com/urfave/cli/v2"
)

var (
	cmdAdd = &cli.Command{
		Name:  "add",
		Usage: "Add tasks to your TODO list",
		UsageText: "Enclose each task in quotes, such as\n\n" +
			"task add \"do dishes\" \"wash clothes\" ...",
		Action: addTasks,
	}
	cmdComplete = &cli.Command{
		Name:  "do",
		Usage: "Mark tasks on your TODO list as complete",
		UsageText: "For each task that you'd like to mark as complete" +
			", pass the task number as displayed by the \"list\" " +
			"command. For example\n\ntasks do 1 6 15\n\nwill mark" +
			" as complete the 1st, 6th and 15th tasks displayed" +
			" by the \"list\" command.",
		Action: completeTasks,
	}
	cmdListIncomplete = &cli.Command{
		Name:   "list",
		Usage:  "List all of your incomplete tasks",
		Action: listIncompleteTasks,
	}
	cmdListCompleted = &cli.Command{
		Name:   "completed",
		Usage:  "List all of your completed tasks since 24h ago",
		Action: listCompletedTasks,
	}
	cmdRemove = &cli.Command{
		Name:  "rm",
		Usage: "Delete incomplete tasks from your TODO list",
		UsageText: "For each task that you'd like to delete, pass the" +
			" task number as displayed by the \"list\" command. " +
			"For example\n\ntasks rm 4 9\n\nwill delete the 4th " +
			"and 9th tasks displayed by the \"list\" command.",
		Action: removeTasks,
	}
)

const (
	dataDirName = ".tasks"
	dbFileName  = "tasks.db"
)

var taskStore *store.TaskStore

func main() {
	dbFilepath, err := getDBFilepath()
	if err != nil {
		log.Fatal(err)
	}

	taskStore, err = store.NewTaskStore(*dbFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer taskStore.Close()

	app := &cli.App{
		Name:      "task",
		Usage:     "task is a CLI for managing your TODOs.",
		UsageText: "task [command]",
		Commands: []*cli.Command{
			cmdAdd,
			cmdComplete,
			cmdListIncomplete,
			cmdListCompleted,
			cmdRemove,
		},
	}
	app.Run(os.Args)
}

func getDBFilepath() (*string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("Could not find user home directory")
	}

	dataDirFilepath := filepath.Join(home, dataDirName)
	if _, err := os.Stat(dataDirFilepath); os.IsNotExist(err) {
		if err := os.Mkdir(dataDirFilepath, 0755); err != nil {
			return nil, fmt.Errorf(
				"Could not create data directory")
		}
	}

	dbFilepath := filepath.Join(dataDirFilepath, dbFileName)
	return &dbFilepath, nil
}

func addTasks(ctx *cli.Context) error {
	args := ctx.Args()
	addedTasks := make([]*store.Task, args.Len())
	for i := 0; i < args.Len(); i++ {
		task := &store.Task{Description: args.Get(i)}
		err := taskStore.AddTask(task)
		if err != nil {
			log.Fatalln(err)
		}
		addedTasks[i] = task
	}

	fmt.Println("Added the following tasks:")
	for _, task := range addedTasks {
		fmt.Printf("- %s\n", task.Description)
	}
	return nil
}

func listIncompleteTasks(ctx *cli.Context) error {
	incompleteTasks, err := taskStore.GetIncompleteTasks()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("You have the following incomplete tasks:")
	for i, task := range incompleteTasks {
		fmt.Printf("%d. %s\n", i+1, task.Description)
	}
	return nil
}

func listCompletedTasks(ctx *cli.Context) error {
	fullDayAgo := time.Now().Add(-time.Hour * 24)
	completedTasks, err := taskStore.GetCompletedTasks(fullDayAgo)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("You have finished the following tasks since 24h ago:")
	for _, task := range completedTasks {
		fmt.Printf("- %s\n", task.Description)
	}
	return nil
}

func completeTasks(ctx *cli.Context) error {
	args := ctx.Args()
	taskNums := make([]int, args.Len())
	for i := 0; i < args.Len(); i++ {
		argInt, err := strconv.Atoi(args.Get(i))
		if err != nil {
			log.Fatalln(err)
		}
		taskNums[i] = argInt
	}

	completedTasks, err := taskStore.CompleteTasks(taskNums)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Marked the following tasks as complete:")
	for _, task := range completedTasks {
		fmt.Printf("- %s\n", task.Description)
	}
	return nil
}

func removeTasks(ctx *cli.Context) error {
	args := ctx.Args()
	taskNums := make([]int, args.Len())
	for i := 0; i < args.Len(); i++ {
		argInt, err := strconv.Atoi(args.Get(i))
		if err != nil {
			log.Fatalln(err)
		}
		taskNums[i] = argInt
	}

	removedTasks, err := taskStore.RemoveTasks(taskNums)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Deleted the following incomplete tasks:")
	for _, task := range removedTasks {
		fmt.Printf("- %s\n", task.Description)
	}
	return nil
}
