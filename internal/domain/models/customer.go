package models

import (
	"errors"
	"slices"
	"strings"
	"time"
)

type Customer struct {
	ID                            string     `db:"id" json:"id"`
	Name                          string     `db:"name" json:"name"`
	Nickname                      *string    `db:"nickname" json:"nickname,omitempty"`
	Email                         string     `db:"email" json:"email"`
	PasswordHash                  *string    `db:"password_hash" json:"password_hash,omitempty"`
	BirthDate                     time.Time  `db:"birth_date" json:"birth_date"`
	PhoneNumber                   *string    `db:"phone_number" json:"phone_number,omitempty"`
	CPF                           string     `db:"cpf" json:"cpf"`
	Gender                        string     `db:"gender" json:"gender"`
	Address                       string     `db:"address" json:"address"`
	AddressNumber                 string     `db:"address_number" json:"address_number"`
	AddressComplement             *string    `db:"address_complement" json:"address_complement,omitempty"`
	AddressNeighborhood           string     `db:"address_neighborhood" json:"address_neighborhood"`
	AddressCity                   string     `db:"address_city" json:"address_city"`
	AddressState                  string     `db:"address_state" json:"address_state"`
	AddressZipCode                string     `db:"address_zip_code" json:"address_zip_code"`
	ProviderOrigin                string     `db:"provider_origin" json:"provider_origin"`
	ExternalID                    *string    `db:"external_id" json:"external_id,omitempty"`
	ProfileImageURL               *string    `db:"profile_image_url" json:"profile_image_url,omitempty"`
	Interests                     *string    `db:"interests" json:"interests,omitempty"`
	HowHeardAboutUs               *string    `db:"how_heard_about_us" json:"how_heard_about_us,omitempty"`
	PreferredCommunicationChannel string     `db:"preferred_comunication_channel" json:"preferred_communication_channel"`
	CreatedAt                     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt                     time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt                     *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

func (c *Customer) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(c.Email) == "" {
		return errors.New("email is required")
	}
	if c.PasswordHash != nil && strings.TrimSpace(*c.PasswordHash) == "" {
		return errors.New("password_hash is required")
	}
	if c.BirthDate.IsZero() {
		return errors.New("birth_date is required")
	}
	if strings.TrimSpace(c.CPF) == "" {
		return errors.New("cpf is required")
	}
	if strings.TrimSpace(c.Gender) == "" {
		return errors.New("gender is required")
	}
	if strings.TrimSpace(c.Address) == "" {
		return errors.New("address is required")
	}
	if strings.TrimSpace(c.AddressNumber) == "" {
		return errors.New("address_number is required")
	}
	if strings.TrimSpace(c.AddressNeighborhood) == "" {
		return errors.New("address_neighborhood is required")
	}
	if strings.TrimSpace(c.AddressCity) == "" {
		return errors.New("address_city is required")
	}
	if strings.TrimSpace(c.AddressState) == "" {
		return errors.New("address_state is required")
	}
	if strings.TrimSpace(c.AddressZipCode) == "" {
		return errors.New("address_zip_code is required")
	}
	if strings.TrimSpace(c.ProviderOrigin) == "" {
		return errors.New("provider_origin is required")
	}

	if prefCom := strings.TrimSpace(c.PreferredCommunicationChannel); prefCom != "" {
		if !isValidCommunicationChannel(prefCom) {
			return errors.New("preferred_communication_channel is invalid")
		}
	} else {
		return errors.New("preferred_communication_channel is required")
	}

	return nil
}

func isValidCommunicationChannel(channel string) bool {
	validChannels := []string{"email", "sms", "whatsapp", "no-preference"}
	return slices.Contains(validChannels, channel)
}
