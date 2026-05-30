package handlers

import (
	"currency/database"
	"currency/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Сервер валют работает! Добро пожаловать!")
}

var Users = make(map[int]models.User)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	database.DB.Create(&user)

	database.DB.Exec(
		`
		INSERT INTO wallets(user_id, currency_id, balance)
		VALUES
		($1, 1, 0),
		($1, 2, 0),
		($1, 3, 0)
		`,
		user.ID,
	)

	fmt.Printf(
		"Пользователь сохранён: %s %s | ID: %d\n",
		user.FirstName,
		user.LastName,
		user.ID,
	)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)
}

func Getusers(w http.ResponseWriter, r *http.Request) {

	var users []models.User

	database.DB.Raw("SELECT * FROM users").Scan(&users)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(users)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	var user models.User
	database.DB.Raw(
		"SELECT * FROM users WHERE email = ? AND password = ?",
		req.Email,
		req.Password,
	).Scan(&user)
	if user.ID == 0 {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}
	fmt.Printf("Пользователь вошёл: %s %s | Email: %s\n", user.FirstName, user.LastName, user.Email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
