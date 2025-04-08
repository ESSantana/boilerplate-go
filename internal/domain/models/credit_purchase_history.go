package models

import (
	"errors"
	"slices"
	"strings"
	"time"
)

type CreditPurchaseHistory struct {
	ID                   string     `db:"id" json:"id"`
	CustomerID           string     `db:"customer_id" json:"customer_id"`
	CreditAmount         int        `db:"credit_amount" json:"credit_amount"`
	PaymentValue         string     `db:"payment_value" json:"payment_value"`
	PaymentMethod        string     `db:"payment_method" json:"payment_method"`
	PaymentStatus        string     `db:"payment_status" json:"payment_status"`
	PaymentVendor        string     `db:"payment_vendor" json:"payment_vendor"`
	PaymentTransactionID string     `db:"payment_transaction_id" json:"payment_transaction_id"`
	CreatedAt            time.Time  `db:"created_at" json:"created_at"`
	UpdateAt             time.Time  `db:"update_at" json:"update_at"`
	DeletedAt            *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

func (c *CreditPurchaseHistory) Validate() error {
	if strings.TrimSpace(c.CustomerID) == "" {
		return errors.New("customer_id is required")
	}
	if c.CreditAmount <= 0 {
		return errors.New("credit_amount must be greater than 0")
	}
	if strings.TrimSpace(c.PaymentValue) == "" {
		return errors.New("payment_value is required")
	}
	if strings.TrimSpace(c.PaymentMethod) == "" {
		return errors.New("payment_method is required")
	}
	if strings.TrimSpace(c.PaymentVendor) == "" {
		return errors.New("payment_vendor is required")
	}
	if strings.TrimSpace(c.PaymentTransactionID) == "" {
		return errors.New("payment_transaction_id is required")
	}
	if ps := strings.TrimSpace(c.PaymentStatus); ps != "" {
		if !isValidPaymentStatus(c.PaymentStatus) {
			return errors.New("payment_status is invalid")
		}
	} else {
		return errors.New("payment_status is required")
	}

	return nil
}

func isValidPaymentStatus(status string) bool {
	validStatuses := []string{"pending", "paid", "canceled", "failed", "refunded"}
	return slices.Contains(validStatuses, status)
}
