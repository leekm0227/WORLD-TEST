package main

import (
	"flag"
	server "go-server/src"
)

// GOGC=150 go run . -test=100
func main() {
	port := flag.String("port", "8888", "port number")
	flag.Parse()
	server.Run(*port)
}
