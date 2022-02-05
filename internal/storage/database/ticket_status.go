package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apiErr "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
)

// TicketStatusDB - interface that dabase connection must implement
type TicketStatusDB interface {
	CreateStatus(ctx context.Context, status *model.Status) error
	GetStatusByID(ctx context.Context, statusID *model.StatusID) (*model.Status, error)
	UpdateStatus(ctx context.Context, status *model.Status) error
	ListAllStatus(ctx context.Context) ([]*model.Status, error)
	DeleteStatus(ctx context.Context, StatusID *model.StatusID) (bool, error)
}


const createStatusQuery = `INSERT INTO ticket_statuses (
	 name, weight
	)
	VALUES (
		 :name, :weight
		)
		RETURNING status_id`

func (d *database) CreateStatus(ctx context.Context, status *model.Status) (err error) {

	rows, err := d.conn.NamedQueryContext(ctx, createStatusQuery, status)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" && pqError.Constraint == "ticket_statuses_name" {
				err = apiErr.ErrStatusExists
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
	if err := rows.Scan(&status.ID); err != nil {
		err = errors.Wrap(err, "Could not get the Status ID")
	}
	return
}

const getStatusByIDQuery = `
	SELECT status_id, name, weight, created_at, updated_at, deleted_at
	FROM ticket_statuses
	WHERE status_id = $1 
	AND deleted_at is NULL`

func (d *database) GetStatusByID(ctx context.Context, statusID *model.StatusID) (*model.Status, error) {

	status := model.Status{}
	if err := d.conn.GetContext(ctx, &status, getStatusByIDQuery, statusID); err != nil {

		return nil, err
	}
	return &status, nil

}

const updateStatusQuery = `
	UPDATE ticket_statuses
	SET 
		status_id = :status_id,
		name = :name,
		weight = :weight,
		updated_at = NOW()
	WHERE status_id = :status_id 
	AND deleted_at is NULL`

func (d *database) UpdateStatus(ctx context.Context, status *model.Status) error {

	//println(*status.PasswordHash)
	result, err := d.conn.NamedExecContext(ctx, updateStatusQuery, status)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("Status Not found")
	}
	return nil
}

const listAllStatusQuery = `
	SELECT status_id, name, weight, created_at, updated_at, deleted_at
	FROM ticket_statuses
	WHERE deleted_at is NULL
	ORDER BY weight ASC`

func (d *database) ListAllStatus(ctx context.Context) ([]*model.Status, error) {

	categories := []*model.Status{}
	if err := d.conn.SelectContext(ctx, &categories, listAllStatusQuery); err != nil {
		return nil, errors.Wrap(err, "could not get statuss")
	}
	return categories, nil
}

const deleteStatusQuery = `
	UPDATE ticket_statuses
	SET deleted_at = NOW()
	WHERE status_id = $1 
	AND deleted_at is NULL`

func (d *database) DeleteStatus(ctx context.Context, StatusID *model.StatusID) (bool, error) {

	result, err := d.conn.ExecContext(ctx, deleteStatusQuery, StatusID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
