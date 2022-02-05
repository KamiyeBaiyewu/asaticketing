package env

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-pg/pg/v9"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/config"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/lib/email"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/lib/enforcer"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/storage"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/storage/database"
	"github.com/sirupsen/logrus"
)

// get environment Variable
var ()

// Env - holds utilities needed by routes
type Env struct {
	DB       database.Database
	Storage  storage.Storage
	Config   *config.Info
	Enforcer *casbin.CachedEnforcer
	inboundMail *email.InboundMail

	/* MISC */
	casbinDB *pg.DB
}

// Boot - Initializes environment variables using environent variables
func Boot() *Env {

	// Load the viper configrations
	cfg := config.LoadConfig()

	// Connecct to the Database
	db, err := database.New(cfg)
	if err != nil {
		logrus.Fatal(err.Error())
	}

/* 	inboundMail , err := email.Init(db)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	go inboundMail.Listen() */
	// Initialize the Enforcer => casbin
	enforcer, casbinDB := enforcer.Init(cfg)

	env := &Env{
		DB:       db,
		Enforcer: enforcer,
		casbinDB: casbinDB,
		Config:   cfg,
		//inboundMail: inboundMail,
	}

	return env
}

// Close close all the necessary
func (e *Env) Close() {
	//e.Storage.Close()
	e.casbinDB.Close()
	e.inboundMail.Close()
}

// ReloadPolicies - reloads all the casbin policies stored into memory
func (e *Env) ReloadPolicies() {
	e.Enforcer.InvalidateCache()
	e.Enforcer.LoadPolicy()
}

// AddPolicy - Adds policy to cassbin cached policies and reloads all the policies
func (e *Env) AddPolicy(role string, object string, action string) (saved bool, err error) {
	saved, err = e.Enforcer.AddPolicy(role, object, action)
	if err != nil {
		return
	}
	return
}

// RemovePolicy - removes policy from cassbin cached policies and reloads all the policies
func (e *Env) RemovePolicy(role string, object string, action string) (removed bool, err error) {
	removed, err = e.Enforcer.RemovePolicy(role, object, action)
	if err != nil {
		return
	}
	return
}
