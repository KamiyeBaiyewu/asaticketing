package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// ClosingRemarkID the identifier for a ticket
type ClosingRemarkID string

// NilClosingRemarkID an empty ClosingRemarkID
var NilClosingRemarkID ClosingRemarkID

// ClosingRemark - represents ClosingRemarks
type ClosingRemark struct {
	ID       ClosingRemarkID `json:"id,omitempty" db:"remark_id"`
	UserID   UserID          `json:"-" db:"closed_by"`
	ClosedBy *User           `json:"closed_by,omitempty"`
	TicketID TicketID        `json:"ticket_id,omitempty" db:"ticket_id"`
	Ticket   *Ticket         `json:"ticket,omitempty"`
	CauseID  CauseID         `json:"cause_id,omitempty" db:"cause_id"`
	Cause    *Cause          `json:"cause,omitempty"`
	Remark   *string         `json:"remark,omitempty" db:"remark"`

	CreatedAt *time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"  db:"deleted_at"`
}

// Decode - UserParameters to JSON
func (cr *ClosingRemark) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&cr)
}

// Verify -  ensures required variables are present
func (cr *ClosingRemark) Verify() error {

	if cr.UserID == NilUserID {
		return errors.New("User is required")
	}
	if cr.TicketID == NilTicketID {
		return errors.New("Ticket is required")
	}
	if cr.CauseID == NilCauseID {
		return errors.New("Cause is required")
	}

	/* 	if cr.Remark == nil || (cr.Remark != nil && len(*cr.Remark) == 0) {
		cr.Remark = func() *string {s:= ""; return &s}()
	} */

	return nil
}

// UpdateValues is used to update empty values
func (cr *ClosingRemark) UpdateValues(nv *ClosingRemark) { //nv means new values
	// Avoid updating the same values
	if cr == nv {
		return
	}

	if nv.CauseID != NilCauseID {
		cr.CauseID = nv.CauseID
	}

	if nv.Remark != nil {
		if len(*nv.Remark) != 0 {
			cr.Remark = nv.Remark
		}
	}

}
