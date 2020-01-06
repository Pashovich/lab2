package main

import (
	"flag"
	"sync"

	"./client"
	"./server"
)

var (
	port         = flag.Int("port", 9090, "port")
	clientNumber = flag.Int("n", 2, "clientNumber")
	serverStatus = flag.Bool("server", false, "start server")
	waitGroup    sync.WaitGroup
)

func main() {
	flag.Parse()
	if *serverStatus {
		waitGroup.Add(1)
		go server.StartServer(*port, *clientNumber)
	} else {
		waitGroup.Add(1)
		go client.StartClient(*port)
	}
	waitGroup.Wait()
}
