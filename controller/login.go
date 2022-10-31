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
	"srun/model"
	"time"
)

func LoginBegin(c *gin.Context) {
	var err error
	var user model.User
	user, err = model.GetLoginUser(c.Query("username")) // Find the user
	if err != nil {
		fail(c, err)
		return
	}

	var opts []webauthn.LoginOption

	// Updating the AuthenticatorSelection options.
	// See the struct declarations for values
	allowList := make([]protocol.CredentialDescriptor, 1)
	webAuthnCredentials := user.WebAuthnCredentials()
	for k, v := range webAuthnCredentials {
		allowList[k] = protocol.CredentialDescriptor{
			CredentialID: v.ID,                                     // 允许认证的凭据ID
			Type:         protocol.PublicKeyCredentialType,         // 允许认证的类型 公钥认证
			Transport:    cfg.FD.Login.AllowCredentials.Transports, // 允许的认证器类型
		}
	}
	if len(allowList) > 0 {
		opts = append(opts, webauthn.WithAllowedCredentials(allowList))
	}

	if cfg.FD.Login.Timeout != 0 {
		opts = append(opts, func(cco *protocol.PublicKeyCredentialRequestOptions) {
			cco.Timeout = int(cfg.FD.Login.Timeout)
		})
	}

	if cfg.FD.Login.UserVerification != "" {
		opts = append(opts, webauthn.WithUserVerification(cfg.FD.Login.UserVerification))
	}

	// Handle next steps
	var options *protocol.CredentialAssertion
	var sessionData *webauthn.SessionData
	if len(opts) > 0 {
		options, sessionData, err = cfg.WAWeb.BeginLogin(&user, opts...)
	} else {
		options, sessionData, err = cfg.WAWeb.BeginLogin(&user)
	}
	// handle errors if present
	if err != nil {
		fail(c, err)
		return
	}
	// store the sessionData values
	var marshal []byte
	marshal, err = json.Marshal(sessionData)
	if err != nil {
		fail(c, err)
		return
	}
	redis.GetRds().Set("session:"+c.Query("username"), marshal, time.Minute)
	success(c, options) // return the options generated
	// options.publicKey contain our registration options
}

func LoginFinish(c *gin.Context) {
	user, err := model.GetLoginUser(c.Query("username")) // Get the user
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
	var cred model.Credential
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
	success(c)
}
