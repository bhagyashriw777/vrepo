// Package main is the entry point for the application.
package main

import (
	"os"

	"github.com/bhagyashriw777/vrepo/agent/agent"
)

func main() {
	apps := os.Args
	//remove default path from arguments
	apps = apps[1:]
	agent.Init(apps)
}
