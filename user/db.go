package service

import (
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

// User : This is the user model
type User struct {
	gorm.Model

	FirstName string `gorm:"type:varchar(100);"`
	LastName  string `gorm:"type:varchar(100);"`
	Email     string `gorm:"type:varchar(100);unique_index"`
	Avatar    string `gorm:"type:text"`
}

// Migrate : This is the db migrate function for
// Users
func (us *UserService) Migrate() error {
	glog.Info("Creating users table")
	err := us.DB.AutoMigrate(&User{}).Error
	if err != nil {
		glog.Info(err)
	}
	return nil
}
