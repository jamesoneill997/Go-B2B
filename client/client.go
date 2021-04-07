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
	args := &server.Customer{
		ID:       1,
		Password: "password",
	}

	var reply string
	err = client.Call("Customer.CreateCustomer", args, &reply)

	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Printf("Customer created, id: %d\n", args.ID)

}
