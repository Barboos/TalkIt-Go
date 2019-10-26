package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"time"
)

type Relationship struct {
	ID         uuid.UUID 		`json:"id" db:"id"`
	CreatedAt  time.Time 		`json:"created_at" db:"created_at"`
	UpdatedAt  time.Time 		`json:"updated_at" db:"updated_at"`
	FollowerID uuid.UUID   	`json:"follower_id" db:"follower_id"`
	FollowedID uuid.UUID   	`json:"followed_id" db:"followed_id"`
}

// String is not required by pop and may be deleted
func (r Relationship) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Relationships is not required by pop and may be deleted
type Relationships []Relationship

// String is not required by pop and may be deleted
func (r Relationships) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Relationship) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Relationship) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Relationship) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
