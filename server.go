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
	"sort"
	"strconv"
	"sync"
)

//Customer struct stores customer info
type Customer struct {
	ID       int    `json: id`
	Password string `json: password`
}

//Order struct stores order info
type Order struct {
	ID         int  `json: id`
	ProductID  int  `json: productID`
	Quantity   int  `json: quantity`
	CustomerID int  `json: customerID`
	Date       Date `json: date`
}

//Product struct stores product info
type Product struct {
	ID              int    `json: id`
	Name            string `json: name`
	Quantity        int    `json: quantity`
	RestockDate     int    `json: restockDate`
	RestockQuantity int    `json: restockQuantity`
}

//Date struct is used to store dates for orders
type Date struct {
	D string `json: d`
	M string `json: m`
	Y string `json: y`
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

func (cust *Customer) ListProducts(customerDetails *Customer, response *[]Product) error {
	mu.Lock()
	defer mu.Unlock()
	var listings []Product

	products, err := json.Marshal(readProducts())
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(products), &listings); err != nil {
		return err
	}

	*response = listings

	return nil
}

// func prettifyProduct(json []byte) string {

// }

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

func (cust *Customer) MakeOrder(orderDetails *Order, response *string) error {
	mu.Lock()
	defer mu.Unlock()
	orderDetails.ID = len(readOrders()) + 1

	fmt.Println("Order request received.")

	orders := readOrders()
	orders = append(orders, *orderDetails)

	jsonOrders, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("data/orders.json", jsonOrders, 0644)

	if err != nil {
		return err
	}

	file, _ := ioutil.ReadFile("data/orders.json")

	ords := getCurrentOrders(2, file)
	for _, ord := range ords {
		fmt.Println(ord.Date)
	}

	*response = fmt.Sprintf("Order successfully! Your Order ID is %d", orderDetails.ID)
	fmt.Printf("Order %d created successfully\n", orderDetails.ID)
	return nil

}

/*GetProjections will provide 6 month availability to client*/
// func (*Customer) GetProjections(product *Product, response *string) {

// }

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

//getCurrentStock takes in a product ID and returns the current level of stock
func getCurrentStock(id int) int {
	file, _ := ioutil.ReadFile("data/products.json")
	data := []Product{}
	_ = json.Unmarshal([]byte(file), &data)

	return 0

}

//getCurrentOrders takes in a product ID and returns the current orders for that product
func getCurrentOrders(id int, file []byte) []Order {

	data := []Order{}
	newData := []Order{}
	_ = json.Unmarshal([]byte(file), &data)

	//remove orders for other products
	for _, ord := range data {
		if ord.ProductID == id {
			newData = append(newData, ord)
		}
	}
	sort.Slice(newData, func(i, j int) bool {

		iStamp, err := strconv.Atoi(fmt.Sprintf("%s%s%s", newData[i].Date.Y, newData[i].Date.M, newData[i].Date.D))
		jStamp, err := strconv.Atoi(fmt.Sprintf("%s%s%s", newData[j].Date.Y, newData[j].Date.M, newData[j].Date.D))

		if err != nil {
			log.Fatal(err)
		}

		return iStamp < jStamp
	})

	return newData
}

//getCurrentRestock takes in a product ID and returns the day of the month that that product is restocked
func getCurrentRestock(id int) int {
	return 0
}

//helper function for removing slice element
func remove(orderSlice []Order, index int) []Order {
	return append(orderSlice[:index], orderSlice[index+1:]...)
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
