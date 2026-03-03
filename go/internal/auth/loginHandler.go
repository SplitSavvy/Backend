package auth

// import (
// 	"encoding/json"
// 	"net/http"
// 	"splitsavvy/internal/database"
// 	"splitsavvy/internal/password"
// )

// type LoginRequest struct {
// 	Identifier string `json:"identifier"`
// 	Password   string `json:"password"`
// }

// func HandleLogin(w http.ResponseWriter, r *http.Request) {
// 	var req LoginRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid Request Format", http.StatusBadRequest)
// 		return
// 	}
// 	if len(req.Identifier) == 0 {
// 		http.Error(w, "Identifier is required", http.StatusBadRequest)
// 		return
// 	}

// 	if len(req.Password) == 0{
// 		http.Error(w, "Password can't be empty", http.StatusBadRequest)
// 		return
// 	}

// 	var (
// 		passwordHash  string
// 		username      string
// 		firstName     string
// 		phoneVerified bool
// 	)

// 	isPhoneLogin := false
// 	var query string

// 	if req.Identifier[0] == '+' {
// 		isPhoneLogin = true
// 		query = `SELECT hashed_password, username, first_name, phone_number_verified
//                  FROM users WHERE phone_number = $1`
// 	} else {
// 		query = `SELECT hashed_password, username, first_name, phone_number_verified
//                  FROM users WHERE username = $1`
// 	}

// 	err := database.DB.QueryRow(r.Context(), query, req.Identifier).Scan(&passwordHash, &username, &firstName, &phoneVerified)
// 	if err != nil {
// 		http.Error(w, "User not found", http.StatusUnauthorized)
// 		return
// 	}

// 	if isPhoneLogin && !phoneVerified {
// 		http.Error(w, "Phone number not verified. Please login with Username.", http.StatusForbidden)
// 		return
// 	}

// 	match, err := password.ComparePasswordAndHash(req.Password, passwordHash)
// 	if err != nil {
// 		http.Error(w, "Error verifying credentials", http.StatusInternalServerError)
// 		return
// 	}

// 	if !match {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 		return
// 	}

// 	response := map[string]string{
// 		"username":   username,
// 		"first_name": firstName,
// 		"status":     "success",
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }
