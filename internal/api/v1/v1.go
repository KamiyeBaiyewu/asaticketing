package v1

import (
	"github.com/gorilla/mux"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/middlewares"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/env"
)

// LoadRoutes helps create version 1 subrouter
func LoadRoutes(router *mux.Router, env *env.Env, authorizer *middlewares.Authorizer) {


	v1Router := router.PathPrefix("/api/v1").Subrouter()

	// has to come before the users and the roles endpoint because of router links
	loadClosingRemark(v1Router, env, authorizer)
	loadNotesAPI(v1Router, env, authorizer)
	loadRolesAPI(v1Router, env, authorizer)
	loadObject(v1Router, env, authorizer)
	loadPolicyAPI(v1Router, env, authorizer)
	loadUsersRolesAPI(v1Router, env, authorizer)
	loadRolesAPI(v1Router, env, authorizer)
	
	// Tickets
	loadTicketAPI(v1Router, env, authorizer)
	loadTicketCause(v1Router, env, authorizer)
	loadTicketCategory(v1Router, env, authorizer)
	loadTicketPriority(v1Router, env, authorizer)
	loadTicketSource(v1Router, env, authorizer)
	loadTicketStatus(v1Router, env, authorizer)
	loadUserAPI(v1Router, env, authorizer)

	loadSLA(v1Router, env, authorizer)
	// (v1Router, env, authorizer)
	// (v1Router, env, authorizer)
	// (v1Router, env, authorizer)

	//Contacts
	loadContactAPI(v1Router, env, authorizer)
}
