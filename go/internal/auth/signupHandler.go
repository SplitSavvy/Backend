package auth

import (
	"encoding/json"
	"net/http"
	"regexp"
	"splitsavvy/internal/database"
	"splitsavvy/internal/password"
	"strings"
)

type SignupRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

var (
	userRegex  = regexp.MustCompile(`^[a-z0-9_.]{4,16}$`)
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Format", http.StatusBadRequest)
		return
	}

	if len(req.Username) == 0 || len(req.FirstName) == 0 || len(req.LastName) == 0 || len(req.Email) == 0 || len(req.Password) == 0 {
		http.Error(w, "Fields Can Not Be Empty", http.StatusBadRequest)
		return
	}

	var exists bool
	req.Username = strings.ToLower(req.Username)
	req.Email = strings.ToLower(req.Email)

	if !userRegex.MatchString(req.Username) {
		http.Error(w, "Invalid Username", http.StatusBadRequest)
		return
	}

	if !emailRegex.MatchString(req.Email) {
		http.Error(w, "Invalid Email", http.StatusBadRequest)
		return
	}

	err := database.DB.QueryRow(r.Context(), "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&exists)

	if err != nil {
		http.Error(w, "Internal DB Error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Username Already Taken", http.StatusConflict)
		return
	}

	err = database.DB.QueryRow(r.Context(), "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)

	if err != nil {
		http.Error(w, "Internal DB Error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email Already Registered", http.StatusConflict)
		return
	}

	hashedPwd, err := password.CreateHash(req.Password, password.DefaultParams)

	if err != nil {
		http.Error(w, "Error While Hashing Password", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO users (username, first_name, last_name, email, hashed_password) VALUES($1, $2, $3, $4, $5)`

	_, err = database.DB.Exec(r.Context(), query, req.Username, req.FirstName, req.LastName, req.Email, hashedPwd)

	if err != nil {
		http.Error(w, "Failed To Create User", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "User created successfully",
	})
}
