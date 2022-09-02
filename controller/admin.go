package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"srun/logic"
)

// CheckLicense license 授权验证
func CheckLicense(c *gin.Context) {
	//
}

// 后端登录

// 配置项管理

func AdminLogin() {

}

func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(logic.IdentityKey)
	c.JSON(200, gin.H{
		"userID":   claims[logic.IdentityKey],
		"userName": user.(*logic.User).UserName,
		"text":     "Hello World.",
	})
}
