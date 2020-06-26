package service

import (
	"apidootoday/config"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinService :
type GinService struct {
	AuthHandler *AuthHandler
}

// NewGinService :
func NewGinService(authHandler *AuthHandler) *GinService {
	return &GinService{
		AuthHandler: authHandler,
	}
}

// Run is the function to run the gin server
func (g *GinService) Run() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/login", g.AuthHandler.Login)
	r.POST("/refresh", g.AuthHandler.Refresh)

	r.Run(fmt.Sprintf(":%d", config.ServerPort)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
