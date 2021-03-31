package structs

import "time"

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
