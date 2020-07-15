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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Run is the function to run the gin server
func (g *GinService) Run() {
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	}
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

		v1.GET("/plans",
			g.AuthHandler.AuthMiddleware,
			g.AuthHandler.GetPlans,
		)

		v1.GET("/user",
			g.AuthHandler.AuthMiddleware,
			g.AuthHandler.GetUser,
		)

		v1.POST("/task",
			g.AuthHandler.AuthMiddleware,
			g.TaskHandler.CreateTask,
		)

		v1.POST("/task/:task_id",
			g.AuthHandler.AuthMiddleware,
			g.TaskHandler.UpdateTask,
		)

		v1.GET("/task/:task_id",
			g.AuthHandler.AuthMiddleware,
			g.TaskHandler.GetTask,
		)

		v1.DELETE("/task/:task_id",
			g.AuthHandler.AuthMiddleware,
			g.TaskHandler.DeleteTask,
		)

		v1.GET("/tasks",
			g.AuthHandler.AuthMiddleware,
			g.TaskHandler.GetTasks,
		)

		v1.POST("/column",
			g.AuthHandler.AuthMiddleware,
			g.TaskHandler.CreateColumn,
		)

		v1.POST("/column/:col_id",
			g.AuthHandler.AuthMiddleware,
			g.TaskHandler.UpdateColumn,
		)

		v1.DELETE("/column/:col_id",
			g.AuthHandler.AuthMiddleware,
			g.TaskHandler.DeleteColumn,
		)

		v1.GET("/column/:col_id",
			g.AuthHandler.AuthMiddleware,
			g.TaskHandler.GetColumn,
		)

		v1.GET("/columns",
			g.AuthHandler.AuthMiddleware,
			g.TaskHandler.GetColumns,
		)

		v1.POST("/repos",
			g.AuthHandler.AuthMiddleware,
			g.TaskHandler.ReposTask,
		)
	}
	r.Run(fmt.Sprintf(":%d", config.ServerPort)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
