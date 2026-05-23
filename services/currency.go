package services

import (
	"currency/database"
	"encoding/json"
	"fmt"
	"net/http"
)

// type Currency struct {
// }

type CurrencyRequest struct {
	UserId  int     `json:"user_id"`
	FromId  int     `json:"from"`
	ToId    int     `json:"to"`
	Balance float64 `json:"balance"`
	Amount  float64 `json:"amount"`
}

type CurrencyResponse struct {
	Result float64 `json:"result"`
}

func CurrencyHandler(w http.ResponseWriter, r *http.Request) {

	var req CurrencyRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var balance float64
	database.DB.Raw("SELECT balance FROM users WHERE id = ?", req.UserId).Scan(&balance)
	fmt.Printf("User balance: %f\n", balance)
	w.Header().Set("Content-Type", "application/json")

	// Проверяем, достаточно ли средств для конвертации
	if balance < req.Amount {
		http.Error(w, "Недостаточно средств для конвертации", http.StatusBadRequest)
		return
	}

	var coefficient float64

	database.DB.Raw(
		"SELECT coefficient FROM exchanges_values WHERE from_id = ? AND to_id = ?",
		req.FromId,
		req.ToId,
	).Scan(&coefficient)
	fmt.Printf("Exchange coefficient: %f\n", coefficient)

	result := req.Amount * coefficient
	fmt.Printf("Converted amount: %f\n", result)

	balance -= req.Amount
	database.DB.Exec(
		"UPDATE users SET balance = ? WHERE id = ?",
		balance,
		req.UserId,
	)

	fmt.Printf("Updated user balance: %f\n", balance)

	response := CurrencyResponse{Result: result}
	json.NewEncoder(w).Encode(response)
}

// { Пример запроса на конвертацию валюты:
//   "user_id": 1,
//   "from": 1,
//   "to": 2,
//   "amount": 145
// }
