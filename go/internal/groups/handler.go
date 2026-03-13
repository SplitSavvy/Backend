package groups

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

type AddToGroupRequest struct {
	UserId      string `json:"user_id"`
	RequesterId string `json:"requester_id"`
}

type GroupMember struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username *string   `json:"username"`
	IsGhost  bool      `json:"is_ghost"`
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

func (h *Handler) AddToGroup(w http.ResponseWriter, r *http.Request) {
	var req AddToGroupRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		sendJSONError(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	uuidToParse := chi.URLParam(r, "id")

	if uuidToParse == "" {
		sendJSONError(w, "group_id is required", http.StatusBadRequest)
		return
	}

	groupId, err := uuid.Parse(uuidToParse)

	if err != nil {
		sendJSONError(w, "UUID parsing failed", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	var someDummyVariable uuid.UUID

	query := `
			SELECT user_id FROM group_members WHERE user_id = $1 AND group_id = $2;
	`
	err = h.DB.QueryRow(ctx, query, req.RequesterId, groupId).Scan(&someDummyVariable)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			sendJSONError(w, "The group member not authorized", http.StatusForbidden)
			return
		}
		sendJSONError(w, "Database Error", http.StatusInternalServerError)
		return
	}

	addQuery := `INSERT INTO group_members(group_id, user_id, invited_by) VALUES($1, $2, $3)`

	_, err = h.DB.Exec(ctx, addQuery, groupId, req.UserId, req.RequesterId)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			sendJSONError(w, "User is already in the group", http.StatusConflict)
			return
		}
		sendJSONError(w, "Database Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"message": "Added succesfully",
		"status":  http.StatusCreated,
	}
	json.NewEncoder(w).Encode(response)

}

func (h *Handler) GetMembers(w http.ResponseWriter, r *http.Request) {
	uuidToParse := chi.URLParam(r, "id")

	if uuidToParse == "" {
		sendJSONError(w, "group_id is required", http.StatusBadRequest)
		return
	}

	groupId, err := uuid.Parse(uuidToParse)

	if err != nil {
		sendJSONError(w, "UUID parsing failed", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	results := []GroupMember{}

	query := `
        SELECT 
            u.id, 
            u.name, 
            u.username, 
            u.is_ghost
        FROM users u
        INNER JOIN group_members gm 
            ON u.id = gm.user_id
        WHERE gm.group_id = $1 
            AND gm.removed_on IS NULL;
    `

	rows, err := h.DB.Query(ctx, query, groupId)
	if err != nil {
		sendJSONError(w, "Database Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var member GroupMember

		err := rows.Scan(
			&member.ID,
			&member.Name,
			&member.Username,
			&member.IsGhost,
		)

		if err != nil {
			sendJSONError(w, "Error parsing database results", http.StatusInternalServerError)
			return
		}

		results = append(results, member)
	}

	if err = rows.Err(); err != nil {
		sendJSONError(w, "Database iteration error", http.StatusInternalServerError)
		return
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
