package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"srun/controller"
	"srun/logger"
	"srun/logic"
	"time"
)

func Setup() *gin.Engine {
	if viper.GetString("app.mode") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// GetLimiters 限制GET请求次数
	GET := tollbooth.NewLimiter(100, &limiter.ExpirableOptions{
		ExpireJobInterval: time.Second,
	})
	// POSTLimiters 限制POST请求次数
	POST := tollbooth.NewLimiter(100, &limiter.ExpirableOptions{
		ExpireJobInterval: time.Second,
	})

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.Use(cors.Default())

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

	r.POST("/register/begin", controller.Limit(POST), controller.Begin)
	r.POST("/register/finish", controller.Finish)
	r.GET("/register/user-exists", controller.Limit(GET), controller.UserExists)
	r.GET("/login/begin", controller.LoginBegin)
	r.POST("/login/finish", controller.LoginFinish)

	r.GET("/admin/login", authMiddleware.LoginHandler)

	admin := r.Group("admin")
	admin.GET("/test", controller.Test)
	admin.GET("/all-cfg", controller.AllCfg)
	admin.POST("/set-login-trans", controller.SetLoginTransports)
	admin.Use(authMiddleware.MiddlewareFunc())

	admin.GET("/refresh_token", authMiddleware.RefreshHandler)
	admin.GET("/logout", authMiddleware.LogoutHandler)
	admin.GET("/hello", controller.HelloHandler)

	return r
}
