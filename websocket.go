package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Упростим проверку для локальной разработки
}
var clients = make(map[*websocket.Conn]string)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Проверка поддержки WebSocket
	if !websocket.IsWebSocketUpgrade(r) {
		http.Error(w, "WebSocket not supported", http.StatusUpgradeRequired)
		return
	}

	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Token required", http.StatusUnauthorized)
		return
	}
	fmt.Println("Received token:", tokenString)

	username, err := validateJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Сохраняем подключение
	clients[conn] = username
	defer delete(clients, conn)

	// Уведомляем всех пользователей
	broadcast(fmt.Sprintf("User %s has joined the chat!", username))

	// Чтение сообщений
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("User %s disconnected: %v\n", username, err)
			break
		}
		broadcast(fmt.Sprintf("%s: %s", username, string(message)))
	}
}

func broadcast(message string) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("Write error:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
