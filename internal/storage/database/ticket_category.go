package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apiErr "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
)

// TicketCategoryDB - interface that dabase connection must implement
type TicketCategoryDB interface {
	CreateCategory(ctx context.Context, category *model.Category) error
	GetCategoryByID(ctx context.Context, categoryID *model.CategoryID) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) error
	ListAllCategories(ctx context.Context) ([]*model.Category, error)
	DeleteCategory(ctx context.Context, CategoryID *model.CategoryID) (bool, error)
}


const createCategoryQuery = `
	INSERT INTO ticket_categories (
	 name, description,weight
	)
	VALUES (
		 :name, :description, :weight
		)
		RETURNING category_id`

func (d *database) CreateCategory(ctx context.Context, category *model.Category) (err error) {

	rows, err := d.conn.NamedQueryContext(ctx, createCategoryQuery, category)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" && pqError.Constraint == "ticket_category_name" {
				err = apiErr.ErrCategoryExists
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
	if err := rows.Scan(&category.ID); err != nil {
		err = errors.Wrap(err, "Could not get the Category ID")
	}
	return
}

const getCategoryByIDQuery = `
	SELECT category_id, name, description, weight, created_at, updated_at, deleted_at
	FROM ticket_categories
	WHERE category_id = $1 
	AND deleted_at is NULL`

func (d *database) GetCategoryByID(ctx context.Context, categoryID *model.CategoryID) (*model.Category, error) {

	category := model.Category{}
	if err := d.conn.GetContext(ctx, &category, getCategoryByIDQuery, categoryID); err != nil {
		return nil, apiErr.ErrNotFound
	}
	return &category, nil

}

const updateCategoryQuery = `
	UPDATE ticket_categories
	SET 
		name = :name,
		description = :description,
		weight = :weight,
		updated_at = NOW()
	WHERE category_id = :category_id
	AND deleted_at is NULL`

func (d *database) UpdateCategory(ctx context.Context, category *model.Category) error {

	//println(*category.PasswordHash)
	result, err := d.conn.NamedExecContext(ctx, updateCategoryQuery, category)
	if err != nil {
	

		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("Category Not found")
	}
	return nil
}

const listAllCategoriesQuery = `
	SELECT category_id, name, description, weight, created_at, updated_at, deleted_at
	FROM ticket_categories
	WHERE deleted_at is NULL
	ORDER BY weight ASC`

func (d *database) ListAllCategories(ctx context.Context) ([]*model.Category, error) {

	categories := []*model.Category{}
	if err := d.conn.SelectContext(ctx, &categories, listAllCategoriesQuery); err != nil {
		return nil, errors.Wrap(err, "could not get categorys")
	}
	return categories, nil
}

const deleteCategoryQuery = `
	UPDATE ticket_categories
	SET deleted_at = NOW()
	WHERE category_id = $1 AND deleted_at is NULL`

func (d *database) DeleteCategory(ctx context.Context, CategoryID *model.CategoryID) (bool, error) {

	result, err := d.conn.ExecContext(ctx, deleteCategoryQuery, CategoryID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
