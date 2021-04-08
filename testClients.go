package main

import (
	"log"
	"net/rpc"
	"sync"

	"github.com/jamesoneill997/Go-B2B/structs"
)

func runClient(cust structs.Customer, wg *sync.WaitGroup) {
	var (
		reply string
	)

	cust.Password = "password"

	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	err = client.Call("Customer.CreateCustomer", cust, &reply)
	wg.Done()
}

func main() {
	wg := new(sync.WaitGroup)
	work := make(chan structs.Customer)

	var customers []structs.Customer

	//create 15 customers
	for i := 0; i < 15; i++ {
		c := structs.Customer{
			Password: "testPassword",
		}
		customers = append(customers, c)
	}

	for j := 0; j < 15; j++ {
		go func(w chan structs.Customer) {
			c := <-w
			wg.Add(1)
			runClient(c, wg)
		}(work)
	}

	for _, c := range customers {
		work <- c
	}

	wg.Wait()
}
