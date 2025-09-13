package models

import (
	"database/sql"
)

type Users struct {
	ID                  string         `json:"id,omitempty" db:"id,omitempty"`
	FirstName           string         `json:"first_name"  db:"first_name"`
	LastName            string         `json:"last_name" db:"last_name"`
	Email               string         `json:"email" db:"email"`
	UserName            string         `json:"username" db:"username"`
	Phone               string         `json:"phone" db:"phone"`
	BirthDate           string         `json:"birth_date" db:"birth_date"`
	Password            string         `json:"password" db:"password"`
	PasswordChangedAt   sql.NullString `json:"password_changed_at" db:"password_changed_at"`
	UserCreatedAt       sql.NullString `json:"user_created_at,omitempty" db:"user_created_at,omitempty"`
	PasswordResetCode   sql.NullString `json:"password_reset_token,omitempty" db:"password_reset_token,omitempty"`
	PasswordCodeExpires sql.NullString `json:"password_token_expires,omitempty" db:"password_token_expires,omitempty"`
	InactiveStatus      bool           `json:"inactive_status,omitempty" db:"inactive_status,omitempty"`
	Role                string         `json:"role,omitempty" db:"role,omitempty"`
}

type Address struct {
	ID     string `json:"id,omitempty" db:"id,omiempty"`
	Street string `json:"street" db:"street"`
	City   string `json:"city" db:"city"`
	UserID int    `json:"user_id" db:"user_id"`
	User   *Users `json:"-"`
}

type UpdatePasswordModel struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UpdatePasswordRespone struct {
	Token          string `json:"token"`
	PasswordUpdate string `json:"password_update"`
}
