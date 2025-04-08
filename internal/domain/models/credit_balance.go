package models

import (
	"errors"
	"strings"
	"time"
)

type CreditBalance struct {
	ID           string     `db:"id" json:"id"`
	CustomerID   string     `db:"customer_id" json:"customer_id"`
	CreditAmount int        `db:"credit_amount" json:"credit_amount"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdateAt     time.Time  `db:"update_at" json:"update_at"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

func (cb *CreditBalance) Validate() error {
	if strings.TrimSpace(cb.CustomerID) == "" {
		return errors.New("customer_id is required")
	}
	if cb.CreditAmount < 0 {
		return errors.New("credit_amount can't be negative")
	}

	return nil
}
