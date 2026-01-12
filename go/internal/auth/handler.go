package auth

import "net/http"

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	login := "Welcome to the page";
	w.Write([]byte("Logged in: " + login))
}