package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/imdario/tagopher"
)

func listAction(c *cli.Context) {
	paths := c.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}
	for _, path := range paths {
		files, err := tagopher.List(path)
		if err != nil {
			fmt.Printf("%s: %s: No such file or directory\n", tagopher.TAG, path)
		} else {
			if len(paths) > 1 {
				fmt.Printf("%s:\n", path)
			}
			for _, file := range files {
				fmt.Println(file)
			}
			if len(paths) > 1 {
				fmt.Println()
			}
		}
	}
}
