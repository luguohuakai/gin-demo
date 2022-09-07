package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"srun/cfg"
	"srun/dao/mysql"
	"srun/dao/redis"
	"srun/model"
	"time"
)

func Begin(c *gin.Context) {
	user, err := model.GetUser(c.Query("username"), c.PostForm("password")) // Find or create the new user
	if err != nil {
		fail(c, err)
		return
	}

	// Updating the AuthenticatorSelection options.
	// See the struct declarations for values
	authSelect := protocol.AuthenticatorSelection{
		AuthenticatorAttachment: protocol.Platform,                // platform：表示仅接受平台内置的、无法移除的认证器，如手机的指纹识别设备 cross-platform：表示仅接受外部认证器，如 USB Key
		RequireResidentKey:      protocol.ResidentKeyUnrequired(), // 是否要求将私钥钥永久存储于认证器中 // 设置为true可实现无用户名登录
		UserVerification:        protocol.VerificationDiscouraged, // 依赖方不关心用户验证
	}

	// Updating the ConveyancePreference options.
	// See the struct declarations for values
	//conveyancePref := protocol.PreferNoAttestation // 如果你没有高安全需求（如银行交易等），请不要向认证器索取证明，即将 attestation 设置为 "none" 对于普通身份认证来说，要求证明不必要的，且会有浏览器提示打扰到用户
	//extension := protocol.AuthenticationExtensions{
	//	"uvm":          true,       // 要求认证器返回用户进行验证的方法
	//	"txAuthSimple": "你正在注册....", // 在认证器上显示与交易有关的简短消息
	//}

	// Handle next steps

	//options, sessionData, err := cfg.WAWeb.BeginRegistration(&user, webauthn.WithAuthenticatorSelection(authSelect), webauthn.WithConveyancePreference(conveyancePref), webauthn.WithExtensions(extension))
	options, sessionData, err := cfg.WAWeb.BeginRegistration(&user, webauthn.WithAuthenticatorSelection(authSelect))
	// handle errors if present
	if err != nil {
		fail(c, err)
		return
	}
	// store the sessionData values
	marshal, err := json.Marshal(sessionData)
	if err != nil {
		fail(c, err)
		return
	}
	redis.GetRds().Set("session:"+c.Query("username"), marshal, time.Minute)
	success(c, options) // return the options generated
	// options.publicKey contain our registration options
}

func Finish(c *gin.Context) {
	user, err := model.GetUser(c.Query("username")) // Get the user
	if err != nil {
		fail(c, err)
		return
	}
	// Get the session data stored from the function above
	var sessionData webauthn.SessionData
	bt, err := redis.GetRds().Get("session:" + c.Query("username")).Bytes()
	if err != nil {
		fail(c, err)
		return
	}
	err = json.Unmarshal(bt, &sessionData)
	if err != nil {
		fail(c, err)
		return
	}
	// using gorilla/sessions it could look like this
	//sessionData := store.Get(r, "registration-session")
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(c.Request.Body)
	if err != nil {
		fmt.Println(err.Error())
		fail(c, err)
		return
	}
	credential, err := cfg.WAWeb.CreateCredential(&user, sessionData, parsedResponse)
	// Handle validation or input errors
	if err != nil {
		fail(c, err)
		return
	}

	// If creation was successful, store the credential object
	err = user.AddCredential(*credential)
	if err != nil {
		fail(c, err)
		return
	}

	err = mysql.GetDB().Model(&user).Update(model.User{Status: 2}).Error
	if err != nil {
		fail(c, err)
		return
	}
	success(c, returnNoData(http.StatusOK, "注册成功")) // Handle next steps
}

func UserExists(c *gin.Context) {
	if c.Query("username") == "" {
		fail(c, errors.New("username can not be empty"))
		return
	}
	if err := model.UserIsWebAuthn(c.Query("username")); err == nil {
		success(c)
	} else {
		if gorm.IsRecordNotFoundError(err) {
			fail(c, err, returnNoData(4001, "user not register"))
			return
		}
		fail(c, err)
		return
	}
}
