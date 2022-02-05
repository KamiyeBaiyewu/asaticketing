package enforcer

import (
	"fmt"
	"strings"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	casbin "github.com/casbin/casbin/v2"
	"github.com/go-pg/pg/v9"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/config"
	"github.com/sirupsen/logrus"
)

var (
	modelPath  string
	policyPath string
	dbTable    string
	dbUser     string
	dbSecret   string
	dbName     string
	dbHost     string
	dbPort     int
	dbParam    string
)

// Init - Initializes the Enforcer
func Init(cfg *config.Info) (enforcer *casbin.CachedEnforcer, db *pg.DB) {

	dbUser = cfg.Database.Username
	dbSecret = cfg.Database.Password
	dbHost = cfg.Database.Hostname
	dbPort = cfg.Database.Port
	dbName = cfg.Database.Database
	dbParam = cfg.Database.Parameter

	// Casbin specific configurations
	modelPath = cfg.Casbin.Model
	policyPath = cfg.Casbin.Policy
	dbTable = cfg.Casbin.Table

	// Build parameters
	param := dbParam

	// If parameter is specified, add a question mark
	// Don't add one if a question mark is already there
	if len(dbParam) > 0 && !strings.HasPrefix(dbParam, "?") {
		param = "?" + dbParam

	}
	// "postgresql://username:password@postgres:5432/database?sslmode=disable"
	dbURL := fmt.Sprintf("postgresql://%v:%v@%v:%d/%v%v", dbUser, dbSecret, dbHost, dbPort, dbName, param)

	opts, _ := pg.ParseURL(dbURL)

	db = pg.Connect(opts)
	//defer db.Close()

	a, _ := pgadapter.NewAdapterByDB(db, pgadapter.WithTableName(dbTable))

	// enforcer, err := casbin.NewEnforcer(*modelPath, a)
	enforcer, err := casbin.NewCachedEnforcer(modelPath, a)
	if err != nil {
		logrus.Fatal(err)
	}
	// Load the policy from DB.
	err = enforcer.LoadPolicy()
	if err != nil {
		logrus.Fatal(err)
	}

	enforcer.EnableAutoSave(true)

	logrus.Infof("Enforcer has loaded the Policies")
	saved, err := addDefaultPolicies(enforcer)
	if saved {
		logrus.Info("New policies have been saved")
	}
	if enforcer == nil {
		panic("the enforcer is nil")
	}
	return

}

func addDefaultPolicies(enforcer *casbin.CachedEnforcer) (saved bool, err error) {
	// Now add policies to the one already stored
	saved, err = enforcer.AddPolicy("user_type_admin", "user_type_agent", "create")
	saved, err = enforcer.AddPolicy("user_type_admin", "user_type_agent", "view")
	saved, err = enforcer.AddPolicy("user_type_admin", "user_type_agent", "update")
	saved, err = enforcer.AddPolicy("user_type_admin", "user_type_agent", "delete")

	/* admin users on users */
	saved, err = enforcer.AddPolicy("user_type_admin", "user_type_user", "create")
	saved, err = enforcer.AddPolicy("user_type_admin", "user_type_user", "view")
	saved, err = enforcer.AddPolicy("user_type_admin", "user_type_user", "update")
	saved, err = enforcer.AddPolicy("user_type_admin", "user_type_user", "delete")

	/* agent users on users */
	saved, err = enforcer.AddPolicy("user_type_agent", "user_type_user", "create")
	saved, err = enforcer.AddPolicy("user_type_agent", "user_type_user", "view")
	saved, err = enforcer.AddPolicy("user_type_agent", "user_type_user", "update")
	saved, err = enforcer.AddPolicy("user_type_agent", "user_type_user", "delete")

	// Add admin policies
	// policy
	
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "policy", enforcer)
	// saved, err = addPolicyForAllAction("admin", "policy", enforcer)
	// user
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "user", enforcer)
	// saved, err = addPolicyForAllAction("admin", "user", enforcer)
	// role
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "role", enforcer)
	// saved, err = addPolicyForAllAction("admin", "role", enforcer)
	// object
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "object", enforcer)
	// saved, err = addPolicyForAllAction("admin", "object", enforcer)
	// ticket_category
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "ticket_category", enforcer)
	// saved, err = addPolicyForAllAction("admin", "ticket_category", enforcer)
	// ticket_priority
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "ticket_priority", enforcer)
	// saved, err = addPolicyForAllAction("admin", "ticket_priority", enforcer)
	// ticket_status
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "ticket_status", enforcer)
	// saved, err = addPolicyForAllAction("admin", "ticket_status", enforcer)
	// ticket_source
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "ticket_source", enforcer)
	// saved, err = addPolicyForAllAction("admin", "ticket_source", enforcer)
	// ticket_sla
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "ticket_sla", enforcer)
	// saved, err = addPolicyForAllAction("admin", "ticket_sla", enforcer)
	// ticket
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "ticket", enforcer)
	// saved, err = addPolicyForAllAction("admin", "ticket", enforcer)
	// note
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "note", enforcer)
	// saved, err = addPolicyForAllAction("admin", "note", enforcer)
	// ticket_cause
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "ticket_cause", enforcer)
	// saved, err = addPolicyForAllAction("admin", "ticket_cause", enforcer)
	// closing_remark
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "closing_remark", enforcer)
	saved, err = addPolicyForAllAction("64e5b10d-7d23-4ee5-b386-8c65a99bbb78", "contact", enforcer)
	// saved, err = addPolicyForAllAction("admin", "closing_remark", enforcer)
	// closed_ticket
	saved, err = addPolicyForAllAction("3b59d805-f72b-4704-b97c-6a3a53de7e81", "closed_ticket", enforcer)
	// saved, err = addPolicyForAllAction("admin", "closed_ticket", enforcer)

	
	
	return
}
func addPolicyForAllAction(role, object string, enforcer *casbin.CachedEnforcer) (saved bool, err error) {

	saved, err = addPolictForAction("create", role, object, enforcer)
	saved, err = addPolictForAction("view", role, object, enforcer)
	saved, err = addPolictForAction("list", role, object, enforcer)
	saved, err = addPolictForAction("update", role, object, enforcer)
	saved, err = addPolictForAction("delete", role, object, enforcer)
	return
}

func addPolictForAction(action, role, object string, enforcer *casbin.CachedEnforcer) (saved bool, err error) {

	return enforcer.AddPolicy(role, object, action)
}

func removePolicyForAllAction(role, object string, enforcer *casbin.CachedEnforcer) (saved bool, err error) {

	saved, err = removePolictForAction("create", role, object, enforcer)
	saved, err = removePolictForAction("view", role, object, enforcer)
	saved, err = removePolictForAction("list", role, object, enforcer)
	saved, err = removePolictForAction("update", role, object, enforcer)
	saved, err = removePolictForAction("delete", role, object, enforcer)
	return
}

func removePolictForAction(action, role, object string, enforcer *casbin.CachedEnforcer) (saved bool, err error) {

	return enforcer.RemovePolicy(role, object, action)
}
