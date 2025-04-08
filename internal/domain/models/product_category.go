package models

import (
	"errors"
	"strings"
	"time"
)

type ProductCategory struct {
	ID          string     `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

func (pc *ProductCategory) Validate() error {
	if strings.TrimSpace(pc.Name) == "" {
		return errors.New("name é obrigatório")
	}

	return nil
}
