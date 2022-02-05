package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/sirupsen/logrus"
)

// InboundEmaiiDB - Holds all the information to connect to the email server
type InboundEmaiiDB interface {
	GetInboundMail(ctx context.Context) (*model.InboudMail, error)
	UpdateInboudMail(ctx context.Context, inboundMail *model.InboudMail) error
}

const getInboundMailQuery = `
	SELECT 
	email_id, "name", status, address, email_user, 
	email_secret, port, secured, mailbox, is_primary, 
	last_seq, last_synced,poll_period, created_by, delete_seen, 
	created_at, updated_at, deleted_at
	FROM 
	inbound_emails;
`

func (d *database) GetInboundMail(ctx context.Context) (*model.InboudMail, error) {

	var inboudMail model.InboudMail
	if err := d.conn.GetContext(ctx, &inboudMail, getInboundMailQuery); err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
			}).Info()
		}

		return nil, err
	}

	return &inboudMail, nil
}

const updateInboudMailQuery = `
	UPDATE inbound_emails
	SET
	name=:name, 
	status=:status, 
	address=:address, 
	email_user=:email_user, 
	email_secret=:email_secret, 
	port=:port,
	secured=:secured, 
	mailbox=:mailbox, 
	is_primary=:is_primary, 
	last_seq=:last_seq, 
	last_synced=:last_synced, 
	delete_seen=:delete_seen, 
	updated_at=NOW()
	WHERE email_id = :email_id
	AND deleted_at is null
`

func (d *database) UpdateInboudMail(ctx context.Context, inboundMail *model.InboudMail) error {

	result, err := d.conn.NamedExecContext(ctx, updateInboudMailQuery, inboundMail)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}
