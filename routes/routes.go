package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"srun/controller"
	"srun/logger"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/demo", controller.Demo)

	return r
}
