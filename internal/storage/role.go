package storage

import (
	"context"
	"errors"

	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
)

var (
	// ErrRoleExists - User already exists in the database
	ErrRoleExists = errors.New("Role already exists")

	// ErrRoleIDExists - Role with the ID already exists
	ErrRoleIDExists = errors.New("Role with ID already exists")

	// ErrRoleAllreadyGranted - tells the user that the role has already been granted to the user
	ErrRoleAllreadyGranted = errors.New("The User has already been granted the role")
)

func (s *storage) CreateRole(ctx context.Context, role *model.Role) error {
	return s.db.CreateRole(ctx, role)
}
func (s *storage) GetRoleByID(ctx context.Context, roleID *model.RoleID) (*model.Role, error) {
	return s.db.GetRoleByID(ctx, roleID)
}

func (s *storage) UpdateRole(ctx context.Context, role *model.Role) error {
	return s.db.UpdateRole(ctx, role)
}
func (s *storage) DeleteRole(ctx context.Context, roleID *model.RoleID) (bool, error) {
	return s.db.DeleteRole(ctx, roleID)
}
func (s *storage) ListAllRoles(ctx context.Context) ([]*model.Role, error) {
	return s.db.ListAllRoles(ctx)
}

func (s *storage) GrantRole(ctx context.Context, userID *model.UserID, roleID *model.RoleID) error {
	return s.db.GrantRole(ctx, userID, roleID)
}
func (s *storage) RevokeRole(ctx context.Context, userID *model.UserID, roleID *model.RoleID) error {
	return s.db.RevokeRole(ctx, userID, roleID)
}
func (s *storage) GetRoleByUser(ctx context.Context, userID *model.UserID) ([]*model.Role, error) {
	return s.db.GetRoleByUser(ctx, userID)
}
func (s *storage) GetRoleUsers(ctx context.Context, roleID *model.RoleID) ([]*model.User, error) {
	return s.db.GetRoleUsers(ctx, roleID)
}

func (s *storage) IsPrimaryRole(ctx context.Context, userID *model.UserID, roleID *model.RoleID) (bool, error) {
	return s.db.IsPrimaryRole(ctx, userID, roleID)
}
