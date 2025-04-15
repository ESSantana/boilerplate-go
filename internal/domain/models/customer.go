package models

import (
	"time"
)

type Customer struct {
	ID                            string     `db:"id" json:"id"`
	Name                          string     `db:"name" json:"name" validate:"required"`
	Nickname                      *string    `db:"nickname" json:"nickname,omitempty"`
	Email                         string     `db:"email" json:"email" validate:"required,email"`
	PasswordHash                  *string    `db:"password_hash" json:"password_hash,omitempty" validate:"required_if_not_match_format=ID uuid"`
	BirthDate                     time.Time  `db:"birth_date" json:"birth_date" validate:"required"`
	PhoneNumber                   *string    `db:"phone_number" json:"phone_number,omitempty"`
	CPF                           string     `db:"cpf" json:"cpf" validate:"required"`
	Gender                        string     `db:"gender" json:"gender" validate:"required"`
	Address                       string     `db:"address" json:"address" validate:"required"`
	AddressNumber                 string     `db:"address_number" json:"address_number" validate:"required"`
	AddressComplement             *string    `db:"address_complement" json:"address_complement,omitempty"`
	AddressNeighborhood           string     `db:"address_neighborhood" json:"address_neighborhood" validate:"required"`
	AddressCity                   string     `db:"address_city" json:"address_city" validate:"required"`
	AddressState                  string     `db:"address_state" json:"address_state" validate:"required"`
	AddressZipCode                string     `db:"address_zip_code" json:"address_zip_code" validate:"required"`
	ProviderOrigin                string     `db:"provider_origin" json:"provider_origin" validate:"required_if_match_format=ID uuid"`
	ExternalID                    *string    `db:"external_id" json:"external_id,omitempty"`
	ProfileImageURL               *string    `db:"profile_image_url" json:"profile_image_url,omitempty"`
	Interests                     *string    `db:"interests" json:"interests,omitempty"`
	HowHeardAboutUs               *string    `db:"how_heard_about_us" json:"how_heard_about_us,omitempty"`
	PreferredCommunicationChannel string     `db:"preferred_comunication_channel" json:"preferred_communication_channel" validate:"required,oneof=email sms whatsapp no-preference"`
	CreatedAt                     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt                     time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt                     *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
