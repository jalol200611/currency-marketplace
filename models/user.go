package models

type User struct {
	ID        int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	FirstName string `json:"first_name" gorm:"column:first_name"`
	LastName  string `json:"last_name" gorm:"column:last_name"`
	Email     string `json:"email" gorm:"column:email"`
	Phone     string `json:"phone" gorm:"column:phone"`
	Password  string `json:"password" gorm:"column:password"`
	// Wallet_Id int    `json:"wallet_id" gorm:"column:wallet_id"`
}

type Currency struct {
	ID       int     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Code     string  `json:"code" gorm:"column:code"`
	Name     string  `json:"name" gorm:"column:name"`
	Symbol   string  `json:"symbol" gorm:"column:symbol"`
	Exchange float64 `json:"exchange" gorm:"column:exchange"`
	IsCrypto bool    `json:"crypto" gorm:"column:crypto"`
}
type TopUpRequest struct {
	UserId     int     `json:"user_id" gorm:"column:user_id"`
	CurrencyId int     `json:"currency_id" gorm:"column:currency_id"`
	Amount     float64 `json:"amount"  gorm:"column:amount"`
}
