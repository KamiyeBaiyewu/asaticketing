package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// SLAID is the identifier for a ticket
type SLAID string

// NilSLAID is an empty SLAID
var NilSLAID SLAID

//SLA - represents Tickets SLA
type SLA struct {
	ID          SLAID      `json:"id,omitempty" db:"agreement_id"`
	Name        *string    `json:"name,omitempty" db:"name"`
	GracePeriod *int       `json:"grace_period,omitempty" db:"grace_period"`
	Weight      *int       `json:"weight,omitempty" db:"weight"`
	CreatedAt   *time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"  db:"deleted_at"`
}

// Decode - SLA to JSON
func (s *SLA) Decode(reader io.Reader) error {
	if err := json.NewDecoder(reader).Decode(&s); err != nil {
		return err
	}
	return nil
}

// Verify -  Verify the SLA values
func (s *SLA) Verify() error {

	if s.Name == nil || (s.Name != nil && len(*s.Name) == 0) {
		return errors.New("Name is required")
	}

	if s.Weight == nil || (s.Weight != nil && *s.Weight >= 10)  {
		s.Weight = func() *int { b := 10; return &b }()
	}else if (s.Weight != nil && *s.Weight <= 1){
		s.Weight = func() *int { b := 1; return &b }()
	}
	if s.GracePeriod == nil || (s.GracePeriod != nil && *s.GracePeriod == 0) {
		return errors.New("Grace Period is required")
	}

	return nil
}

// UpdateValues is used to update empty values
func (s *SLA) UpdateValues(nv *SLA) { //nv means new values
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
	if nv.GracePeriod != nil {
		if *nv.GracePeriod != 0 {
			s.GracePeriod = nv.GracePeriod
		}
	}

}
