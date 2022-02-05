package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apiErr "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
)

// TicketPriorityDB - interface that dabase connection must implement
type TicketPriorityDB interface {
	CreatePriority(ctx context.Context, priority *model.Priority) error
	GetPriorityByID(ctx context.Context, priorityID *model.PriorityID) (*model.Priority, error)
	UpdatePriority(ctx context.Context, priority *model.Priority) error
	ListAllPriorities(ctx context.Context) ([]*model.Priority, error)
	DeletePriority(ctx context.Context, PriorityID *model.PriorityID) (bool, error)
}


const createPriorityQuery = `INSERT INTO ticket_priorities (
	 name, weight
	)
	VALUES (
		 :name, :weight
		)
		RETURNING priority_id`

func (d *database) CreatePriority(ctx context.Context, priority *model.Priority) (err error) {

	rows, err := d.conn.NamedQueryContext(ctx, createPriorityQuery, priority)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" && pqError.Constraint == "ticket_priorities_name" {
				err = apiErr.ErrPriorityExists
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
	if err := rows.Scan(&priority.ID); err != nil {
		err = errors.Wrap(err, "Could not get the Priority ID")
	}
	return
}

const getPriorityByIDQuery = `
	SELECT priority_id, name, weight, created_at, updated_at, deleted_at
	FROM ticket_priorities
	WHERE priority_id = $1 
	AND deleted_at is NULL`

func (d *database) GetPriorityByID(ctx context.Context, priorityID *model.PriorityID) (*model.Priority, error) {
	priority := model.Priority{}
	if err := d.conn.GetContext(ctx, &priority, getPriorityByIDQuery, priorityID); err != nil {

		return nil, err
	}
	return &priority, nil

}

const updatePriorityQuery = `
	UPDATE ticket_priorities
	SET 
		priority_id = :priority_id,
		name = :name,
		weight = :weight,
		updated_at = NOW()
	WHERE priority_id = :priority_id 
	AND deleted_at is NULL`

func (d *database) UpdatePriority(ctx context.Context, priority *model.Priority) error {
	//println(*priority.PasswordHash)
	result, err := d.conn.NamedExecContext(ctx, updatePriorityQuery, priority)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("Priority Not found")
	}
	return nil
}

const listAllPrioritiesQuery = `
	SELECT priority_id, name, weight, created_at, updated_at, deleted_at
	FROM ticket_priorities
	WHERE deleted_at is NULL
	ORDER BY weight ASC`

func (d *database) ListAllPriorities(ctx context.Context) ([]*model.Priority, error) {
	categories := []*model.Priority{}
	if err := d.conn.SelectContext(ctx, &categories, listAllPrioritiesQuery); err != nil {
		return nil, errors.Wrap(err, "could not get prioritys")
	}
	return categories, nil
}

const deletePriorityQuery = `
	UPDATE ticket_priorities
	SET deleted_at = NOW()
	WHERE priority_id = $1 
	AND deleted_at is NULL`

func (d *database) DeletePriority(ctx context.Context, PriorityID *model.PriorityID) (bool, error) {
	result, err := d.conn.ExecContext(ctx, deletePriorityQuery, PriorityID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
