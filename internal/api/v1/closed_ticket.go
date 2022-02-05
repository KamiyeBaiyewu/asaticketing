package v1

import (
	"context"
	"fmt"
	"log"
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

//ClosedTicketAPI - holds the ticket endpoints for the tickets
type ClosedTicketAPI struct {
	env *env.Env
	db  database.Database
}

// Load help create a subrouter for the tickets
func loadClosedTicketAPI(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	closedTicketAPI := &ClosedTicketAPI{env: env,
		db: env.DB,
	}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("GET", "/closed_tickets/", closedTicketAPI.List, authorizer.ObjAuthorize("closed_ticket", "view")),                  //retrieves all closed tickets
		newAPIEndpoint("GET", "/closed_tickets/{ticketID}", closedTicketAPI.Get, authorizer.ObjAuthorize("closed_ticket", "view")),         //retrieves a ticket using its ID
		newAPIEndpoint("DELETE", "/closed_tickets/{ticketID}", closedTicketAPI.Delete, authorizer.ObjAuthorize("closed_ticket", "delete")), //retrieves a ticket using its ID

	}
	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Get -  retreives ticket information
func (api *ClosedTicketAPI) Get(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "[API-Gateway] -> closedTicketAPI.Get()")

	vars := mux.Vars(r)
	ticketID := model.TicketID(vars["ticketID"])

	ctx := r.Context()

	ticket, err := api.db.GetTicketByID(ctx, &ticketID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching ticket TicketID: %v", ticketID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}
	if err := api.getTicketProps(ctx, ticket); err != nil {
		logger.WithError(err).Error()
		utils.WriteError(w, http.StatusNotFound, err, nil)
		return
	}
	logger.WithField("TicketID", ticketID).Debug("Get Ticket Complete")

	utils.WriteJSON(w, http.StatusOK, ticket)
}

// List - List all the tickets
// GET - /tickets
// Permission Admin
func (api *ClosedTicketAPI) List(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> closedTicketAPI.List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})
	log.Printf("Query => %+v\n", r.URL.Query())
	ctx := r.Context()

	tickets, err := api.db.ListAllTickets(ctx)
	if err != nil {
		errMessage := fmt.Sprintf("Error retreiving all the tickets")
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error retreiving all the tickets", nil)
		return

	}
	// add all the properties for the tickets
	for index := range tickets {
		if err := api.getTicketProps(ctx, tickets[index]); err != nil {
			logger.WithError(err).Error()
			utils.WriteError(w, http.StatusNotFound, err, nil)
			return
		}
	}
	logger.Info("Tickets List Returned")

	utils.WriteJSON(w, http.StatusOK, &tickets)

}

// Delete - Deletes a ticket
// DELETE - /tickets/{userID}
func (api *ClosedTicketAPI) Delete(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> closedTicketAPI.Delete()")

	vars := mux.Vars(r)
	ticketID := model.TicketID(vars["ticketID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"UserID":   ticketID,
		"pricipal": principal,
	})

	ctx := r.Context()

	deleted, err := api.db.DeleteTicket(ctx, &ticketID)
	if err != nil {
		errMessage := fmt.Sprintf("Error deleting ticket: %v", ticketID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if deleted {

		logger.Info("Ticket Deleted")
	}

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}

func (api *ClosedTicketAPI) getTicketProps(ctx context.Context, ticket *model.Ticket) (err error) {

	ticket.Category, err = api.env.DB.GetCategoryByID(ctx, &ticket.CategoryID)
	if err != nil {
		logrus.WithError(err).Warn("Error category of ticket")
		return
	}
	ticket.Priority, err = api.db.GetPriorityByID(ctx, &ticket.PriorityID)
	if err != nil {
		logrus.WithError(err).Warn("Error priority of ticket")
		return
	}
	ticket.Status, err = api.db.GetStatusByID(ctx, &ticket.StatusID)
	if err != nil {
		logrus.WithError(err).Warn("Error status of ticket")
		return
	}
	ticket.SLA, err = api.db.GetSLAByID(ctx, &ticket.SLAID)
	if err != nil {
		logrus.WithError(err).Warn("Error SLA of ticket")
		return
	}
	ticket.Source, err = api.db.GetSourceByID(ctx, &ticket.SourceID)
	if err != nil {
		logrus.WithError(err).Warn("Error source of ticket")
		return
	}

	ticket.CreatedBy, err = api.db.GetUserByID(ctx, &ticket.UserID)
	if err != nil {
		logrus.WithError(err).Warn("Error user id of the person who created the ticket")
		return
	}
	
	/*
		ticket.ClosingRemark, err = api.db.ClosingRemark(ctx, &ticket.ID)
		if err != nil {
			logrus.WithError(err).Warn("Error retrieving the closing remark")
			return
		}else {
			//fetch the root cause
			ticket.ClosingRemark.Cause, err = api.db.GetCauseByID(ctx, &ticket.ClosingRemark.CauseID)
			if err != nil {
				logrus.WithError(err).Warn("Error retrieving root cause of ticket")
				return
			}
			ticket.ClosingRemark.ClosedBy, err = api.db.GetUserByID(ctx, &ticket.ClosingRemark.UserID)
			if err != nil {
				logrus.WithError(err).Warn("Error retrieving who closed the ticket")
				return
			}
		} */

	// Remove all the unecessary
	// category
	ticket.CategoryID = model.NilCategoryID
	ticket.Category.Description = nil
	ticket.Category.Weight = nil
	ticket.Category.CreatedAt = nil
	ticket.Category.UpdatedAt = nil
	ticket.Category.DeletedAt = nil

	ticket.PriorityID = model.NilPriorityID
	ticket.Priority.Weight = nil
	ticket.Priority.CreatedAt = nil
	ticket.Priority.UpdatedAt = nil
	ticket.Priority.DeletedAt = nil

	ticket.StatusID = model.NilStatusID
	ticket.Status.Weight = nil
	ticket.Status.CreatedAt = nil
	ticket.Status.UpdatedAt = nil
	ticket.Status.DeletedAt = nil

	ticket.SLAID = model.NilSLAID
	ticket.SLA.Weight = nil
	ticket.SLA.CreatedAt = nil
	ticket.SLA.UpdatedAt = nil
	ticket.SLA.DeletedAt = nil

	ticket.SourceID = model.NilSourceID
	ticket.Source.Weight = nil
	ticket.Source.CreatedAt = nil
	ticket.Source.UpdatedAt = nil
	ticket.Source.DeletedAt = nil

	ticket.UserID = model.NilUserID
	ticket.CreatedBy.Name = func() *string {
		s := fmt.Sprintf("%s %s", *ticket.CreatedBy.Firstname, *ticket.CreatedBy.Lastname)
		return &s
	}()
	ticket.CreatedBy.Firstname = nil
	ticket.CreatedBy.Lastname = nil
	ticket.CreatedBy.Email = nil
	ticket.CreatedBy.Type = nil
	ticket.CreatedBy.RoleID = model.NilRoleID
	ticket.CreatedBy.PasswordHash = nil
	ticket.CreatedBy.IsActive = nil
	ticket.CreatedBy.IsSystem = nil
	ticket.CreatedBy.CreatedAt = nil
	ticket.CreatedBy.UpdatedAt = nil
	ticket.CreatedBy.DeletedAt = nil

	return
}

func (api *ClosedTicketAPI) getNoteProps(ctx context.Context, note *model.Note) (err error) {

	note.CreatedBy, err = api.db.GetUserByID(ctx, &note.UserID)
	if err != nil {
		logrus.WithError(err).Warn("Error fetching the user who created the ticket")
		return
	}

	note.CreatedBy.Name = func() *string {
		s := fmt.Sprintf("%s %s", *note.CreatedBy.Firstname, *note.CreatedBy.Lastname)
		return &s
	}()

	note.CreatedBy.Firstname = nil
	note.CreatedBy.Lastname = nil
	note.CreatedBy.Type = nil
	note.CreatedBy.CreatedAt = nil
	note.CreatedBy.UpdatedAt = nil
	note.CreatedBy.RoleID = model.NilRoleID
	note.TicketID = model.NilTicketID

	return nil
}
