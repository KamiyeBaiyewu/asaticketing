package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// NoteID is the identifier for the note
type NoteID string

// NilNoteID is an empty NoteID
var NilNoteID NoteID

// Note - represents User Notes
type Note struct {
	// Note Note `json:"note" db:"note"`
	ID        NoteID     `json:"id,omitempty" db:"note_id"`
	Note      *string    `json:"note,omitempty" db:"note"`
	TicketID  TicketID   `json:"ticket_id,omitempty" db:"ticket_id"`
	UserID    UserID     `json:"-" db:"created_by"`
	CreatedBy *User      `json:"created_by,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"  db:"deleted_at"`
}

// Decode - UserParameters to JSON
func (n *Note) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&n)
}

//Verify all fields before create or update
func (n *Note) Verify() error {

	if n.Note == nil || (n.Note != nil && len(*n.Note) == 0) {
		return errors.New("Note is required")
	}
	if n.TicketID == NilTicketID {
		return errors.New("Ticket is required")
	}

	// Ensure we know who crated the note
	if n.UserID == NilUserID {
		return errors.New("User is Required")
	}

	return nil
}

// UpdateValues is used to update empty values
func (n *Note) UpdateValues(nv *Note) { //nv means new values
	// Avoid updating the same values
	if n == nv {
		return
	}
	if nv.Note != nil || len(*nv.Note) != 0 {
		n.Note = nv.Note
	}

}
