package main

import (
	"github.com/imdario/tagopher"
	"github.com/codegangsta/cli"
)

func initAction(c *cli.Context) {
	err := tagopher.Init()
	if err != nil {
		panic(err)
	}
}
