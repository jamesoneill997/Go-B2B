package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"strconv"
	"strings"
	"time"

	"github.com/jamesoneill997/Go-B2B/structs"
)

func runClient(custs <-chan structs.Customer, results chan<- string) {
	var (
		reply        string
		availability []string
		// products []structs.Product
		// commands = []string{"ListProducts"}
	)

	for cust := range custs {
		time.Sleep(time.Duration(numberGenerator(100, 3000)) * time.Millisecond)
		cust.Password = "password"

		client, err := rpc.DialHTTP("tcp", "localhost:1234")
		if err != nil {
			log.Fatal("dialing: ", err)
		}

		err = client.Call("Customer.CreateCustomer", cust, &reply)

		//extract id
		id, err := strconv.Atoi(strings.Fields(reply)[len(strings.Fields(reply))-1])
		cust.ID = id

		fmt.Println(reply)

		order := structs.Order{
			ProductID:  numberGenerator(1, 4),
			Quantity:   numberGenerator(1, 150),
			CustomerID: cust.ID,
			Date:       dateGenerator(),
		}

		err = client.Call("Customer.MakeOrder", order, &reply)
		fmt.Println(reply)

		//extract order id
		orderID, err := strconv.Atoi(strings.Fields(reply)[len(strings.Fields(reply))-1])
		order.ID = orderID

		err = client.Call("Customer.ListOrders", cust.ID, &availability)
		fmt.Println(availability)

		err = client.Call("Customer.CancelOrder", orderID, &reply)
		fmt.Println(reply)

		err = client.Call("Customer.ListOrders", cust.ID, &availability)
		fmt.Println(availability)

		// err = client.Call("Customer.GetProjections", 2, &availability)
		// fmt.Println(availability)

		if err != nil {
			fmt.Println(err)
		}

		results <- "Done"
	}

}

func numberGenerator(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return (rand.Intn(max-min+1) + min)
}

func dateGenerator() structs.Date {
	d := strconv.Itoa(numberGenerator(1, 28))
	m := strconv.Itoa(numberGenerator(5, 12))
	dInt, _ := strconv.Atoi(d)
	mInt, _ := strconv.Atoi(m)

	if dInt < 10 {
		d = fmt.Sprintf("0%s", d)
	}

	if mInt < 10 {
		m = fmt.Sprintf("0%s", m)
	}

	date := structs.Date{
		D: d,
		M: m,
		Y: "2021",
	}
	return date
}

func main() {
	maxRoutines := 15
	numJobs := 10
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

	//TODO refactor to use sync.WaitGroup
	for a := 1; a < numJobs; a++ {
		<-results
	}
}
