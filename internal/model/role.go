package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// RoleID is the identifier for the role
type RoleID string

// NilRoleID is an empty RoleID
var NilRoleID RoleID

// Role - represents User Roles
type Role struct {
	// Role Role `json:"role" db:"role"`
	ID          RoleID     `json:"id,omitempty" db:"role_id"`
	Name        *string    `json:"name" db:"name"`
	Description *string    `json:"description,omitempty" db:"description"`
	CreatedAt   *time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"  db:"deleted_at"`
}

// Decode - UserParameters to JSON
func (r *Role) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&r)
}

//Verify all fields before create or update
func (r *Role) Verify() error {

	if r.Name == nil || (r.Name != nil && len(*r.Name) == 0) {
		return errors.New("Name is required")
	}

	return nil
}

// UpdateValues is used to update empty values
func (r *Role) UpdateValues(nv *Role) { //nv means new values
	// Avoid updating the same values
	if r == nv {
		return
	}
	if nv.Name != nil || len(*nv.Name) != 0 {
		r.Name = nv.Name
	}
	if nv.Description != nil || len(*nv.Description) != 0 {
		r.Description = nv.Description
	}

}
