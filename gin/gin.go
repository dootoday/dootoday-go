package gin

import (
	"apidootoday/config"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Run is the function to run the gin server
func Run() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(fmt.Sprintf(":%d", config.ServerPort)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
