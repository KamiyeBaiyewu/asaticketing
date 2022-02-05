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

// ClosingRemarkAPI - structure holds handlers for Tickets Categories
type ClosingRemarkAPI struct {
	db  database.Database
	env *env.Env
}

// Load help create a subrouter for the closingRemarks
func loadClosingRemark(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	api := &ClosingRemarkAPI{env: env, db: env.DB}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/closing_remarks", api.Create, authorizer.ObjAuthorize("closing_remark", "create")),
		newAPIEndpoint("GET", "/closing_remarks/{closingRemarkID}", api.Get, authorizer.ObjAuthorize("closing_remark", "view")), //retrieves a closingRemark using its ID
		newAPIEndpoint("GET", "/closing_remarks", api.List, authorizer.ObjAuthorize("closing_remark", "list")),             //retrieves all the closingRemarks

		newAPIEndpoint("PATCH", "/closing_remarks/{closingRemarkID}", api.Update, authorizer.ObjAuthorize("closing_remark", "update")),  //updates a closingRemark using its ID
		newAPIEndpoint("DELETE", "/closing_remarks/{closingRemarkID}", api.Delete, authorizer.ObjAuthorize("closing_remark", "delete")), //delete a closingRemark using its ID

	}

	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new ClosingRemark
func (api *ClosingRemarkAPI) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	principal := middlewares.GetPrincipal(r)
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "closing_remark.go -> ClosingRemarkApi.Create()")

	//Load parameters
	var closingRemark model.ClosingRemark

	if err := closingRemark.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}


	closingRemark.UserID = principal.UserID
	if err := closingRemark.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"closingRemark": *closingRemark.Remark,
	})
	err := api.db.CreateClosingRemark(ctx, &closingRemark)
	if err != nil {
		logger.WithError(err).Error("Creating closingRemark")
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	createdClosingRemark, err := api.db.GetClosingRemarkByID(ctx, &closingRemark.ID)
	if err != nil {
		logger.WithError(err).Warn(err.Error())
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, &createdClosingRemark)
}

// Get -  retreives closingRemark information
func (api *ClosingRemarkAPI) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//principal := middlewares.GetPrincipal(r)
	logger := logrus.WithField("func", "closing_remark.go -> ClosingRemarkApi.Get()")

	vars := mux.Vars(r)
	closingRemarkID := model.ClosingRemarkID(vars["closingRemarkID"])

	//Get all the closingRemark fields from the database to ensure all is well
	closingRemark, err := api.db.GetClosingRemarkByID(ctx, &closingRemarkID)
	if err != nil {
		errMessage := fmt.Sprintf("Retrieving closingRemark ID: %v", closingRemarkID)
		logger.WithError(err).Error(errMessage)
		utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	logger.WithField("closingRemarkID", closingRemarkID).Debug("Get ClosingRemark Complete")

	utils.WriteJSON(w, http.StatusOK, closingRemark)
}

// Update - Updated ClosingRemark  Details
// PATCH - /closingRemarks/{closingRemarkID}
// Permission MemberIsTarget, Admin
func (api *ClosingRemarkAPI) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "closing_remark.go -> Update()")

	vars := mux.Vars(r)
	closingRemarkID := model.ClosingRemarkID(vars["closingRemarkID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"ClosingRemarkID": closingRemarkID,
		"pricipal":   principal,
	})

	var closingRemark model.ClosingRemark

	// Decode Parameters
	if err := closingRemark.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"closingRemark": *closingRemark.Remark,
	})

	logger = logger.WithField("closingRemarkID", closingRemarkID)

	storedClosingRemark, err := api.db.GetClosingRemarkByID(ctx, &closingRemarkID)
	if err != nil {
		errMessage := fmt.Sprintf("Retrieving closingRemark ID: %v", closingRemarkID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	storedClosingRemark.UpdateValues(&closingRemark)
	// now update the database values
	err = api.db.UpdateClosingRemark(ctx, storedClosingRemark)
	if err != nil {
		logger.WithError(err).Warn("Error updating closingRemark.")
		utils.WriteError(w, http.StatusInternalServerError, err, map[string]string{
			"requestID": ctx.Value("correlationid").(string),
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"ClosingRemark Updated": true,
	})

	utils.WriteJSON(w, http.StatusOK, storedClosingRemark)

}

// List - List all the closingRemarks
// GET - /closingRemarks
// Permission Admin
func (api *ClosingRemarkAPI) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "closing_remark.go -> List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	closingRemarks, err := api.db.ListAllRemarks(ctx)
	if err != nil {
		logger.WithError(err).Error("Retreiving all ticket closingRemarks")
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return

	}

	logger.Info("Categories Returned")

	utils.WriteJSON(w, http.StatusOK, &closingRemarks)

}

// Delete - Deletes a closingRemark
// DELETE - /closingRemarks/{closingRemarkID}/closingRemarks/{closingRemarkID}
// Permission MemberIsTarget
func (api *ClosingRemarkAPI) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "closing_remark.go -> Delete()")

	vars := mux.Vars(r)
	closingRemarkID := model.ClosingRemarkID(vars["closingRemarkID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"ClosingRemarkID": closingRemarkID,
		"pricipal":   principal,
	})

	deleted, err := api.db.DeleteClosingRemark(ctx, &closingRemarkID)
	if err != nil {

		logger.WithError(err).Error("deleteing closingRemark")
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	logger.Info("ClosingRemark Deleted")

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}
