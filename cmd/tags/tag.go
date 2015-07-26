package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/imdario/tagopher"
	"os"
	"strings"
)

func cleanTag(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func collectTags(value []string) (tags []string) {
	const SEPARATOR = ","
	for _, v := range value {
		if strings.Contains(v, SEPARATOR) {
			t := strings.Split(v, SEPARATOR)
			for _, e := range t {
				tags = append(tags, cleanTag(e))
			}
		} else {
			tags = append(tags, cleanTag(v))
		}
	}
	return
}

func applyTags(c *cli.Context, flag string) (err error) {
	value := c.StringSlice(flag)
	if len(value) == 0 {
		return
	}
	tags := collectTags(value)
	for _, file := range c.Args() {
		switch (flag) {
		case "add":
			err = tagopher.AddTags(file, tags...)
		case "remove":
			err = tagopher.RemoveTags(file, tags...)
		case "rename":
			err = tagopher.RenameTags(file, tags...)
		}
	}
	return
}

func tagAction(c *cli.Context) {
	for _, flag := range []string{"add", "remove", "rename"} {
		if err := applyTags(c, flag); err != nil {
			if os.IsNotExist(err) {
				fmt.Println(err)
				break
			}
			panic(err)
		}
	}
}
