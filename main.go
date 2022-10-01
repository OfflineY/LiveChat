package main

import (
	// ./modules
	"OnlineChat/modules"
)

var Version string = "2.1"

func main() {
	modules.Bootstrap(Version)
}

// main.go
