package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apiErr  "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
)

// ObjectDB - interface that dabase connection must implement
type ObjectDB interface {
	CreateObject(ctx context.Context, user *model.Object) error
	GetObjectByID(ctx context.Context, userID *model.ObjectID) (*model.Object, error)
	UpdateObject(ctx context.Context, user *model.Object) error
	ListAllObjects(ctx context.Context) ([]*model.Object, error)
	DeleteObject(ctx context.Context, ObjectID *model.ObjectID) (bool, error)
}


const createObjectQuery = `
	INSERT INTO objects (
	 name, description, created_by
	)
	VALUES (
		 :name, :description, :created_by
		)
		RETURNING object_id`

func (d *database) CreateObject(ctx context.Context, user *model.Object) (err error) {
	rows, err := d.conn.NamedQueryContext(ctx, createObjectQuery, user)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" && pqError.Constraint == "objects_name" {
				err = apiErr.ErrObjectExists
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
	if err := rows.Scan(&user.ID); err != nil {
		err = errors.Wrap(err, "Could not get the Object ID")
	}
	return
}

const getObjectByIDQuery = `
	SELECT object_id, name, description,is_standard,created_by,  created_at, updated_at, deleted_at
	FROM objects
	WHERE object_id = $1 
	AND deleted_at is NULL`

func (d *database) GetObjectByID(ctx context.Context, userID *model.ObjectID) (*model.Object, error) {

	user := model.Object{}
	if err := d.conn.GetContext(ctx, &user, getObjectByIDQuery, userID); err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
			}).Info()
		}

		return nil, err
	}
	return &user, nil

}

const updateObjectQuery = `
	UPDATE objects
	SET 
		name = :name,
		description = :description,
		updated_at = NOW()
	WHERE object_id = :object_id
	AND deleted_at is NULL`

func (d *database) UpdateObject(ctx context.Context, user *model.Object) error {

	
	result, err := d.conn.NamedExecContext(ctx, updateObjectQuery, user)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("Object Not found")
	}
	return nil
}

const listAllObjectsQuery = `
	SELECT object_id, name, description,is_standard, created_by, created_at, updated_at, deleted_at
	FROM objects
	WHERE deleted_at is NULL`

func (d *database) ListAllObjects(ctx context.Context) ([]*model.Object, error) {

	object := []*model.Object{}
	if err := d.conn.SelectContext(ctx, &object, listAllObjectsQuery); err != nil {
		return nil, errors.Wrap(err, "could not get users")
	}
	return object, nil
}

const deleteObjectQuery = `
	UPDATE objects
	SET deleted_at = NOW()
	WHERE object_id = $1 AND deleted_at is NULL`

func (d *database) DeleteObject(ctx context.Context, ObjectID *model.ObjectID) (bool, error) {

	result, err := d.conn.ExecContext(ctx, deleteObjectQuery, ObjectID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
