package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// JobApplication is used by pop to map your job_applications database table to your go code.
type JobApplication struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	NewParams   nulls.String `json:"new_params" db:"new_params"`
	OrigParams  nulls.String `json:"orig_params" db:"orig_params"`
	ProcessedAt nulls.Time   `json:"processed_at" db:"processed_at"`
	Response    nulls.String `json:"response" db:"response"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (j JobApplication) String() string {
	jj, _ := json.Marshal(j)
	return string(jj)
}

// JobApplications is not required by pop and may be deleted
type JobApplications []JobApplication

// String is not required by pop and may be deleted
func (j JobApplications) String() string {
	jj, _ := json.Marshal(j)
	return string(jj)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (j *JobApplication) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (j *JobApplication) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (j *JobApplication) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
