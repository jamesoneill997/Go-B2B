package main

import (
	"log"
	"net/rpc"

	"github.com/jamesoneill997/Go-B2B/structs"
)

func runClient(custs <-chan structs.Customer, results chan<- string) {
	var (
		reply string
	)

	for cust := range custs {
		cust.Password = "password"

		client, err := rpc.DialHTTP("tcp", "localhost:1234")
		if err != nil {
			log.Fatal("dialing: ", err)
		}

		err = client.Call("Customer.CreateCustomer", cust, &reply)
		results <- "Done"
	}
}

func main() {
	maxRoutines := 15
	numJobs := 100
	jobs := make(chan structs.Customer, numJobs)
	results := make(chan string, numJobs)

	defer close(jobs)

	var customers []structs.Customer

	//create 100 customers
	for i := 0; i < numJobs; i++ {
		c := structs.Customer{
			Password: "testPassword",
		}
		customers = append(customers, c)

	}

	for i := 0; i < maxRoutines; i++ {
		go runClient(jobs, results)
	}

	for _, cust := range customers {
		jobs <- cust
	}

	for a := 1; a < numJobs; a++ {
		<-results
	}
}
