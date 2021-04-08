package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
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
	ID              int    `json: id`
	Name            string `json: name`
	Quantity        int    `json: quantity`
	RestockDate     int    `json: restockDate`
	RestockQuantity int    `json: restockQuantity`
}

type Response struct {
	Message string
}

var mu sync.Mutex

func (cust *Customer) CreateCustomer(customerDetails *Customer, response *string) error {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("Create customer request received.")

	customers := readCustomers()
	customerDetails.ID = generateID()

	customers = append(customers, *customerDetails)

	jsonCustomers, err := json.Marshal(customers)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("data/customers.json", jsonCustomers, 0644)

	if err != nil {
		return err
	}

	*response = fmt.Sprintf("Account created successfully! Your customer ID is %d", customerDetails.ID)
	fmt.Printf("Customer %d created successfully\n", customerDetails.ID)
	return nil
}

/*Login takes customer login details and returns whether login has been successful*/
func (cust *Customer) Login(customerDetails *Customer, response *string) error {
	registeredCustomers := readCustomers()
	enteredID := customerDetails.ID
	enteredPassword := customerDetails.Password

	for _, cust := range registeredCustomers {
		if enteredID == cust.ID {
			if err := checkPassword(enteredPassword, cust.Password); err != nil {
				*response = "Invalid Credentials"
				return err
			}

			*response = "Login successful"
			return nil
		}
	}

	*response = "Customer not found."
	return errors.New(*response)
}

/*				Unexported Helper functions				*/

//checks a users password
func checkPassword(enteredPassword string, correctPassword string) error {
	if enteredPassword == correctPassword {
		return nil
	}
	return errors.New("Invalid credentials")
}

/*generates unique ID for user*/
func generateID() int {
	return len(readCustomers()) + 1
}

//gets list of all current customers
func readCustomers() []Customer {
	file, _ := ioutil.ReadFile("data/customers.json")
	data := []Customer{}
	_ = json.Unmarshal([]byte(file), &data)

	return data
}

//gets list of all current orders
func readOrders() []Order {
	file, _ := ioutil.ReadFile("data/orders.json")
	data := []Order{}
	_ = json.Unmarshal([]byte(file), &data)

	return data
}

//gets list of all current products
func readProducts() []Product {
	file, _ := ioutil.ReadFile("data/products.json")
	data := []Product{}
	_ = json.Unmarshal([]byte(file), &data)

	return data
}

/*							Main 							*/
func main() {
	customer := new(Customer)
	rpc.Register(customer)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		http.Serve(l, nil)

	}
}
