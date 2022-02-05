package model

import (
	"encoding/json"
	"fmt"
	"io"
)

// Credentials - used for the login Api
type Credentials struct {
	SessionData

	// Username/Password Login
	Email    string `json:"email"`
	Password string `json:"password"`

	// DeviceID DeviceID `json:"deviceID"`
	// Google and Facebook Login coming soon
}

// Principal is an authenticated entity
type Principal struct {
	UserID UserID `json:"userID,omitemoty"`
	Name   string `json:"name,omitempty"`
	Role   string `json:"role,omitempty"`
	Type   string `json:"type,omitempty"`
}

// Decode - Credentials to JSON
func (c *Credentials) Decode(reader io.Reader) error {
	if err := json.NewDecoder(reader).Decode(&c); err != nil {
		return err
	}
	return nil
}

// NilPricipal is an uninitialized principal
var NilPricipal Principal

func (p *Principal) String() string {
	if p.UserID != "" {
		return fmt.Sprintf("UserID[%s]", p.UserID)
	}

	return "(none)"
}
