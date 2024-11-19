package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/login", login)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
