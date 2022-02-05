package model

import "time"

// InboudMailID is the identifier for the student
type InboudMailID string

// NilInboudMailID is an empty InboudMailID
var NilInboudMailID InboudMailID

// InboudMail is a structure that represents InboudMail Object
type InboudMail struct {
	ID          InboudMailID `json:"id,omitempty" db:"email_id"`
	Name        *string      `json:"name,omitempty" db:"name"`
	Status      *string      `json:"status,omitempty" db:"status"`
	Address     *string      `json:"address,omitempty" db:"address"`
	EmailUser   *string      `json:"email_user,omitempty" db:"email_user"`
	EmailSecret *string      `json:"email_secret,omitempty" db:"email_secret"`
	Port        *int         `json:"port,omitempty" db:"port"`
	Secured     *bool        `json:"secured,omitempty" db:"secured"`
	Mailbox     *string      `json:"mailbox,omitempty" db:"mailbox"`
	IsPrimary   *bool        `json:"is_primary,omitempty" db:"is_primary"`
	LastSeq     *int         `json:"last_seq,omitempty" db:"last_seq"`
	PollPeriod  *int         `json:"poll_period,omitempty" db:"poll_period"`
	LastSynced  *time.Time   `json:"last_synced,omitempty"  db:"last_synced"`
	UserID      UserID       `json:"created_by,omitempty" db:"created_by"`
	DeleteSeen  *bool        `json:"delete_seen,omitempty" db:"delete_seen"`
	CreatedAt   *time.Time   `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt   *time.Time   `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty"  db:"deleted_at"`
}
