package models

import (
	"errors"
	"strings"
	"time"
)

type Product struct {
	ID              string     `db:"id" json:"id"`
	ProviderID      string     `db:"provider_id" json:"provider_id"`
	CategoryID      string     `db:"category_id" json:"category_id"`
	Name            string     `db:"name" json:"name"`
	Description     *string    `db:"description" json:"description,omitempty"`
	CreditCost      int        `db:"credit_cost" json:"credit_cost"`
	AverageDuration int        `db:"average_duration" json:"average_duration"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdateAt        time.Time  `db:"update_at" json:"update_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

func (p *Product) Validate() error {
	if strings.TrimSpace(p.ProviderID) == "" {
		return errors.New("provider_id is required")
	}
	if strings.TrimSpace(p.CategoryID) == "" {
		return errors.New("category_id is required")
	}
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("name is required")
	}
	if p.CreditCost <= 0 {
		return errors.New("credit_cost must be greater than 0")
	}
	if p.AverageDuration <= 0 {
		return errors.New("average_duration must be greater than 0")
	}

	return nil
}
