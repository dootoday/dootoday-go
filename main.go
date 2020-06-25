package main

import (
	"apidootoday/gin"
	"apidootoday/gorm"

	"github.com/golang/glog"
)

func main() {
	_, err := gorm.InitDB()
	if err != nil {
		glog.Fatal("Having trouble connecting the database", err)
	}

	gin.Run()
}
