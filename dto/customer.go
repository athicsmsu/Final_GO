package dto

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"pass"`
}

type CustomerDto struct {
	CustomerID  int    `json:"id"`
	FirstName   string `json:"fname"`
	LastName    string `json:"lname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone"`
	Address     string `json:"address"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
