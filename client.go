package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"github.com/jamesoneill997/Go-B2B/structs"
)

func main() {

	var (
		userHasAccount string
		reply          string
		customer       structs.Customer
		cmds           string
		products       []structs.Product
		availability   []string
	)

	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	login(client, userHasAccount, customer, reply)
	customer.ID, _ = strconv.Atoi(strings.Split(reply, " ")[len(reply)])

	fmt.Println("Welcome to go-B2B, type 'help' for a full list of commands")

	for {
		fmt.Print(">> ")
		fmt.Scanf("%v", &cmds)
		commands := strings.Split(cmds, " ")

		switch commands[0] {
		case "help":
			printHelp()
		case "login":
			login(client, userHasAccount, customer, reply)
		case "logout":
			logout()
		case "placeorder":
			var productID int
			var quantity int
			var requestedDate string

			fmt.Print("\n>> Please enter the product ID: ")
			fmt.Scanf("%d", &productID)
			fmt.Print("\n>> Please enter the quantity that you would like to order: ")
			fmt.Scanf("%d", &quantity)
			fmt.Print("\n>> Please enter the date that you would like to place the order for (Note strict formatting) dd/mm/yyyy: ")
			fmt.Scanf("%s", &requestedDate)

			tmpDate := strings.Split(requestedDate, "/")

			date := structs.Date{
				D: tmpDate[0],
				M: tmpDate[1],
				Y: tmpDate[2],
			}

			order := structs.Order{
				ProductID:  productID,
				CustomerID: customer.ID,
				Quantity:   quantity,
				Date:       date,
			}

			placeOrder(client, customer, order, reply)

		case "listproducts":
			listProducts(client, customer, &products)

		case "listorders":
			listOrders(client, customer, &availability)

		case "cancelorder":
			cancelOrder(client, customer, reply)

		case "showavailability":
			var withDate string
			stockReq := structs.StockRequest{}
			var requestedDate string

			fmt.Print("\n>>Would you like to know availability for a specific date (y/n)")
			fmt.Scanf("%s", &withDate)
			if withDate == "y" {
				fmt.Print("\n>> Please enter the product ID: ")
				fmt.Scanf("%d", &stockReq.ProductID)

				fmt.Print("\n>> Please enter the date that you would like to check availability for (Note strict formatting) dd/mm/yyyy: ")
				fmt.Scanf("%s", &requestedDate)

				tmpDate := strings.Split(requestedDate, "/")

				date := structs.Date{
					D: tmpDate[0],
					M: tmpDate[1],
					Y: tmpDate[2],
				}

				stockReq.Date = date

				showAvailabilityWithDate(client, stockReq, &reply)
			} else {
				fmt.Print("\n>> Please enter the product ID: ")
				fmt.Scanf("%d", &stockReq.ProductID)

				showAvailabilityWithoutDate(client, stockReq.ProductID, &availability)
			}
		}

	}

}

func printHelp() {
	validCommands := []structs.Command{
		{
			Name:        "help",
			Description: "Print a list of all commands and their usage",
		},
		{
			Name:        "login",
			Description: "Login as a user, this command will prompt you for ID and password",
		},
		{
			Name:        "logout",
			Description: "End current session",
		},
		{
			Name:        "placeorder",
			Description: "Places an order takes 3 parameters, productID, quantity and date",
		},
		{
			Name:        "listproducts",
			Description: "Print a list of all available products",
		},
		{
			Name:        "listorders",
			Description: "Print a list of all of your orders",
		},
		{
			Name:        "cancelorder",
			Description: "Cancel an order by ID",
		},
		{
			Name:        "showavailability",
			Description: "Print a list of availability for a given product, optional parameter Date will show availability on that date. No parameter will show availability for next 6 months",
		},
	}
	fmt.Println("go-B2B Assistance. Please see list of valid commands below")
	fmt.Println("Command                                                 Description")
	for _, cmd := range validCommands {
		fmt.Printf("%s%s%s\n", cmd.Name, strings.Repeat(" ", 40-len(cmd.Name)), cmd.Description)
	}
}

func login(client *rpc.Client, userHasAccount string, customer structs.Customer, reply string) {
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

		err := client.Call("Customer.CreateCustomer", customer, &reply)
		fmt.Println(reply)

		if err != nil {
			log.Fatal(err)
		}

	}
}

func logout() {
	fmt.Println("Thank you for shopping with go-B2B, see you again soon!")
	os.Exit(1)
}

func placeOrder(client *rpc.Client, customer structs.Customer, order structs.Order, reply string) {
	err := client.Call("Customer.MakeOrder", order, &reply)
	fmt.Println(reply)

	//extract order id
	orderID, err := strconv.Atoi(strings.Fields(reply)[len(strings.Fields(reply))-1])
	order.ID = orderID

	if err != nil {
		log.Fatal(err)
	}
}

func listProducts(client *rpc.Client, customer structs.Customer, reply *[]structs.Product) {
	err := client.Call("Customer.ListProducts", customer, reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ID             Name                Quantity")
	for _, prod := range *reply {
		fmt.Printf("%d %s %s %s %d\n", prod.ID, strings.Repeat(" ", 10), prod.Name, strings.Repeat(" ", 20-len(prod.Name)), prod.Quantity)
	}

}

func listOrders(client *rpc.Client, customer structs.Customer, availability *[]string) {
	err := client.Call("Customer.ListOrders", customer.ID, &availability)
	fmt.Println(*availability)
	if err != nil {
		log.Fatal(err)
	}

}

func cancelOrder(client *rpc.Client, customer structs.Customer, reply string) {
	var id int
	fmt.Print("Please enter the ID of the order that you would like to cancel: ")
	fmt.Scanf("%d", &id)

	err := client.Call("Customer.CancelOrder", id, &reply)
	fmt.Println(reply)

	if err != nil {
		fmt.Println(err)
	}

}

func showAvailabilityWithDate(client *rpc.Client, stockReq structs.StockRequest, reply *string) {
	err := client.Call("Customer.GetProjectionFromDate", stockReq, &reply)
	fmt.Println(*reply)

	if err != nil {
		fmt.Println(err)
	}
}

func showAvailabilityWithoutDate(client *rpc.Client, id int, availability *[]string) {
	err := client.Call("Customer.GetProjections", 2, &availability)
	fmt.Println(*availability)

	if err != nil {
		fmt.Println(err)
	}
}
