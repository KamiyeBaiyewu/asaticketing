package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/auth"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/utils"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/env"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/storage/database"
	"github.com/sirupsen/logrus"
)

// Authorizer - middleware for verifying token
type Authorizer struct {
	enforcer *casbin.CachedEnforcer
	db       database.Database
	//rolesCache    gcache.Cache //Have to enable cache to reduce the ammount of database calls
}

// NewAuthorizer - creates instance of Autorization Middleware
func NewAuthorizer(env *env.Env) *Authorizer {

	authorizer := &Authorizer{
		enforcer: env.Enforcer,
		db:       env.DB,
	}

	// TODO: uncommenet code when cache is added
	/* 	authorizer.rolesCache = gcache.New(200).LRU().LoaderExpireFunc(func(key interface{}) (interface{}, *time.Duration, error) {
		userID := key.(model.UserID)
		roles, err := authorizer.db.GetRoleNamesForUser(context.Background(), &userID)
		if err != nil {
			return nil, nil, err
		}
		expire := time.Duration(env.Config.Authorizer.CacheExpiration) * time.Second
		return roles, &expire, nil
	}).Build() */

	return authorizer
}

// ObjAuthorize - checks if a role can access a certain resource
func (a *Authorizer) ObjAuthorize(obj, act string) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		logger := logrus.WithField("func", "[Authorizer.ObjAuthorize]")
		fn := func(w http.ResponseWriter, r *http.Request) {
			// service.NewBasicJwtauthService()
			//get request context
			ctx := r.Context()

			// extract the token from the request
			token, err := getToken(r)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"Error": err.Error(),
					"Op":    "retrieving token",
				}).Error()

				utils.WriteError(w, http.StatusUnauthorized, "Invalid Token", nil)
				return
			}

			principal, err := a.getPrincipal(ctx, token)
			if err != nil {
				logrus.Debug(err.Error())
				utils.WriteError(w, http.StatusUnauthorized, err.Error(), nil)
				return
			}
			logger.WithFields(
				logrus.Fields{
					"principal": principal,
				})

			ctx = addDetailsToContext(ctx, principal)

			authorized := false

			if strings.Contains(principal.Role, "|") { //means multiple roles present

				roles := strings.Split(principal.Role, "|")
				for _, role := range roles {
					//println("P => Checking Role: ", role, " Object: ", obj, " Action: ", act)
					// casbin enforce
					authorized, err = a.enforcer.Enforce(role, obj, act)
					if err != nil {
						logrus.Debug(err.Error())
						utils.WriteError(w, http.StatusUnauthorized, "Error during Authorization", nil)
						return
					}
					if authorized {
						break
					}
				}
			} else {
				//println("S => Checking Role: ", principal.Role, " Object: ", obj, " Action: ", act)
				// casbin enforce
				authorized, err = a.enforcer.Enforce(principal.Role, obj, act)
				if err != nil {
					logrus.Debug(err.Error())
					utils.WriteError(w, http.StatusUnauthorized, "Error during Authorization", nil)
					return
				}
			}

			if authorized {

				ctx = context.WithValue(ctx, principalContextKey, *principal)
				//put the context in the http.request context and make sure you take it out at http handlers
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			} else {
				logger.WithFields(
					logrus.Fields{
						"Role":   string(principal.Role),
						"Object": obj,
						"Action": act,
					},
				).Info()
				utils.WriteError(w, http.StatusUnauthorized, "unauthorized", nil)
				return
			}

		}

		return http.HandlerFunc(fn)
	}
}

func (a *Authorizer) enforce(role, object, action string) (authorized bool, err error) {
	authorized, err = a.enforcer.Enforce(role, object, action)
	return
}

// UserOnUserType -  helps to check the kind of proncipal user can carry out action on another type of user
func (a *Authorizer) UserOnUserType(principalUserType, objectUserType, action string) (authorized bool, err error) {

	principalUserType = fmt.Sprintf("user_type_%v", principalUserType)
	objectUserType = fmt.Sprintf("user_type_%v", objectUserType)
	authorized, err = a.enforcer.Enforce(principalUserType, objectUserType, action)
	return
}
func addDetailsToContext(ctx context.Context, principal *model.Principal) (valContext context.Context) {
	//put desired data in the context
	valContext = ctx
	valContext = context.WithValue(valContext, "userid", string(principal.UserID))
	valContext = context.WithValue(valContext, "role", string(principal.Role))

	return
}

// getPrincipal -  checks the the token submitted is valid
func (a *Authorizer) getPrincipal(ctx context.Context, accessToken string) (*model.Principal, error) {

	principal, err := auth.VerifyToken(accessToken)
	if err == nil {

		// TODO: uncommenet code when cache is added

		/* 	roles, err := a.rolesCache.Get(principal.UserID)
		   	if err != nil {
		   		logrus.WithError(err).Info("Retrieving IDs of primary and auxilary role")
		   		return nil, err
		   	} */

		roles, err := a.db.GetRoleNamesForUser(ctx, &principal.UserID)
		if err != nil {
			logrus.WithError(err).Info("Retrieving IDs of primary and auxilary role")
			return nil, err
		}
		principal.Role = roles
	}
	return principal, err

}
