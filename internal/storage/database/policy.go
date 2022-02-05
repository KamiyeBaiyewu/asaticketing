package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apiErr  "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
)

// PolicyDB - holds all the data for policies and actions
type PolicyDB interface {
	CreatePolicy(ctx context.Context, userPolicy *model.Policy) (err error)
	GetPolicyByID(ctx context.Context, policyID *model.PolicyID) (*model.Policy, error)
	ListAllPolicies(ctx context.Context) ([]*model.Policy, error)
	DeletePolicy(ctx context.Context, policyID *model.PolicyID) (bool, error)

	// MISC
	ListPermitedObjectActions(ctx context.Context, userID *model.UserID) ([]*model.Action, error)
	ListPermitedSystemActions(ctx context.Context, userID *model.UserID) ([]*model.Action, error)
}

const createPolicyQuery = `
		INSERT INTO object_policies (
			 role_id, object_id, action, created_by
			)
			VALUES (
				 :role_id, :object_id,  :action,  :created_by
				)
				RETURNING policy_id`

func (d *database) CreatePolicy(ctx context.Context, userPolicy *model.Policy) (err error) {
	rows, err := d.conn.NamedQueryContext(ctx, createPolicyQuery, userPolicy)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Constraint == "unique_policy" {
				err = apiErr.ErrPolicyExists
				return
			}
			switch pqError.Code.Name() {
			case "invalid_text_representation":
				return apiErr.ErrInvalidValues
			case "unique_violation":
				if pqError.Constraint == "object_policies_unique" {
					return apiErr.ErrPolicyExists
				}
			// One of the Foreign key ID is missing
			case "foreign_key_violation":
				switch pqError.Constraint {
				case "object_policy_role_id_fkey":
					return errors.New("Role does not exist")
				case "object_policy_object_id_fkey":
					return errors.New("Object does not exist")
				// case "object_policy_action_id_fkey":
				// 	return errors.New("Action does not exist")
				case "object_policy_created_by_fkey":
					return errors.New("User does not exist")
				}
			}

			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
			}).Info()
			return errors.New("Error Creating Policy")
		}
		logrus.WithError(err).Error("This error was not caught")
		return errors.New("Error Creating Policy")
	}

	rows.Next()
	if err := rows.Scan(&userPolicy.ID); err != nil {
		err = errors.Wrap(err, "Could not get the Policy ID")
	}
	return
}

const getPolicyByIDQuery = `
	SELECT  objp.policy_id , objp.role_id ,objp.object_id, objp.created_by, objp.is_standard, objp.created_at, objp.updated_at, objp.deleted_at 
	FROM object_policies objp 
	WHERE objp.policy_id = $1
	AND objp.deleted_at IS NULL
`

func (d *database) GetPolicyByID(ctx context.Context, policyID *model.PolicyID) (*model.Policy, error) {
	userPolicy := model.Policy{}
	if err := d.conn.GetContext(ctx, &userPolicy, getPolicyByIDQuery, policyID); err != nil {
		println(err.Error())
		return nil, err
	}
	return &userPolicy, nil
}

const listAllPoliciesQuery = `
	SELECT  objp.policy_id , objp.role_id ,objp.object_id, objp.created_by, objp.is_standard, objp.created_at, objp.updated_at, objp.deleted_at 
	FROM object_policies objp 
	WHERE objp.deleted_at IS NULL 
`

func (d *database) ListAllPolicies(ctx context.Context) ([]*model.Policy, error) {

	userPolicies := []*model.Policy{}
	if err := d.conn.SelectContext(ctx, &userPolicies, listAllPoliciesQuery); err != nil {
		if pqError, ok := err.(*pq.Error); ok {

			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
			}).Info()
		}
		return nil, err //errors.New("No policies found")
	}
	return userPolicies, nil
}

const deletePolicyQuery = `
	UPDATE object_policies
	SET deleted_at = NOW()
	WHERE policy_id = $1 AND deleted_at is NULL;
	`

func (d *database) DeletePolicy(ctx context.Context, policyID *model.PolicyID) (bool, error) {

	result, err := d.conn.ExecContext(ctx, deletePolicyQuery, policyID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}
	return true, nil
}

/*
	MISC
*/
const listPermitedObjectActionsQuery = `
	SELECT upper(concat(v1,'_',v2)) FROM casbin_rules cr 
	WHERE v0 IN 
	(SELECT r.role_id::text 
	from users u
	INNER JOIN roles r on r.role_id = u.role_id 
	WHERE u.user_id = $1 
	AND u.deleted_at is NULL
	AND r.deleted_at is NULL
	UNION 
	SELECT r.role_id::text 
	from roles r
	INNER JOIN users_roles ur 
	ON ur.role_id = r.role_id 
	WHERE ur.user_id = $1)
`

func (d *database) ListPermitedObjectActions(ctx context.Context, userID *model.UserID) ([]*model.Action, error) {

	PermitedActions := []*model.Action{}
	if err := d.conn.SelectContext(ctx, &PermitedActions, listPermitedObjectActionsQuery, userID); err != nil {
		return nil, errors.Wrap(err, "could not get policies")
	}
	return PermitedActions, nil
}


const listPermitedSystemActionsQuery = `
	SELECT upper(so."name") FROM system_objects so  
	INNER JOIN system_object_policies sop 
	ON sop.object_id = so.object_id 
	WHERE sop.role_id IN 
	(SELECT r.role_id 
	from users u
	INNER JOIN roles r on r.role_id = u.role_id 
	WHERE u.user_id = $1 
	AND u.deleted_at is NULL
	AND r.deleted_at is NULL
	UNION 
	SELECT r.role_id 
	from roles r
	INNER JOIN users_roles ur 
	ON ur.role_id = r.role_id 
	WHERE ur.user_id = $1)
`

func (d *database) ListPermitedSystemActions(ctx context.Context, userID *model.UserID) ([]*model.Action, error) {

	PermitedActions := []*model.Action{}
	if err := d.conn.SelectContext(ctx, &PermitedActions, listPermitedSystemActionsQuery, userID); err != nil {
		return nil, errors.Wrap(err, "could not get policies")
	}
	return PermitedActions, nil
}

