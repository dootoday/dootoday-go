package gorm

import (
	"apidootoday/config"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/golang/glog"
	
	// for the driver I guess
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	// DB is the connection handle for the db
	DB *gorm.DB
)

// InitDB : This Initializes the first db, and exports it to be passed around
func InitDB() (*gorm.DB, error) {
	// attempt to open a new connection to the db
	glog.Info("Opening a new connection to the db...")
	connStr := fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", 
		config.DbUsername, config.DbPassword, config.DbHostName, config.DbName,
	)
	db, err := gorm.Open(config.DbDriver, connStr);
	if err != nil {
		return db, err
	}
	defer db.Close()
	return db, err
}
