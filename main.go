package main

import (
	"github.com/jamesoneill997/Go-B2B/client"
	"github.com/jamesoneill997/Go-B2B/server"
)

func main() {
	server.StartServer()
	client.StartClient()
}
