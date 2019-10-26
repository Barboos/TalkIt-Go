package models

import (
	"strings"
	"time"

	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                   	uuid.UUID  						`json:"id" db:"id"`
	CreatedAt            	time.Time  						`json:"created_at" db:"created_at"`
	UpdatedAt            	time.Time  						`json:"updated_at" db:"updated_at"`
	Email                	string     						`json:"email" db:"email"`
	Name                 	string     						`json:"name" db:"name"`
	PasswordHash         	string     						`json:"password_hash" db:"password_hash"`
	Microposts           	Microposts 						`has_many:"microposts" order_by:"created_at desc" db:"-"`
	Password             	string     						`json:"-" db:"-"`
	PasswordConfirmation 	string     						`json:"-" db:"-"`
}

// Create wraps up the pattern of encrypting the password and
// running validations. Useful when writing tests.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	u.PasswordHash = string(ph)
	return tx.ValidateAndCreate(u)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.EmailIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
		// check to see if the email address is already taken:
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email)
				if u.ID != uuid.Nil {
					q = q.Where("id != ?", u.ID)
				}
				b, err = q.Exists(u)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
	), err
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	u.Email = strings.ToLower(u.Email)
	if len(u.Password) > 0 {
		ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}
		u.PasswordHash = string(ph)
	}
	return validate.Validate(
		// check to see if the email address is already taken:
		&validators.StringsMatch{Name: "Password", Field: u.Password,
			Field2:  u.PasswordConfirmation,
			Message: "Password does not match confirmation"},
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email)
				if u.ID != uuid.Nil {
					q = q.Where("id != ?", u.ID)
				}
				b, err = q.Exists(u)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), err
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (m Users) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}
