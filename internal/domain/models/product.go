package models

import (
	"time"

	"github.com/application-ellas/ellas-backend/internal/domain/constants"
	"github.com/application-ellas/ellas-backend/internal/domain/errors"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID                string     `db:"id" json:"id"`
	ServiceProviderID string     `db:"service_provider_id" json:"service_provider_id"`
	Name              string     `db:"name" json:"name"`
	CategoryID        string     `db:"category_id" json:"category_id"`
	Description       string     `db:"description" json:"description"`
	Price             string     `db:"price" json:"price"`
	Duration          int        `db:"duration" json:"duration"`
	Latitude          float64    `db:"latitude" json:"latitude"`
	Longitude         float64    `db:"longitude" json:"longitude"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at" json:"deleted_at"`
}

func (product *Product) Validate(validationTyep constants.ModelValidationType) error {
	if product.ServiceProviderID == "" {
		return errors.NewValidationError("service provider id is required")
	}
	if product.Name == "" {
		return errors.NewValidationError("product name is required")
	}
	if product.CategoryID == "" {
		return errors.NewValidationError("product category id is required")
	}
	if product.Price == "" {
		return errors.NewValidationError("product price is required")
	}
	if product.Duration <= 0 {
		return errors.NewValidationError("product duration must be greater than 0")
	}
	if product.Latitude < -90 || product.Latitude > 90 {
		return errors.NewValidationError("product latitude must be between -90 and 90")
	}
	if product.Longitude < -180 || product.Longitude > 180 {
		return errors.NewValidationError("product longitude must be between -180 and 180")
	}

	if validationTyep == constants.ValidationTypeUpdate {
		if product.ID == "" {
			return errors.NewValidationError("product id is required")
		}
	}

	return nil
}

func (product *Product) GetPrice() decimal.Decimal {
	price, _ := decimal.NewFromString(product.Price)
	return price
}
