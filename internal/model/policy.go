package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// Action - holds action  value
type Action string

// NilAction is an empty Action
var NilAction Action

// PolicyID is the identifier for the policy
type PolicyID string

// NilPolicyID is an empty PolicyID
var NilPolicyID PolicyID

// Policy - represents User Policies
type Policy struct {
	ID       PolicyID `json:"id,omitempty" db:"policy_id"`
	RoleID   RoleID   `json:"role_id,omitempty" db:"role_id"`
	ObjectID ObjectID `json:"object_id,omitempty" db:"object_id"`
	Action   *string  `json:"action,omitempty" db:"action"`
	UserID     UserID     `json:"-" db:"created_by"` //Can also be used to retrieve
	IsStandard *bool      `json:"-" db:"is_standard"`
	CreatedAt  *time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"  db:"deleted_at"`

	/* MISC */
	CreatedBy *User   `json:"created_by,omitempty" `
	Role      *Role   `json:"role,omitempty" `
	Object    *Object `json:"object,omitempty"`
}

// Decode - UserParameters to JSON
func (p *Policy) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&p)
}

// Verify ensures all values are safe
func (p *Policy) Verify() error {
	if p.RoleID == NilRoleID {
		return errors.New("Role is required")
	}
	if p.ObjectID == NilObjectID {
		return errors.New("Oject is Required")
	}
	if p.Action == nil || (p.Action != nil && len(*p.Action) == 0) {
		return errors.New("Action is required")
	}
	if p.UserID == NilUserID {
		return errors.New("User is required")
	}
	return nil
}
