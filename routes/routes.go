package routes

import (
	"currency/handlers"
	"currency/services"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {

	r := mux.NewRouter()

	// HOME

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// AUTH

	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")

	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// USERS

	r.HandleFunc("/users", handlers.Getusers).Methods("GET")

	// WALLETS

	r.HandleFunc("/wallets", services.GetWallets).Methods("GET")

	// TOP UP

	r.HandleFunc("/topUp", services.TopUpHandler).Methods("POST")

	// CONVERT

	r.HandleFunc("/convert", services.CurrencyHandler).Methods("POST")

	// MARKET

	r.HandleFunc("/market/order", services.CreateMarketOrder).Methods("POST")

	r.HandleFunc("/market/buy", services.BuyMarketOrder).Methods("POST")

	r.HandleFunc("/markets", services.GetMarketOrders).Methods("GET")

	// TRANSACTIONS

	r.HandleFunc("/transactions", services.GetTransactions).Methods("GET")
	r.HandleFunc(
		"/rates",
		services.GetRates,
	).Methods("GET")
	r.HandleFunc(
		"/market/cancel",
		services.CancelMarketOrder,
	).Methods("POST")
	return r
}
