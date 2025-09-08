// order-service/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Order struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ProductName string `json:"product_name"`
}

var (
	orders = map[int]Order{}
	nextOrderID = 1
	mu sync.Mutex // Мьютекс для защиты доступа к orders и nextOrderID
)

// Здесь мы имитируем обращение к User Service.
// В реальном приложении использовался бы HTTP-клиент.
func isUserExists(userID int) bool {
	// В реальном приложении здесь был бы HTTP-запрос к user-service
	// Для простоты, пока будем считать, что пользователь существует, если его ID >= 1.
	// Позже мы можем добавить реальный вызов.
	return userID >= 1
}

func main() {
	http.HandleFunc("/orders", handleOrders)
	http.HandleFunc("/orders/user/", handleUserOrders) // Для получения заказов пользователя

	fmt.Println("Order Service listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mu.Lock()
		orderList := make([]Order, 0, len(orders))
		for _, order := range orders {
			orderList = append(orderList, order)
		}
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orderList)

	case http.MethodPost:
		var newOrder Order
		err := json.NewDecoder(r.Body).Decode(&newOrder)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Проверяем, существует ли пользователь
		if !isUserExists(newOrder.UserID) {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}

		mu.Lock()
		newOrder.ID = nextOrderID
		orders[newOrder.ID] = newOrder
		nextOrderID++
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newOrder)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleUserOrders(w http.ResponseWriter, r *http.Request) {
	// Извлекаем User ID из URL, например, /orders/user/123
	idStr := r.URL.Path[len("/orders/user/"):]
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		mu.Lock()
		userOrderList := []Order{}
		for _, order := range orders {
			if order.UserID == userID {
				userOrderList = append(userOrderList, order)
			}
		}
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userOrderList)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
