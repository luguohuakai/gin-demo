package logic

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"srun/dao/mysql"
	"srun/model"
	"time"
)

// 登录表单
type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var IdentityKey = "username"

func JWT() *jwt.GinJWTMiddleware {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "srun webauthn admin",
		Key:         []byte("srunsoft"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.Admin); ok {
				return jwt.MapClaims{
					IdentityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			//fmt.Println(fmt.Sprintf("%#v", claims))
			if v, ok := claims[IdentityKey].(string); ok {
				return &model.Admin{
					Username: v,
				}
			} else {
				fmt.Println("token error")
				zap.L().Error("token error")
				return nil
			}
		},
		// 验证器 用户身份验证
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var data login
			if err := c.ShouldBindJSON(&data); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			// : 校验数据库中的账密
			var admin model.Admin
			if err := mysql.GetDB().First(&admin, "username = ?", data.Username).Error; err != nil {
				return nil, err
			}
			if SHA1(data.Password) == admin.Password {
				return &admin, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		// 配置授权人 仅在身份验证成功后调用
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*model.Admin); ok {
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
