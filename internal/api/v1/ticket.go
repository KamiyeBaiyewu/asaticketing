package v1

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/middlewares"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/responses"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/utils"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/env"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/storage/database"
	"github.com/sirupsen/logrus"
)

//TicketAPI - holds the ticket endpoints for the tickets
type TicketAPI struct {
	env *env.Env
	db  database.Database
}

// Load help create a subrouter for the tickets
func loadTicketAPI(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	ticketsAPI := &TicketAPI{env: env,
		db: env.DB,
	}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/tickets", ticketsAPI.Create, authorizer.ObjAuthorize("ticket", "create")),
		newAPIEndpoint("GET", "/tickets/{ticketID}", ticketsAPI.Get, authorizer.ObjAuthorize("ticket", "view")), //retrieves a ticket using its ID
		newAPIEndpoint("GET", "/tickets", ticketsAPI.List, authorizer.ObjAuthorize("ticket", "list")),           //retrieves all the ticjets

		newAPIEndpoint("PATCH", "/tickets/{ticketID}", ticketsAPI.Update, authorizer.ObjAuthorize("ticket", "update")),  //updates a ticket using its ID
		newAPIEndpoint("DELETE", "/tickets/{ticketID}", ticketsAPI.Delete, authorizer.ObjAuthorize("ticket", "delete")), //delete a ticket using its ID

		newAPIEndpoint("POST", "/tickets/{ticketID}/close", ticketsAPI.Close, authorizer.ObjAuthorize("ticket", "update")),               //adds a note to a ticket
		newAPIEndpoint("GET", "/tickets/{ticketID}/notes", ticketsAPI.ListNotes, authorizer.ObjAuthorize("ticket", "update")),              //retrieves all the notes for a ticket
		newAPIEndpoint("POST", "/tickets/{ticketID}/notes", ticketsAPI.AddNote, authorizer.ObjAuthorize("ticket", "update")),               //adds a note to a ticket
		newAPIEndpoint("DELETE", "/tickets/{ticketID}/notes/{noteID}", ticketsAPI.DeleteNote, authorizer.ObjAuthorize("ticket", "update")), //deletes a note for a ticket
	
	
		/* 
		newAPIEndpoint("GET", "/tickets/{ticketID}", ticketsAPI.Get, authorizer.ObjAuthorize("ticket", "view")), //retrieves a ticket using its ID
		newAPIEndpoint("GET", "/tickets/{ticketID}", ticketsAPI.Get, authorizer.ObjAuthorize("ticket", "view")), //retrieves a ticket using its ID
		 */
	}
	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new Ticket
func (api *TicketAPI) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> TicketsApi.Create()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	// Decode parameters
	var ticket model.Ticket
	if err := ticket.Decode(r.Body); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		err = errors.New("Error with submited values")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Get the userID from the Token
	ticket.UserID = principal.UserID

	if err := ticket.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// log.Printf("UserParameter => %+v\n", userParameters)
	logger = logger.WithFields(logrus.Fields{
		"TicketID":           ticket.ID,
		"Ticket Subject":     *ticket.Subject,
		"Ticket Description": *ticket.Description,
	})

	// Get the SLA
	sla, err := api.db.GetSLAByID(ctx, &ticket.SLAID)

	if err != nil {
		logger.WithError(err).Warn("")
		utils.WriteError(w, http.StatusConflict, err, nil)
		return
	}

	// set the due date using the grace period in the SLA
	ticket.DueDate = func() *time.Time { t := time.Now(); t = t.Add(time.Duration(*sla.GracePeriod) * time.Hour); return &t }()
	// set the userID
	ticket.UserID = principal.UserID
	if err := api.db.CreateTicket(ctx, &ticket); err != nil {
		logger.WithError(err).Warn("")
		utils.WriteError(w, http.StatusConflict, err, nil)
		return
	}

	createdTicket, err := api.db.GetTicketByID(ctx, &ticket.ID)
	if err != nil {
		logger.WithError(err).Error()
		utils.WriteError(w, http.StatusConflict, err, nil)
		return
	}

	if err := api.getTicketProps(ctx, createdTicket); err != nil {
		logger.WithError(err).Error()
		utils.WriteError(w, http.StatusNotFound, err, nil)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, &createdTicket)
}

// Get -  retreives ticket information
func (api *TicketAPI) Get(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "[API-Gateway] -> TicketsApi.Get()")

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
func (api *TicketAPI) List(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> TicketsApi.List()")

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

// Update - Updated Ticket  Details
// PATCH - /tickets/{ticketID}
func (api *TicketAPI) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> TicketsApi.Update()")

	vars := mux.Vars(r)
	ticketID := model.TicketID(vars["ticketID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"TicketID": ticketID,
		"pricipal": principal,
	})

	// Decode parameters
	var ticket model.Ticket

	if err := ticket.Decode(r.Body); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ticket.ID = ticketID
	storedticket, err := api.db.GetTicketByID(ctx, &ticketID)
	if err != nil {
		errMessage := fmt.Sprintf("Retrieving ticket ID: %v", ticketID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err, nil)
		return
	}
	storedticket.UpdateValues(&ticket)
	logger = logger.WithField("TicketID", ticketID)

	err = api.db.UpdateTicket(ctx, storedticket)
	if err != nil {
		errMessage := fmt.Sprintf("Error updating ticket TicketID: %v", ticketID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	utils.WriteJSON(w, http.StatusOK, storedticket)

}

// Delete - Deletes a ticket
// DELETE - /tickets/{userID}
func (api *TicketAPI) Delete(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> TicketsApi.Delete()")

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

// ListNotes - returns all the notes for a ticket
func (api *TicketAPI) ListNotes(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> TicketsApi.ListNotes()")

	principal := middlewares.GetPrincipal(r)

	vars := mux.Vars(r)
	ticketID := model.TicketID(vars["ticketID"])

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
		"TicketID": ticketID,
	})

	
	ctx := r.Context()

	notes, err := api.db.ListAllTicketNotes(ctx, &ticketID)
	if err != nil {
		errMessage := fmt.Sprintf("Error retreiving all the tickets")
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error retreiving all the tickets", nil)
		return

	}

	for index := range notes {
		err = api.getNoteProps(ctx, notes[index])
		if err != nil {
			logger.WithError(err).Warn("Adding note properties")
			utils.WriteError(w, http.StatusInternalServerError, "Error retreiving all the tickets", nil)
			return
		}
	}

	logger.Info("Ticket Notes Returned")

	utils.WriteJSON(w, http.StatusOK, &notes)
}

// AddNote - adds note to a ticket
func (api *TicketAPI) AddNote(w http.ResponseWriter, r *http.Request) {

	logger := logrus.WithField("func", "[API-Gateway] -> TicketsApi.AddNote()")

	principal := middlewares.GetPrincipal(r)

	vars := mux.Vars(r)
	ticketID := model.TicketID(vars["ticketID"])

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
		"TicketID": ticketID,
	})

	var note model.Note
	if err := note.Decode(r.Body); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()

	note.UserID = principal.UserID
	note.TicketID = ticketID
	if err := note.Verify(); err != nil {
		logger.WithError(err).Warn("Some field is missing")
		utils.WriteError(w, http.StatusBadRequest, "Not all fields were found", map[string]string{
			"error": err.Error(),
		})
		return

	}

	if err := api.db.CreateNote(ctx, &note); err != nil {

		logger.WithError(err).Warn("Error creating note")

		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	createdNote, err := api.db.GetNoteByID(ctx, &note.ID)
	if err != nil {
		logger.WithError(err).Warn("Error creating note")
		utils.WriteError(w, http.StatusConflict, "Error creating note", nil)
		return
	}

	err = api.getNoteProps(ctx, createdNote)
	if err != nil {
		logger.WithError(err).Warn("Adding note properties")
		utils.WriteError(w, http.StatusInternalServerError, "Error retreiving all the tickets", nil)
		return
	}
	logger.Info("Note created")
	utils.WriteJSON(w, http.StatusCreated, &createdNote)
}

// Close - closes a ticket
func (api *TicketAPI) Close(w http.ResponseWriter, r *http.Request) {

	logger := logrus.WithField("func", "[API-Gateway] -> TicketsApi.DeleteNote()")

	principal := middlewares.GetPrincipal(r)

	vars := mux.Vars(r)
	ticketID := model.TicketID(vars["ticketID"])

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
		"TicketID": ticketID,
	})

	ctx := r.Context()

	//Load parameters
	var closingRemark model.ClosingRemark

	if err := closingRemark.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	closingRemark.UserID = principal.UserID
	closingRemark.TicketID = ticketID
	if err := closingRemark.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}

	
	err := api.db.CreateClosingRemark(ctx, &closingRemark)
	if err != nil {
		logger.WithError(err).Error("Creating closingRemark")
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	createdClosingRemark, err := api.db.GetClosingRemarkByID(ctx, &closingRemark.ID)
	if err != nil {
		logger.WithError(err).Warn("Retrieving remark for recently closed ticket")
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}
	// update the closed value on the database
	_, err = api.db.CloseTicket(ctx, &ticketID)
	if err != nil {
		logger.WithError(err).Warn("Updating closed at to now")
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, &createdClosingRemark)

}

// DeleteNote - deletes note from a ticekt
func (api *TicketAPI) DeleteNote(w http.ResponseWriter, r *http.Request) {

	logger := logrus.WithField("func", "[API-Gateway] -> TicketsApi.DeleteNote()")

	principal := middlewares.GetPrincipal(r)

	vars := mux.Vars(r)
	ticketID := model.TicketID(vars["ticketID"])
	noteID := model.NoteID(vars["noteID"])

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
		"TicketID": ticketID,
		"noteID":   noteID,
	})

	ctx := r.Context()

	deleted, err := api.db.DeleteTicketNote(ctx, &ticketID, &noteID)
	if err != nil {
		errMessage := fmt.Sprintf("Error deleting ticket note: %v", noteID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if deleted {

		logger.Info("Note Deleted")
	}

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})
}

func (api *TicketAPI) getTicketProps(ctx context.Context, ticket *model.Ticket) (err error) {

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

func (api *TicketAPI) getNoteProps(ctx context.Context, note *model.Note) (err error) {

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
