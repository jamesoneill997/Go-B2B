package main

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/jamesoneill997/Go-B2B/structs"
)

func main() {
	var userHasAccount string
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing: ", err)
	}
	// Synchronous call
	args := &structs.Customer{
		Password: "password",
	}

	fmt.Print("Welcome to B2B-CLI, Do you already have an account? (y/n): ")
	fmt.Scanf("%s", &userHasAccount)

	if userHasAccount == "y" {
		fmt.Println("Welcome back!")
		//login prompt
	} else {
		fmt.Println("Nice to meet you!")
		//create account prompt
	}

	var reply string
	err = client.Call("Customer.CreateCustomer", args, &reply)
	fmt.Println(reply)

	if err != nil {
		log.Fatal("Error:", err)
	}

}
