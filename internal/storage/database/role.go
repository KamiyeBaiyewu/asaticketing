package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apiErr  "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
)

var (
	
)

// RoleDB - Interface holds all the methods for storing and retrieving roles
type RoleDB interface {
	CreateRole(ctx context.Context, role *model.Role) error
	GetRoleByID(ctx context.Context, roleID *model.RoleID) (*model.Role, error)

	UpdateRole(ctx context.Context, role *model.Role) error
	DeleteRole(ctx context.Context, roleID *model.RoleID) (bool, error)
	ListAllRoles(ctx context.Context) ([]*model.Role, error)


}

const createRoleQuery = `
		INSERT INTO roles (
			 name, description
			)
			VALUES (
				 :name, :description
				)
				RETURNING role_id`

func (d *database) CreateRole(ctx context.Context, userRole *model.Role) (err error) {

	rows, err := d.conn.NamedQueryContext(ctx, createRoleQuery, userRole)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Constraint == "unique_role" {
				err = apiErr.ErrRoleExists
				return
			}
			if pqError.Code.Name() == "unique_violation" && pqError.Constraint == "roles_pkey" {
				err = apiErr.ErrRoleIDExists
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
	if err := rows.Scan(&userRole.ID); err != nil {
		err = errors.Wrap(err, "Could not get the Role ID")
	}
	return
}

const getRoleByIDQuery = `
	SELECT role_id, name, description, created_at, updated_at, deleted_at
	from roles
	WHERE role_id = $1`

func (d *database) GetRoleByID(ctx context.Context, roleID *model.RoleID) (*model.Role, error) {
	userRole := model.Role{}
	if err := d.conn.GetContext(ctx, &userRole, getRoleByIDQuery, roleID); err != nil {

		return nil, err
	}
	return &userRole, nil
}

const updateRoleQuery = `
		update roles
		SET name = :name,
		description = :description,
		updated_at = NOW()
		WHERE role_id = :role_id`

func (d *database) UpdateRole(ctx context.Context, role *model.Role) error {
	result, err := d.conn.NamedExecContext(ctx, updateRoleQuery, role)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("Role Not found")
	}

	return nil
}

const deleteRoleQuery = `
	update roles
	SET deleted_at = NOW()
	WHERE role_id = $1 AND deleted_at is NULL;
	`

func (d *database) DeleteRole(ctx context.Context, roleID *model.RoleID) (bool, error) {

	result, err := d.conn.ExecContext(ctx, deleteRoleQuery, roleID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}
	return true, nil
}

const listAllRolesQuery = `
	SELECT role_id, name, description, created_at, updated_at, deleted_at
	from roles
	WHERE deleted_at is NULL;
`

func (d *database) ListAllRoles(ctx context.Context) ([]*model.Role, error) {
	userRoles := []*model.Role{}
	if err := d.conn.SelectContext(ctx, &userRoles, listAllRolesQuery); err != nil {
		return nil, errors.Wrap(err, "could not get roles")
	}
	return userRoles, nil
}
