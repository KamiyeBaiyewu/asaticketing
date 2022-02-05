package model

import (
	"errors"
	"time"

	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserID is the identifier for the student
type UserID string

// NilUserID is an empty UserID
var NilUserID UserID

var (
	userTypes = []string{"admin", "agent", "user"}
)

// User is a structure that represents User Object
type User struct {
	ID           UserID     `json:"id,omitempty" db:"user_id"`
	Name         *string    `json:"name,omitempty"`
	Firstname    *string    `json:"firstname,omitempty" db:"firstname"`
	Lastname     *string    `json:"lastname,omitempty" db:"lastname"`
	Email        *string    `json:"email,omitempty" db:"email"`
	Type         *string    `json:"type,omitempty" db:"user_type"`
	RoleID       RoleID     `json:"role_id,omitempty" db:"role_id"`
	PasswordHash *[]byte    `json:"-" db:"password_hash"`
	IsActive     *bool      `json:"-" db:"is_active"`
	IsSystem     *bool      `json:"-" db:"is_system"`
	CreatedAt    *time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"  db:"deleted_at"`

	// MISC
	Role            *Role     `json:"role,omitempty"` //Primary role
	Password        *string   `json:"password,omitempty"`
	ConfirmPassword *string   `json:"confirm_password,omitempty"`
	Actions         []*Action `json:"actions,omitempty"`
	Roles           []*Role   `json:"roles,omitempty"`
}

// SetPassword accepts password string and creates hash
func (u *User) SetPassword(password string) (err error) {

	hash, err := HashPassword(password)
	if err != nil {
		return
	}
	u.PasswordHash = &hash

	return
}

//Verify all fields before create or update
func (u *User) Verify() error {

	if u.Email == nil || (u.Email != nil && len(*u.Email) == 0) {
		return errors.New("Email is required")
	}
	if u.Firstname == nil || (u.Firstname != nil && len(*u.Firstname) == 0) {
		return errors.New("Firstname is required")
	}
	if u.Lastname == nil || (u.Lastname != nil && len(*u.Lastname) == 0) {
		return errors.New("Lastname is required")
	}
	if u.Type == nil || (u.Type != nil && len(*u.Type) == 0) {
		return errors.New("Type is required")
	} else if !utils.ItemExists(userTypes, *u.Type) {
		println(*u.Type)
		return errors.New("Invalid user type")
	}
	if u.RoleID == NilRoleID {
		return errors.New("Role is required")
	}

	return nil
}

// CheckPassword  verifies the user's password
func (u *User) CheckPassword(password string) error {
	if u.PasswordHash != nil && len(*u.PasswordHash) == 0 {
		return errors.New("password not set")
	}
	return bcrypt.CompareHashAndPassword(*u.PasswordHash, []byte(password))
}

// HashPassword helps turn password in plain text into bcrypt hash
func HashPassword(password string) ([]byte, error) {

	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// UpdateValues is used to update empty values
func (u *User) UpdateValues(nv *User) { //nv means new values
	// Avoid updating the same values
	if u == nv {
		return
	}
	if nv.PasswordHash != nil || len(*nv.PasswordHash) != 0 {
		u.PasswordHash = nv.PasswordHash
	}

}
