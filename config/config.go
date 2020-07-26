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
	// AccessTokenSecret : This is the jwt secret for access token
	AccessTokenSecret string

	// RefreshTokenSecret : This is the jwt secret for refresh token
	RefreshTokenSecret string

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

	// RPApiKey : API Key for RazorPay
	RPApiKey string

	// RPApiSecret : API Key for RazorPay
	RPApiSecret string

	// FrontendBase : This is the front-end base
	FrontendBase string

	// BackendBase : This is the back-end base
	BackendBase string

	// DooTodayLogo : This is the logo URL
	DooTodayLogo string

	// DooTodayName : This is the name
	DooTodayName string

	// DooTodayDesc : This is the description
	DooTodayDesc string

	// RedisHost : This is the Redis host
	RedisHost string

	// RedisPort : This is the RedisPort
	RedisPort string

	// RedisPass : This is the RedisPass
	RedisPass string

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

	RPApiKey = os.Getenv("RP_API_KEY")
	RPApiSecret = os.Getenv("RP_API_SECRET")

	AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	RefreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")

	FrontendBase = os.Getenv("FRONT_END_BASE")
	BackendBase = os.Getenv("BACK_END_BASE")
	DooTodayLogo = os.Getenv("DOO_TODAY_LOGO")
	DooTodayName = os.Getenv("DOO_TODAY_NAME")
	DooTodayDesc = os.Getenv("DOO_TODAY_DESC")

	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	RedisPass = os.Getenv("REDIS_PASSWORD")

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

	accessTokenSecret := viper.GetString("accessTokenSecret")
	refreshTokenSecret := viper.GetString("refreshTokenSecret")
	serverHostName := viper.GetString("serverHostName")
	serverPort := viper.GetString("serverPort")
	dbDriver := viper.GetString("dbDriver")
	dbHostname := viper.GetString("dbHostName")
	dbPassword := viper.GetString("dbPassword")
	dbPort := viper.GetString("dbPort")
	dbUser := viper.GetString("dbUsername")
	dbName := viper.GetString("dbName")
	dbTestDBName := viper.GetString("testDbName")
	rpAPIKey := viper.GetString("razorpayKey")
	rpAPISecret := viper.GetString("razorpaySecret")
	frontEndBase := viper.GetString("frontEndBase")
	backeEndBase := viper.GetString("backeEndBase")
	dooTodayLogo := viper.GetString("dooTodayLogo")
	dooTodayName := viper.GetString("dooTodayName")
	dooTodayDesc := viper.GetString("dooTodayDesc")
	redisHost := viper.GetString("redisHost")
	redisPort := viper.GetString("redisPort")
	redisPassword := viper.GetString("redisPassword")

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
	os.Setenv("RP_API_KEY", rpAPIKey)
	os.Setenv("RP_API_SECRET", rpAPISecret)
	os.Setenv("ACCESS_TOKEN_SECRET", accessTokenSecret)
	os.Setenv("REFRESH_TOKEN_SECRET", refreshTokenSecret)
	os.Setenv("FRONT_END_BASE", frontEndBase)
	os.Setenv("BACK_END_BASE", backeEndBase)
	os.Setenv("DOO_TODAY_LOGO", dooTodayLogo)
	os.Setenv("DOO_TODAY_NAME", dooTodayName)
	os.Setenv("DOO_TODAY_DESC", dooTodayDesc)
	os.Setenv("REDIS_HOST", redisHost)
	os.Setenv("REDIS_PORT", redisPort)
	os.Setenv("REDIS_PASSWORD", redisPassword)
	glog.Info("setEnvironmentVariablesFromConfig: Config finished reading in settings from file.")
}
