package models

import "time"

type Role = string

const (
	RoleAdmin           Role = "admin"
	RoleCustomer        Role = "customer"
	RoleServiceProvider Role = "service_provider"
)

type UserRole struct {
	ID        string     `db:"id" json:"id"`
	UserID    string     `db:"user_id" json:"user_id"`
	Role      Role       `db:"role" json:"role"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}
