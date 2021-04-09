package main

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/jamesoneill997/Go-B2B/structs"
	"github.com/urfave/cli"
)

func main() {
	var (
		userHasAccount string
		reply          string
		customer       structs.Customer
		cmd            string
		cmdArr         []string
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
	app := cli.NewApp()
	app.Name = "B2B-CLI"
	app.Usage = "A B2B ordering system created in Go!"

	orderFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "product",
			Usage: "Specify desired Product",
		},
		cli.StringFlag{
			Name:  "date",
			Usage: "Specify date of order (format = dd/mm/yyyy)",
		},
		cli.StringFlag{
			Name:  "time",
			Usage: "Specify time of order (24hr format = 14:30)",
		},
	}

	credFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "Customer ID",
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "Customer password",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "logout",
			Usage: "Ends current session",
			Action: func(c *cli.Context) error {
				log.Fatal("Goodbye!")
				return err
			},
		},

		{
			Name:  "login",
			Usage: "Prompts user login.",
			Flags: credFlags,
		},
		{
			Name:  "order",
			Usage: "Order a specified product",
			Flags: orderFlags,
		},
		{
			Name:  "listproducts",
			Usage: "List all available products",
			Action: func(c *cli.Context) error {
				err := client.Call("Customer.ListProducts", customer, &reply)
				fmt.Println(&reply)
				return err
			},
		},
		{
			Name:  "availability",
			Usage: "List all availability, optional flags allow product and date specification",
			Flags: orderFlags,
		},
		{
			Name:  "listorders",
			Usage: "List all orders for current user",
		},
		{
			Name:  "cancelorder",
			Usage: "Cancel order given an ID, run without flags will cancel all orders for current user",
			Flags: orderFlags,
		},
	}

	if err != nil {
		log.Fatal("Error:", err)
	}

	//CLI
	app.Run([]string{"help"})

	for {
		fmt.Print(">> ")
		fmt.Scanf("%s", &cmd)
		cmdArr = append(cmdArr, cmd)
		err := app.Run(cmdArr)
		if err != nil {
			log.Fatal(err)
		}
	}

}
