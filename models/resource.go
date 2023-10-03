package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Resource is used by pop to map your resources database table to your go code.
type Resource struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (r Resource) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Resources is not required by pop and may be deleted
type Resources []Resource

// String is not required by pop and may be deleted
func (r Resources) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Resource) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: r.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Resource) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Resource) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
