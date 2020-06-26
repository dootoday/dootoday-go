package service

import (
	"errors"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

// UserService : this is the user service
// point of interaction for the outside world
type UserService struct {
	DB           *gorm.DB
	jwtService   IJWTService
	gauthservice IGoogleService
}

// NewUserService : is an instantiator for user service
func NewUserService(
	db *gorm.DB,
	jwt IJWTService,
	gauth IGoogleService,
) *UserService {
	return &UserService{
		DB:           db,
		jwtService:   jwt,
		gauthservice: gauth,
	}
}

var (
	// EmailNotVerifiedError :
	EmailNotVerifiedError = errors.New("Email not verified")
)

// Login :
func (us *UserService) Login(idToken string) (uint, error) {
	info, err := us.gauthservice.VerifyIDToken(idToken)
	if err != nil {
		return uint(0), err
	}
	if !info.VerifiedEmail {
		return uint(0), EmailNotVerifiedError
	}
	fname, _ := us.jwtService.GetInfoFromToken(
		idToken, "given_name",
	)
	lname, _ := us.jwtService.GetInfoFromToken(
		idToken, "family_name",
	)
	avatar, _ := us.jwtService.GetInfoFromToken(
		idToken, "picture",
	)
	user := User{
		FirstName: fname,
		Email:     info.Email,
		GoogleID:  info.UserId,
		LastName:  lname,
		Avatar:    avatar,
	}
	userExists, err := us.UserExists(&user)
	if err != nil {
		return uint(0), err
	}
	if !userExists {
		glog.Info("User does not exists")
		// Creating a new user
		err := us.Create(&user)
		if err != nil {
			return uint(0), err
		}
		glog.Info("New user created for ", user.Email)
		// TODO : create initial task lists and tasks
	}
	glog.Info("User login ", user.Email)
	return user.ID, nil
}
