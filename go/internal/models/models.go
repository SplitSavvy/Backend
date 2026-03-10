package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Username    *string    `json:"username,omitempty" db:"username"`
	PhoneNumber *string    `json:"phone_number,omitempty" db:"phone_number"`
	IsGhost     bool       `json:"is_ghost" db:"is_ghost"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type UserCredential struct {
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Group struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	CreatedBy uuid.UUID  `json:"created_by" db:"created_by"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type GroupMember struct {
	GroupID   uuid.UUID  `json:"group_id" db:"group_id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Role      string     `json:"role" db:"role"`
	AddedOn   time.Time  `json:"added_on" db:"added_on"`
	RemovedOn *time.Time `json:"removed_on,omitempty" db:"removed_on"`
}
