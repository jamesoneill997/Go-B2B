package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

//Customer struct stores customer info
type Customer struct {
	ID       int    `json: id`
	Password string `json: password`
}

//Order struct stores order info
type Order struct {
	ID         int       `json: id`
	ProductID  int       `json: productID`
	CustomerID int       `json: customerID`
	Date       time.Time `json: date`
}

//Product struct stores product info
type Product struct {
	ID              int       `json: id`
	Name            string    `json: name`
	Quantity        int       `json: quantity`
	RestockDate     time.Time `json: restockDate`
	RestockQuantity int       `json: restockQuantity`
}

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func readCustomers() []Customer {
	file, _ := ioutil.ReadFile("data/customers.json")
	data := []Customer{}

	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (cust *Customer) CreateCustomer(customerDetails *Customer, response *string) error {
	customers := readCustomers()
	customers = append(customers, *customerDetails)

	jsonCustomers, err := json.Marshal(customers)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("data/customers.json", jsonCustomers, 0644)

	if err != nil {
		return err
	}

	return nil
}

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func StartServer() {
	customer := new(Customer)
	rpc.Register(customer)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
