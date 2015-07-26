package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/imdario/tagopher"
	"os"
)

var (
	initCommand = cli.Command{
		Name:   "init",
		Usage:  fmt.Sprintf("create an empty %s repository", tagopher.TAG_NAME),
		Action: initAction,
	}
	tagCommand = cli.Command{
		Name:   "tag",
		Usage:  "add, remove or replace tags attached to given files (flags can be used multiple times or once with comma-separated tags)",
		Action: tagAction,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name: "add, a",
				Value: &cli.StringSlice{},
			},
			cli.StringSliceFlag{
				Name: "remove, r",
				Value: &cli.StringSlice{},
			},
			cli.StringSliceFlag{
				Name: "rename, e",
				Value: &cli.StringSlice{},
			},
		},
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
		for _, alias := range command.Aliases {
			if alias == os.Args[1] {
				known = true
				break
			}
		}
	}
	if !known {
		tmp := []string{os.Args[0], tagCommand.Name}
		os.Args = append(tmp, os.Args[1:]...)
	}
	app.Run(os.Args)
}
