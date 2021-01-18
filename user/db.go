package service

import (
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

// User : This is the user model
type User struct {
	gorm.Model

	FirstName       string `gorm:"type:varchar(100);"`
	LastName        string `gorm:"type:varchar(100);"`
	Email           string `gorm:"type:varchar(100);"`
	GoogleID        string `gorm:"index:googleid"`
	Avatar          string `gorm:"type:text"`
	TimeZoneOffset  int    `gorm:"type:smallint"`
	AllowAutoUpdate bool   `gorm:"default:0"`
}

// Migrate : This is the db migrate function for
// Users
func (us *UserService) Migrate() error {
	glog.Info("Creating users table")
	err := us.DB.AutoMigrate(&User{}).Error
	if err != nil {
		glog.Info(err)
	}
	// Drop unique key on email
	err = us.DB.Exec(`ALTER TABLE users DROP INDEX uix_users_email;`).Error
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
		}
		return false, err
	}
	return true, nil
}

// GetUserByID :
func (us *UserService) GetUserByID(userID uint) (*User, error) {
	var user User
	err := us.DB.Where("id=?", userID).Find(&user).Error
	return &user, err
}

// UpdateUser :
func (us *UserService) UpdateUser(user *User) error {
	return us.DB.Model(user).Update(user).Error
}

// GetUsersByTimeZoneOffset :
func (us *UserService) GetUsersByTimeZoneOffset(offset int) ([]User, error) {
	users := []User{}
	err := us.DB.Where("time_zone_offset=?", offset).Find(&users).Error
	return users, err
}
