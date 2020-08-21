package main

import (
	ginservice "apidootoday/gin"
	gauthservice "apidootoday/googleauth"
	"apidootoday/gorm"
	jwtservice "apidootoday/jwtservice"
	orderservice "apidootoday/orderservice"
	rdclient "apidootoday/redisclient"
	subscriptionservice "apidootoday/subscription"
	ts "apidootoday/taskservice"
	userservice "apidootoday/user"
	"context"

	"apidootoday/config"

	"github.com/golang/glog"
)

func main() {
	db, err := gorm.InitDB()
	if err != nil {
		glog.Fatal("Having trouble connecting the database", err)
	}
	db.LogMode(config.Debug)
	tokenService := jwtservice.NewTokenService()
	gauthService := gauthservice.NewGoogleAuthService()
	us := userservice.NewUserService(db, tokenService, gauthService)
	subscription := subscriptionservice.NewSubscriptionService(db)
	order := orderservice.NewOrderService(db)
	taskdbservice := ts.NewTaskDBService(db)
	recurringtaskservice := ts.NewRecurringTaskService(taskdbservice)
	taskservice := ts.NewTaskService(taskdbservice, recurringtaskservice)

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
	redisClient := rdclient.NewRedisClient(
		context.Background(),
		config.RedisHost,
		config.RedisPort,
		config.RedisPass,
		0,
		string(config.Environment),
	)
	authHandlers := ginservice.NewAuthHandler(
		us, tokenService, gauthService, taskservice, subscription,
	)
	userHandler := ginservice.NewUserHandler(us)
	taskHandlers := ginservice.NewTaskHandler(taskservice, recurringtaskservice, redisClient)
	subscriptionHandler := ginservice.NewSubscriptionHandler(subscription, order, us)
	ginService := ginservice.NewGinService(
		authHandlers, taskHandlers,
		subscriptionHandler, userHandler,
	)
	// Run gin
	ginService.Run()

	defer db.Close()
}
