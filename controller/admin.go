package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"srun/cfg"
	"srun/logic"
)

// CheckLicense license 授权验证
//func CheckLicense(c *gin.Context) {
//	//
//}

// AllCfg 获取全部配置
func AllCfg(c *gin.Context) {
	success(c, cfg.FD)
}

// SetLoginTransports 设置transports (usb nfc ble internal)
func SetLoginTransports(c *gin.Context) {
	var ec cfg.ExcludeCredentials
	var err error
	err = c.ShouldBindJSON(&ec)
	if err != nil {
		fail(c, err)
		return
	}
	cfg.VP.Set("register.exclude_credentials.transports", ec.Transports)
	err = cfg.VP.WriteConfig()
	if err != nil {
		fail(c, err)
		return
	}
	success(c)
}

type User struct {
	Name  []string `json:"name,omitempty" binding:"required,max=6,inArray=11 22,omitempty"`
	Email string   `json:"email" binding:"email"`
}

func Test(c *gin.Context) {
	var u User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		fail(c, err)
		return
	}
	success(c)
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
