package main

import (
	"medusa/src/api/webserver"
	"medusa/src/api/worker"
	"os"
)

func main() {
	a := os.Args[2]
	switch a {
	case "worker":
		worker.Boot()
	case "web":
		webserver.Boot()
	default:
		panic("Cannot identify the server type")
	}
}
