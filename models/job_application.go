package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// internal struct to hold the job parameters.
type application struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ProjectID     uuid.UUID `json:"project_id" db:"project_id"`
	ClientID      uuid.UUID `json:"client_id" db:"client_id"`
	RuntimeID     uuid.UUID `json:"runtime_id" db:"runtime_id"`
	DatabaseID    uuid.UUID `json:"database_id" db:"database_id"`
	EnvironmentID uuid.UUID `json:"environment_id" db:"environment_id"`
	Repository    string    `json:"repository" db:"repository"`
	Branch        string    `json:"branch" db:"branch"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// JobApplication is used by pop to map your job_applications database table to your go code.
type JobApplication struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	Type        string       `json:"type" db:"type"`
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

// NewParamApplication returns the application struct from the new_params field.
func (j *JobApplication) NewParamApplication() (*Application, error) {
	targetApp := &application{}
	err := json.Unmarshal([]byte(j.NewParams.String), &targetApp)
	if err != nil {
		return nil, errors.New("Failed to unmarshal the job parameters. " + err.Error() + " " + j.NewParams.String)
	}
	return j.targetAppApplication(targetApp)
}

// OrigParamApplication returns the application struct from the orig_params field.
func (j *JobApplication) OrigParamApplication() (*Application, error) {
	targetApp := &application{}
	err := json.Unmarshal([]byte(j.OrigParams.String), &targetApp)
	if err != nil {
		return nil, errors.New("Failed to unmarshal the job parameters. " + err.Error() + " " + j.OrigParams.String)
	}
	return j.targetAppApplication(targetApp)
}

// targetAppApplication returns the application struct frome the given application.
func (j *JobApplication) targetAppApplication(targetApp *application) (*Application, error) {
	// execute the create project job.
	project := Project{}
	DB.Where("id = ?", targetApp.ProjectID).First(&project)
	client := Client{}
	DB.Where("id = ?", targetApp.ClientID).First(&client)
	dbtype := Dbtype{}
	DB.Where("id = ?", targetApp.DatabaseID).First(&dbtype)
	runtime := Runtime{}
	DB.Where("id = ?", targetApp.RuntimeID).First(&runtime)
	environment := Environment{}
	DB.Where("id = ?", targetApp.EnvironmentID).First(&environment)

	hosts := Hosts{}
	DB.Where("environment_id = ?", environment.ID).All(&hosts)
	// create an application struct to pass to the server.
	app := &Application{
		Project:     &project,
		Client:      &client,
		Database:    &dbtype,
		Runtime:     &runtime,
		Environment: &environment,
		Repository:  targetApp.Repository,
		Branch:      targetApp.Branch,
	}
	return app, nil
}
