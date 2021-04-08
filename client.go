package main

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/jamesoneill997/Go-B2B/structs"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing: ", err)
	}
	// Synchronous call
	args := &structs.Customer{
		Password: "password",
	}

	var reply string
	err = client.Call("Customer.CreateCustomer", args, &reply)

	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Printf("Customer created, id: %d\n", args.ID)

}
