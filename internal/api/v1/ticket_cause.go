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

// TicketCauseAPI - structure holds handlers for Tickets Causes
type TicketCauseAPI struct {
	db  database.Database
	env *env.Env
}

// Load help create a subrouter for the causes
func loadTicketCause(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	api := &TicketCauseAPI{env: env, db: env.DB}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/ticket_causes", api.Create, authorizer.ObjAuthorize("ticket_cause", "create")),
		newAPIEndpoint("GET", "/ticket_causes/{causeID}", api.Get, authorizer.ObjAuthorize("ticket_cause", "view")), //retrieves a cause using its ID
		newAPIEndpoint("GET", "/ticket_causes", api.List, authorizer.ObjAuthorize("ticket_cause", "list")),             //retrieves all the causes

		newAPIEndpoint("PATCH", "/ticket_causes/{causeID}", api.Update, authorizer.ObjAuthorize("ticket_cause", "update")),  //updates a cause using its ID
		newAPIEndpoint("DELETE", "/ticket_causes/{causeID}", api.Delete, authorizer.ObjAuthorize("ticket_cause", "delete")), //delete a cause using its ID

	}

	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new Cause
func (api *TicketCauseAPI) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	principal := middlewares.GetPrincipal(r)

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_cause.go -> CauseApi.Create()")

	//Load parameters
	var cause model.Cause

	if err := cause.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := cause.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}
	
	// Add the userID from the principal
	cause.UserID = principal.UserID

	logger = logger.WithFields(logrus.Fields{
		"cause": *cause.Name,
	})
	err := api.db.CreateCause(ctx, &cause)
	if err != nil {
		logger.WithError(err).Error("Creating cause")
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	createdCause, err := api.db.GetCauseByID(ctx, &cause.ID)
	if err != nil {
		logger.WithError(err).Warn(err.Error())
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, &createdCause)
}

// Get -  retreives cause information
func (api *TicketCauseAPI) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//principal := middlewares.GetPrincipal(r)
	logger := logrus.WithField("func", "ticket_cause.go -> CauseApi.Get()")

	vars := mux.Vars(r)
	causeID := model.CauseID(vars["causeID"])

	//Get all the cause fields from the database to ensure all is well
	cause, err := api.db.GetCauseByID(ctx, &causeID)
	if err != nil {
		errMessage := fmt.Sprintf("Retrieving cause ID: %v", causeID)
		logger.WithError(err).Error(errMessage)
		utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	logger.WithField("causeID", causeID).Debug("Get Cause Complete")

	utils.WriteJSON(w, http.StatusOK, cause)
}

// Update - Updated Cause  Details
// PATCH - /causes/{causeID}
// Permission MemberIsTarget, Admin
func (api *TicketCauseAPI) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_cause.go -> Update()")

	vars := mux.Vars(r)
	causeID := model.CauseID(vars["causeID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"CauseID": causeID,
		"pricipal":   principal,
	})

	var cause model.Cause

	// Decode Parameters
	if err := cause.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"cause": *cause.Name,
	})

	logger = logger.WithField("causeID", causeID)

	storedCause, err := api.db.GetCauseByID(ctx, &causeID)
	if err != nil {
		errMessage := fmt.Sprintf("Retrieving cause ID: %v", causeID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	storedCause.UpdateValues(&cause)
	// now update the database values
	err = api.db.UpdateCause(ctx, storedCause)
	if err != nil {
		logger.WithError(err).Warn("Error updating cause.")
		utils.WriteError(w, http.StatusInternalServerError, err, map[string]string{
			"requestID": ctx.Value("correlationid").(string),
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"Cause Updated": true,
	})

	utils.WriteJSON(w, http.StatusOK, storedCause)

}

// List - List all the causes
// GET - /causes
// Permission Admin
func (api *TicketCauseAPI) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_cause.go -> List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	causes, err := api.db.ListAllCauses(ctx)
	if err != nil {
		logger.WithError(err).Error("Retreiving all ticket causes")
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return

	}

	logger.Info("Causes Returned")

	utils.WriteJSON(w, http.StatusOK, &causes)

}

// Delete - Deletes a cause
// DELETE - /causes/{causeID}/causes/{causeID}
// Permission MemberIsTarget
func (api *TicketCauseAPI) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_cause.go -> Delete()")

	vars := mux.Vars(r)
	causeID := model.CauseID(vars["causeID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"CauseID": causeID,
		"pricipal":   principal,
	})

	deleted, err := api.db.DeleteCause(ctx, &causeID)
	if err != nil {

		logger.WithError(err).Error("deleteing cause")
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	logger.Info("Cause Deleted")

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}
