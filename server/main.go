package main

import (
	"AAA/src/server"
	"AAA/src/test"
	"flag"
)

// go run . -test=10
func main() {
	port := flag.String("port", "8888", "port number")
	dummySize := flag.Int("test", 0, "dummy size")
	flag.Parse()

	if *dummySize > 0 {
		test.Run(*port, *dummySize)
	}

	server.Run(*port)
}
