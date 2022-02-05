package v1

import (
	"context"
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

// NotesAPI - structure holds rest endpoints for notes
type NotesAPI struct {
	db database.Database
}

// Load help create a subrouter for the notes
func loadNotesAPI(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	notesAPI := &NotesAPI{db: env.DB}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/notes", notesAPI.Create, authorizer.ObjAuthorize("note", "create")),
		newAPIEndpoint("GET", "/notes/{noteID}", notesAPI.Get, authorizer.ObjAuthorize("note", "view")), //retrieves a note using its ID
		newAPIEndpoint("GET", "/notes", notesAPI.List, authorizer.ObjAuthorize("note", "list")),         //retrieves all the notes

		newAPIEndpoint("PATCH", "/notes/{noteID}", notesAPI.Update, authorizer.ObjAuthorize("note", "update")),  //updates a user using its ID
		newAPIEndpoint("DELETE", "/notes/{noteID}", notesAPI.Delete, authorizer.ObjAuthorize("note", "delete")), //delete a user using its ID
	}
	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new Note
func (api *NotesAPI) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> NotesApi.Create()")

	principal := middlewares.GetPrincipal(r)
	// Decode parameters
	var note model.Note
	if err := note.Decode(r.Body); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	note.UserID = principal.UserID
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

	logger.Info("Note created")

	utils.WriteJSON(w, http.StatusCreated, &createdNote)
}

// Get -  retreives note information
func (api *NotesAPI) Get(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "[API-Gateway] -> NotesApi.Get()")

	vars := mux.Vars(r)
	noteID := model.NoteID(vars["noteID"])

	ctx := r.Context()

	note, err := api.db.GetNoteByID(ctx, &noteID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching note NoteID: %v", noteID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	api.getNoteProps(ctx, note)
	logger.WithField("NoteID", noteID).Debug("Get Note Complete")

	utils.WriteJSON(w, http.StatusOK, note)
}

// List - List all the notes
// GET - /notes
// Permission Admin
func (api *NotesAPI) List(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> NotesApi.List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	ctx := r.Context()

	notes, err := api.db.ListAllNotes(ctx)
	if err != nil {
		errMessage := fmt.Sprintf("Error retreiving all the users")
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error retreiving all the notes", nil)
		return

	}

	// add all the note properties
	for index := range notes {
		api.getNoteProps(ctx,notes[index])
	}
	logger.Info("Notes List Returned")

	utils.WriteJSON(w, http.StatusOK, &notes)

}

// Update - Updated Note  Details
// PATCH - /notes/{noteID}
func (api *NotesAPI) Update(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> NotesApi.Update()")

	vars := mux.Vars(r)
	noteID := model.NoteID(vars["noteID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"NoteID":   noteID,
		"pricipal": principal,
	})

	// Decode parameters
	var userNote model.Note

	if err := userNote.Decode(r.Body); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()
	logger = logger.WithFields(logrus.Fields{
		"NotesID": noteID,
	})

	savedNote, err := api.db.GetNoteByID(ctx, &noteID)
	if err != nil {

		errMessage := fmt.Sprintf("Error getting note with ID: %v", noteID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusConflict, "Error getting note", nil)
		return
	}

	savedNote.UpdateValues(&userNote)

	// now update the database values
	if err := api.db.UpdateNote(ctx, &userNote); err != nil {
		errMessage := fmt.Sprintf("Error updating note NoteID: %v", noteID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error updating note", nil)
		return
	}
	logger.Info("Note Updated")
	utils.WriteJSON(w, http.StatusOK, &responses.ActUpdated{
		Updated: true,
	})

}

// Delete - Deletes a note
// DELETE - /notes/{userID}
func (api *NotesAPI) Delete(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> NotesApi.Delete()")

	vars := mux.Vars(r)
	noteID := model.NoteID(vars["noteID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"UserID":   noteID,
		"pricipal": principal,
	})

	ctx := r.Context()

	deleted, err := api.db.DeleteNote(ctx, &noteID)
	if err != nil {
		errMessage := fmt.Sprintf("Error deleting user: %v", noteID)
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

func (api *NotesAPI) getNoteProps(ctx context.Context, note *model.Note)(err error) {

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

	return nil
}
