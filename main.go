package main

import (
	"apidootoday/gin"
	"apidootoday/gorm"
	jwtservice "apidootoday/jwtservice"
	userservice "apidootoday/user"
	"fmt"
	"time"

	"github.com/golang/glog"
)

func main() {
	t := jwtservice.GetRefreshToken(uint(1))
	ts := jwtservice.NewTokenService(t, jwtservice.RefreshTokenType)
	fmt.Println(ts.GetUserID())
	fmt.Println(ts.IsTokenValid())
	time.Sleep(time.Second * 11)
	fmt.Println(ts.IsTokenValid())
}

func poxyMain() {
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
