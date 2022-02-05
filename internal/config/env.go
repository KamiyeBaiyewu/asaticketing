package config

import (
	"github.com/spf13/viper"
)

// LoadEnvConfig - helps to load the configuration from enviroment varaibles and set the default values
func loadEnvConfig() (vCfg *viper.Viper) {
	vCfg = viper.New()

	vCfg.BindEnv("app_version", "APP_VERSION")
	vCfg.SetDefault("app_version", "vHEAD")
	vCfg.BindEnv("data_directory", "DATA_DIRECTORY")
	vCfg.SetDefault("data_directory", "")
	vCfg.BindEnv("http_addr", "HTTP_ADDR")
	vCfg.SetDefault("http_addr", ":8000")
	

	// Database parameters
	vCfg.BindEnv("db.driver", "DB_DRIVER")
	vCfg.SetDefault("db.driver", "postgres")
	vCfg.BindEnv("db.user", "DB_USER")
	vCfg.SetDefault("db.user", "postgres")
	vCfg.BindEnv("db.secret", "DB_SECRET")
	vCfg.SetDefault("db.secret", "Peaceg419")
	vCfg.BindEnv("db.host", "DB_HOST")
	vCfg.SetDefault("db.host", "db")
	vCfg.BindEnv("db.port", "DB_PORT")
	vCfg.SetDefault("db.port", 5432)
	vCfg.BindEnv("db.name", "DB_NAME")
	vCfg.SetDefault("db.name", "asa_tickets_app")
	vCfg.BindEnv("db.param", "DB_PARAM")
	vCfg.SetDefault("db.param", "sslmode=disable")
	vCfg.BindEnv("db.timeout", "DB_TIMEOUT")
	vCfg.SetDefault("db.timeout", 2000)
	vCfg.BindEnv("db.url", "DB_URL")
	// vCfg.SetDefault("db.url", "postgres://postgres:Peaceg419@db:5432/asa_tickets_app?sslmode=disable")

	// Load the enviroment variables
	vCfg.BindEnv("casbin.model", "CASBIN_MODEL")
	vCfg.SetDefault("casbin.model", "./config/model.conf")
	vCfg.BindEnv("casbin.table", "CASBIN_TABLE")
	vCfg.SetDefault("casbin.table", "casbin_rules")
	vCfg.BindEnv("casbin.policy", "CASBIN_POLICY")
	vCfg.SetDefault("casbin.policy", "./config/policy.csv")

	// authorizer
	vCfg.BindEnv("authorizer.cache_exp", "AUTHORIZER_CACHE_EXP")
	vCfg.SetDefault("authorizer.cache_exp", 180)
	

	return
}
