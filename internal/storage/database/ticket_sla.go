package database

import (
	"context"

	"github.com/lib/pq"
	apiErr "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// SLADB - interface that dabase connection must implement
type SLADB interface {
	CreateSLA(ctx context.Context, sla *model.SLA) error
	GetSLAByID(ctx context.Context, slaID *model.SLAID) (*model.SLA, error)
	UpdateSLA(ctx context.Context, sla *model.SLA) error
	ListAllSLA(ctx context.Context) ([]*model.SLA, error)
	DeleteSLA(ctx context.Context, SLAID *model.SLAID) (bool, error)
}

const createSLAQuery = `INSERT INTO ticket_slas (
	 name, weight,grace_period
	)
	VALUES (
		 :name, :weight, :grace_period
		)
		RETURNING agreement_id`

func (d *database) CreateSLA(ctx context.Context, sla *model.SLA) (err error) {

	rows, err := d.conn.NamedQueryContext(ctx, createSLAQuery, sla)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" && pqError.Constraint == "ticket_slas_name" {
				err = apiErr.ErrSLAExists
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
	if err := rows.Scan(&sla.ID); err != nil {
		err = errors.Wrap(err, "Could not get the SLA ID")
	}
	return
}

const getSLAByIDQuery = `
	SELECT agreement_id, name, weight,grace_period, created_at, updated_at, deleted_at
	FROM ticket_slas
	WHERE agreement_id = $1 
	AND deleted_at is NULL`

func (d *database) GetSLAByID(ctx context.Context, slaID *model.SLAID) (*model.SLA, error) {

	sla := model.SLA{}
	if err := d.conn.GetContext(ctx, &sla, getSLAByIDQuery, slaID); err != nil {

		return nil, apiErr.ErrNotFound
	}
	return &sla, nil

}

const updateSLAQuery = `
	UPDATE ticket_slas
	SET 
		agreement_id = :agreement_id,
		name = :name,
		grace_period = :grace_period,
		weight = :weight,
		updated_at = NOW()
	WHERE agreement_id = :agreement_id 
	AND deleted_at is NULL`

func (d *database) UpdateSLA(ctx context.Context, sla *model.SLA) error {

	result, err := d.conn.NamedExecContext(ctx, updateSLAQuery, sla)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("SLA Not found")
	}
	return nil
}

const listAllSLAQuery = `
	SELECT agreement_id, name,grace_period, weight, created_at, updated_at, deleted_at
	FROM ticket_slas
	WHERE deleted_at is NULL
	ORDER BY weight ASC`

func (d *database) ListAllSLA(ctx context.Context) ([]*model.SLA, error) {

	categories := []*model.SLA{}
	if err := d.conn.SelectContext(ctx, &categories, listAllSLAQuery); err != nil {
		return nil, errors.Wrap(err, "could not get slas")
	}
	return categories, nil
}

const deleteSLAQuery = `
	UPDATE ticket_slas
	SET deleted_at = NOW()
	WHERE agreement_id = $1 
	AND deleted_at is NULL`

func (d *database) DeleteSLA(ctx context.Context, SLAID *model.SLAID) (bool, error) {

	result, err := d.conn.ExecContext(ctx, deleteSLAQuery, SLAID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
