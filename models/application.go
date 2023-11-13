package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Application is used by pop to map your applications database table to your go code.
type Application struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	ProjectID     uuid.UUID    `json:"project_id" db:"project_id"`
	Project       *Project     `json:"project" belongs_to:"project"`
	ClientID      uuid.UUID    `json:"client_id" db:"client_id"`
	Client        *Client      `json:"client" belongs_to:"client"`
	RuntimeID     uuid.UUID    `json:"runtime_id" db:"runtime_id"`
	Runtime       *Runtime     `json:"runtime" belongs_to:"runtime"`
	DatabaseID    uuid.UUID    `json:"database_id" db:"database_id"`
	Database      *Dbtype      `json:"database" belongs_to:"dbtype"`
	EnvironmentID uuid.UUID    `json:"environment_id" db:"environment_id"`
	Environment   *Environment `json:"environment" belongs_to:"environment"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (a Application) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Applications is not required by pop and may be deleted
type Applications []Application

// String is not required by pop and may be deleted
func (a Applications) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Application) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: a.ProjectID, Name: "ProjectID"},
		&validators.UUIDIsPresent{Field: a.ClientID, Name: "ClientID"},
		&validators.UUIDIsPresent{Field: a.RuntimeID, Name: "RuntimeID"},
		&validators.UUIDIsPresent{Field: a.DatabaseID, Name: "DatabaseID"},
		&validators.UUIDIsPresent{Field: a.EnvironmentID, Name: "EnvironmentID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Application) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Application) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
