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
	// ErrEmailNotVerified :
	ErrEmailNotVerified = errors.New("Email not verified")
)

// Login : returns userID, isNewUser, error
func (us *UserService) Login(idToken string) (uint, bool, error) {
	info, err := us.gauthservice.VerifyIDToken(idToken)
	if err != nil {
		return uint(0), false, err
	}
	if !info.VerifiedEmail {
		return uint(0), false, ErrEmailNotVerified
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
		return uint(0), false, err
	}
	isNewUser := false
	if !userExists {
		glog.Info("User does not exists")
		// Creating a new user
		err := us.Create(&user)
		if err != nil {
			return uint(0), false, err
		}
		glog.Info("New user created for ", user.Email)
		isNewUser = true
	}
	glog.Info("User login ", user.Email)
	return user.ID, isNewUser, nil
}

// UpdateUserTimeZoneOffset :
func (us *UserService) UpdateUserTimeZoneOffset(userID uint, offset int) error {
	user, err := us.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.TimeZoneOffset = offset
	err = us.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAutoTaskMove :
func (us *UserService) UpdateAutoTaskMove(userID uint, allow bool) error {
	user, err := us.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.AllowAutoUpdate = allow
	err = us.UpdateAllowAutoTaskMove(user)
	if err != nil {
		return err
	}
	return nil
}
