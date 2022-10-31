package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
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
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page Not Found Error"})
	})

	// 前端静态资源
	r.StaticFile("/", "/srun3/www/webauthn/index.html")
	r.Static("/js", "/srun3/www/webauthn/js")
	r.Static("/css", "/srun3/www/webauthn/css")
	r.Static("/image", "/srun3/www/webauthn/image")
	r.Static("/assets", "/srun3/www/webauthn/assets")
	r.Static("/icons", "/srun3/www/webauthn/icons")

	// webauthn注册/登录
	r.POST("/register/begin", logic.CheckLicenseMiddleware, controller.Limit(POST), controller.Begin)
	r.POST("/register/finish", controller.Finish)
	r.GET("/login/begin", logic.CheckLicenseMiddleware, controller.Limit(GET), controller.LoginBegin)
	r.POST("/login/finish", controller.LoginFinish)

	// 后台登录
	r.POST("/admin/login", authMiddleware.LoginHandler)

	// 后台接口
	admin := r.Group("admin")
	admin.Use(authMiddleware.MiddlewareFunc())
	admin.GET("/all-cfg", controller.AllCfg)
	admin.POST("/set-login-trans", controller.SetLoginTransports)
	admin.POST("/set-attestation", controller.SetAttestation)
	admin.POST("/set-timeout", controller.SetTimeout)
	admin.POST("/set-attachment", controller.SetAttachment)
	admin.POST("/set-require-resident-key", controller.SetRequireResidentKey)
	admin.POST("/set-user-verification", controller.SetUserVerification)
	admin.GET("/get-user", controller.GetUser)
	admin.DELETE("/del-user", logic.CheckLicenseMiddleware, controller.DelUser)
	admin.DELETE("/batch-del-user", logic.CheckLicenseMiddleware, controller.BatchDelUser)
	admin.GET("/get-sso", controller.GetSso)
	admin.POST("/edit-sso", controller.EditSso)
	admin.GET("/get-north", controller.GetNorth)
	admin.POST("/edit-north", controller.EditNorth)
	admin.POST("/active", controller.Limit(POST), controller.Active)
	admin.GET("/license-status", controller.LicenseStatus)
	admin.GET("/get-license", controller.GetLicense)

	admin.GET("/refresh_token", authMiddleware.RefreshHandler)
	admin.GET("/logout", authMiddleware.LogoutHandler)

	return r
}
