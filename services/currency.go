package services

import (
	"encoding/json"
	"net/http"
)

func CurrencyHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"USD": 12700,
		"EUR": 13800,
		"BTC": 68000,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
