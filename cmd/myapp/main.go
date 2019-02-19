package main

import (
	"github.com/nikhil-github/api-cab-data/pkg/wiring"
	"log"
)

var (
	gitCommit = "nil" // Git commit hash, set during build time
	version   = "nil" // Version is the build version of the application's source code
)

func main() {
	var cfg wiring.Config
	a := wiring.App{
		Config:    cfg,
		GitCommit: gitCommit,
		Version:   version,
	}
	a.Run()
	log.Println("started app")
}
