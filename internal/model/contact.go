package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

//ContactID the identifier for a contact
type ContactID string

// NilContactID an empty ContactID
var NilContactID ContactID

/*
contact_id uuid PRIMARY KEY,
    firstname TEXT,
    lastname TEXT,
    phone_no INTEGER,
    email text NOT NULL,
    created_by uuid REFERENCES users,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
	deleted_at TIMESTAMP WITH TIME ZONE NOT NULL
*/

// Contact - represents Contacts
type Contact struct {
	ID        ContactID `json:"id,omitempty" db:"contact_id"`
	Firstname *string   `json:"firstname,omitempty" db:"firstname"`
	Lastname  *string   `json:"lastname,omitempty" db:"lastname"`
	PhoneNo   *string   `json:"phoneno,omitempty" db:"phone_no"`
	Email     *string   `json:"email,omitempty" db:"email"`

	UserID UserID     `json:"-"  db:"created_by"`
	CreatedAt *time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"  db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"  db:"deleted_at"`


	/* MISC */
	CreatedBy *User `json:"created_by,omitempty"`
}

// Decode - UserParameters to JSON
func (c *Contact) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&c)
}

// check this again

// Verify -  ensures required variables are present
func (c *Contact) Verify() error {

	if c.Firstname == nil {
		return errors.New("firstname is required")
	}

	if c.Lastname == nil {
		return errors.New("lastname is required")
	}
	if c.UserID == NilUserID {
		return errors.New("The creator is missing")
	}

	return nil
}

// UpdateValues is used to update empty values
func (c *Contact) UpdateValues(nv *Contact) { //nv means new values
	// Avoid updating the same values
	if c == nv {
		return
	}

 if nv.Firstname != nil{
	 c.Firstname = nv.Firstname
 }
 if nv.Lastname != nil{
	 c.Lastname = nv.Lastname
 }
 if nv.PhoneNo != nil{
	 c.PhoneNo = nv.PhoneNo
}
 if nv.Email != nil{
	 c.Email = nv.Email	
}

}
