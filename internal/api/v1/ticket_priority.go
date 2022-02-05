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

// TicketPriorityAPI - structure holds handlers for Tickets Priorities
type TicketPriorityAPI struct {
	env *env.Env
	db  database.Database
}

// Load help create a subrouter for the priorities
func loadTicketPriority(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	api := &TicketPriorityAPI{env: env,
		db: env.DB,
	}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/ticket_priorities", api.Create, authorizer.ObjAuthorize("ticket_priority", "create")),
		newAPIEndpoint("GET", "/ticket_priorities/{priorityID}", api.Get, authorizer.ObjAuthorize("ticket_priority", "view")), //retrieves a priority using its ID
		newAPIEndpoint("GET", "/ticket_priorities", api.List, authorizer.ObjAuthorize("ticket_priority", "list")),             //retrieves all the priorities

		newAPIEndpoint("PATCH", "/ticket_priorities/{priorityID}", api.Update, authorizer.ObjAuthorize("ticket_priority", "update")),  //updates a priority using its ID
		newAPIEndpoint("DELETE", "/ticket_priorities/{priorityID}", api.Delete, authorizer.ObjAuthorize("ticket_priority", "delete")), //delete a priority using its ID

	}

	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new Priority
func (api *TicketPriorityAPI) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_priority.go -> PriorityApi.Create()")

	//Load parameters
	var priority model.Priority

	if err := priority.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := priority.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}
	logger = logger.WithFields(logrus.Fields{
		"priority": *priority.Name,
	})

	if err := api.db.CreatePriority(ctx, &priority); err != nil {
		logger.WithError(err).Warn(err.Error())
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}
	createdPriority, err := api.db.GetPriorityByID(ctx, &priority.ID)

	if err != nil {
		logger.WithError(err).Warn(err.Error())
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, &createdPriority)
}

// Get -  retreives priority information
func (api *TicketPriorityAPI) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logrus.WithField("func", "ticket_priority.go -> PriorityApi.Get()")

	vars := mux.Vars(r)
	priorityID := model.PriorityID(vars["priorityID"])

	//Get all the priority fields from the database to ensure all is well
	priority, err := api.db.GetPriorityByID(ctx, &priorityID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching priority PriorityID: %v", priorityID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	logger.WithField("priorityID", priorityID).Debug("Get Priority Complete")

	utils.WriteJSON(w, http.StatusOK, priority)
}

// Update - Updated Priority  Details
// PATCH - /priorities/{priorityID}
// Permission MemberIsTarget, Admin
func (api *TicketPriorityAPI) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_priority.go -> Update()")

	vars := mux.Vars(r)
	priorityID := model.PriorityID(vars["priorityID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"PriorityID": priorityID,
		"pricipal":   principal,
	})

	var priority model.Priority

	// Decode Parameters
	if err := priority.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"priority": *priority.Name,
	})

	logger = logger.WithField("priorityID", priorityID)

	storedPriority, err := api.db.GetPriorityByID(ctx, &priorityID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching priority PriorityID: %v", priorityID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	storedPriority.UpdateValues(&priority)
	// now update the database values
	err = api.db.UpdatePriority(ctx, storedPriority)
	if err != nil {
		logger.WithError(err).Error("Updating priority")
		utils.WriteError(w, http.StatusInternalServerError, err, map[string]string{
			"requestID": ctx.Value("correlationid").(string),
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"Priority Updated": true,
	})

	utils.WriteJSON(w, http.StatusOK, storedPriority)

}

// List - List all the priorities
// GET - /priorities
// Permission Admin
func (api *TicketPriorityAPI) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_priority.go -> List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	priorities, err := api.db.ListAllPriorities(ctx)
	if err != nil {
		logger.WithError(err).Warn("Retreiving all ticket priorities")
		utils.WriteError(w, http.StatusInternalServerError, "Error retreiving all the priorities", nil)
		return

	}

	logger.Info("Priorities Returned")

	utils.WriteJSON(w, http.StatusOK, &priorities)

}

// Delete - Deletes a priority
// DELETE - /priorities/{priorityID}/priorities/{priorityID}
// Permission MemberIsTarget
func (api *TicketPriorityAPI) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "ticket_priority.go -> Delete()")

	vars := mux.Vars(r)
	priorityID := model.PriorityID(vars["priorityID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"PriorityID": priorityID,
		"pricipal":   principal,
	})

	deleted, err := api.db.DeletePriority(ctx, &priorityID)
	if err != nil {
		errMessage := fmt.Sprintf("Error deleting priority ID: %v", priorityID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	logger.Info("Priority Deleted")

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}
