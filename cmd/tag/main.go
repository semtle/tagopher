package main

import (
	"github.com/codegangsta/cli"
	"github.com/imdario/tagopher"
	"os"
)

var (
	initCommand = cli.Command{
		Name:   "init",
		Usage:  "create an empty Tagopher repository",
		Action: initAction,
	}
	tagCommand = cli.Command{
		Name:   "tag",
		Usage:  "add, remove or replace tags attached to given files",
		Action: tagAction,
	}
	listCommand = cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "list current repository's path entries",
		Action:  listAction,
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "tag"
	app.Usage = "organize files with tags"
	app.Commands = []cli.Command{
		initCommand,
		tagCommand,
		listCommand,
	}
	known := false
	for _, command := range app.Commands {
		if command.Name == os.Args[1] {
			known = true
			break
		}
	}
	if !known {
		os.Args = append([]string{tagopher.TAG}, os.Args...)
	}
	app.Run(os.Args)
}
