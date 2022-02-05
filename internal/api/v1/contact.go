package v1

import (
	"context"
	"errors"
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

//ContactAPI - holds the contact endpoints
type ContactAPI struct {
	env *env.Env
	db  database.Database
}

// Load help create a subrouter for the contacts
func loadContactAPI(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	contactsAPI := &ContactAPI{env: env,
		db: env.DB,
	}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/contacts", contactsAPI.Create, authorizer.ObjAuthorize("contact", "create")),
		newAPIEndpoint("GET", "/contacts/{contactID}", contactsAPI.Get, authorizer.ObjAuthorize("contact", "view")), //retrieves a contactt using its ID
		newAPIEndpoint("GET", "/contacts", contactsAPI.List, authorizer.ObjAuthorize("contact", "list")),           //retrieves all the contacts

		newAPIEndpoint("PATCH", "/contacts/{contactID}", contactsAPI.Update, authorizer.ObjAuthorize("contact", "update")),  //updates a contact using its ID
		newAPIEndpoint("DELETE", "/contacts/{contactID}", contactsAPI.Delete, authorizer.ObjAuthorize("contact", "delete")), //delete a contact using its ID
}
	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new contact
func (api *ContactAPI) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> ContactsApi.Create()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	// Decode parameters
	var contact model.Contact
	if err := contact.Decode(r.Body); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		err = errors.New("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Get the userID from the Token
	contact.UserID = principal.UserID

	if err := contact.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"Contact Firstname":     *contact.Firstname,
		"Contact Lastname": *contact.Lastname,
		})

	if err := api.db.CreateContact(ctx, &contact); err != nil {
		logger.WithError(err).Warn("")
		utils.WriteError(w, http.StatusConflict, err, nil)
		return
	}

	createdContact, err := api.db.GetContactByID(ctx, &contact.ID)
	if err != nil {
		logger.WithError(err).Error()
		utils.WriteError(w, http.StatusConflict, err, nil)
		return
	}
	
	if err := api.getContactProps(ctx, createdContact); err != nil {
		logger.WithError(err).Error()
		utils.WriteError(w, http.StatusNotFound, err, nil)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, &createdContact)
}

// Get -  retreives contact information
func (api *ContactAPI) Get(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "[API-Gateway] -> ContactsApi.Get()")

	vars := mux.Vars(r)
	contactID := model.ContactID(vars["contactID"])

	ctx := r.Context()

	contact, err := api.db.GetContactByID(ctx, &contactID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching contact ContactID: %v", contactID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}
	if err := api.getContactProps(ctx, contact); err != nil {
		logger.WithError(err).Error()
		utils.WriteError(w, http.StatusNotFound, err, nil)
		return
	}
	logger.WithField("ContactID", contactID).Debug("Contact found")

	utils.WriteJSON(w, http.StatusOK, contact)
}

// List - List all the contacts
// GET - /contacts
// Permission Admin
func (api *ContactAPI) List(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> ContactsApi.List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	ctx := r.Context()

	contacts, err := api.db.ListAllContacts(ctx)
	if err != nil {
		errMessage := fmt.Sprintf("Error retreiving all the contacts")
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error retreiving all the contacts", nil)
		return

	}
	// add all the properties for the contacts
	for index := range contacts {
		if err := api.getContactProps(ctx, contacts[index]); err != nil {
			logger.WithError(err).Error()
			utils.WriteError(w, http.StatusNotFound, err, nil)
			return
		}
	}
	logger.Info("Contacts List Returned")

	utils.WriteJSON(w, http.StatusOK, &contacts)

}

// Update - Updated Contact  Details
// PATCH - /contacts/{contactID}
func (api *ContactAPI) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> ContactsApi.Update()")

	vars := mux.Vars(r)
	contactID := model.TicketID(vars["contactID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"ContactID": contactID,
		"principal": principal,
	})

	// Decode parameters
	var contact model.Contact

	if err := contact.Decode(r.Body); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	contact.ID = model.ContactID(contactID)
	storedcontact, err := api.db.GetContactByID(ctx, &contact.ID)
	if err != nil {
		errMessage := fmt.Sprintf("Retrieving contact ID: %v", &contactID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err, nil)
		return
	}
	storedcontact.UpdateValues(&contact)
	logger = logger.WithField("ContactID", contactID)

	err = api.db.UpdateContact(ctx, storedcontact)
	if err != nil {
		errMessage := fmt.Sprintf("Error updating contact ContactID: %v", contactID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	utils.WriteJSON(w, http.StatusOK, storedcontact)

}

// Delete - Deletes a contact
// DELETE - /contacts/{userID}
func (api *ContactAPI) Delete(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> ContactsApi.Delete()")

	vars := mux.Vars(r)
	contactID := model.ContactID(vars["contactID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"UserID":   contactID,
		"pricipal": principal,
	})

	ctx := r.Context()

	deleted, err := api.db.DeleteContact(ctx, &contactID)
	if err != nil {
		errMessage := fmt.Sprintf("Error deleting contact: %v", contactID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if deleted {

		logger.Info("Contact Deleted")
	}

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}




func (api *ContactAPI) getContactProps(ctx context.Context, contact *model.Contact) (err error) {

	contact.CreatedBy, _ = api.db.GetUserByID(ctx, &contact.UserID)
	return
}

