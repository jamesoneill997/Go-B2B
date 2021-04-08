package main

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/jamesoneill997/Go-B2B/structs"
)

func main() {
	var (
		userHasAccount string
		reply          string
		customer       structs.Customer
	)

	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	fmt.Print("Welcome to B2B-CLI, Do you already have an account? (y/n): ")
	fmt.Scanf("%s", &userHasAccount)

	if userHasAccount == "y" {
		fmt.Println("Welcome back!")
		fmt.Print("Please enter your ID: ")
		fmt.Scanf("%d", &customer.ID)

		fmt.Print("Please enter your Password: ")
		fmt.Scanf("%s", &customer.Password)

		err := client.Call("Customer.Login", customer, &reply)

		if err != nil {
			log.Fatal(err)
		}

	} else {
		fmt.Println("Nice to meet you!")
		fmt.Print("Please enter the password that you would like to use: ")
		fmt.Scanf("%s", &customer.Password)

		err = client.Call("Customer.CreateCustomer", customer, &reply)
	}

	fmt.Println(reply)

	if err != nil {
		log.Fatal("Error:", err)
	}

}
