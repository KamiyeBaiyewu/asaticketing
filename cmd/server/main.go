package main

import (
	"net/http"

	"github.com/lilkid3/ASA-Ticket/Backend/internal/api"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/config"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/env"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load the Environment Varaiables
	//flag.Parse()

	// Load the viper configrations
	//cfg := config.LoadConfig()
	logrus.SetLevel(logrus.DebugLevel)

	env := env.Boot()

	logrus.WithField("version", config.Version).Debug("Starting Server")

	// Close all Connections
	defer env.Close()

	// Create a new Router
	handler, err := api.NewRouter(env)
	if err != nil {
		logrus.WithError(err).Fatal("Error while building Router")
	}

	server := http.Server{
		Handler: handler,
		Addr:    env.Config.HTTPAddr,
	}
	logrus.WithFields(logrus.Fields{
		"HTTP Address": env.Config.HTTPAddr,
	}).Info("Starting Web Server")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Error("Server failed")
	}

}
