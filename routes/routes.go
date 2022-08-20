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

	r.StaticFile("/", "index.html")

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/demo", controller.Demo)

	r.GET("/register/begin", controller.Begin)
	r.POST("/register/finish", controller.Finish)
	r.GET("/login/begin", controller.LoginBegin)
	r.POST("/login/finish", controller.LoginFinish)

	return r
}
