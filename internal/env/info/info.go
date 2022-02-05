// Package info reads the application settings.
package info

import (
	"encoding/json"

	"github.com/lilkid3/ASA-Ticket/Backend/internal/env/jsonconfig"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/env/storage"
)

// *****************************************************************************
// Application Settings
// *****************************************************************************

// Info structures the application settings.
type Info struct {
	Storage storage.Info `json:"PostgreSQL"`
	path    string
}

// Path returns the env.json path
func (c *Info) Path() string {

	return c.path
}

// ParseJSON unmarshals bytes to structs
func (c *Info) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

// New returns a instance of the application settings.
func New(path string) *Info {
	return &Info{
		path: path,
	}
}

// LoadJSONConfig reads the JSON configuration file.
func LoadJSONConfig(configFile string) (*Info, error) {
	// Create a new configuration with the path to the file
	config := New(configFile)

	// Load the configuration file
	err := jsonconfig.Load(configFile, config)

	// Return the configuration
	return config, err
}

// Close close all the necessary
func (e *Info) Close() {
	// e.DB.Close()
	// e.jaegerCloser.Close()
	// e.SDClient.SDClient.Deregister()
}
