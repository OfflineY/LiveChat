package main

import (
	// ./modules
	"OnlineChat/modules"
)

var Version string = "2.0.0 BETA "

func main() {
	modules.Bootstrap(Version)
}

// main.go
