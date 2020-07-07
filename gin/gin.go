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
	TaskHandler *TaskHandler
}

// NewGinService :
func NewGinService(
	authHandler *AuthHandler,
	taskHandler *TaskHandler,
) *GinService {
	return &GinService{
		AuthHandler: authHandler,
		TaskHandler: taskHandler,
	}
}

// CORSMiddleware :
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Run is the function to run the gin server
func (g *GinService) Run() {
	r := gin.Default()
	r.Use(CORSMiddleware())
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

		v1.GET("/user",
			g.AuthHandler.AuthMiddleware,
			g.AuthHandler.GetUser,
		)

		v1.POST("/createtask",
			g.AuthHandler.AuthMiddleware,
			g.TaskHandler.CreateTask,
		)
	}
	r.Run(fmt.Sprintf(":%d", config.ServerPort)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
