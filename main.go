package main

import (
	"net/http"

	"github.com/twaananen/boulderlog/internal/handlers"
)

func main() {
	http.HandleFunc("/auth/status", handlers.AuthStatus)
	http.HandleFunc("/auth/login", handlers.Login)
	http.HandleFunc("/auth/logout", handlers.Logout)
	http.ListenAndServe(":8080", nil)
}
