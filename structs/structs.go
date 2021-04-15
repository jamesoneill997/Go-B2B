package structs

type Customer struct {
	ID       int    `json: id`
	Password string `json: password`
}

//Product struct stores product info
type Product struct {
	ID              int    `json: id`
	Name            string `json: name`
	Quantity        int    `json: quantity`
	RestockDate     int    `json: restockDate`
	RestockQuantity int    `json: restockQuantity`
}

//Order struct stores order info
type Order struct {
	ID         int  `json: id`
	ProductID  int  `json: productID`
	Quantity   int  `json: quantity`
	CustomerID int  `json: customerID`
	Date       Date `json: date`
}

//Date struct is used to store dates for orders
type Date struct {
	D string `json: d`
	M string `json: m`
	Y string `json: y`
}

//StockRequest is used to store data used to query for stock availability
type StockRequest struct {
	Date      Date
	ProductID int
}

//Command type used for storing commands from CLI
type Command struct {
	Name        string
	Args        string
	Description string
}
