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

// TicketStatusAPI - structure holds handlers for Tickets Categories
type TicketStatusAPI struct {
	env *env.Env
	db  database.Database
}

// Load help create a subrouter for the status
func loadTicketStatus(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	api := &TicketStatusAPI{env: env, db: env.DB}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/ticket_statuses", api.Create, authorizer.ObjAuthorize("ticket_status", "create")),
		newAPIEndpoint("GET", "/ticket_statuses/{statusID}", api.Get, authorizer.ObjAuthorize("ticket_status", "view")), //retrieves a status using its ID
		newAPIEndpoint("GET", "/ticket_statuses", api.List, authorizer.ObjAuthorize("ticket_status", "list")),           //retrieves all the status

		newAPIEndpoint("PATCH", "/ticket_statuses/{statusID}", api.Update, authorizer.ObjAuthorize("ticket_status", "update")),  //updates a status using its ID
		newAPIEndpoint("DELETE", "/ticket_statuses/{statusID}", api.Delete, authorizer.ObjAuthorize("ticket_status", "delete")), //delete a status using its ID

	}

	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new Status
func (api *TicketStatusAPI) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_status.go -> StatusApi.Create()")

	//Load parameters
	var status model.Status

	if err := status.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// log.Printf("StatusParameter => %+v\n", status)
	logger = logger.WithFields(logrus.Fields{
		"status": *status.Name,
	})

	if err := status.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := api.db.CreateStatus(ctx, &status); err != nil {
		logger.WithError(err).Warn(err.Error())
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}
	createdStatus, err := api.db.GetStatusByID(ctx, &status.ID)
	if err != nil {
		logger.WithError(err).Warn(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
	}
	utils.WriteJSON(w, http.StatusCreated, &createdStatus)
}

// Get -  retreives status information
func (api *TicketStatusAPI) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := logrus.WithField("func", "ticket_status.go -> StatusApi.Get()")

	vars := mux.Vars(r)
	statusID := model.StatusID(vars["statusID"])

	//Get all the status fields from the database to ensure all is well
	status, err := api.db.GetStatusByID(ctx, &statusID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching status StatusID: %v", statusID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	logger.WithField("statusID", statusID).Debug("Get Status Complete")

	utils.WriteJSON(w, http.StatusOK, status)
}

// Update - Updated Status  Details
// PATCH - /status/{statusID}
// Permission MemberIsTarget, Admin
func (api *TicketStatusAPI) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_status.go -> Update()")

	vars := mux.Vars(r)
	statusID := model.StatusID(vars["statusID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"StatusID": statusID,
		"pricipal": principal,
	})

	var status model.Status

	// Decode Parameters
	if err := status.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"status": *status.Name,
	})

	logger = logger.WithField("statusID", statusID)

	storedStatus, err := api.db.GetStatusByID(ctx, &statusID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching status StatusID: %v", statusID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	storedStatus.UpdateValues(&status)
	// now update the database values

	if err := api.db.UpdateStatus(ctx, storedStatus); err != nil {
		logger.WithError(err).Warn("Error updating status.")
		utils.WriteError(w, http.StatusInternalServerError, "Error updating status.", map[string]string{
			"requestID": ctx.Value("correlationid").(string),
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"Status Updated": true,
	})

	utils.WriteJSON(w, http.StatusOK, status)

}

// List - List all the status
// GET - /status
// Permission Admin
func (api *TicketStatusAPI) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_status.go -> List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	status, err := api.db.ListAllStatus(ctx)
	if err != nil {
		logger.WithError(err).Warn("Retreiving all status")
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return

	}

	logger.Info("Categories Returned")

	utils.WriteJSON(w, http.StatusOK, &status)

}

// Delete - Deletes a status
// DELETE - /status/{statusID}/status/{statusID}
// Permission MemberIsTarget
func (api *TicketStatusAPI) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_status.go -> Delete()")

	vars := mux.Vars(r)
	statusID := model.StatusID(vars["statusID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"StatusID": statusID,
		"pricipal": principal,
	})

	deleted, err := api.db.DeleteStatus(ctx, &statusID)
	if err != nil {
		errMessage := fmt.Sprintf("Error deleting status: %v", statusID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	logger.Info("Status Deleted")

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}
