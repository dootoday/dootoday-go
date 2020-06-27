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

	v1 := r.Group("/v1")
	{
		v1.POST("/login", g.AuthHandler.Login)
		v1.POST("/refresh", g.AuthHandler.Refresh)

		v1.POST("/apply-promo",
			g.AuthHandler.AuthMiddleware,
			g.AuthHandler.ApplyPromo,
		)
	}

	r.Run(fmt.Sprintf(":%d", config.ServerPort)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
