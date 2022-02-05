package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apiErr "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
)

// TicketSourceDB - interface that dabase connection must implement
type TicketSourceDB interface {
	CreateSource(ctx context.Context, source *model.Source) error
	GetSourceByID(ctx context.Context, sourceID *model.SourceID) (*model.Source, error)
	UpdateSource(ctx context.Context, source *model.Source) error
	ListAllSources(ctx context.Context) ([]*model.Source, error)
	DeleteSource(ctx context.Context, SourceID *model.SourceID) (bool, error)
}



const createSourceQuery = `INSERT INTO ticket_sources (
	 name, weight
	)
	VALUES (
		 :name, :weight
		)
		RETURNING source_id`

func (d *database) CreateSource(ctx context.Context, source *model.Source) (err error) {
	rows, err := d.conn.NamedQueryContext(ctx, createSourceQuery, source)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" && pqError.Constraint == "ticket_sources_name" {
				err = apiErr.ErrSourceExists
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
	if err := rows.Scan(&source.ID); err != nil {
		err = errors.Wrap(err, "Could not get the Source ID")
	}
	return
}

const getSourceByIDQuery = `
	SELECT source_id, name, weight, created_at, updated_at, deleted_at
	FROM ticket_sources
	WHERE source_id = $1 
	AND deleted_at is NULL`

func (d *database) GetSourceByID(ctx context.Context, sourceID *model.SourceID) (*model.Source, error) {
	source := model.Source{}
	if err := d.conn.GetContext(ctx, &source, getSourceByIDQuery, sourceID); err != nil {

		return nil, err
	}
	return &source, nil

}

const updateSourceQuery = `
	UPDATE ticket_sources
	SET 
		source_id = :source_id,
		name = :name,
		weight = :weight,
		updated_at = NOW()
	WHERE source_id = :source_id 
	AND deleted_at is NULL`

func (d *database) UpdateSource(ctx context.Context, source *model.Source) error {

	result, err := d.conn.NamedExecContext(ctx, updateSourceQuery, source)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("Source Not found")
	}
	return nil
}

const listAllSourcesQuery = `
	SELECT source_id, name, weight, created_at, updated_at, deleted_at
	FROM ticket_sources
	WHERE deleted_at is NULL
	ORDER BY weight ASC`

func (d *database) ListAllSources(ctx context.Context) ([]*model.Source, error) {

	categories := []*model.Source{}
	if err := d.conn.SelectContext(ctx, &categories, listAllSourcesQuery); err != nil {
		return nil, errors.Wrap(err, "could not get sources")
	}
	return categories, nil
}

const deleteSourceQuery = `
	UPDATE ticket_sources
	SET deleted_at = NOW()
	WHERE source_id = $1 
	AND deleted_at is NULL`

func (d *database) DeleteSource(ctx context.Context, SourceID *model.SourceID) (bool, error) {

	result, err := d.conn.ExecContext(ctx, deleteSourceQuery, SourceID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
