package controller

import (
	"encoding/json"
	"fmt"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"srun/cfg"
	"srun/dao/mysql"
	"srun/dao/redis"
	"srun/logic"
	"srun/model"
	"time"
)

func AdminLoginBegin(c *gin.Context) {
	user, err := model.GetLoginAdmin(c.Query("username")) // Find the user
	if err != nil {
		fail(c, err)
		return
	}

	// Updating the AuthenticatorSelection options.
	// See the struct declarations for values
	allowList := make([]protocol.CredentialDescriptor, 1)
	webAuthnCredentials := user.WebAuthnCredentials()
	for k, v := range webAuthnCredentials {
		allowList[k] = protocol.CredentialDescriptor{
			//CredentialID: credentialToAllowID, // 允许认证的凭据ID
			CredentialID: v.ID,                             // 允许认证的凭据ID
			Type:         protocol.PublicKeyCredentialType, // 允许认证的类型 公钥认证
			Transport: []protocol.AuthenticatorTransport{
				protocol.USB,
				protocol.Internal,
				protocol.NFC,
				protocol.BLE,
			}, // 允许的认证器类型
		}
	}

	// Handle next steps

	options, sessionData, err := cfg.WAWeb.BeginLogin(&user, webauthn.WithAllowedCredentials(allowList), webauthn.WithUserVerification(protocol.VerificationDiscouraged))
	//options, sessionData, err := cfg.WAWeb.BeginLogin(&user)
	//options, sessionData, err := cfg.WAWeb.BeginLogin(&user, webauthn.WithUserVerification(protocol.VerificationPreferred))
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
	redis.GetRds().Set("session-admin:"+c.Query("username"), marshal, time.Minute)
	success(c, options) // return the options generated
	// options.publicKey contain our registration options
}

func AdminLoginFinish(c *gin.Context) {
	user, err := model.GetLoginAdmin(c.Query("username")) // Get the user
	if err != nil {
		fail(c, err)
		return
	}
	// Get the session data stored from the function above
	var sessionData webauthn.SessionData
	bt, err := redis.GetRds().Get("session-admin:" + c.Query("username")).Bytes()
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
	//sessionData := store.Get(r, "login-session")
	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(c.Request.Body)
	if err != nil {
		fail(c, err)
		return
	}
	parsedResponse.Response.CollectedClientData.Origin = fmt.Sprintf("%s://%s:%d", viper.GetString("app.protocol"), viper.GetString("app.host"), viper.GetInt("app.port"))
	credential, err := cfg.WAWeb.ValidateLogin(&user, sessionData, parsedResponse)
	// Handle validation or input errors
	if err != nil {
		fail(c, err)
		return
	}
	var cred model.AdminCredential
	if err = mysql.GetDB().First(&cred, "uid = ? and cid = ?", user.ID, credential.ID).Error; err == nil {
		err := cred.UpdateCredential(*credential)
		if err != nil {
			fail(c, err)
			return
		}
	} else {
		fail(c, err)
		return
	}
	// If login was successful, handle next steps
	//success(c)
	//return
	// todo: 调用jwt的登录
	// 模拟将账号密码写入 请求 post 的 json 中
	marshal, err := json.Marshal(user)
	if err != nil {
		fail(c, err)
		return
	}
	_, err = c.Request.Body.Read(marshal)
	if err != nil {
		fail(c, err)
		return
	}
	logic.JWT().LoginHandler(c)
}
