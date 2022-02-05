package database

import (
	"context"
	"strings"

	"github.com/lib/pq"
	apiErr "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// UserRoleDB - interfaces holds methods for users the relation between users and roles
type UserRoleDB interface {
	GrantRole(ctx context.Context, userID *model.UserID, roleID *model.RoleID) error
	RevokeRole(ctx context.Context, userID *model.UserID, roleID *model.RoleID) error
	GetRoleByUser(ctx context.Context, userID *model.UserID) ([]*model.Role, error)
	GetRoleUsers(ctx context.Context, roleID *model.RoleID) ([]*model.User, error)

	IsPrimaryRole(ctx context.Context, userID *model.UserID, roleID *model.RoleID) (bool, error)
	GetRoleNamesForUser(ctx context.Context, userID *model.UserID) (string, error)
}

const grantUserRoleQuery = `
	INSERT INTO users_roles (user_id, role_id)
	VALUES ($1, $2)`

func (d *database) GrantRole(ctx context.Context, userID *model.UserID, roleID *model.RoleID) error {
	//TODO: only grant role is the user does not already have have the grant
	if _, err := d.conn.ExecContext(ctx, grantUserRoleQuery, userID, roleID); err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" && pqError.Constraint == "users_roles_user" {
				return apiErr.ErrRoleAllreadyGranted
			}
			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
			}).Info()
		}
		return errors.Wrap(err, "could not grant user role")
	}
	return nil
}

const revokeUserRoleQuery = `
	UPDATE users_roles
	SET deleted_at = NOW()
	WHERE user_id = $1 
	AND role_id = $2
	AND deleted_at is NULL`

func (d *database) RevokeRole(ctx context.Context, userID *model.UserID, roleID *model.RoleID) error {
	result, err := d.conn.ExecContext(ctx, revokeUserRoleQuery, userID, roleID)
	if err != nil {
		return errors.Wrap(err, "could not revoke user role")
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("No matching user and role found")
	}
	return nil
}

const getRolesByUserQuery = `
	SELECT ro.name  ,ro.role_id 
	FROM users_roles ur 
	INNER JOIN roles ro 
	ON ro.role_id = ur.role_id
	WHERE ur.user_id = $1
	AND ur.deleted_at is NULL
	AND ro.deleted_at is NULL
	UNION 
	SELECT ro.name ,ro.role_id 
	from roles ro
	INNER JOIN users us 
	ON us.role_id = ro.role_id 
	WHERE us.user_id = $1
`

func (d *database) GetRoleByUser(ctx context.Context, userID *model.UserID) ([]*model.Role, error) {
	// TODO: Also add the user's primary role
	var roles []*model.Role
	if err := d.conn.SelectContext(ctx, &roles, getRolesByUserQuery, userID); err != nil {
		return nil, errors.Wrap(err, "could not get the roles for user")
	}
	return roles, nil
}

const getRolesUsersQuery = `
	SELECT us.user_id,us.firstname, us.lastname, us.email 
	FROM users_roles ur 
	INNER JOIN users us ON us.user_id = ur.user_id
	INNER JOIN roles ro ON ro.role_id = ur.role_id
	WHERE ur.role_id = $1 
	AND ur.deleted_at is NULL
	AND us.deleted_at is NULL
	AND ro.deleted_at is NULL
	UNION
	SELECT us.user_id,us.firstname, us.lastname, us.email
	from users us 
	WHERE us.role_id = $1`

func (d *database) GetRoleUsers(ctx context.Context, roleID *model.RoleID) ([]*model.User, error) {
	var users []*model.User
	if err := d.conn.SelectContext(ctx, &users, getRolesUsersQuery, roleID); err != nil {
		return nil, errors.Wrap(err, "could not get the users for a role")
	}
	return users, nil
}

/*
const getRolesByUserIDQuery = `
	SELECT role
	FROM users_roles
	WHERE user_id = $1`

 func (d *database) GetRoleByUser(ctx context.Context, userID model.UserID) ([]*model.Role, error) {
	var roles []*model.Role
	if err := d.conn.SelectContext(ctx, &roles, getRolesByUserIDQuery, userID); err != nil {
		return nil, errors.Wrap(err, "could not get user roles")
	}
	return roles, nil
} */

const isPrimaryRoleQuery = `select exists(
	select 1 from users  
	where 
	user_id = $1 AND role_id = $2
	AND deleted_at is NULL
	)`

func (d *database) IsPrimaryRole(ctx context.Context, userID *model.UserID, roleID *model.RoleID) (bool, error) {

	var exist bool
	if err := d.conn.GetContext(ctx, &exist, isPrimaryRoleQuery, userID, roleID); err != nil {
		return false, err
	}

	return exist, nil
}

const getRoleNamesForUserQuery = `
	SELECT ro.role_id
	from users us
	INNER JOIN roles ro on ro.role_id = us.role_id
	WHERE us.user_id = $1 
	AND us.deleted_at is NULL
	AND ro.deleted_at is NULL
	UNION 
	SELECT ro.role_id
	from roles ro
	INNER JOIN users_roles ur 
	ON ur.role_id = ro.role_id 
	WHERE ur.user_id = $1
	AND ur.deleted_at is NULL
	AND ro.deleted_at is NULL
	`

func (d *database) GetRoleNamesForUser(ctx context.Context, userID *model.UserID) (string, error) {
	var userRoles []string
	var userRole string
	if err := d.conn.SelectContext(ctx, &userRoles, getRoleNamesForUserQuery, userID); err != nil {

		return "", err
	}
	if len(userRoles) > 0 {

		if len(userRoles) > 1 {

			userRole = strings.Join(userRoles, "|")
		} else {
			userRole = userRoles[0]
		}
	}
	return userRole, nil

}
