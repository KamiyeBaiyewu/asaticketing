package database

import (
	"context"

	"github.com/lib/pq"
	apiErr "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ContactsDB - holds the interface to the Contacts
type ContactsDB interface {
	CreateContact(ctx context.Context, contact *model.Contact) (err error)
	GetContactByID(ctx context.Context, contactID *model.ContactID) (*model.Contact, error)
	ListAllContacts(ctx context.Context) ([]*model.Contact, error)
	UpdateContact(ctx context.Context, contact *model.Contact) error
	DeleteContact(ctx context.Context, contactID *model.ContactID) (bool, error)


}

//CONFIRM THE DB INSERTIONS
const createContactQuery = `
		INSERT INTO contacts (
		firstname,lastname,phone_no,email,created_by
		)
			VALUES (
				:firstname,:lastname,:phone_no,:email,:created_by
				)
				RETURNING contact_id`

func (d *database) CreateContact(ctx context.Context, contact *model.Contact) (err error) {
	rows, err := d.conn.NamedQueryContext(ctx, createContactQuery, contact)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Constraint == "unique_contact" {
				err = apiErr.ErrTicketExists
				return
			}
			// One of the Foreign key ID is missing
		
			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
				"Error":          err,
			}).Info()
			return apiErr.ErrCreatingContact
		}
		logrus.WithError(err).Error()
		return apiErr.ErrCreatingContact
	}

	rows.Next()
	if err := rows.Scan(&contact.ID); err != nil {
		err = errors.Wrap(err, "Could not get the Contact ID")
	}
	return
}



// correct the  db insertions
const getContactByIDQuery = `
	SELECT ct.contact_id, ct.firstname, ct.lastname, ct.phone_no,ct.email ,ct.created_by, ct.created_at, ct.deleted_at
	FROM contacts ct
	WHERE ct.contact_id = $1
	AND ct.deleted_at IS NULL
`

func (d *database) GetContactByID(ctx context.Context, contactID *model.ContactID) (*model.Contact, error) {
	contact := model.Contact{}
	if err := d.conn.GetContext(ctx, &contact, getContactByIDQuery, contactID); err != nil {

		logrus.WithError(err).Error()
		return nil, apiErr.ErrNotFound
	}
	return &contact, nil
}

// correct db insertions

const listAllContactsQuery = `
	SELECT ct.contact_id, ct.firstname, ct.lastname, ct.phone_no,ct.email ,ct.created_by, ct.created_at, ct.deleted_at
	FROM contacts ct
	WHERE ct.deleted_at IS NULL`

func (d *database) ListAllContacts(ctx context.Context) ([]*model.Contact, error) {
	contacts := []*model.Contact{}
	if err := d.conn.SelectContext(ctx, &contacts, listAllContactsQuery); err != nil {
		println(err.Error())
		return nil, err
	}
	return contacts, nil
}

const updateContactQuery = `
		UPDATE contacts
		SET firstname = :firstname,
		lasttname = :lastname,
		phone_no = :phone_no,
		email = :email,
		updated_at = NOW()
		WHERE contact_id = :contact_id
		AND deleted_at is null`

func (d *database) UpdateContact(ctx context.Context, contact *model.Contact) error {
	result, err := d.conn.NamedExecContext(ctx, updateContactQuery, contact)
	if err != nil {
		println("PQERROR => ", err.Error())
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {


			
			}

			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
			}).Info()
			return apiErr.ErrUpdatingContact
		}
		return apiErr.ErrUpdatingContact
	}

	rows, err := result.RowsAffected()

	if err != nil || rows == 0 {
		return errors.New("Contact Not found")
	}

	return nil
}

const deleteContactQuery = `
	UPDATE contacts
	SET deleted_at = NOW()
	WHERE contact_id = $1 
	AND deleted_at is NULL;
	`

func (d *database) DeleteContact(ctx context.Context, contactID *model.ContactID) (bool, error) {
	result, err := d.conn.ExecContext(ctx, deleteContactQuery, contactID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}
	return true, nil
}
