package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apiErr "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
)

// TicketCauseDB - interface that dabase connection must implement
type TicketCauseDB interface {
	CreateCause(ctx context.Context, cause *model.Cause) error
	GetCauseByID(ctx context.Context, causeID *model.CauseID) (*model.Cause, error)
	UpdateCause(ctx context.Context, cause *model.Cause) error
	ListAllCauses(ctx context.Context) ([]*model.Cause, error)
	DeleteCause(ctx context.Context, CauseID *model.CauseID) (bool, error)
}


const createCauseQuery = `
	INSERT INTO ticket_causes (
	 name, description,weight,created_by
	)
	VALUES (
		 :name, :description, :weight, :created_by
		)
		RETURNING cause_id`

func (d *database) CreateCause(ctx context.Context, cause *model.Cause) (err error) {

	rows, err := d.conn.NamedQueryContext(ctx, createCauseQuery, cause)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" && pqError.Constraint == "ticket_causes_name" {
				err = apiErr.ErrCauseExists
				return
			}

			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
			}).Info()
		}
		return
	}

	rows.Next()
	if err := rows.Scan(&cause.ID); err != nil {
		err = errors.Wrap(err, "Could not get the Cause ID")
	}
	return
}

const getCauseByIDQuery = `
	SELECT cause_id, name, description, weight,is_standard,created_by, created_at, updated_at, deleted_at
	FROM ticket_causes
	WHERE cause_id = $1 
	AND deleted_at is NULL`

func (d *database) GetCauseByID(ctx context.Context, causeID *model.CauseID) (*model.Cause, error) {

	cause := model.Cause{}
	if err := d.conn.GetContext(ctx, &cause, getCauseByIDQuery, causeID); err != nil {

		return nil, err
	}
	return &cause, nil

}

const updateCauseQuery = `
	UPDATE ticket_causes
	SET 
		name = :name,
		description = :description,
		weight = :weight,
		updated_at = NOW()
	WHERE cause_id = :cause_id
	AND deleted_at is NULL`

func (d *database) UpdateCause(ctx context.Context, cause *model.Cause) error {

	//println(*cause.PasswordHash)
	result, err := d.conn.NamedExecContext(ctx, updateCauseQuery, cause)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("Cause Not found")
	}
	return nil
}

const listAllCausesQuery = `
	SELECT cause_id, name, description, weight,is_standard,created_by, created_at, updated_at, deleted_at
	FROM ticket_causes
	WHERE deleted_at is NULL
	ORDER BY weight ASC`

func (d *database) ListAllCauses(ctx context.Context) ([]*model.Cause, error) {

	causes := []*model.Cause{}
	if err := d.conn.SelectContext(ctx, &causes, listAllCausesQuery); err != nil {
		return nil, apiErr.ErrInternal
	}
	return causes, nil
}

const deleteCauseQuery = `
	UPDATE ticket_causes
	SET deleted_at = NOW()
	WHERE cause_id = $1 AND deleted_at is NULL`

func (d *database) DeleteCause(ctx context.Context, CauseID *model.CauseID) (bool, error) {

	result, err := d.conn.ExecContext(ctx, deleteCauseQuery, CauseID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
