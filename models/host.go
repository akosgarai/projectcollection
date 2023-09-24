package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Host is used by pop to map your hosts database table to your go code.
type Host struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	Name          string       `json:"name" db:"name"`
	IP            string       `json:"ip" db:"ip"`
	EnvironmentID uuid.UUID    `json:"environment_id" db:"environment_id"`
	Environment   *Environment `json:"environment" belongs_to:"environment"`
	SSHUser       string       `json:"ssh_user" db:"ssh_user"`
	SSHPort       int          `json:"ssh_port" db:"ssh_port"`
	SSHKey        string       `json:"ssh_key" db:"ssh_key"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (h Host) String() string {
	jh, _ := json.Marshal(h)
	return string(jh)
}

// Hosts is not required by pop and may be deleted
type Hosts []Host

// String is not required by pop and may be deleted
func (h Hosts) String() string {
	jh, _ := json.Marshal(h)
	return string(jh)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (h *Host) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: h.Name, Name: "Name"},
		&validators.StringIsPresent{Field: h.IP, Name: "IP"},
		&validators.StringIsPresent{Field: h.SSHUser, Name: "SSHUser"},
		&validators.IntIsPresent{Field: h.SSHPort, Name: "SSHPort"},
		&validators.StringIsPresent{Field: h.SSHKey, Name: "SSHKey"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (h *Host) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (h *Host) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
