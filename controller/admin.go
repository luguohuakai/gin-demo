package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"srun/cfg"
	"srun/dao/mysql"
	"srun/logic"
	"srun/model"
	"strconv"
	"strings"
)

type User struct {
	Name  []string `json:"name,omitempty" binding:"required,max=6,inArray=11 22,omitempty"`
	Email string   `json:"email" binding:"email"`
}

// AllCfg 获取全部配置
func AllCfg(c *gin.Context) {
	success(c, cfg.VP.AllSettings())
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
	cfg.VP.Set("login.allow_credentials.transports", ec.Transports)
	err = cfg.VP.WriteConfig()
	if err != nil {
		fail(c, err)
		return
	}
	success(c)
}

// SetAttestation 设置 Attestation conveyance preference
func SetAttestation(c *gin.Context) {
	var r cfg.Register
	var err error
	err = c.ShouldBindJSON(&r)
	if err != nil {
		fail(c, err)
		return
	}
	cfg.VP.Set("register.attestation", r.Attestation)
	err = cfg.VP.WriteConfig()
	if err != nil {
		fail(c, err)
		return
	}
	success(c)
}

// SetTimeout 设置超时时间(毫秒)
func SetTimeout(c *gin.Context) {
	var r cfg.Register
	var err error
	err = c.ShouldBindJSON(&r)
	if err != nil {
		fail(c, err)
		return
	}
	cfg.VP.Set("register.timeout", r.Timeout)
	cfg.VP.Set("login.timeout", r.Timeout)
	err = cfg.VP.WriteConfig()
	if err != nil {
		fail(c, err)
		return
	}
	success(c)
}

// SetAttachment 设置 Authenticator selection 之 Authenticator Attachment
func SetAttachment(c *gin.Context) {
	var as cfg.AuthenticatorSelection
	var err error
	err = c.ShouldBindJSON(&as)
	if err != nil {
		fail(c, err)
		return
	}
	cfg.VP.Set("register.authenticator_selection.authenticator_attachment", as.AuthenticatorAttachment)
	err = cfg.VP.WriteConfig()
	if err != nil {
		fail(c, err)
		return
	}
	success(c)
}

// SetRequireResidentKey 设置 Authenticator selection 之 Require resident key (username-less flow)
func SetRequireResidentKey(c *gin.Context) {
	var as cfg.AuthenticatorSelection
	var err error
	err = c.ShouldBindJSON(&as)
	if err != nil {
		fail(c, err)
		return
	}
	cfg.VP.Set("register.authenticator_selection.require_resident_key", as.RequireResidentKey)
	err = cfg.VP.WriteConfig()
	if err != nil {
		fail(c, err)
		return
	}
	success(c)
}

// SetUserVerification 设置 Authenticator selection 之 User verification required (also for authentication)
func SetUserVerification(c *gin.Context) {
	var as cfg.AuthenticatorSelection
	var err error
	err = c.ShouldBindJSON(&as)
	if err != nil {
		fail(c, err)
		return
	}
	cfg.VP.Set("register.authenticator_selection.user_verification", as.UserVerification)
	cfg.VP.Set("login.user_verification", as.UserVerification)
	err = cfg.VP.WriteConfig()
	if err != nil {
		fail(c, err)
		return
	}
	success(c)
}

// todo: 设置 CredProtect Extension

// GetUser 获取已注册用户列表
func GetUser(c *gin.Context) {
	var qu model.QueryUser
	if err := c.ShouldBindQuery(&qu); err != nil {
		fail(c, err)
	} else {
		if userLst, total, err := qu.GetUserLst(); err != nil {
			fail(c, err)
		} else {
			list(c, userLst, total)
		}
	}
}

// DelUser 删除用户
func DelUser(c *gin.Context) {
	if err := mysql.GetDB().Unscoped().Delete(&model.User{}, "id = ?", c.Query("id")).Error; err != nil {
		fail(c, err)
	} else {
		// 删除凭据
		mysql.GetDB().Unscoped().Delete(&model.Credential{}, "uid = ?", c.Query("id"))
		success(c)
	}
}

// BatchDelUser 批量删除用户
func BatchDelUser(c *gin.Context) {
	idsArr := strings.Split(c.Query("ids"), ",")
	var idsArrInt []int
	for _, v := range idsArr {
		if id, err := strconv.Atoi(v); err != nil {
			fail(c, err)
			return
		} else {
			idsArrInt = append(idsArrInt, id)
		}
	}
	if err := mysql.GetDB().Unscoped().Delete(&model.User{}, idsArrInt).Error; err != nil {
		fail(c, err)
	} else {
		// 删除凭据
		mysql.GetDB().Unscoped().Delete(&model.Credential{}, "uid in ?", idsArrInt)
		success(c)
	}
}

// GetSso 获取sso配置
func GetSso(c *gin.Context) {
	success(c, viper.Get("sso"))
}

// EditSso 修改sso配置
func EditSso(c *gin.Context) {
	type sso struct {
		Url    string `json:"url,omitempty" binding:"url"`
		Secret string `json:"secret,omitempty"`
	}
	var s sso
	err := c.ShouldBindJSON(&s)
	if err != nil {
		fail(c, err)
		return
	} else {
		viper.Set("sso", s)
		err := viper.WriteConfig()
		if err != nil {
			fail(c, err)
			return
		} else {
			success(c)
		}
	}
}

// GetNorth 获取北向接口配置
func GetNorth(c *gin.Context) {
	file := "/etc/northbound.ini"
	vp := viper.New()
	vp.SetConfigFile(file)
	vp.SetConfigName("northbound.ini")
	vp.SetConfigType("ini")
	vp.AddConfigPath("etc")
	err := vp.ReadInConfig()
	if err != nil {
		fail(c, err)
		return
	}
	success(c, vp.AllSettings())
}

// EditNorth 修改北向接口配置
func EditNorth(c *gin.Context) {
	var err error
	type north struct {
		Protocol    string `json:"protocol,omitempty" binding:"oneof=http https"`
		InterfaceIp string `json:"interface_ip,omitempty" binding:"ip"`
		Port        int    `json:"port,omitempty"`
	}
	var n north
	err = c.ShouldBindJSON(&n)
	if err != nil {
		fail(c, err)
		return
	}

	file := "/etc/northbound.ini"
	vp := viper.New()
	vp.SetConfigFile(file)
	vp.SetConfigName("northbound.ini")
	vp.SetConfigType("ini")
	vp.AddConfigPath("etc")
	err = vp.ReadInConfig()
	if err != nil {
		fail(c, err)
		return
	}

	if n.Protocol != "" {
		vp.Set("protocol", n.Protocol)
	}
	if n.InterfaceIp != "" {
		vp.Set("interface_ip", n.InterfaceIp)
	}
	if n.Port != 0 {
		vp.Set("port", n.Port)
	}
	err = vp.WriteConfig()
	if err != nil {
		fail(c, err)
		return
	}
	success(c)
}

// Active 激活license
func Active(c *gin.Context) {
	var err error
	var as logic.AuthServer
	err = c.ShouldBindJSON(&as)
	if err != nil {
		fail(c, err)
		return
	}

	var a logic.Auth
	err, a = logic.ParseLicense(as.License)
	if err != nil {
		fail(c, err)
		return
	}

	viper.Set("auth.days", a.Days)
	viper.Set("auth.name", a.Name)
	viper.Set("auth.project", a.Project)
	viper.Set("auth.apply_time", a.ApplyTime)
	viper.Set("auth_server.license", as.License)
	err = viper.WriteConfig()
	if err != nil {
		fail(c, err)
		return
	}

	success(c)
}

// LicenseStatus 查询软件授权状态
func LicenseStatus(c *gin.Context) {
	err := logic.CheckLicense()
	if err != nil {
		fail(c, err, returnData(viper.GetString("auth.num"), http.StatusUnauthorized))
		return
	}

	success(c, returnNoData(http.StatusOK, "已授权"))
}

// GetLicense 查询当前license
func GetLicense(c *gin.Context) {
	type auth struct {
		Auth logic.Auth `json:"auth" mapstructure:"auth"`
	}
	var a auth
	err := viper.Unmarshal(&a)
	if err != nil {
		fail(c, err)
		return
	}
	success(c, a.Auth)
}

// 后台生成app_id app_secret

// 使用app_id和app_secret获取token
