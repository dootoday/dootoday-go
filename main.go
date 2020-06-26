package main

import (
	"apidootoday/gin"
	"apidootoday/gorm"
	userservice "apidootoday/user"

	"github.com/golang/glog"
)

func main() {
	db, err := gorm.InitDB()
	if err != nil {
		glog.Fatal("Having trouble connecting the database", err)
	}

	us := userservice.NewUserService(db)
	err = us.Migrate()
	if err != nil {
		glog.Fatal("Having some problem with migration", err)
	}
	defer db.Close()
	gin.Run()
}
