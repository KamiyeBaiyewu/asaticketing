package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// NilTicketID the identifier for a ticket
type TicketID string

// NilTicketID an empty TicketID
var NilTicketID TicketID

// Ticket - represents Tickets
type Ticket struct {
	ID          TicketID   `json:"id,omitempty" db:"ticket_id"`
	Subject     *string    `json:"subject,omitempty" db:"subject"`
	Description *string    `json:"description,omitempty" db:"description"`
	Code        *int       `json:"number,omitempty" db:"number"`
	UserID      UserID     `json:"-" db:"created_by"`
	CreatedBy   *User      `json:"created_by,omitempty"`
	CategoryID  CategoryID `json:"category_id,omitempty" db:"category_id"`
	Category    *Category  `json:"category,omitempty"`
	StatusID    StatusID   `json:"status_id,omitempty" db:"status_id"`
	Status      *Status    `json:"status,omitempty"`
	PriorityID  PriorityID `json:"priority_id,omitempty" db:"priority_id"`
	Priority    *Priority  `json:"priority,omitempty"`
	SLAID       SLAID      `json:"sla_id,omitempty" db:"sla_id"`
	SLA         *SLA       `json:"sla,omitempty"`
	SourceID    SourceID   `json:"source_id,omitempty" db:"source_id"`
	Source      *Source    `json:"source,omitempty"`

	DueDate    *time.Time `json:"deadline,omitempty"  db:"deadline"`
	ClosedAt   *time.Time `json:"closed_at,omitempty"  db:"closed_at"`
	CreatedAt  *time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"  db:"deleted_at"`
	AssignedID UserID     `json:"assigned_id,omitempty" db:"assigned_to"`
	AssignedTo *User      `json:"assigned_to,omitempty"`

	// Users Are represent the people assigned to the ticket
	Users []*User `json:"users,omitempty"`

	
	// Helpful for retrieving Tickets fromt the database

	/* MISC */
	Grace *int `json:"-" db:"grace"`

	// For closed tickets
	ClosingRemark *ClosingRemark `json:"closing_remark,omitempty"`
}

// Decode - UserParameters to JSON
func (t *Ticket) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&t)
}

// Verify -  ensures required variables are present
func (t *Ticket) Verify() error {

	if t.UserID == NilUserID {
		return errors.New("User is required")
	}
	if t.CategoryID == NilCategoryID {
		return errors.New("Category is required")
	}
	if t.StatusID == NilStatusID {
		return errors.New("Status is required")
	}
	if t.PriorityID == NilPriorityID {
		return errors.New("Priority is required")
	}
	if t.SourceID == NilSourceID {
		return errors.New("Source is required")
	}
	if t.SLAID == NilSLAID {
		return errors.New("SLA is required")
	}

	if t.Subject == nil || (t.Subject != nil && len(*t.Subject) == 0) {
		return errors.New("Subject is required")
	}

	if t.Description == nil || (t.Description != nil && len(*t.Description) == 0) {
		return errors.New("Description is required")
	}

	return nil
}

// UpdateValues is used to update empty values
func (t *Ticket) UpdateValues(nv *Ticket) { //nv means new values
	// Avoid updating the same values
	if t == nv {
		return
	}

	/* 	if t.UserID != NilUserID {
		t.UserID = nv.UserID
	} */
	if nv.CategoryID != NilCategoryID {
		t.CategoryID = nv.CategoryID
	}
	if nv.StatusID != NilStatusID {
		t.StatusID = nv.StatusID
	}
	if nv.PriorityID != NilPriorityID {
		t.PriorityID = nv.PriorityID
	}
	if nv.SourceID != NilSourceID {
		t.SourceID = nv.SourceID
	}
	/* if nv.SLAID != NilSLAID {
		t.SLAID = nv.SLAID
	} */

	if nv.AssignedID != NilUserID {
		t.AssignedID = nv.AssignedID
	}

	if nv.Description != nil {
		if len(*nv.Description) != 0 {
			t.Description = nv.Description
		}
	}

}
