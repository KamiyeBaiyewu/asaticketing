package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// PriorityID is the identifier for a ticket
type PriorityID string

// NilPriorityID is an empty PriorityID
var NilPriorityID PriorityID

//Priority - represents Tickets Priority
type Priority struct {
	ID        PriorityID `json:"id,omitempty" db:"priority_id"`
	Name      *string    `json:"name,omitempty" db:"name"`
	Weight    *int       `json:"weight,omitempty" db:"weight"`
	CreatedAt *time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"  db:"deleted_at"`
}

// Decode - Priority to JSON
func (p *Priority) Decode(reader io.Reader) error {
	if err := json.NewDecoder(reader).Decode(&p); err != nil {
		return err
	}
	return nil
}

// Verify -  Verify the Priority values
func (p *Priority) Verify() error {

	if p.Name == nil || (p.Name != nil && len(*p.Name) == 0) {
		return errors.New("Name is required")
	}
	if p.Weight == nil || (p.Weight != nil && *p.Weight >= 10)  {
		p.Weight = func() *int { b := 10; return &b }()
	}else if (p.Weight != nil && *p.Weight <= 1){
		p.Weight = func() *int { b := 1; return &b }()
	}

	return nil
}

// UpdateValues is used to update empty values
func (p *Priority) UpdateValues(nv *Priority) { //nv means new values
	// Avoid updating the same values
	if p == nv {
		return
	}

	if nv.Name != nil {
		if len(*nv.Name) != 0 {
			p.Name = nv.Name
		}
	}
	if nv.Weight != nil {

		if *nv.Weight >= 10 {
			p.Weight = func() *int { b := 10; return &b }()
		} else if *nv.Weight <= 1 {
			p.Weight = func() *int { b := 1; return &b }()

		} else {

			p.Weight = nv.Weight
		}

	
	}

}
