package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	apiErr  "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
	"github.com/sirupsen/logrus"
)

// UserDB - interface that dabase connection must implement
type UserDB interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, userID *model.UserID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	ListAllUsers(ctx context.Context) ([]*model.User, error)
	DeleteUser(ctx context.Context, UserID *model.UserID) (bool, error)
}

var (
	// ErrUserExists - User already exists in the database
	ErrUserExists = errors.New("user with that email exists")
)

const createUserQuery = `
	INSERT INTO users (
		firstname, lastname, email, password_hash, user_type, role_id
		)
		VALUES (
			:firstname, :lastname, :email, :password_hash, :user_type, :role_id
			)
			RETURNING user_id`

func (d *database) CreateUser(ctx context.Context, user *model.User) (err error) {

	rows, err := d.conn.NamedQueryContext(ctx, createUserQuery, user)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "not_null_violation":
				switch pqError.Column {
				case "user_type":
					return apiErr.ErrUserTypeRequired
				}
			case "unique_violation":
				switch pqError.Constraint{
				case "user_email":
					return apiErr.ErrEmailAlreadyExists
				}
			// One of the Foreign key ID is missing
			case "foreign_key_violation":
				switch pqError.Constraint {
				case "users_role_id_fkey":
					return apiErr.ErrRoleNotExist
			}
		}
			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
			}).Info()
		}
		logrus.WithFields(logrus.Fields{
			"Erros":err.Error(),
		})
		return apiErr.ErrCreatingUser
	}

	rows.Next()
	if err := rows.Scan(&user.ID); err != nil {
		err = errors.Wrap(err, "Could not get the User ID")
	}
	return
}

const getUserByIDQuery = `
	SELECT user_id, firstname, lastname, email,role_id, password_hash, user_type, is_active, is_system, created_at, updated_at, deleted_at
	FROM users
	WHERE user_id = $1`

func (d *database) GetUserByID(ctx context.Context, userID *model.UserID) (*model.User, error) {
	user := model.User{}
	if err := d.conn.GetContext(ctx, &user, getUserByIDQuery, userID); err != nil {
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

const getUserByEmailQuery = `
	SELECT user_id, firstname, lastname, email,role_id, password_hash, user_type, is_active, is_system, created_at, updated_at, deleted_at
	FROM users
	WHERE email = $1 AND deleted_at is NULL`

func (d *database) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := model.User{}
	if err := d.conn.GetContext(ctx, &user, getUserByEmailQuery, email); err != nil {
		// logrus.WithError(err)
		return nil, err
	}
	return &user, nil

}

const updateUserQuery = `
		UPDATE users
		SET user_id = :user_id,
		 firstname = :firstname,
		 lastname = :lastname,
		 email = :email,
		 password_hash = :password_hash,
		 role_id = :role_id,
		 user_type = :user_type,
		 is_active = :is_active,
		updated_at = NOW()
		WHERE user_id = :user_id;
`

func (d *database) UpdateUser(ctx context.Context, user *model.User) error {

	result, err := d.conn.NamedExecContext(ctx, updateUserQuery, user)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("User Not found")
	}
	return nil
}

const listAllUsersQuery = `
	SELECT  user_id, firstname, lastname, email,role_id, password_hash, user_type, is_active, is_system, created_at, updated_at, deleted_at
	FROM users
	WHERE deleted_at is NULL;
`

func (d *database) ListAllUsers(ctx context.Context) ([]*model.User, error) {
	categories := []*model.User{}
	if err := d.conn.SelectContext(ctx, &categories, listAllUsersQuery); err != nil {
		return nil, errors.Wrap(err, "could not get users")
	}
	return categories, nil
}

const deleteUserQuery = `
	UPDATE users
	SET deleted_at = NOW()
	WHERE user_id = $1 AND deleted_at is NULL;
	`

func (d *database) DeleteUser(ctx context.Context, UserID *model.UserID) (bool, error) {
	result, err := d.conn.ExecContext(ctx, deleteUserQuery, UserID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
