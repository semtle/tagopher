package main

import (
	"github.com/codegangsta/cli"
	"github.com/imdario/tagopher"
)

func initAction(c *cli.Context) {
	err := tagopher.Init()
	if err != nil {
		panic(err)
	}
}
