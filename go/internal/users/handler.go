package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"unicode/utf8"

	"splitsavvy/internal/password"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

type UserSearchResult struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username *string   `json:"username"`
	InGroup  bool      `json:"in_group"`
}

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var pgErr *pgconn.PgError
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJSONError(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	tx, err := h.DB.Begin(ctx)

	if err != nil {
		sendJSONError(w, "Failed to Connect to DB", http.StatusInternalServerError)
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
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				switch pgErr.ConstraintName {
				case "users_username_key":
					sendJSONError(w, "Username already exist", http.StatusConflict)
				case "users_phone_number_key":
					sendJSONError(w, "Phone number already exist", http.StatusConflict)
				case "unique_ghost_contact":
					sendJSONError(w, "Ghost user already exisit", http.StatusConflict)
				default:
					sendJSONError(w, "Conflict", http.StatusConflict)
				}
				return
			}
		}
		sendJSONError(w, "Could not create user profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !req.IsGhost {
		if req.Email == nil || req.Password == nil {
			sendJSONError(w, "Email and Password is mandatory", http.StatusBadRequest)
			return
		}
		hashedPassword, err := password.CreateHash(*req.Password, password.DefaultParams)
		if err != nil {
			sendJSONError(w, "Failed to secure password", http.StatusInternalServerError)
			return
		}
		credQuery := `INSERT INTO user_credentials (user_id, email, password_hash)
            VALUES ($1, $2, $3)
		`
		_, err = tx.Exec(ctx, credQuery, newUserID, *req.Email, hashedPassword)

		if err != nil {
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" && pgErr.ConstraintName == "user_credentials_email_key" {
					sendJSONError(w, "Email ID already exist", http.StatusConflict)
					return
				}
			}
			sendJSONError(w, "Couldn't Save Credentials"+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		sendJSONError(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	response := map[string]interface{}{
		"message": "User created successfully",
		"id":      newUserID,
		"status":  http.StatusCreated,
	}
	json.NewEncoder(w).Encode(response)

}

func (h *Handler) SearchUser(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	groupId := r.URL.Query().Get("group_id")
	requesterId := r.URL.Query().Get("requester_id")
	if utf8.RuneCountInString(searchTerm) < 3 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
		return
	}

	if groupId == "" {
		sendJSONError(w, "group_id is required", http.StatusBadRequest)
		return
	}

	if requesterId == "" {
		sendJSONError(w, "group_id is required", http.StatusBadRequest)
		return
	}

	dbSearchTerm := searchTerm + "%"
	ctx := context.Background()

	query := `
            SELECT 
                u.id AS user_id,
                u.name,
                u.username,
                CASE 
                    WHEN gm.user_id IS NOT NULL THEN TRUE
                    ELSE FALSE
                END AS in_group
            FROM users u
            LEFT JOIN group_members gm 
                ON u.id = gm.user_id 
                AND gm.group_id = $2
                AND gm.removed_on IS NULL
            WHERE 
                (u.name ILIKE $1 OR u.username ILIKE $1)
                AND (u.is_ghost = FALSE OR u.created_by = $3)
            LIMIT 5;
            `

	rows, err := h.DB.Query(ctx, query, dbSearchTerm, groupId, requesterId)
	if err != nil {
		sendJSONError(w, "Database error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var results []UserSearchResult

	for rows.Next() {
		var user UserSearchResult
		err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.InGroup)
		if err != nil {
			sendJSONError(w, "Error reading user data", http.StatusInternalServerError)
			return
		}
		results = append(results, user)
	}

	if err = rows.Err(); err != nil {
		sendJSONError(w, "Error iterating users", http.StatusInternalServerError)
		return
	}

	if results == nil {
		results = []UserSearchResult{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)

}

func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  message,
		"status": statusCode,
	})
}
