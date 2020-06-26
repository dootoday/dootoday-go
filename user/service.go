package service

import "github.com/jinzhu/gorm"

// UserService : this is the user service
// point of interaction for the outside world
type UserService struct {
	DB *gorm.DB
}

// NewUserService : is an instantiator for user service
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB: db,
	}
}
