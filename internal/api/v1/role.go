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
	"github.com/lilkid3/ASA-Ticket/Backend/internal/storage"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/storage/database"
	"github.com/sirupsen/logrus"
)

// RolesAPI - structure holds rest endpoints for roles
type RolesAPI struct {
	db database.Database
}

// Load help create a subrouter for the roles
func loadRolesAPI(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	rolesAPI := &RolesAPI{db: env.DB}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/roles", rolesAPI.Create),//, authorizer.ObjAuthorize("role", "create")),
		newAPIEndpoint("GET", "/roles/{roleID}", rolesAPI.Get),//,  authorizer.ObjAuthorize("role", "view")), //retrieves a role using its ID
		newAPIEndpoint("GET", "/roles", rolesAPI.List),//, authorizer.ObjAuthorize("role", "list")),         //retrieves all the roles

		newAPIEndpoint("PATCH", "/roles/{roleID}", rolesAPI.Update),//, authorizer.ObjAuthorize("role", "update")),  //updates a user using its ID
		newAPIEndpoint("DELETE", "/roles/{roleID}", rolesAPI.Delete),//,  authorizer.ObjAuthorize("role", "delete")), //delete a user using its ID
	}
	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new Role
func (api *RolesAPI) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> RolesApi.Create()")

	// Decode parameters
	var role model.Role
	if err := role.Decode(r.Body); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := role.Verify(); err != nil {
		logger.WithError(err).Warn("Some field is missing")
		utils.WriteError(w, http.StatusBadRequest, "Not all fields were found", map[string]string{
			"error": err.Error(),
		})
		return

	}
	logger = logger.WithFields(logrus.Fields{
		"Role": role.Name,
	})

	if err := api.db.CreateRole(ctx, &role); err != nil {

		switch err {
		case storage.ErrRoleExists:
			logger.WithError(err).Warn("Role Already Exists")
			utils.WriteError(w, http.StatusConflict, err.Error(), nil)
			return
		case storage.ErrRoleIDExists:
			logger.WithError(err).Warn("RoleID Already Exists")
			utils.WriteError(w, http.StatusConflict, err.Error(), nil)
			return
		default:
			logger.WithError(err).Warn("Error creating role")

			utils.WriteError(w, http.StatusConflict, err.Error(), nil)
			return
		}
	}

	createdRole, err := api.db.GetRoleByID(ctx, &role.ID)
	if err != nil {
		logger.WithError(err).Warn("Error creating role")
		utils.WriteError(w, http.StatusConflict, "Error creating role", nil)
		return
	}

	logger.Info("Role created")

	utils.WriteJSON(w, http.StatusCreated, &createdRole)
}

// Get -  retreives role information
func (api *RolesAPI) Get(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "[API-Gateway] -> RolesApi.Get()")

	vars := mux.Vars(r)
	roleID := model.RoleID(vars["roleID"])

	ctx := r.Context()

	role, err := api.db.GetRoleByID(ctx, &roleID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching role RoleID: %v", roleID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	logger.WithField("RoleID", roleID).Debug("Get Role Complete")

	utils.WriteJSON(w, http.StatusOK, role)
}

// List - List all the roles
// GET - /roles
// Permission Admin
func (api *RolesAPI) List(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> RolesApi.List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	ctx := r.Context()

	roles, err := api.db.ListAllRoles(ctx)
	if err != nil {
		errMessage := fmt.Sprintf("Error retreiving all the users")
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error retreiving all the roles", nil)
		return

	}

	logger.Info("Roles List Returned")

	utils.WriteJSON(w, http.StatusOK, &roles)

}

// Update - Updated Role  Details
// PATCH - /roles/{roleID}
func (api *RolesAPI) Update(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> RolesApi.Update()")

	vars := mux.Vars(r)
	roleID := model.RoleID(vars["roleID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"RoleID":   roleID,
		"pricipal": principal,
	})

	// Decode parameters
	var userRole model.Role

	if err := userRole.Decode(r.Body); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()
	logger = logger.WithFields(logrus.Fields{
		"RolesID": roleID,
		"Role":    userRole.Name,
	})

	savedRole, err := api.db.GetRoleByID(ctx, &roleID)
	if err != nil {

		errMessage := fmt.Sprintf("Error getting role with ID: %v", roleID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusConflict, "Error getting role", nil)
		return
	}

	savedRole.UpdateValues(&userRole)

	// now update the database values
	if err := api.db.UpdateRole(ctx, &userRole); err != nil {
		errMessage := fmt.Sprintf("Error updating role RoleID: %v", roleID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error updating role", nil)
		return
	}
	logger.Info("Role Updated")
	utils.WriteJSON(w, http.StatusOK, &responses.ActUpdated{
		Updated: true,
	})

}

// Delete - Deletes a role
// DELETE - /roles/{userID}
func (api *RolesAPI) Delete(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> RolesApi.Delete()")

	vars := mux.Vars(r)
	roleID := model.RoleID(vars["roleID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"UserID":   roleID,
		"pricipal": principal,
	})

	ctx := r.Context()

	deleted, err := api.db.DeleteRole(ctx, &roleID)
	if err != nil {
		errMessage := fmt.Sprintf("Error deleting user: %v", roleID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if deleted {

		logger.Info("Role Deleted")
	}

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}
