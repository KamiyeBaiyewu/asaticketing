package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// CauseID is the identifier for a ticket
type CauseID string

// NilCauseID is an empty CauseID
var NilCauseID CauseID

//Cause - represents Tickets Cause
type Cause struct {
	ID          CauseID `json:"id,omitempty" db:"cause_id"`
	Name        *string    `json:"name,omitempty" db:"name"`
	Description *string    `json:"description,omitempty" db:"description"`
	Weight      *int       `json:"weight,omitempty" db:"weight"`
	UserID      UserID     `json:"-" db:"created_by"`
	IsStandard  *bool      `json:"-" db:"is_standard"`
	CreatedAt   *time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"  db:"deleted_at"`
}

// Verify -  Verify the Cause values
func (c *Cause) Verify() error {

	if c.Name == nil || (c.Name != nil && len(*c.Name) == 0) {
		return errors.New("Name is required")
	}

	if c.Description == nil || (c.Description != nil && len(*c.Description) == 0) {
		c.Description = func() *string { s := ""; return &s }()
	}
	if c.Weight == nil || (c.Weight != nil && *c.Weight >= 10) {
		c.Weight = func() *int { b := 10; return &b }()
	} else if c.Weight != nil && *c.Weight <= 1 {
		c.Weight = func() *int { b := 1; return &b }()
	}

	return nil
}

// Decode - Cause to JSON
func (c *Cause) Decode(reader io.Reader) error {
	if err := json.NewDecoder(reader).Decode(&c); err != nil {
		return err
	}
	return nil
}

// UpdateValues is used to update empty values
func (c *Cause) UpdateValues(nv *Cause) { //nv means new values
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
	if nv.Weight != nil {
		if *nv.Weight != 0 {
			if *nv.Weight >= 10 {
				c.Weight = func() *int { b := 10; return &b }()
			} else if *nv.Weight <= 1 {
				c.Weight = func() *int { b := 1; return &b }()

			} else {

				c.Weight = nv.Weight
			}
		}
	}

}
