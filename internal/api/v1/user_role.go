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

// UserRoleAPI - structure holds rest endpoints for user's roles
type UserRoleAPI struct {
	db database.Database
}

// Load help create a subrouter for the roles
func loadUsersRolesAPI(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	userRolesAPI := &UserRoleAPI{db: env.DB}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/users/{userID}/role/{roleID}", userRolesAPI.GrantRole, authorizer.ObjAuthorize("role", "create")),
		newAPIEndpoint("DELETE", "/users/{userID}/role/{roleID}", userRolesAPI.RevokeRole, authorizer.ObjAuthorize("role", "delete")),
		newAPIEndpoint("GET", "/roles/{roleID}/users", userRolesAPI.RoleUsers, authorizer.ObjAuthorize("role", "list")),
		newAPIEndpoint("GET", "/users/{userID}/roles", userRolesAPI.UsersRoles, authorizer.ObjAuthorize("role", "list")),

		// List the User Roles
	}
	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// GrantRole -  Handler grants a user a role
func (api *UserRoleAPI) GrantRole(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := logrus.WithField("func", "[API-Gateway] -> UsersRolesApi.GrantRole()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	roleID := model.RoleID(vars["roleID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(
		logrus.Fields{
			"UserID":    userID,
			"RoleID":    roleID,
			"principal": principal,
		})

	// first ensure the user does not have this role as his primary role
	isPrimaryRole, err := api.db.IsPrimaryRole(ctx, &userID, &roleID)
	if err != nil {
		logger.WithError(err).Warn("Error granting role")
		utils.WriteError(w, http.StatusInternalServerError, "Error granting role", nil)
		return
	}
	if isPrimaryRole { //return an error if this role is user's primary role
		utils.WriteError(w, http.StatusConflict, "Primary role of user", nil)
		return
	}
	err = api.db.GrantRole(ctx, &userID, &roleID)
	if err != nil {
			errMessage := fmt.Sprintf("Error granting role")
			logger.WithError(err).Warn(errMessage)
			utils.WriteError(w, http.StatusInternalServerError, err, nil)
			return
	
	}

	logger.Info("Role Granted")

	utils.WriteJSON(w, http.StatusOK, &responses.ActGranted{
		Granted: true,
	})
}

// RevokeRole -  DELETES the role if it exists
func (api *UserRoleAPI) RevokeRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := logrus.WithField("func", "[API-Gateway] -> UsersRolesApi.GrantRole()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	roleID := model.RoleID(vars["roleID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(
		logrus.Fields{
			"UserID":    userID,
			"RoleID":    roleID,
			"principal": principal,
		})

	// first ensure the user does not have this role as his primary role
	isPrimaryRole, err := api.db.IsPrimaryRole(ctx, &userID, &roleID)
	if err != nil {
		logger.WithError(err).Warn("Error revoke role")
		utils.WriteError(w, http.StatusInternalServerError, "Error revoking role", nil)
		return
	}
	if isPrimaryRole { //return an error if this role is user's primary role
		utils.WriteError(w, http.StatusConflict, "Primary role of user", nil)
		return
	}

	err = api.db.RevokeRole(ctx, &userID, &roleID)
	if err != nil {
		errMessage := "Error revoking role"
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error revoking role", nil)
		return
	}

	logger.Info("Role Revoked")
	utils.WriteJSON(w, http.StatusCreated, &responses.ActRevoked{
		Revoked: true,
	})
}

// RoleUsers - Fetches all the users associated with a role
func (api *UserRoleAPI) RoleUsers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := logrus.WithField("func", "[API-Gateway] -> UsersRolesApi.RoleUsers()")

	vars := mux.Vars(r)
	roleID := model.RoleID(vars["roleID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(
		logrus.Fields{
			"roleID":    roleID,
			"principal": principal,
		})

	users, err := api.db.GetRoleUsers(ctx, &roleID)
	if err != nil {
		errMessage := "Error retreiving users"
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error retrieving users", nil)
		return
	}

	logger.Info("Users Returned")

	utils.WriteJSON(w, http.StatusOK, &users)

}

// UsersRoles - Fetched all the roles associated with a user
func (api *UserRoleAPI) UsersRoles(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := logrus.WithField("func", "[API-Gateway] -> UsersRolesApi.UsersRoles()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(
		logrus.Fields{
			"userID":    userID,
			"principal": principal,
		})

	roles, err := api.db.GetRoleByUser(ctx, &userID)
	if err != nil {

		errMessage := "Error retreiving roles"
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error revoking roles", nil)
		return
	}

	logger.Info("Roles Returned")

	utils.WriteJSON(w, http.StatusOK, &roles)
}
