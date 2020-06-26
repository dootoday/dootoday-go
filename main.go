package main

import (
	ginservice "apidootoday/gin"
	gauthservice "apidootoday/googleauth"
	"apidootoday/gorm"
	jwtservice "apidootoday/jwtservice"
	subscriptionservice "apidootoday/subscription"
	userservice "apidootoday/user"

	"github.com/golang/glog"
)

func main() {
	db, err := gorm.InitDB()
	if err != nil {
		glog.Fatal("Having trouble connecting the database", err)
	}

	tokenService := jwtservice.NewTokenService()
	gauthService := gauthservice.NewGoogleAuthService()

	us := userservice.NewUserService(
		db, tokenService, gauthService,
	)

	subscription := subscriptionservice.NewSubscriptionService(db)

	// Table migrations
	err = us.Migrate()
	if err != nil {
		glog.Fatal("Having some problem with migration", err)
	}
	err = subscription.Migrate()
	if err != nil {
		glog.Fatal("Having some problem with migration", err)
	}

	authHandlers := ginservice.NewAuthHandler(
		us, tokenService, gauthService,
	)
	ginService := ginservice.NewGinService(authHandlers)
	// Run gin
	ginService.Run()

	defer db.Close()
}
