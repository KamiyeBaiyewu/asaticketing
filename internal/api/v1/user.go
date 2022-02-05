package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/auth"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/middlewares"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/requests"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/responses"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/utils"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/env"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/storage/database"

	"github.com/sirupsen/logrus"
)

// UserAPI - structure holds privies rest for users
type UserAPI struct {
	db database.Database
}

// Load help create a subrouter for the users
func loadUserAPI(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {

	userAPI := &UserAPI{db: env.DB}

	apiEndpoint := []apiEndpoint{

		newAPIEndpoint("POST", "/users", userAPI.Create),
		newAPIEndpoint("GET", "/users/{userID}", userAPI.Get, authorizer.ObjAuthorize("user", "view")), //retrieves a user using its ID
		newAPIEndpoint("GET", "/users", userAPI.List, authorizer.ObjAuthorize("user", "list")),         //retrieves all the users

		newAPIEndpoint("PATCH", "/users/{userID}", userAPI.Update, authorizer.ObjAuthorize("user", "update")),  //updates a user using its ID
		newAPIEndpoint("DELETE", "/users/{userID}", userAPI.Delete, authorizer.ObjAuthorize("user", "delete")), //delete a user using its ID
		// ----- AUTHORIZATION -----
		newAPIEndpoint("POST", "/login", userAPI.Login),
		// ----- TOKENS -----
		newAPIEndpoint("POST", "/refresh", userAPI.RefreshToken),
	}
	for _, api := range apiEndpoint {

		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}

}

// Create - Creates a new User
func (api *UserAPI) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "user -> user.go -> UserApi.Create()")

	//Load parameters
	var userParameters requests.UserParameters

	if err := userParameters.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := userParameters.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// ensure the password is available
	if !utils.ValidatePassword(userParameters.Password) {
		utils.WriteError(w, http.StatusBadRequest, "Invalid Password", map[string]string{
			"error": ("Password: 8 Characters conatining one lower, one upper, one number, one character"),
		})
		return
	}

	// log.Printf("UserParameter => %+v\n", userParameters)
	logger = logger.WithFields(logrus.Fields{
		"email": *userParameters.Email,
	})
	hashed, err := model.HashPassword(userParameters.Password)
	if err != nil {
		logger.WithError(err).Warn("Could not hash password.")
		utils.WriteError(w, http.StatusInternalServerError, "could not hash password", map[string]string{
			"requestID": ctx.Value("correlationid").(string),
		})
		return
	}

	newUser := &model.User{
		Firstname:    userParameters.Firstname,
		Lastname:     userParameters.Lastname,
		Email:        userParameters.Email,
		Type:         userParameters.Type,
		RoleID:       model.RoleID(userParameters.RoleID),
		PasswordHash: &hashed,
	}

	ctx = r.Context()

	err = api.db.CreateUser(ctx, newUser)
	if err != nil {
		logger.WithError(err).Error()
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}
	createdUser, err := api.db.GetUserByID(ctx, &newUser.ID)
	if err != nil {
		logger.WithError(err).Warn("Error fetching newly created user")
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	// get the user's role
	role, err := api.db.GetRoleByID(ctx, &createdUser.RoleID)
	if err != nil {
		logger.WithError(err).Warn("Error fetching the roles of newly created user")
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}
	createdUser.Role = role
	// remove all that is not needed from the endpoint
	createdUser.RoleID = model.NilRoleID

	api.writeToTokenResponse(ctx, w, http.StatusCreated, createdUser, userParameters.DeviceID, true)

}

// Login - Validates User Credentials
func (api *UserAPI) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "user -> user.go -> UserApi.Login()")

	var credentials model.Credentials

	if err := credentials.Decode(r.Body); err != nil {
		logger.WithError(err).Warn("Could not decode credentials")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
	}

	logger = logger.WithFields(logrus.Fields{
		"email": credentials.Email,
	})

	if err := credentials.SessionData.Verify(); err != nil {
		logger.WithError(err).Warn("Not all fields found")
		utils.WriteError(w, http.StatusBadRequest, "Not all fields were found", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// get the user by email
	// User comes from the User Service
	user, err := api.db.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		logger.WithError(err).Warn("Error logging in")
		utils.WriteError(w, http.StatusBadRequest, "Invalid email or password", nil)
		return

	}
	//checking if password is correct
	if err := user.CheckPassword(credentials.Password); err != nil {
		logger.WithError(err).Warn("Error logging in")
		utils.WriteError(w, http.StatusBadRequest, "Invalid email or password", nil)
		return
	}

	// Add the User permitted Actions on objects
	objActions, err := api.db.ListPermitedObjectActions(ctx, &user.ID)
	if err != nil {
		logger.WithError(err).Warn("Error fetching object actions")
		utils.WriteError(w, http.StatusInternalServerError, "Currently unable to login", nil)
		return
	}
	// Add the User permitted Actions on system objects
	sysActions, err := api.db.ListPermitedSystemActions(ctx, &user.ID)
	if err != nil {
		logger.WithError(err).Warn("Error fetching system actions")
		utils.WriteError(w, http.StatusInternalServerError, "Currently unable to login", nil)
		return
	}
	actions := append(sysActions, objActions...)
	user.Actions = actions
	logger.WithField("userID", user.ID).Debug("User logged in")

	// get the user's role
	role, err := api.db.GetRoleByID(ctx, &user.RoleID)
	if err != nil {
		logger.WithError(err).Warn("Error fetching the roles of newly created user")
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}
	user.Role = role

	// remove all that is not needed from the endpoint
	user.RoleID = model.NilRoleID

	api.writeToTokenResponse(ctx, w, http.StatusOK, user, credentials.DeviceID, true)

}

// Get -  retreives user information
func (api *UserAPI) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	principal := middlewares.GetPrincipal(r)
	logger := logrus.WithField("func", "user -> user.go -> UserApi.Get()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])

	logger = logger.WithFields(logrus.Fields{
		"UserID":   userID,
		"pricipal": principal,
	})

	//Get all the user fields from the database to ensure all is well
	user, err := api.db.GetUserByID(ctx, &userID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching user UserID: %v", userID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusNotFound, "User not found", nil)
		return
	}

	logger.WithField("userID", userID).Debug("Get User Complete")

	utils.WriteJSON(w, http.StatusOK, user)
}

// Update - Updated User  Details
// PATCH - /users/{userID}
// Permission MemberIsTarget, Admin
func (api *UserAPI) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "user.go -> Update()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"UserID":   userID,
		"pricipal": principal,
	})
	var userRequest requests.UserParameters

	// Decode Parameters
	if err := userRequest.Decode(r.Body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := userRequest.Verify(); err != nil {
		logger.WithError(err).Warn("Error with submitted values")
		utils.WriteError(w, http.StatusBadRequest, "Error with submitted values", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"Request email":  *userRequest.Email,
		"Request userID": userRequest.ID,
		"userID":         userID,
	})

	storedUser, err := api.db.GetUserByID(ctx, &userID)
	if err != nil {
		errMessage := fmt.Sprintf("Error fetching user UserID: %v", userID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusConflict, err, nil)
		return
	}

	// fill the values that are not being updated
	if len(userRequest.Password) > 0 {
		if userRequest.Password != userRequest.ConfirmPassword {
			utils.WriteError(w, http.StatusInternalServerError, "Passwords don't match", map[string]string{
				"requestID": ctx.Value("correlationid").(string),
			})
			return
		}
		if err = storedUser.SetPassword(userRequest.Password); err != nil {
			errMessage := fmt.Sprintf("Error setting password UserID: %v", userID)
			logger.WithError(err).Warn(errMessage)
			utils.WriteError(w, http.StatusInternalServerError, "Error setting password ", map[string]string{
				"requestID": ctx.Value("correlationid").(string),
			})
			return
		}
	}
	storedUser.UpdateValues(&userRequest.User)

	// now update the database values
	err = api.db.UpdateUser(ctx, storedUser)
	if err != nil {
		logger.WithError(err).Warn("Error updating user.")
		utils.WriteError(w, http.StatusInternalServerError, "Error updating user.", map[string]string{
			"requestID": ctx.Value("correlationid").(string),
		})
		return
	}

	logger.Info("User Updated")
	utils.WriteJSON(w, http.StatusOK, storedUser)

}

// List - List all the users
// GET - /users
// Permission Admin
func (api *UserAPI) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "user.go -> List()")

	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"pricipal": principal,
	})

	users, err := api.db.ListAllUsers(ctx)
	if err != nil {
		errMessage := fmt.Sprintf("Error retreiving all the users")
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, "Error retreiving all the users", nil)
		return

	}

	logger.Info("Users Returned")

	utils.WriteJSON(w, http.StatusOK, &users)

}

// Delete - Deletes a user
// DELETE - /users/{userID}/users/{userID}
// Permission MemberIsTarget
func (api *UserAPI) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Show function name in error logs to track errors faster
	logger := logrus.WithField("func", "user.go -> Delete()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := middlewares.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"UserID":   userID,
		"pricipal": principal,
	})

	deleted, err := api.db.DeleteUser(ctx, &userID)
	if err != nil {
		errMessage := fmt.Sprintf("Error deleting user: %v", userID)
		logger.WithError(err).Warn(errMessage)
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	logger.Info("User Deleted")

	utils.WriteJSON(w, http.StatusOK, &responses.ActDeleted{
		Deleted: deleted,
	})

}

// RefreshToken - accepts a refresh Token to return a new access token
func (api *UserAPI) RefreshToken(w http.ResponseWriter, r *http.Request) {}

// writeToTokenResponse - generates access and refresh tokens, returnes them to user, Refresh token is stored in database as session
func (api *UserAPI) writeToTokenResponse(ctx context.Context, w http.ResponseWriter, status int, user *model.User, deviceID model.DeviceID, cookie bool) {

	// Issue the toke
	fullName := fmt.Sprintf("%s %s", *user.Firstname, *user.Lastname)
	principal := model.Principal{UserID: user.ID, Name: fullName, Role: *user.Role.Name, Type: *user.Type}
	tokens, err := auth.IssueToken(principal)
	if err != nil || tokens == nil {
		logrus.WithError(err).Warn("Error issuing token")
		utils.WriteError(w, http.StatusUnauthorized, "Error Issuing Token", nil)
		return
	}

	session := &model.Session{
		UserID:       user.ID,
		DeviceID:     deviceID,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.RefreshTokenExpiresAt,
	}

	if err := api.db.SaveRefreshToken(ctx, session); err != nil {
		logrus.WithError(err).Warn("Error issuing token")
		utils.WriteError(w, http.StatusUnauthorized, "Error Issuing Token", nil)
		return
	}

	// write token response:
	tokenResponse := responses.TokenResponse{
		User:   user,
		Tokens: tokens,
	}
	if cookie {
		// TODO:
	}

	utils.WriteJSON(w, status, tokenResponse)

}
