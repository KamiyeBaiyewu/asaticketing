package v1

import (
	"encoding/json"
	"net/http"

	"github.com/lilkid3/ASA-Ticket/Backend/internal/config"
	"github.com/sirupsen/logrus"
)

//API for returning Version

// ServerVersion represents the server version
type ServerVersion struct {
	Version string `json:"version"`
}

// VersionHandler helps to write current server version to request
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var err error
	versionJSON, err := json.Marshal(ServerVersion{
		Version: config.Version,
	})
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(versionJSON); err != nil {
		logrus.WithError(err).Debug("Error writing version")
	}

}
