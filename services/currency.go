package services

import (
	"encoding/json"
	"net/http"
)

type CurrencyRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type CurrencyResponse struct {
	Result float64 `json:"result"`
}

func CurrencyHandler(w http.ResponseWriter, r *http.Request) {

	var req CurrencyRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var result float64

	switch {
	case req.From == "USD" && req.To == "RUB":
		result = req.Amount * 90

	case req.From == "RUB" && req.To == "USD":
		result = req.Amount / 90

	case req.From == "USD" && req.To == "EUR":
		result = req.Amount * 0.92

	case req.From == "EUR" && req.To == "USD":
		result = req.Amount / 0.92

	default:
		http.Error(w, "Неподдерживаемая валюта", http.StatusBadRequest)
		return
	}

	data := CurrencyResponse{
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// {
//   "from": "USD",
//   "to": "RUB",				//Пример конвертации валюты из USD в RUB
//   "amount": 100
// }
