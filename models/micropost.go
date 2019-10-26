package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"time"
)

type Micropost struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	User      User      `belongs_to:"user" db:"-"`
	UserID 	uuid.UUID `json:"user_id" db:"user_id"`
}

// String is not required by pop and may be deleted
func (m Micropost) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Microposts is not required by pop and may be deleted
type Microposts []Micropost

// String is not required by pop and may be deleted
func (m Microposts) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *Micropost) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.Title, Name: "Title"},
		&validators.StringIsPresent{Field: m.Content, Name: "Content"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *Micropost) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *Micropost) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
