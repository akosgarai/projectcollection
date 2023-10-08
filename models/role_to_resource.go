package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// RoleToResource is used by pop to map your role_to_resources database table to your go code.
type RoleToResource struct {
	ID         uuid.UUID `json:"id" db:"id"`
	RoleID     uuid.UUID `json:"role_id" db:"role_id"`
	ResourceID uuid.UUID `json:"resource_id" db:"resource_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (r RoleToResource) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// RoleToResources is not required by pop and may be deleted
type RoleToResources []RoleToResource

// String is not required by pop and may be deleted
func (r RoleToResources) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *RoleToResource) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *RoleToResource) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *RoleToResource) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
