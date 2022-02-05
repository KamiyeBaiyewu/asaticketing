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

// TicketSourceAPI - structure holds handlers for Tickets Sources
type TicketSourceAPI struct {
	env *env.Env
	db  database.Database
}

// Load help create a subrouter for the sources
func loadTicketSource(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	api := &TicketSourceAPI{env: env,
		db: env.DB,
	}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/ticket_sources", api.Create, authorizer.ObjAuthorize("ticket_source", "create")),
		newAPIEndpoint("GET", "/ticket_sources/{sourceID}", api.Get, authorizer.ObjAuthorize("ticket_source", "view")), //retrieves a source using its ID
		newAPIEndpoint("GET", "/ticket_sources", api.List, authorizer.ObjAuthorize("ticket_source", "list")),           //retrieves all the sources

		newAPIEndpoint("PATCH", "/ticket_sources/{sourceID}", api.Update, authorizer.ObjAuthorize("ticket_source", "update")),  //updates a source using its ID
		newAPIEndpoint("DELETE", "/ticket_sources/{sourceID}", api.Delete, authorizer.ObjAuthorize("ticket_source", "delete")), //delete a source using its ID

	}

	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new Source
func (api *TicketSourceAPI) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_source.go -> SourceApi.Create()")

	//Load parameters
	var source model.Source

	if err := source.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// log.Printf("SourceParameter => %+v\n", source)
	logger = logger.WithFields(logrus.Fields{
		"source": *source.Name,
	})

	if err := source.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := api.db.CreateSource(ctx, &source); err != nil {
		logger.WithError(err).Warn(err.Error())
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	createdSource, err := api.db.GetSourceByID(ctx, &source.ID)
	if err != nil {
		logger.WithError(err).Warn(err.Error())
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, &createdSource)
}

// Get -  retreives source information
func (api *TicketSourceAPI) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := logrus.WithField("func", "ticket_source.go -> SourceApi.Get()")

	vars := mux.Vars(r)
	sourceID := model.SourceID(vars["sourceID"])

	//Get all the source fields from the database to ensure all is well
	source, err := api.db.GetSourceByID(ctx, &sourceID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching source SourceID: %v", sourceID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err, nil)
		return
	}

	logger.WithField("sourceID", sourceID).Debug("Get Source Complete")

	utils.WriteJSON(w, http.StatusOK, source)
}

// Update - Updated Source  Details
// PATCH - /sources/{sourceID}
// Permission MemberIsTarget, Admin
func (api *TicketSourceAPI) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_source.go -> Update()")

	vars := mux.Vars(r)
	sourceID := model.SourceID(vars["sourceID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"SourceID": sourceID,
		"pricipal": principal,
	})

	var source model.Source

	// Decode Parameters
	if err := source.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithField("sourceID", sourceID)

	storedSource, err := api.db.GetSourceByID(ctx, &sourceID)
	if err != nil {
		errMessage := fmt.Sprintf("Retrieving source ID: %v", sourceID)
		logger.WithError(err).Error(errMessage)
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}
	logger = logger.WithFields(logrus.Fields{
		"source": *storedSource.Name,
	})

	storedSource.UpdateValues(&source)
	// now update the database values

	if err := api.db.UpdateSource(ctx, storedSource); err != nil {
		logger.WithError(err).Warn("Error updating source.")
		utils.WriteError(w, http.StatusInternalServerError, "Error updating source.", map[string]string{
			"requestID": ctx.Value("correlationid").(string),
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"Source Updated": true,
	})

	utils.WriteJSON(w, http.StatusOK, source)

}

// List - List all the sources
// GET - /sources
// Permission Admin
func (api *TicketSourceAPI) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_source.go -> List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	sources, err := api.db.ListAllSources(ctx)
	if err != nil {
		logger.WithError(err).Error("Retreiving all the sources")
		utils.WriteError(w, http.StatusInternalServerError, "Error retreiving all the sources", nil)
		return

	}

	logger.Info("Sources Returned")

	utils.WriteJSON(w, http.StatusOK, &sources)

}

// Delete - Deletes a source
// DELETE - /sources/{sourceID}/sources/{sourceID}
// Permission MemberIsTarget
func (api *TicketSourceAPI) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_source.go -> Delete()")

	vars := mux.Vars(r)
	sourceID := model.SourceID(vars["sourceID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"SourceID": sourceID,
		"pricipal": principal,
	})

	deleted, err := api.db.DeleteSource(ctx, &sourceID)
	if err != nil {
		errMessage := fmt.Sprintf("Error deleting source ID: %v", sourceID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	logger.Info("Source Deleted")

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}
