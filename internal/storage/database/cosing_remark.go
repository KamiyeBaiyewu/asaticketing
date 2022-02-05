package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apiErr "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
)

// ClosingRemarkDB - interface that dabase connection must implement
type ClosingRemarkDB interface {
	CreateClosingRemark(ctx context.Context, closingRemark *model.ClosingRemark) error
	GetClosingRemarkByID(ctx context.Context, closingRemarkID *model.ClosingRemarkID) (*model.ClosingRemark, error)
	UpdateClosingRemark(ctx context.Context, closingRemark *model.ClosingRemark) error
	ListAllRemarks(ctx context.Context) ([]*model.ClosingRemark, error)
	DeleteClosingRemark(ctx context.Context, ClosingRemarkID *model.ClosingRemarkID) (bool, error)
}


const createClosingRemarkQuery = `
	INSERT INTO ticket_closing_remarks (
		ticket_id, cause_id, closed_by, remark
	)
	VALUES (
		:ticket_id, :cause_id, :closed_by, :remark
		)
		RETURNING remark_id`

func (d *database) CreateClosingRemark(ctx context.Context, closingRemark *model.ClosingRemark) (err error) {

	rows, err := d.conn.NamedQueryContext(ctx, createClosingRemarkQuery, closingRemark)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" && pqError.Constraint == "ticket_closing_remarks_unique" {
				err = apiErr.ErrClosingRemarkExists
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
	if err := rows.Scan(&closingRemark.ID); err != nil {
		err = errors.Wrap(err, "Could not get the ClosingRemark ID")
	}
	return
}

const getClosingRemarkByIDQuery = `
	SELECT remark_id, ticket_id, cause_id, closed_by, remark, created_at, updated_at, deleted_at
	FROM ticket_closing_remarks
	WHERE remark_id = $1 
	AND deleted_at is NULL`

func (d *database) GetClosingRemarkByID(ctx context.Context, closingRemarkID *model.ClosingRemarkID) (*model.ClosingRemark, error) {

	closingRemark := model.ClosingRemark{}
	if err := d.conn.GetContext(ctx, &closingRemark, getClosingRemarkByIDQuery, closingRemarkID); err != nil {
		return nil, apiErr.ErrNotFound
	}
	return &closingRemark, nil

}

const updateClosingRemarkQuery = `
	UPDATE ticket_closing_remarks
	SET 
	 cause_id = :cause_id,
	 remark = :remark,
	updated_at = NOW()
	WHERE remark_id = :remark_id
	AND deleted_at is NULL`

func (d *database) UpdateClosingRemark(ctx context.Context, closingRemark *model.ClosingRemark) error {

	//println(*closingRemark.PasswordHash)
	result, err := d.conn.NamedExecContext(ctx, updateClosingRemarkQuery, closingRemark)
	if err != nil {
	

		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("ClosingRemark Not found")
	}
	return nil
}

const listAllRemarksQuery = `
	SELECT remark_id, ticket_id, cause_id, closed_by, remark, created_at, updated_at, deleted_at
	FROM ticket_closing_remarks
	WHERE deleted_at is NULL`

func (d *database) ListAllRemarks(ctx context.Context) ([]*model.ClosingRemark, error) {

	categories := []*model.ClosingRemark{}
	if err := d.conn.SelectContext(ctx, &categories, listAllRemarksQuery); err != nil {
		return nil, errors.Wrap(err, "could not get closingRemarks")
	}
	return categories, nil
}

const deleteClosingRemarkQuery = `
	UPDATE ticket_closing_remarks
	SET deleted_at = NOW()
	WHERE remark_id = $1 AND deleted_at is NULL`

func (d *database) DeleteClosingRemark(ctx context.Context, ClosingRemarkID *model.ClosingRemarkID) (bool, error) {

	result, err := d.conn.ExecContext(ctx, deleteClosingRemarkQuery, ClosingRemarkID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
