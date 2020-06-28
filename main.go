package main

import (
	// ginservice "apidootoday/gin"
	gauthservice "apidootoday/googleauth"
	"apidootoday/gorm"
	jwtservice "apidootoday/jwtservice"
	orderservice "apidootoday/orderservice"
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
	us := userservice.NewUserService(db, tokenService, gauthService)
	subscription := subscriptionservice.NewSubscriptionService(db)
	order := orderservice.NewOrderService(db)

	// Table migrations
	err = us.Migrate()
	if err != nil {
		glog.Fatal("Having some problem with migration", err)
	}
	err = subscription.Migrate()
	if err != nil {
		glog.Fatal("Having some problem with migration", err)
	}
	err = order.Migrate()
	if err != nil {
		glog.Fatal("Having some problem with migrating orders", err)
	}

	order.CreateNewOrder(2, 1, 2000)

	// authHandlers := ginservice.NewAuthHandler(
	// 	us, tokenService, gauthService, subscription,
	// )
	// ginService := ginservice.NewGinService(authHandlers)
	// // Run gin
	// ginService.Run()

	defer db.Close()
}
