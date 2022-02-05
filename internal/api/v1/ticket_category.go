package v1

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/middlewares"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/responses"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/utils"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/env"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/storage/database"
	"github.com/sirupsen/logrus"
)

// TicketCategoryAPI - structure holds handlers for Tickets Categories
type TicketCategoryAPI struct {
	db  database.Database
	env *env.Env
}

// Load help create a subrouter for the categories
func loadTicketCategory(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	api := &TicketCategoryAPI{env: env, db: env.DB}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/ticket_categories", api.Create, authorizer.ObjAuthorize("ticket_category", "create")),
		newAPIEndpoint("GET", "/ticket_categories/{categoryID}", api.Get, authorizer.ObjAuthorize("ticket_category", "view")), //retrieves a category using its ID
		newAPIEndpoint("GET", "/ticket_categories", api.List, authorizer.ObjAuthorize("ticket_category", "list")),             //retrieves all the categories

		newAPIEndpoint("PATCH", "/ticket_categories/{categoryID}", api.Update, authorizer.ObjAuthorize("ticket_category", "update")),  //updates a category using its ID
		newAPIEndpoint("DELETE", "/ticket_categories/{categoryID}", api.Delete, authorizer.ObjAuthorize("ticket_category", "delete")), //delete a category using its ID

	}

	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new Category
func (api *TicketCategoryAPI) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_category.go -> CategoryApi.Create()")

	//Load parameters
	var category model.Category

	if err := category.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}



	if err := category.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"category": *category.Name,
	})
	err := api.db.CreateCategory(ctx, &category)
	if err != nil {
		logger.WithError(err).Error("Creating category")
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	createdCategory, err := api.db.GetCategoryByID(ctx, &category.ID)
	if err != nil {
		logger.WithError(err).Warn(err.Error())
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, &createdCategory)
}

// Get -  retreives category information
func (api *TicketCategoryAPI) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//principal := middlewares.GetPrincipal(r)
	logger := logrus.WithField("func", "ticket_category.go -> CategoryApi.Get()")

	vars := mux.Vars(r)
	categoryID := model.CategoryID(vars["categoryID"])

	//Get all the category fields from the database to ensure all is well
	category, err := api.db.GetCategoryByID(ctx, &categoryID)
	if err != nil {
		errMessage := fmt.Sprintf("Retrieving category ID: %v", categoryID)
		logger.WithError(err).Error(errMessage)
		utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	logger.WithField("categoryID", categoryID).Debug("Get Category Complete")

	utils.WriteJSON(w, http.StatusOK, category)
}

// Update - Updated Category  Details
// PATCH - /categories/{categoryID}
// Permission MemberIsTarget, Admin
func (api *TicketCategoryAPI) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_category.go -> Update()")

	vars := mux.Vars(r)
	categoryID := model.CategoryID(vars["categoryID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"CategoryID": categoryID,
		"pricipal":   principal,
	})

	var category model.Category

	// Decode Parameters
	if err := category.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"category": *category.Name,
	})

	logger = logger.WithField("categoryID", categoryID)

	storedCategory, err := api.db.GetCategoryByID(ctx, &categoryID)
	if err != nil {
		errMessage := fmt.Sprintf("Retrieving category ID: %v", categoryID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	storedCategory.UpdateValues(&category)
	// now update the database values
	err = api.db.UpdateCategory(ctx, storedCategory)
	if err != nil {
		logger.WithError(err).Warn("Error updating category.")
		utils.WriteError(w, http.StatusInternalServerError, err, map[string]string{
			"requestID": ctx.Value("correlationid").(string),
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"Category Updated": true,
	})

	utils.WriteJSON(w, http.StatusOK, storedCategory)

}

// List - List all the categories
// GET - /categories
// Permission Admin
func (api *TicketCategoryAPI) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_category.go -> List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	categories, err := api.db.ListAllCategories(ctx)
	if err != nil {
		logger.WithError(err).Error("Retreiving all ticket categories")
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return

	}

	logger.Info("Categories Returned")

	utils.WriteJSON(w, http.StatusOK, &categories)

}

// Delete - Deletes a category
// DELETE - /categories/{categoryID}/categories/{categoryID}
// Permission MemberIsTarget
func (api *TicketCategoryAPI) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_category.go -> Delete()")

	vars := mux.Vars(r)
	categoryID := model.CategoryID(vars["categoryID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"CategoryID": categoryID,
		"pricipal":   principal,
	})

	deleted, err := api.db.DeleteCategory(ctx, &categoryID)
	if err != nil {

		logger.WithError(err).Error("deleteing category")
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	logger.Info("Category Deleted")

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}
