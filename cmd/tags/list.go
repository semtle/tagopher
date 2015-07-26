package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	"github.com/imdario/tagopher"
	"github.com/imdario/termtable"
	"os"
	"strings"
	"time"
)

func asTermtableRows(files []tagopher.File) (rows [][]string) {
	for _, file := range files {
		rows = append(rows, asTermtableRow(file))
	}
	return
}

func asTermtableRow(file tagopher.File) (row []string) {
	row = append(row, fmt.Sprintf("%s", file.Mode()))
	row = append(row, humanize.Bytes(uint64(file.Size())))
	row = append(row, file.ModTime().Format(time.Stamp))
	row = append(row, file.String())
	tags := strings.Join(file.Tags(), ", ")
	if tags == "" {
		tags = "-"
	}
	row = append(row, tags)
	return
}

func listAction(c *cli.Context) {
	paths := c.Args()
	if len(paths) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		paths = []string{cwd}
	}
	for _, path := range paths {
		files, err := tagopher.List(path)
		if err != nil {
			fmt.Printf("%s\n", err)
		} else {
			if len(paths) > 1 {
				fmt.Printf("%s:\n", path)
			}
			t := termtable.NewTable(asTermtableRows(files), &termtable.TableOptions{
				Padding: 2,
			})
			fmt.Println(t.Render())
			if len(paths) > 1 {
				fmt.Println()
			}
		}
	}
}
