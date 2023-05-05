package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/shaurya947/gophercises-task/store"
	"github.com/urfave/cli/v2"
)

var (
	cmdAdd = &cli.Command{
		Name:  "add",
		Usage: "Add a new task to your TODO list",
	}
	cmdDo = &cli.Command{
		Name:  "do",
		Usage: "Mark a task on your TODO list as complete",
	}
	cmdList = &cli.Command{
		Name:  "list",
		Usage: "List all of your incomplete tasks",
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
		Commands:  []*cli.Command{cmdAdd, cmdDo, cmdList},
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
