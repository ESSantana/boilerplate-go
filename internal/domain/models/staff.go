package models

import (
	"errors"
	"slices"
	"strings"
	"time"
)

type Staff struct {
	ID                  string     `db:"id" json:"id"`
	Name                string     `db:"name" json:"name"`
	Nickname            *string    `db:"nickname" json:"nickname,omitempty"`
	Email               string     `db:"email" json:"email"`
	PasswordHash        *string    `db:"password_hash" json:"password_hash,omitempty"`
	StaffRole           string     `db:"staff_role" json:"staff_role"`
	BirthDate           time.Time  `db:"birth_date" json:"birth_date"`
	PhoneNumber         *string    `db:"phone_number" json:"phone_number,omitempty"`
	Document            string     `db:"document" json:"document"`
	DocumentType        string     `db:"document_type" json:"document_type"`
	Address             string     `db:"address" json:"address"`
	AddressNumber       string     `db:"address_number" json:"address_number"`
	AddressComplement   *string    `db:"address_complement" json:"address_complement,omitempty"`
	AddressNeighborhood string     `db:"address_neighborhood" json:"address_neighborhood"`
	AddressCity         string     `db:"address_city" json:"address_city"`
	AddressState        string     `db:"address_state" json:"address_state"`
	AddressZipCode      string     `db:"address_zip_code" json:"address_zip_code"`
	ProviderOrigin      string     `db:"provider_origin" json:"provider_origin"`
	ExternalID          *string    `db:"external_id" json:"external_id,omitempty"`
	ProfileImageURL     *string    `db:"profile_image_url" json:"profile_image_url,omitempty"`
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt           *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

func (s *Staff) Validate() error {

	if strings.TrimSpace(s.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(s.Email) == "" {
		return errors.New("email is required")
	}
	if s.BirthDate.IsZero() {
		return errors.New("birth_date is required")
	}
	if strings.TrimSpace(s.Document) == "" {
		return errors.New("document is required")
	}
	if !isValidDocumentType(s.DocumentType) {
		return errors.New("document_type is invalid")
	}
	if strings.TrimSpace(s.Address) == "" {
		return errors.New("address is required")
	}
	if strings.TrimSpace(s.AddressNumber) == "" {
		return errors.New("address_number is required")
	}
	if strings.TrimSpace(s.AddressNeighborhood) == "" {
		return errors.New("address_neighborhood is required")
	}
	if strings.TrimSpace(s.AddressCity) == "" {
		return errors.New("address_city is required")
	}
	if strings.TrimSpace(s.AddressState) == "" {
		return errors.New("address_state is required")
	}
	if strings.TrimSpace(s.AddressZipCode) == "" {
		return errors.New("address_zip_code is required")
	}
	if strings.TrimSpace(s.ProviderOrigin) == "" {
		return errors.New("provider_origin is required")
	}

	if sr := strings.TrimSpace(s.StaffRole); sr != "" {
		if !isValidStaffRole(sr) {
			return errors.New("staff_role is is invalid")
		}
	} else {
		return errors.New("staff_role is is required")
	}

	return nil
}

func isValidStaffRole(role string) bool {
	validRoles := []string{"admin", "manager", "service_provider"}
	return slices.Contains(validRoles, role)
}

func isValidDocumentType(docType string) bool {
	validTypes := []string{"cpf", "cnpj"}
	return slices.Contains(validTypes, docType)
}
