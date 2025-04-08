package models

import (
	"errors"
	"strings"
	"time"
)

type CustomerEvent struct {
	ID         string     `db:"id" json:"id"`
	CustomerID string     `db:"customer_id" json:"customer_id"`
	EventType  string     `db:"event_type" json:"event_type"`
	EventData  *string    `db:"event_data" json:"event_data,omitempty"`
	Latitude   float64    `db:"latitude" json:"latitude"`
	Longitude  float64    `db:"longitude" json:"longitude"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdateAt   time.Time  `db:"update_at" json:"update_at"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

func (ce *CustomerEvent) Validate() error {
	if strings.TrimSpace(ce.CustomerID) == "" {
		return errors.New("customer_id is required")
	}
	if strings.TrimSpace(ce.EventType) == "" {
		return errors.New("event_type is required")
	}

	if ce.Latitude < -90 || ce.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if ce.Longitude < -180 || ce.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}

	return nil
}