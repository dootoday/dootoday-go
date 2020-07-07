package main

import (
	ginservice "apidootoday/gin"
	gauthservice "apidootoday/googleauth"
	"apidootoday/gorm"
	jwtservice "apidootoday/jwtservice"
	orderservice "apidootoday/orderservice"
	subscriptionservice "apidootoday/subscription"
	ts "apidootoday/taskservice"
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
	taskdbservice := ts.NewTaskDBService(db)
	taskservice := ts.NewTaskService(taskdbservice)

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
	err = taskdbservice.Migrate()
	if err != nil {
		glog.Fatal("Having some problem with migrating tasks", err)
	}

	authHandlers := ginservice.NewAuthHandler(
		us, tokenService, gauthService, subscription,
	)
	taskHandlers := ginservice.NewTaskHandler(taskservice)
	ginService := ginservice.NewGinService(authHandlers, taskHandlers)
	// Run gin
	ginService.Run()

	defer db.Close()
}
