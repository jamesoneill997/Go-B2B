package structs

type Customer struct {
	ID       int    `json: id`
	Password string `json: password`
}
