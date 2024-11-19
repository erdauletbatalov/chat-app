package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

type Room struct {
	Name     string
	Clients  map[*websocket.Conn]string
	Messages []Message
}

var rooms = make(map[string]*Room)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Упростим проверку для локальной разработки
}
var clients = make(map[*websocket.Conn]string)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Получаем имя пользователя и комнату из параметров
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Anonymous"
	}

	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		roomName = "general"
	}

	// Создаем комнату, если её нет
	room, exists := rooms[roomName]
	if !exists {
		room = &Room{
			Name:    roomName,
			Clients: make(map[*websocket.Conn]string),
		}
		rooms[roomName] = room
	}

	// Добавляем клиента в комнату
	room.Clients[conn] = username

	// Уведомляем пользователей в комнате
	broadcastToRoom(room, fmt.Sprintf("User %s has joined the room!", username))

	// Отправляем историю сообщений
	for _, msg := range room.Messages {
		conn.WriteJSON(msg)
	}

	// Чтение сообщений
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("User %s disconnected: %v\n", username, err)
			break
		}

		// Добавляем сообщение в историю
		msg.Username = username
		room.Messages = append(room.Messages, msg)

		// Рассылаем сообщение всем в комнате
		broadcastToRoom(room, fmt.Sprintf("%s: %s", msg.Username, msg.Content))
	}

	// Удаляем клиента при выходе
	delete(room.Clients, conn)
	broadcastToRoom(room, fmt.Sprintf("User %s has left the room!", username))
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

func broadcastToRoom(room *Room, message string) {
	for client := range room.Clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("Write error:", err)
			client.Close()
			delete(room.Clients, client)
		}
	}
}
