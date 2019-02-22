package main

import (
	"github.com/nikhil-github/api-cab-data/pkg/wiring"
)

// TODO :- add git version/build number

func main() {
	var cfg wiring.Config
	a := wiring.App{
		Config: cfg,
	}
	a.Run()
}




