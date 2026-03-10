package groups

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	DB *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{DB: db}
}

type CreateGroupRequest struct {
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"id"`
}

func (h *Handler) HandleGroupRequest(w http.ResponseWriter, r *http.Request) {
	var req CreateGroupRequest
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

	var newGroupID uuid.UUID

	groupQuery := `INSERT INTO groups(name, created_by)
	VALUES($1, $2)
	returning id`

	err = tx.QueryRow(ctx, groupQuery,
		req.Name,
		req.CreatedBy,
	).Scan(&newGroupID)

	if err != nil {
		sendJSONError(w, "Error while creating group", http.StatusInternalServerError)
		return
	}

	groupMemberQuery := `INSERT INTO group_members(group_id, user_id)
	VALUES($1, $2)`

	_, err = tx.Exec(ctx, groupMemberQuery, newGroupID, req.CreatedBy)

	if err != nil {
		sendJSONError(w, "Failed to add members", http.StatusInternalServerError)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		sendJSONError(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	response := map[string]interface{}{
		"message": "Group created successfully",
		"id":      newGroupID,
		"status":  http.StatusCreated,
	}
	json.NewEncoder(w).Encode(response)
}

func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  message,
		"status": statusCode,
	})
}
