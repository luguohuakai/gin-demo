package controller

import (
	"encoding/json"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
	"net/http"
	"srun/cfg"
	"srun/dao/mysql"
	"srun/dao/redis"
	"srun/model"
)

func LoginBegin(c *gin.Context) {
	user, err := model.GetUser(c.Query("username")) // Find the user
	if err != nil {
		fail(c, err)
		return
	}

	// Updating the AuthenticatorSelection options.
	// See the struct declarations for values
	//allowList := make([]protocol.CredentialDescriptor, 1)
	//allowList[0] = protocol.CredentialDescriptor{
	//	//CredentialID: credentialToAllowID, // 允许认证的凭据ID
	//	Type: protocol.PublicKeyCredentialType, // 允许认证的类型 公钥认证
	//	//Transport: []protocol.AuthenticatorTransport{
	//	//	protocol.NFC,
	//	//	protocol.Internal,
	//	//	protocol.USB,
	//	//	protocol.BLE,
	//	//}, // 允许的认证器类型
	//}

	// Handle next steps

	//options, sessionData, err := cfg.WAWeb.BeginLogin(&user, webauthn.WithAllowedCredentials(allowList), webauthn.WithUserVerification(protocol.VerificationRequired))
	options, sessionData, err := cfg.WAWeb.BeginLogin(&user, webauthn.WithUserVerification(protocol.VerificationPreferred))
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
	redis.GetRds().Set("session:"+c.Query("username"), marshal, 0)
	success(c, options) // return the options generated
	// options.publicKey contain our registration options
}

func LoginFinish(c *gin.Context) {
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
	//sessionData := store.Get(r, "login-session")
	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(c.Request.Body)
	if err != nil {
		fail(c, err)
		return
	}
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
	// If login was successful, handle next steps
	success(c, returnNoData(http.StatusOK, "登录成功"))
}