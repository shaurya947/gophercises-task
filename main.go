package main

import (
	"os"

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

func main() {
	app := &cli.App{
		Name:      "task",
		Usage:     "task is a CLI for managing your TODOs.",
		UsageText: "task [command]",
		Commands:  []*cli.Command{cmdAdd, cmdDo, cmdList},
	}
	app.Run(os.Args)
}
