package client

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/jamesoneill997/Go-B2B/server"
)

func StartClient() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing: ", err)
	}
	// Synchronous call
	args := &server.Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

}
