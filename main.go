package main

import (
	"apidootoday/config"
	cs "apidootoday/cronservice"
	es "apidootoday/emailservice"
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
	"flag"
	"fmt"

	"github.com/golang/glog"
)

func main() {
	noRunServer := false
	flag.BoolVar(&noRunServer, "no-run-server", false, "Don't startup the server.")

	noMigration := false
	flag.BoolVar(&noMigration, "no-migration", false, "Don't start migration.")

	runBatch := ""
	flag.StringVar(&runBatch, "batch", "", "Specify a batch to run.")

	flag.Parse()

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
	emailService := es.NewEmailService()
	cronService := cs.NewCronService(us, taskservice, emailService)

	if !noMigration {
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
	}

	redisClient := rdclient.NewRedisClient(
		context.Background(),
		config.RedisHost,
		config.RedisPort,
		config.RedisPass,
		0,
		string(config.Environment),
	)

	// Setup handlers
	authHandlers := ginservice.NewAuthHandler(
		us, tokenService, gauthService, taskservice, subscription, emailService,
	)
	userHandler := ginservice.NewUserHandler(us)
	taskHandlers := ginservice.NewTaskHandler(taskservice, recurringtaskservice, redisClient)
	subscriptionHandler := ginservice.NewSubscriptionHandler(subscription, order, us)

	// Run batches
	if runBatch != "" {
		switch runBatch {
		case "move-tasks-to-today":
			err := cronService.MoveTasksToTodayCron()
			if err != nil {
				glog.Error(err)
			}
		case "morning-email-reminder":
			err = cronService.DailyMorningEmailCron()
			if err != nil {
				glog.Error(err)
			}
		case "send-email-test":
			err := emailService.SendWelcomeEmail(
				"sanborn.sen@gmail.com",
				"Sudipta Sen",
				"Sanborn",
			)
			if err != nil {
				glog.Error(err)
			}
		default:
			glog.Fatal(fmt.Sprintf("Unrecognized batch: %s", runBatch))
		}
	}

	glog.Info(noRunServer)

	// Run server
	if !noRunServer {
		ginService := ginservice.NewGinService(
			authHandlers, taskHandlers,
			subscriptionHandler, userHandler,
		)
		// Run gin
		ginService.Run()
	}

	defer db.Close()
}
