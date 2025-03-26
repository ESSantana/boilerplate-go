package models

import "time"

type User struct {
	//SHA1 hash of the user email
	ID             string `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	Email          string `db:"email" json:"email"`
	ProviderOrigin string `db:"provider_origin" json:"provider_origin"`
	// Unique identifier from the external provider
	ExternalID      string     `db:"external_id" json:"external_id"`
	ProfileImageURL string     `db:"profile_image_url" json:"profile_image_url"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at"`
}
