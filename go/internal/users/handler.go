package users

import (
	"context"
	"encoding/json"
	"net/http"

	"splitsavvy/internal/password"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	DB *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{DB: db}
}

type CreateUserRequest struct {
	Name        string     `json:"name"`
	Username    *string    `json:"username"`
	PhoneNumber *string    `json:"phone_number"`
	Email       *string    `json:"email"`
	Password    *string    `json:"password"`
	IsGhost     bool       `json:"is_ghost"`
	CreatedBy   *uuid.UUID `json:"created_by"`
}

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	tx, err := h.DB.Begin(ctx)

	if err != nil {
		http.Error(w, "Failed to Connect to DB", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(ctx)

	var newUserID uuid.UUID

	userQuery := `INSERT INTO users(name, username, phone_number, is_ghost, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`
	err = tx.QueryRow(ctx, userQuery,
		req.Name,
		req.Username,
		req.PhoneNumber,
		req.IsGhost,
		req.CreatedBy,
	).Scan(&newUserID)

	if err != nil {
		http.Error(w, "Could not create user profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !req.IsGhost {
		if req.Email == nil || req.Password == nil {
			http.Error(w, "Email and Password is mandatory", http.StatusBadRequest)
			return
		}
		hashedPassword, err := password.CreateHash(*req.Password, password.DefaultParams)
		if err != nil {
			http.Error(w, "Failed to secure password", http.StatusInternalServerError)
			return
		}
		credQuery := `INSERT INTO user_credentials (user_id, email, password_hash)
            VALUES ($1, $2, $3)
		`
		_, err = tx.Exec(ctx, credQuery, newUserID, *req.Email, hashedPassword)

		if err != nil {
			http.Error(w, "Couldn't Save Credentials"+err.Error(), http.StatusInternalServerError)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	response := map[string]interface{}{
		"message": "User created successfully",
		"id":      newUserID,
	}
	json.NewEncoder(w).Encode(response)

}
