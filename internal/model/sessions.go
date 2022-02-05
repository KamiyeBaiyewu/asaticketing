package model

import "errors"

// DeviceID is an abstract type to represent a specific device
type DeviceID string

// NilDeviceID represents an empty DeviceID
var NilDeviceID DeviceID

// Session represents structure used to store sessions in the database
type Session struct {
	UserID       UserID   `db:"user_id"`
	DeviceID     DeviceID `db:"device_id"`
	RefreshToken string   `db:"refresh_token"`
	ExpiresAt    int64    `db:"expires_at"`
}

// SessionData used to represent data sent in json body with requests
type SessionData struct {
	DeviceID DeviceID `json:"deviceID"`
}

//Verify all fields before create or update
func (s *SessionData) Verify() error {
	
	if len(s.DeviceID) == 0 {
		return errors.New("DeviceID is required")
	}

	return nil
}
