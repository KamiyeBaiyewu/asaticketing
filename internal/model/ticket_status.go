package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// StatusID is the identifier for a ticket
type StatusID string

// NilStatusID is an empty StatusID
var NilStatusID StatusID

//Status - represents Tickets Status
type Status struct {
	ID        StatusID   `json:"id,omitempty" db:"status_id"`
	Name      *string    `json:"name,omitempty" db:"name"`
	Weight    *int       `json:"weight,omitempty" db:"weight"`
	CreatedAt *time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"  db:"deleted_at"`
}

// Decode - Status to JSON
func (s *Status) Decode(reader io.Reader) error {
	if err := json.NewDecoder(reader).Decode(&s); err != nil {
		return err
	}
	return nil
}

// Verify -  Verify the Status values
func (s *Status) Verify() error {

	if s.Name == nil || (s.Name != nil && len(*s.Name) == 0) {
		return errors.New("Status is required")
	}

	if s.Weight == nil || (s.Weight != nil && *s.Weight >= 10)  {
		s.Weight = func() *int { b := 10; return &b }()
	}else if (s.Weight != nil && *s.Weight <= 1){
		s.Weight = func() *int { b := 1; return &b }()
	}

	return nil
}

// UpdateValues is used to update empty values
func (s *Status) UpdateValues(nv *Status) { //nv means new values
	// Avoid updating the same values
	if s == nv {
		return
	}

	if nv.Name != nil {
		if len(*nv.Name) != 0 {
			s.Name = nv.Name
		}
	}
	if nv.Weight != nil {
		if *nv.Weight >= 10 {
			s.Weight = func() *int { b := 10; return &b }()
		} else if *nv.Weight <= 1 {
			s.Weight = func() *int { b := 1; return &b }()

		} else {

			s.Weight = nv.Weight
		}
	}

}
