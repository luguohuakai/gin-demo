package controller

import (
	"encoding/json"
	"fmt"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"srun/cfg"
	"srun/dao/mysql"
	"srun/dao/redis"
	"srun/model"
	"time"
)

func Begin(c *gin.Context) {
	var err error
	var user model.User
	user, err = model.GetUser(c.Query("username"), "begin", c.PostForm("password")) // Find or create the new user
	if err != nil {
		fail(c, err)
		return
	}

	var opts []webauthn.RegistrationOption

	var authSelect protocol.AuthenticatorSelection
	var flag bool
	if cfg.FD.Register.AuthenticatorSelection.AuthenticatorAttachment != "" {
		flag = true
		authSelect.AuthenticatorAttachment = cfg.FD.Register.AuthenticatorSelection.AuthenticatorAttachment
	}
	if cfg.FD.Register.AuthenticatorSelection.RequireResidentKey != "" {
		flag = true
		if cfg.FD.Register.AuthenticatorSelection.RequireResidentKey == "true" {
			authSelect.RequireResidentKey = protocol.ResidentKeyRequired()
		} else {
			authSelect.RequireResidentKey = protocol.ResidentKeyUnrequired()
		}
	}
	if cfg.FD.Register.AuthenticatorSelection.UserVerification != "" {
		flag = true
		authSelect.UserVerification = cfg.FD.Register.AuthenticatorSelection.UserVerification
	}
	if flag {
		opts = append(opts, webauthn.WithAuthenticatorSelection(authSelect))
	}

	// （可选）用于标识要排除的凭证，可以避免同一个用户多次注册同一个认证器。如果用户试图注册相同的认证器，用户代理会抛出 InvalidStateError 错误。数组中的每一项都是一个公钥凭证对象
	var excludeList []protocol.CredentialDescriptor

	for _, v := range user.WebAuthnCredentials() {
		excludeList = append(excludeList, protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: v.ID,
			Transport:    cfg.FD.Register.ExcludeCredentials.Transports,
		})
	}
	if len(excludeList) > 0 {
		opts = append(opts, webauthn.WithExclusions(excludeList))
	}

	if cfg.FD.Register.Attestation != "" {
		opts = append(opts, webauthn.WithConveyancePreference(cfg.FD.Register.Attestation))
	}

	if cfg.FD.Register.Timeout != 0 {
		opts = append(opts, func(cco *protocol.PublicKeyCredentialCreationOptions) {
			cco.Timeout = int(cfg.FD.Register.Timeout)
		})
	}

	//extension := protocol.AuthenticationExtensions{
	//	"uvm":          true,       // 要求认证器返回用户进行验证的方法
	//	"txAuthSimple": "你正在注册....", // 在认证器上显示与交易有关的简短消息
	//}

	var options *protocol.CredentialCreation
	var sessionData *webauthn.SessionData
	if len(opts) > 0 {
		options, sessionData, err = cfg.WAWeb.BeginRegistration(&user, opts...)
	} else {
		options, sessionData, err = cfg.WAWeb.BeginRegistration(&user)
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

func Finish(c *gin.Context) {
	user, err := model.GetUser(c.Query("username"), "finish") // Get the user
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
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(c.Request.Body)
	if err != nil {
		fmt.Println(err.Error())
		fail(c, err)
		return
	}
	parsedResponse.Response.CollectedClientData.Origin = fmt.Sprintf("%s://%s:%d", viper.GetString("app.protocol"), viper.GetString("app.host"), viper.GetInt("app.port"))
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
	success(c, returnNoData(http.StatusOK, "注册成功"))
}
