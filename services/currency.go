package services

import (
	"currency/database"
	"currency/models"
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

	database.DB.Raw(
		`
		SELECT balance
		FROM wallets
		WHERE user_id = $1
		AND currency_id = $2
		`,
		req.UserId,
		req.FromId,
	).Scan(&balance)

	fmt.Printf("Wallet balance: %f\n", balance)

	if balance < req.Amount {
		http.Error(w, "Недостаточно средств", http.StatusBadRequest)
		return
	}

	var coefficient float64

	database.DB.Raw(
		`
		SELECT coefficient
		FROM exchanges_values
		WHERE from_id = $1
		AND to_id = $2
		`,
		req.FromId,
		req.ToId,
	).Scan(&coefficient)

	result := req.Amount * coefficient

	// списываем FROM валюту

	database.DB.Exec(
		`
		UPDATE wallets
		SET balance = balance - $1
		WHERE user_id = $2
		AND currency_id = $3
		`,
		req.Amount,
		req.UserId,
		req.FromId,
	)

	// начисляем TO валюту

	resultDB := database.DB.Exec(
		`
	UPDATE wallets
	SET balance = balance + $1
	WHERE user_id = $2
	AND currency_id = $3
	`,
		result,
		req.UserId,
		req.ToId,
	)

	fmt.Printf("ROWS UPDATED: %d\n", resultDB.RowsAffected)

	if resultDB.RowsAffected == 0 {

		fmt.Println("КОШЕЛЁК НЕ НАЙДЕН")

		database.DB.Exec(
			`
		INSERT INTO wallets(user_id, currency_id, balance)
		VALUES ($1, $2, $3)
		`,
			req.UserId,
			req.ToId,
			result,
		)

		fmt.Println("Создан новый кошелёк")
	}

	response := CurrencyResponse{
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func TopUpHandler(w http.ResponseWriter, r *http.Request) {
	var req models.TopUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var balance float64

	database.DB.Raw(
		"SELECT balance FROM wallets WHERE user_id = ? AND currency_id = ?",
		req.UserId,
		req.CurrencyId,
	).Scan(&balance)
	fmt.Printf("Баланс: %f\n", balance)
	balance = balance + req.Amount
	database.DB.Exec(
		"UPDATE wallets SET balance = ? WHERE user_id = ? AND currency_id = ?",
		balance,
		req.UserId,
		req.CurrencyId,
	)
	fmt.Printf("Ваш баланс пополнен на %f | Баланс после пополнения: %f\n", req.Amount, balance)
	notification := fmt.Sprintf("Ваш баланс пополнен на %f | Баланс после пополнения: %f", req.Amount, balance)
	json.NewEncoder(w).Encode(notification)
}

type MarketOrder struct {
	Id int `json:"id"`

	SellerId   int    `json:"seller_id"`
	SellerName string `json:"seller_name"`

	SellCurrId int `json:"sell_curr_id"`
	BuyCurrId  int `json:"buy_curr_id"`

	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

func CreateMarketOrder(w http.ResponseWriter, r *http.Request) {

	var order MarketOrder

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем баланс пользователя

	var walletBalance float64

	database.DB.Raw(
		"SELECT balance FROM wallets WHERE user_id = $1 AND currency_id = $2",
		order.SellerId,
		order.SellCurrId,
	).Scan(&walletBalance)

	fmt.Printf("Баланс пользователя: %f\n", walletBalance)

	// Проверяем хватает ли денег

	if walletBalance < order.Amount {
		http.Error(w, "Недостаточно средств", http.StatusBadRequest)
		return
	}

	// Блокируем валюту seller

	database.DB.Exec(
		`
	UPDATE wallets
	SET balance = balance - $1
	WHERE user_id = $2
	AND currency_id = $3
	`,
		order.Amount,
		order.SellerId,
		order.SellCurrId,
	)

	// Создаём ордер

	result := database.DB.Exec(
		`
		INSERT INTO markets(
			seller_id,
			sell_curr_id,
			buy_curr_id,
			amount,
			price
		)
		VALUES ($1, $2, $3, $4, $5)
		`,
		order.SellerId,
		order.SellCurrId,
		order.BuyCurrId,
		order.Amount,
		order.Price,
	)

	// Проверяем SQL ошибки

	if result.Error != nil {
		fmt.Println("SQL ERROR:", result.Error)
		http.Error(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	// Лог в терминал

	notificationOrder := fmt.Sprintf(
		"Создан ордер | Seller ID: %d | Sell Currency: %d | Buy Currency: %d | Amount: %f | Price: %f",
		order.SellerId,
		order.SellCurrId,
		order.BuyCurrId,
		order.Amount,
		order.Price,
	)

	fmt.Println(notificationOrder)

	// Ответ клиенту

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notificationOrder)
	json.NewEncoder(w).Encode(order)
}

type BuyRequest struct {
	OrderId  int    `json:"order_id"`
	BuyerId  int    `json:"buyer_id"`
	Password string `json:"password"`
}

func BuyMarketOrder(w http.ResponseWriter, r *http.Request) {

	var req BuyRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	var dbPassword string

	database.DB.Raw(
		`
	SELECT password
	FROM users
	WHERE id = $1
	`,
		req.BuyerId,
	).Scan(&dbPassword)
	if dbPassword != req.Password {

		http.Error(
			w,
			"Неверный пароль",
			http.StatusBadRequest,
		)

		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ищем ордер

	var order MarketOrder

	database.DB.Raw(
		`
SELECT
	id,
	seller_id,
	sell_curr_id,
	buy_curr_id,
	amount,
	price
FROM markets
WHERE id = $1
	`,
		req.OrderId,
	).Scan(&order)
	fmt.Printf("ORDER ID FROM REQUEST: %d\n", req.OrderId)
	fmt.Printf("ORDER FROM DB: %+v\n", order)
	if order.SellerId == 0 {
		http.Error(w, "Ордер не найден", http.StatusBadRequest)
		return
	}
	if order.SellerId == req.BuyerId {

		http.Error(
			w,
			"Нельзя купить свой собственный ордер",
			http.StatusBadRequest,
		)

		return
	}
	fmt.Printf("ORDER AFTER SCAN: %+v\n", order)
	// Проверяем баланс покупателя

	var buyerBalance float64

	database.DB.Raw(
		`
		SELECT balance
		FROM wallets
		WHERE user_id = $1
		AND currency_id = $2
		`,
		req.BuyerId,
		order.BuyCurrId,
	).Scan(&buyerBalance)

	if buyerBalance < order.Price {
		http.Error(w, "Недостаточно средств", http.StatusBadRequest)
		return
	}

	// Снимаем деньги у buyer

	database.DB.Exec(
		`
		UPDATE wallets
		SET balance = balance - $1
		WHERE user_id = $2
		AND currency_id = $3
		`,
		order.Price,
		req.BuyerId,
		order.BuyCurrId,
	)

	// Даём деньги seller

	database.DB.Exec(
		`
		UPDATE wallets
		SET balance = balance + $1
		WHERE user_id = $2
		AND currency_id = $3
		`,
		order.Price,
		order.SellerId,
		order.BuyCurrId,
	)

	// Даём валюту buyer

	database.DB.Exec(
		`
		UPDATE wallets
		SET balance = balance + $1
		WHERE user_id = $2
		AND currency_id = $3
		`,
		order.Amount,
		req.BuyerId,
		order.SellCurrId,
	)

	// // Забираем валюту seller

	// database.DB.Exec(
	// 	`
	// 	UPDATE wallets
	// 	SET balance = balance - $1
	// 	WHERE user_id = $2
	// 	AND currency_id = $3
	// 	`,
	// 	order.Amount,
	// 	order.SellerId,
	// 	order.SellCurrId,
	// )

	// Записываем транзакцию

	resultTransaction := database.DB.Exec(
		`
	INSERT INTO transactions(
		buyer_id,
		seller_id,
		sell_curr_id,
		buy_curr_id,
		amount,
		price
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	`,
		req.BuyerId,
		order.SellerId,
		order.SellCurrId,
		order.BuyCurrId,
		order.Amount,
		order.Price,
	)

	if resultTransaction.Error != nil {
		fmt.Println("TRANSACTION ERROR:", resultTransaction.Error)
	}
	// Удаляем ордер

	database.DB.Exec(
		"DELETE FROM markets WHERE id = $1",
		req.OrderId,
	)

	fmt.Println("Ордер успешно куплен")

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode("Ордер успешно куплен")
}

// ---Подключение к API для получения курсов валют---
type ExchangeApiResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func UpdateCurrencyRates() {

	resp, err := http.Get("https://open.er-api.com/v6/latest/USD")
	if err != nil {
		fmt.Println("API ERROR:", err)
		return
	}

	defer resp.Body.Close()

	var data ExchangeApiResponse

	json.NewDecoder(resp.Body).Decode(&data)

	fmt.Println("Курсы обновлены")

	// Карта:
	// code -> id

	currencies := map[string]int{
		"USD": 1,
		"EUR": 2,
		"RUB": 3,
		"BTC": 4,
		"TON": 5,
		"ETH": 6,
	}

	for code, id := range currencies {

		rate := data.Rates[code]

		fmt.Printf("%s = %f\n", code, rate)

		database.DB.Exec(
			`
			UPDATE exchanges_values
			SET coefficient = $1
			WHERE from_id = $2
			AND to_id = $3
			`,
			rate,
			1, // USD
			id,
		)
	}
}

type Transaction struct {
	Id int `json:"id"`

	BuyerId   int    `json:"buyer_id"`
	BuyerName string `json:"buyer_name"`

	SellerId   int    `json:"seller_id"`
	SellerName string `json:"seller_name"`

	SellCurrId int `json:"sell_curr_id"`
	BuyCurrId  int `json:"buy_curr_id"`

	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`

	CreatedAt string `json:"created_at"`
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {

	var transactions []Transaction

	database.DB.Raw(
		`
	SELECT

		t.id,

		t.buyer_id,
		b.first_name || ' ' || b.last_name AS buyer_name,

		t.seller_id,
		s.first_name || ' ' || s.last_name AS seller_name,

		t.sell_curr_id,
		t.buy_curr_id,

		t.amount,
		t.price,

		t.created_at

	FROM transactions t

	JOIN users b
	ON b.id = t.buyer_id

	JOIN users s
	ON s.id = t.seller_id
	`,
	).Scan(&transactions)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(transactions)
}

func GetMarketOrders(w http.ResponseWriter, r *http.Request) {

	var orders []MarketOrder
	database.DB.Raw(
		`
	SELECT

		m.id,

		m.seller_id,

		u.first_name || ' ' || u.last_name
		AS seller_name,

		m.sell_curr_id,
		m.buy_curr_id,

		m.amount,
		m.price

	FROM markets m

	JOIN users u
	ON u.id = m.seller_id
	`,
	).Scan(&orders)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(orders)
}

type Wallet struct {
	Id         int     `json:"id"`
	UserId     int     `json:"user_id"`
	CurrencyId int     `json:"currency_id"`
	Balance    float64 `json:"balance"`
}

func GetWallets(w http.ResponseWriter, r *http.Request) {

	var wallets []Wallet

	database.DB.Raw(
		"SELECT * FROM wallets",
	).Scan(&wallets)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(wallets)
}

type Rate struct {
	FromId      int     `json:"from_id"`
	ToId        int     `json:"to_id"`
	Coefficient float64 `json:"coefficient"`
}

func GetRates(
	w http.ResponseWriter,
	r *http.Request,
) {

	var rates []Rate

	database.DB.Raw(
		`
	SELECT
		from_id,
		to_id,
		coefficient
	FROM exchanges_values
	`,
	).Scan(&rates)
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(rates)
}

type CancelOrderRequest struct {
	OrderId  int    `json:"order_id"`
	UserId   int    `json:"user_id"`
	Password string `json:"password"`
}

func CancelMarketOrder(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req CancelOrderRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(&req)
	var dbPassword string

	database.DB.Raw(
		`
	SELECT password
	FROM users
	WHERE id = $1
	`,
		req.UserId,
	).Scan(&dbPassword)

	if dbPassword != req.Password {

		http.Error(
			w,
			"Неверный пароль",
			http.StatusBadRequest,
		)

		return
	}
	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	var order MarketOrder

	database.DB.Raw(
		`
		SELECT
			id,
			seller_id,
			sell_curr_id,
			amount
		FROM markets
		WHERE id = $1
		`,
		req.OrderId,
	).Scan(&order)

	if order.Id == 0 {

		http.Error(
			w,
			"Ордер не найден",
			http.StatusBadRequest,
		)

		return
	}

	if order.SellerId != req.UserId {

		http.Error(
			w,
			"Нельзя удалить чужой ордер",
			http.StatusBadRequest,
		)

		return
	}

	// Возвращаем деньги

	database.DB.Exec(
		`
		UPDATE wallets
		SET balance = balance + $1
		WHERE user_id = $2
		AND currency_id = $3
		`,
		order.Amount,
		order.SellerId,
		order.SellCurrId,
	)

	// Удаляем ордер

	database.DB.Exec(
		`
		DELETE FROM markets
		WHERE id = $1
		`,
		req.OrderId,
	)

	json.NewEncoder(w).Encode("Ордер отменён")
}

// { Пример запроса на конвертацию валюты:
//   "user_id": 1,
//   "from": 1,
//   "to": 2,
//   "amount": 145
// }
// { Пример запроса на пополнение баланса:
//   "user_id": 1,
//   "currency_id": 2,
//   "amount": 1000
// }
// { Пример запроса на создание рыночного ордера:
// {
//   "seller_id": 1,
//   "sell_curr_id": 3,
//   "buy_curr_id": 2,
//   "amount": 1,
//   "price": 7500000
// }
