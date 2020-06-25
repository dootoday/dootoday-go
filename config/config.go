package config

import (
	"os"
	"strconv"

	"github.com/golang/glog"
	"github.com/spf13/viper"
)

func init() {
	SetSettingsFromViper()
}

var (
	// ServerPort : port to run gin local server on
	ServerPort int

	// ServerHostName : Hostname to run this server on
	ServerHostName string

	// Debug : Debug mode true / false
	Debug bool

	// DbDriver : Driver for the db
	DbDriver string

	// DbUsername : username for the db
	DbUsername string

	// DbPassword : password for the dn
	DbPassword string

	// DbHostName : host for the db
	DbHostName string

	// DbName : name of the db
	DbName string

	// DbPort : port for the db
	DbPort string

	// Environment : dev environment, production, docker, etc
	Environment AppEnvironment

	// AppEnvironments : array of all app environments
	AppEnvironments = []AppEnvironment{
		AppEnvironmentTesting,
		AppEnvironmentLocal,
		AppEnvironmentStaging,
		AppEnvironmentProduction,
	}
)

// AppEnvironment : string wrapper for environment name
type AppEnvironment string

const (
	// AppEnvironmentTesting : testing env
	AppEnvironmentTesting = AppEnvironment("testing")
	// AppEnvironmentLocal :
	AppEnvironmentLocal = AppEnvironment("local")
	// AppEnvironmentStaging :
	AppEnvironmentStaging = AppEnvironment("staging")
	// AppEnvironmentProduction :
	AppEnvironmentProduction = AppEnvironment("production")
	// AppEnvironmentJenkins : is the jenkins environment
	AppEnvironmentJenkins = AppEnvironment("jenkins")
)

func getEnvironment() AppEnvironment {
	hostEnvironment := os.Getenv("ENVIRONMENT")
	for _, env := range AppEnvironments {
		if env == AppEnvironment(hostEnvironment) {
			Environment = env
			return env
		}
	}

	// set to local config if environment not found
	return AppEnvironmentLocal
}

// SetSettingsFromViper : sets global settings using viper
func SetSettingsFromViper() {
	Environment = getEnvironment()
	glog.Info("We're in our the following environment: ", Environment)

	// SetENV if not in a production environment
	// Check for local
	if Environment != AppEnvironmentProduction && Environment != AppEnvironmentStaging {
		setEnvironmentVariablesFromConfig(Environment)
	}

	if Environment == AppEnvironmentTesting {
		DbName = os.Getenv("TEST_DB_NAME")
	} else {
		DbName = os.Getenv("DB_NAME")
	}
	DbDriver = os.Getenv("DB_DRIVER")
	DbHostName = os.Getenv("DB_HOSTNAME")
	DbUsername = os.Getenv("DB_USERNAME")
	DbPort = os.Getenv("DB_PORT")
	DbPassword = os.Getenv("DB_PASSWORD")
	glog.Info("Db settings: ", DbDriver, " ", DbHostName, " ", DbName)

	Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	ServerHostName = os.Getenv("SERVER_HOSTNAME")
	ServerPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
}

func setEnvironmentVariablesFromConfig(env AppEnvironment) {
	// get and set basePath of project
	baseProjectPath := "/Users/sudiptasen/personal/apidootoday"
	viper.AddConfigPath(baseProjectPath + "/config/")
	viper.SetConfigType("yaml")
	if env == AppEnvironmentLocal || env == AppEnvironmentTesting {
		glog.Info("Reading configuration from localConfig")
		viper.SetConfigName("localConfig")
	}

	if env == AppEnvironmentJenkins {
		glog.Info("Reading from jenkins settings")
		viper.SetConfigName("jenkinsConfig")
	}

	err := viper.ReadInConfig()
	if err != nil {
		glog.Info("Failed reading local settings: ", err)
	}
	debug := viper.GetBool("debug")

	serverHostName := viper.GetString("serverHostName")
	serverPort := viper.GetString("serverPort")
	dbDriver := viper.GetString("dbDriver")
	dbHostname := viper.GetString("dbHostName")
	dbPassword := viper.GetString("dbPassword")
	dbPort := viper.GetString("dbPort")
	dbUser := viper.GetString("dbUsername")
	dbName := viper.GetString("dbName")
	dbTestDBName := viper.GetString("testDbName")

	// Set the OS Environment variables
	os.Setenv("DB_DRIVER", dbDriver)
	os.Setenv("DB_HOSTNAME", dbHostname)
	os.Setenv("DB_USERNAME", dbUser)
	os.Setenv("DB_PORT", dbPort)
	os.Setenv("DB_NAME", dbName)
	os.Setenv("DB_PASSWORD", dbPassword)
	os.Setenv("TEST_DB_NAME", dbTestDBName)
	os.Setenv("DEBUG", strconv.FormatBool(debug))
	os.Setenv("SERVER_HOSTNAME", serverHostName)
	os.Setenv("SERVER_PORT", serverPort)
	glog.Info("setEnvironmentVariablesFromConfig: Config finished reading in settings from file.")
}
