// Package info reads the application settings.
package config

// *****************************************************************************
// Application Settings
// *****************************************************************************

// LoadJSONConfig reads the configuration from viper.
func LoadConfig() *Info {
	// Create a new configuration with the path to the file

	// we would use the ENV varaibles
	vCfg := loadEnvConfig()

	// Set the app version
	//Version = vCfg.GetString("app_version")
	config := &Info{
		AppVersion:    vCfg.GetString("app_version"),
		DataDirectory: vCfg.GetString("data_directory"),
		HTTPAddr:      vCfg.GetString("http_addr"),
		Database: database{
			Driver:    vCfg.GetString("db.driver"),
			Username:  vCfg.GetString("db.user"),
			Password:  vCfg.GetString("db.secret"),
			Database:  vCfg.GetString("db.name"),
			Hostname:  vCfg.GetString("db.host"),
			Port:      vCfg.GetInt("db.port"),
			Parameter: vCfg.GetString("db.param"),
			Timeout:   vCfg.GetInt64("db.timeout"),
			URL:       vCfg.GetString("db.url"),
		},
		Casbin: casbin{
			Model:  vCfg.GetString("casbin.model"),
			Policy: vCfg.GetString("casbin.policy"),
			Table:  vCfg.GetString("casbin.table"),
		},
		Authorizer: authorizer{
			CacheExpiration: vCfg.GetInt("authorizer.cache_exp"),
		},
	}

	// log.Printf("Config => %+v\n\n", config)
	// Return the configuration
	return config
}
