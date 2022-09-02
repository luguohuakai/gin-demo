package logic

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// 登录表单
type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var IdentityKey = "id"

// User model
// todo: 与数据库中的用户模型绑定
type User struct {
	UserName string
}

func JWT() *jwt.GinJWTMiddleware {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "srun webauthn admin",
		Key:         []byte("srunsoft"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					IdentityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[IdentityKey].(string),
			}
		},
		// 验证器 用户身份验证
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVal login
			if err := c.ShouldBind(&loginVal); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVal.Username
			password := loginVal.Password

			// todo: 校验数据库中的账密
			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &User{
					UserName: userID,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		// 配置授权人 仅在身份验证成功后调用
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.UserName == "admin" {
				return true
			}

			return false
		},
		// 配置未授权返回项
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		// 配置token获取位置
		TokenLookup: "header: Authorization, query: token, cookie: jwt", // 从这些地方去找token
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time.
		// You can override it to use another time value.
		// This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	return authMiddleware
}
