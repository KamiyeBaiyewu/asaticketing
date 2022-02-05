package requests

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
)

// UserParameters - structure that represents data the clients sends
type UserParameters struct {
	model.User
	model.SessionData
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// Decode - UserParameters to JSON
func (u *UserParameters) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&u)
}

//Verify all fields before create or update
func (u *UserParameters) Verify() error {
	// Verify user and Session data
	if err := u.User.Verify(); err != nil {
		return err
	}
	if err := u.SessionData.Verify(); err != nil {
		return err
	}

	if len(u.Password) > 0 {

		if u.Password == "" {
			return errors.New("Password field is required")
		}
		if u.ConfirmPassword == "" {
			return errors.New("Confirm Password field is required")
		}

		// Ensure password and Confirm Password Match
		if u.Password != u.ConfirmPassword {

			return errors.New("Passwords don't match")
		}
	}

	return nil
}

// ToUser helps to convert a user request into a user type
func (u *UserParameters) ToUser() *model.User {

	user := &model.User{}
	if u.ID != model.NilUserID {
		user.ID = u.ID
	}
	if u.Firstname != nil {
		user.Firstname = u.Firstname
	}
	if u.Type != nil {
		user.Type = u.Type
	}
	if u.Lastname != nil {
		user.Lastname = u.Lastname
	}
	if u.Email != nil {
		user.Email = u.Email
	}
	if u.IsActive != nil {
		user.IsActive = u.IsActive
	}
	if u.PasswordHash != nil {
		user.PasswordHash = u.PasswordHash
	}
	if u.IsSystem != nil {
		user.IsSystem = u.IsSystem
	}
	if u.RoleID != model.NilRoleID {
		user.RoleID = u.RoleID
	}
	if u.Password != "" {
		user.Password = func() *string { s := u.Password; return &s }()
	}
	if u.ConfirmPassword != "" {
		user.ConfirmPassword = func() *string { s := u.ConfirmPassword; return &s }()
	}

	return user
}
