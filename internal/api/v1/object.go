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

// ObjectAPI - structure holds handlers for Tickets Objects
type ObjectAPI struct {
	db database.Database
}

// Load help create a subrouter for the objects
func loadObject(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	api := &ObjectAPI{db: env.DB}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/objects", api.Create, authorizer.ObjAuthorize("object", "create")),
		newAPIEndpoint("GET", "/objects/{objectID}", api.Get, authorizer.ObjAuthorize("object", "view")), //retrieves a object using its ID
		newAPIEndpoint("GET", "/objects", api.List, authorizer.ObjAuthorize("object", "list")),           //retrieves all the objects

		newAPIEndpoint("PATCH", "/objects/{objectID}", api.Update, authorizer.ObjAuthorize("object", "update")),  //updates a object using its ID
		newAPIEndpoint("DELETE", "/objects/{objectID}", api.Delete, authorizer.ObjAuthorize("object", "delete")), //delete a object using its ID

	}

	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new Object
func (api *ObjectAPI) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "object.go -> ObjectApi.Create()")

	// get the principal from context
	principal := middlewares.GetPrincipal(r)

	//Load parameters
	var object model.Object

	if err := object.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// INsert the principal userID into the Object
	// object.CreatedBy = principal.UserID

	// log.Printf("ObjectParameter => %+v\n", object)
	logger = logger.WithFields(logrus.Fields{
		"object": *object.Name,
		"prinicipal": principal,
	})
	object.CreatedBy = principal.UserID
	if err := object.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}

	err := api.db.CreateObject(ctx, &object)
	if err != nil {
		logger.WithError(err).Warn("Storing noew object")
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}
	createdObject, err := api.db.GetObjectByID(ctx, &object.ID)
	if err != nil {
		logger.WithError(err).Warn("Retreiving newly created object")
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, &createdObject)
}

// Get -  retreives object information
func (api *ObjectAPI) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	principal := middlewares.GetPrincipal(r)

	logger := logrus.WithField("func", "object.go -> ObjectApi.Get()")
	vars := mux.Vars(r)
	objectID := model.ObjectID(vars["objectID"])

	logger.WithFields(
		logrus.Fields{"objectID": objectID,
			"principal": principal}).Debug("Get Object Complete")

	//Get all the object fields from the database to ensure all is well
	object, err := api.db.GetObjectByID(ctx, &objectID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching object ObjectID: %v", objectID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.WriteJSON(w, http.StatusOK, object)
}

// Update - Updated Object  Details
// PATCH - /objects/{objectID}
// Permission MemberIsTarget, Admin
func (api *ObjectAPI) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "object.go -> Update()")

	vars := mux.Vars(r)
	objectID := model.ObjectID(vars["objectID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"ObjectID": objectID,
		"pricipal": principal,
	})

	var object model.Object

	// Decode Parameters
	if err := object.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"object": *object.Name,
	})

	logger = logger.WithField("objectID", objectID)

	storedObject, err := api.db.GetObjectByID(ctx, &objectID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching object ObjectID: %v", objectID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	storedObject.UpdateValues(&object)
	// now update the database values
	err = api.db.UpdateObject(ctx, storedObject)
	if err != nil {
		// TODO: respond to error wiht the correlationid in the context
		utils.WriteError(w, http.StatusInternalServerError, "Temporarily unable to update object.", map[string]string{
			"requestID": ctx.Value("correlationid").(string),
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"Object Updated": true,
	})

	utils.WriteJSON(w, http.StatusOK, storedObject)

}

// List - List all the objects
// GET - /objects
// Permission Admin
func (api *ObjectAPI) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "object.go -> List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	objects, err := api.db.ListAllObjects(ctx)
	if err != nil {
		errMessage := fmt.Sprintf("Error retreiving all the objects")
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error retreiving all the objects", nil)
		return

	}

	logger.Info("Objects Returned")

	utils.WriteJSON(w, http.StatusOK, &objects)

}

// Delete - Deletes a object
// DELETE - /objects/{objectID}/objects/{objectID}
// Permission MemberIsTarget
func (api *ObjectAPI) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "object.go -> Delete()")

	vars := mux.Vars(r)
	objectID := model.ObjectID(vars["objectID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"ObjectID": objectID,
		"pricipal": principal,
	})

	// Get the Object first to ensure is not a standart object
	object, err := api.db.GetObjectByID(ctx, &objectID)
	if err != nil {
		
		logger.WithError(err).Warn("retrieving object")
		utils.WriteError(w, http.StatusNotFound, "Object not found", nil)
		return
	}
	if *object.IsStandard {

		utils.WriteError(w, http.StatusConflict, "Standard Objects can not be deleted", nil)
		return
	}

	deleted, err := api.db.DeleteObject(ctx, &objectID)
	if err != nil {
		errMessage := fmt.Sprintf("Error deleting object: %v", objectID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	logger.Info("Object Deleted")

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}
