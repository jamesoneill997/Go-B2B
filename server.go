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
	"strings"
	"sync"

	"github.com/jamesoneill997/Go-B2B/structs"
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

type StockRequest struct {
	Date      Date
	ProductID int
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

	*response = fmt.Sprintf("Order successfully! Your Order ID is %d", orderDetails.ID)
	fmt.Printf("Order %d created successfully\n", orderDetails.ID)
	return nil

}

/*GetProjections will provide 6 month availability to client*/
func (*Customer) GetProjections(id int, response *[]string) error {
	mu.Lock()
	defer mu.Unlock()

	availability := projectionFetcher(id)

	*response = availability

	return nil

}

/*GetProjectionFromDate takes in a product and a date and returns the availability of that product on that date*/
func (*Customer) GetProjectionFromDate(query StockRequest, response *string) error {
	mu.Lock()
	defer mu.Unlock()

	availability := projectionFetcher(query.ProductID)
	for _, a := range availability {
		s := strings.Split(a, ":")
		d, quantity := strings.Trim(s[0], "{"), strings.Trim(s[1], " ")
		d = strings.Trim(d, "}")

		date := strings.Split(d, " ")

		thisDate := Date{date[0], date[1], date[2]}

		if query.Date.M == thisDate.M {
			*response = fmt.Sprintf("The projected quantity for %v is %s", date, quantity)
			return nil
		}
	}

	*response = "Projection not available"
	return errors.New("Projection not available")

}

/*ListOrders function will list all orders made by the specified customer*/
func (*Customer) ListOrders(custID int, response *[]string) error {
	mu.Lock()
	defer mu.Unlock()
	ords := readOrders()
	custOrders := []Order{}
	strOrders := []string{}

	for _, ord := range ords {
		if ord.CustomerID == custID {
			custOrders = append(custOrders, ord)
			strOrders = append(strOrders, fmt.Sprintf("\n Order ID: %d \n Product ID: %d \n Product Name: %s \n Quantity: %d \n Date: %v\n", ord.ID, ord.ProductID, getProductName(ord.ProductID), ord.Quantity, ord.Date))
		}
	}

	*response = strOrders

	return nil

}

/*CancelOrder takes in an order ID and deletes that order.*/
func (*Customer) CancelOrder(orderID int, response *string) error {
	mu.Lock()
	defer mu.Unlock()

	orders := readOrders()

	for i, ord := range orders {
		if ord.ID == orderID {
			orders = remove(orders, i)
			jsonOrders, err := json.Marshal(orders)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile("data/orders.json", jsonOrders, 0644)

			if err != nil {
				return err
			}
			fmt.Printf("Order %d has been cancelled", orderID)
			*response = "Successfully cancelled order"
			return nil
		}
	}

	*response = "Order not found"
	return errors.New("Order not found")
}

/*

   UNEXPORTED HELPER FUNCTIONS

*/

func checkOrders(month int, ords []Order) []Order {
	ordersThisMonth := []Order{}

	for _, ord := range ords {
		intMonth, _ := strconv.Atoi(ord.Date.M)

		if intMonth == month {
			ordersThisMonth = append(ordersThisMonth, ord)
		}
	}

	return ordersThisMonth
}

func projectionFetcher(id int) []string {
	file, _ := ioutil.ReadFile("data/orders.json")
	ords := getCurrentOrders(id, file)
	availability := []string{}
	restocks := []structs.Date{}
	currStock := getCurrentStock(id)

	restockDay, restockQuantity := getCurrentRestock(id)

	for i := 4; i <= 10; i++ {
		restockDate := structs.Date{
			D: strconv.Itoa(restockDay),
			M: fmt.Sprintf("0%d", i),
			Y: "2021",
		}

		restocks = append(restocks, restockDate)
	}
	for i := 4; i <= 10; i++ {

		ordersThisMonth := checkOrders(i, ords)

		if len(ordersThisMonth) == 0 {
			currStock += restockQuantity
			availability = append(availability, fmt.Sprintf("\n%s: %d", restocks[i-4], currStock))
		} else {

			for _, order := range ordersThisMonth {

				if order.Date.D < strconv.Itoa(restockDay) {
					currStock -= order.Quantity
					availability = append(availability, fmt.Sprintf("\n%s: %d", order.Date, currStock))
				} else {
					restockDate := structs.Date{
						D: strconv.Itoa(restockDay),
						M: strconv.Itoa(i + 3),
						Y: "2021",
					}

					currStock += restockQuantity - order.Quantity
					availability = append(availability, fmt.Sprintf("\n%s: %d", restockDate, currStock))
				}
			}
		}
	}

	return availability
}

//checks a users password
func checkPassword(enteredPassword string, correctPassword string) error {
	if enteredPassword == correctPassword {
		return nil
	}
	return errors.New("Invalid credentials")
}

//get product name given ID
func getProductName(id int) string {
	file, _ := ioutil.ReadFile("data/products.json")
	data := []Product{}
	_ = json.Unmarshal([]byte(file), &data)

	for _, prod := range data {
		if prod.ID == id {
			return prod.Name
		}
	}

	return "Product name not found"
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

	for _, prod := range data {
		if prod.ID == id {
			return prod.Quantity
		}
	}

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

//getCurrentRestock takes in a product ID and returns (restock date, restock quantity)
func getCurrentRestock(id int) (int, int) {
	file, _ := ioutil.ReadFile("data/products.json")
	data := []Product{}
	_ = json.Unmarshal([]byte(file), &data)

	for _, p := range data {
		if id == p.ID {
			return p.RestockDate, p.RestockQuantity
		}
	}

	//no product found, assume that there is no restock
	return 0, 0
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
