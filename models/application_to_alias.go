package models

import (
	"encoding/json"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// ApplicationToAlias is used by pop to map your application_to_aliases database table to your go code.
type ApplicationToAlias struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ApplicationID uuid.UUID `json:"application_id" db:"application_id"`
	AliasID       uuid.UUID `json:"alias_id" db:"alias_id"`
	Alias         Alias     `belongs_to:"alias"`
}

// String is not required by pop and may be deleted
func (a ApplicationToAlias) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// ApplicationToAliases is not required by pop and may be deleted
type ApplicationToAliases []ApplicationToAlias

// String is not required by pop and may be deleted
func (a ApplicationToAliases) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *ApplicationToAlias) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *ApplicationToAlias) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *ApplicationToAlias) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// SelectLabel is used to display the alias name in the select box
func (a ApplicationToAlias) SelectLabel() string {
	return a.Alias.Name
}

// SelectValue is used to display the alias id in the select box
func (a ApplicationToAlias) SelectValue() interface{} {
	return a.AliasID
}
