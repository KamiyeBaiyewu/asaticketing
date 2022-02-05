package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/middlewares"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/responses"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/utils"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/env"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/storage/database"
	"github.com/sirupsen/logrus"
)

// PolicyAPI - structure holds rest endpoints for object policies
type PolicyAPI struct {
	db  database.Database
	env *env.Env
}

// Load help create a subrouter for the policies
func loadPolicyAPI(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	policiesAPI := &PolicyAPI{db: env.DB, env: env}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/policies", policiesAPI.Create, authorizer.ObjAuthorize("policy", "create")),
		newAPIEndpoint("GET", "/policies/{policyID}", policiesAPI.Get, authorizer.ObjAuthorize("policy", "view")),         //retrieves a policy using its ID
		newAPIEndpoint("GET", "/policies", policiesAPI.List, authorizer.ObjAuthorize("policy", "list")),                   //retrieves all the users
		newAPIEndpoint("DELETE", "/policies/{policyID}", policiesAPI.Delete, authorizer.ObjAuthorize("policy", "delete")), //delete a user using its ID

	}
	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new Policy in the system
func (api *PolicyAPI) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> PoliciesApi.Create()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	// Decode parameters
	var policy model.Policy
	if err := policy.Decode(r.Body); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		err = errors.New("Error with submited values")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Get the userID from the Token
	policy.UserID = principal.UserID

	if err := policy.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Retrieve the role from the database
	role, err := api.db.GetRoleByID(ctx, &policy.RoleID)
	if err != nil {
		logger.WithError(err).Warn("Error fetching role")
		utils.WriteError(w, http.StatusNotFound, "Role not found", nil)
		return
	}
	// retrive the Ojects
	object, err := api.db.GetObjectByID(ctx, &policy.ObjectID)

	if err != nil {
		logger.WithError(err).Warn("Error fetching role")
		utils.WriteError(w, http.StatusNotFound, "Role not found", nil)
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"Role":   role.Name,
		"Object": object.Name,
		"Action": policy.Action,
	})

	// println("Code arrived here")
	err = api.db.CreatePolicy(ctx, &policy)
	if err != nil {
		logger.WithError(err).Warn("Creating poicy")
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	createdPolicy, err := api.db.GetPolicyByID(ctx, &policy.ID)
	if err != nil {
		logger.WithError(err).Warn("Retrieving the newly created policy")
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}
	// Now add one polices: One using the role's ID 
	_, err = api.env.AddPolicy(string(role.ID), strings.ToLower(*object.Name), *policy.Action)
	// _, err = api.env.AddPolicy(*role.Name, strings.ToLower(*object.Name), *policy.Action)
	if err != nil {
		// Try to delete the Created Policy
		api.db.DeletePolicy(ctx, &createdPolicy.ID)
		utils.WriteError(w, http.StatusInternalServerError, "Error saving Policy", nil)

		return
	}

	utils.WriteJSON(w, http.StatusCreated, &createdPolicy)

}

// Get -  retreives policy information
func (api *PolicyAPI) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	principal := middlewares.GetPrincipal(r)

	logger := logrus.WithField("func", "polcy.go -> PolicyApi.Get()")

	vars := mux.Vars(r)
	policyID := model.PolicyID(vars["policyID"])
	logger.WithFields(
		logrus.Fields{"principal": principal,
			"policyID": policyID})

	//Get all the policy fields from the database to ensure all is well
	policy, err := api.db.GetPolicyByID(ctx, &policyID)
	if err != nil {
		errMessage := "Error fetching policy"
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err, nil)
		return
	}
	// add all the policy properties
	err = api.getPolicyProps(ctx, policy)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err, nil)
		return
	}
	utils.WriteJSON(w, http.StatusOK, policy)
}

// List - List all the policies
// GET - /policies
// Permission Admin
func (api *PolicyAPI) List(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "[API-Gateway] -> PoliciesApi.List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	ctx := r.Context()

	policies, err := api.db.ListAllPolicies(ctx)
	if err != nil {
		errMessage := fmt.Sprintf("Error retreiving all the users")
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error retreiving all the policies", nil)
		return

	}
	// add all the policy properties for all the policies
	for index := range policies {
		err = api.getPolicyProps(ctx, policies[index])
		if err != nil {
			utils.WriteError(w, http.StatusNotFound, err, nil)
			return
		}
	}

	logger.Info("Policies List Returned")

	utils.WriteJSON(w, http.StatusOK, &policies)

}

// Delete - Deletes a policy
// DELETE - /policies/{userID}
func (api *PolicyAPI) Delete(w http.ResponseWriter, r *http.Request) {
	// Show function name in error logs to track errors faster
	ctx := r.Context()
	logger := logrus.WithField("func", "[API-Gateway] -> PoliciesApi.Delete()")

	vars := mux.Vars(r)
	policyID := model.PolicyID(vars["policyID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"UserID":   policyID,
		"pricipal": principal,
	})

	//Get all the policy fields from the database to ensure all is well
	policy, err := api.db.GetPolicyByID(ctx, &policyID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching policy PolicyID: %v", policyID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, err, nil)
		return
	}

	// Retrieve the role from the database
	role, err := api.db.GetRoleByID(ctx, &policy.RoleID)
	if err != nil {
		logger.WithError(err).Warn("Error fetching role")
		utils.WriteError(w, http.StatusNotFound, "Role not found", nil)
		return
	}
	// retrive the Ojects
	object, err := api.db.GetObjectByID(ctx, &policy.ObjectID)

	if err != nil {
		logger.WithError(err).Warn("Error fetching role")
		utils.WriteError(w, http.StatusNotFound, "Object not found", nil)
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"Role":   role.Name,
		"Object": *object.Name,
		"Action": policy.Action,
	})

	// Now Delete the one created polices: One using the role's ID 
	api.env.RemovePolicy(string(role.ID), strings.ToLower(*object.Name), *policy.Action)
	// api.env.RemovePolicy(*role.Name, strings.ToLower(*object.Name), *policy.Action)

	deleted, err := api.db.DeletePolicy(ctx, &policyID)
	if err != nil {
		errMessage := fmt.Sprintf("Error deleting user: %v", policyID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if deleted {

		logger.Info("Policy Deleted")

	}

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}

func (api *PolicyAPI) getPolicyProps(ctx context.Context, policy *model.Policy) (err error) {

	// get who created it
	policy.CreatedBy, err = api.db.GetUserByID(ctx, &policy.UserID)
	if err != nil {
		logrus.WithError(err).Warn("Error fetching who crated a policy")
		return
	}
	// Get the role
	policy.Role, err = api.db.GetRoleByID(ctx, &policy.RoleID)
	if err != nil {
		logrus.WithError(err).Warn("Error fetching the role of policy")
		return
	}
	// Get the Object
	policy.Object, err = api.db.GetObjectByID(ctx, &policy.ObjectID)
	if err != nil {
		logrus.WithError(err).Warn("Error fetching the object of policy")
		return
	}

	// Remove all the unecessary
	policy.UserID = model.NilUserID
	policy.RoleID = model.NilRoleID
	policy.ObjectID = model.NilObjectID

	// remove Uncesary Properties

	// User
	policy.CreatedBy.RoleID = model.NilRoleID
	policy.CreatedBy.Type = nil
	policy.CreatedBy.Email = nil
	policy.CreatedBy.CreatedAt = nil
	policy.CreatedBy.UpdatedAt = nil
	policy.CreatedBy.DeletedAt = nil

	// Role
	policy.Role.Description = nil
	policy.Role.CreatedAt = nil
	policy.Role.UpdatedAt = nil
	policy.Role.DeletedAt = nil

	// Object
	policy.Object.CreatedBy = model.NilUserID
	policy.Object.Description = nil
	policy.Object.CreatedAt = nil
	policy.Object.UpdatedAt = nil
	policy.Object.DeletedAt = nil

	return
}
