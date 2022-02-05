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

// SLAAPI - structure holds handlers for Tickets Categories
type SLAAPI struct {
	env *env.Env
	db  database.Database
}

// Load help create a subrouter for the sla
func loadSLA(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	api := &SLAAPI{env: env, db: env.DB}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/ticket_slas", api.Create, authorizer.ObjAuthorize("ticket_sla", "create")),
		newAPIEndpoint("GET", "/ticket_slas/{SLAID}", api.Get, authorizer.ObjAuthorize("ticket_sla", "view")), //retrieves a sla using its ID
		newAPIEndpoint("GET", "/ticket_slas", api.List, authorizer.ObjAuthorize("ticket_sla", "list")),        //retrieves all the sla

		newAPIEndpoint("PATCH", "/ticket_slas/{SLAID}", api.Update, authorizer.ObjAuthorize("ticket_sla", "update")),  //updates a sla using its ID
		newAPIEndpoint("DELETE", "/ticket_slas/{SLAID}", api.Delete, authorizer.ObjAuthorize("ticket_sla", "delete")), //delete a sla using its ID

	}

	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new SLA
func (api *SLAAPI) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "sla.go -> SLAAPI.Create()")

	//Load parameters
	var sla model.SLA

	if err := sla.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := sla.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"sla": *sla.Name,
	})
	if err := api.db.CreateSLA(ctx, &sla); err != nil {
		logger.WithError(err).Warn(err.Error())
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}
	createdSLA, err := api.db.GetSLAByID(ctx, &sla.ID)
	if err != nil {
		logger.WithError(err).Warn(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
	}
	utils.WriteJSON(w, http.StatusCreated, &createdSLA)
}

// Get -  retreives sla information
func (api *SLAAPI) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := logrus.WithField("func", "sla.go -> SLAAPI.Get()")

	vars := mux.Vars(r)
	SLAID := model.SLAID(vars["SLAID"])

	//Get all the sla fields from the database to ensure all is well
	sla, err := api.db.GetSLAByID(ctx, &SLAID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching sla ID: %v", SLAID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	logger.WithField("SLAID", SLAID).Debug("Get SLA Complete")

	utils.WriteJSON(w, http.StatusOK, sla)
}

// Update - Updated SLA  Details
// PATCH - /sla/{SLAID}
// Permission MemberIsTarget, Admin
func (api *SLAAPI) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "sla.go -> Update()")

	vars := mux.Vars(r)
	SLAID := model.SLAID(vars["SLAID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"SLAID":    SLAID,
		"pricipal": principal,
	})

	var sla model.SLA

	// Decode Parameters
	if err := sla.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"sla": *sla.Name,
	})

	logger = logger.WithField("SLAID", SLAID)

	storedSLA, err := api.db.GetSLAByID(ctx, &SLAID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching sla ID: %v", SLAID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	storedSLA.UpdateValues(&sla)
	// now update the database values

	if err := api.db.UpdateSLA(ctx, storedSLA); err != nil {
		logger.WithError(err).Warn("Error updating sla.")
		utils.WriteError(w, http.StatusInternalServerError, "Error updating sla.", map[string]string{
			"requestID": ctx.Value("correlationid").(string),
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"SLA Updated": true,
	})

	utils.WriteJSON(w, http.StatusOK, storedSLA)

}

// List - List all the sla
// GET - /sla
// Permission Admin
func (api *SLAAPI) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "sla.go -> List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	sla, err := api.db.ListAllSLA(ctx)
	if err != nil {
		logger.WithError(err).Warn("Retreiving all sla")
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return

	}

	logger.Info("Categories Returned")

	utils.WriteJSON(w, http.StatusOK, &sla)

}

// Delete - Deletes a sla
// DELETE - /sla/{SLAID}/sla/{SLAID}
// Permission MemberIsTarget
func (api *SLAAPI) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "sla.go -> Delete()")

	vars := mux.Vars(r)
	SLAID := model.SLAID(vars["SLAID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"StatusID": SLAID,
		"pricipal": principal,
	})

	deleted, err := api.db.DeleteSLA(ctx, &SLAID)
	if err != nil {
		errMessage := fmt.Sprintf("Error deleting sla: %v", SLAID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	logger.WithFields(logrus.Fields{
		"SLA Deleted": deleted,
	})

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}
