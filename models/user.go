package models

type User struct {
	ID                int        `json:"id"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Email             string     `json:"email"`
	Phone             string     `json:"phone"`
	TrackedCurrencies []Currency `json:"tracked_currencies"`
}

type Currency struct {
	ID       int     `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Exchange float64 `json:"exchange"`
	IsCrypto bool    `json:"crypto"`
}
