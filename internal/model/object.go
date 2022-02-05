package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// ObjectID is the identifier for a ticket
type ObjectID string

// NilObjectID is an empty ObjectID
var NilObjectID ObjectID

//Object - represents System Object
type Object struct {
	ID          ObjectID   `json:"id,omitempty" db:"object_id"`
	Name        *string    `json:"name,omitempty" db:"name"`
	Description *string    `json:"description,omitempty" db:"description"`
	IsStandard  *bool      `json:"-" db:"is_standard"`
	CreatedBy   UserID     `json:"-" db:"created_by"`
	CreatedAt   *time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"  db:"deleted_at"`
}

// Verify -  Verify the Object values
func (c *Object) Verify() error {


	if c.Name == nil || (c.Name != nil && len(*c.Name) == 0) {
		return errors.New("Object is required")
	}
	if c.CreatedBy == NilUserID {
		return errors.New("UserID is required")
	}

	return nil
}

// Decode - Object to JSON
func (c *Object) Decode(reader io.Reader) error {
	if err := json.NewDecoder(reader).Decode(&c); err != nil {
		return err
	}
	return nil
}

// UpdateValues is used to update empty values
func (c *Object) UpdateValues(nv *Object) { //nv means new values
	// Avoid updating the same values
	if c == nv {
		return
	}

	if nv.Name != nil {
		if len(*nv.Name) != 0 {
			c.Name = nv.Name
		}
	}
	if nv.Description != nil {
		if len(*nv.Description) != 0 {
			c.Description = nv.Description
		}
	}

}
