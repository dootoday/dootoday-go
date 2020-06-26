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
	GoogleID  string `gorm:"index:googleid"`
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

// Create :
func (us *UserService) Create(user *User) error {
	glog.Info("Creating users table")
	err := us.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

// UserExists :
func (us *UserService) UserExists(user *User) (bool, error) {
	err := us.DB.Where("google_id=?", user.GoogleID).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
