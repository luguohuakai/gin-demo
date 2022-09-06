package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"srun/controller"
	"srun/logger"
	"srun/logic"
)

func Setup() *gin.Engine {
	if viper.GetString("app.mode") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	authMiddleware := logic.JWT()

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.StaticFile("/", "index.html")
	r.Static("/js", "js")

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/demo", controller.Demo)

	r.POST("/register/begin", controller.Begin)
	r.POST("/register/finish", controller.Finish)
	r.GET("/register/user-exists", controller.UserExists)
	r.GET("/login/begin", controller.LoginBegin)
	r.POST("/login/finish", controller.LoginFinish)

	r.GET("/admin/login", authMiddleware.LoginHandler)

	admin := r.Group("admin")
	admin.Use(authMiddleware.MiddlewareFunc())

	admin.GET("/refresh_token", authMiddleware.RefreshHandler)
	admin.GET("/logout", authMiddleware.LogoutHandler)
	admin.GET("/hello", controller.HelloHandler)

	return r
}
