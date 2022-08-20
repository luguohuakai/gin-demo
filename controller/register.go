package controller

import (
	"encoding/json"
	"fmt"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
	"net/http"
	"srun/cfg"
	"srun/dao/redis"
	"srun/model"
)

func Begin(c *gin.Context) {
	user, err := model.GetUser(c.Query("username")) // Find or create the new user
	if err != nil {
		fail(c, err)
		return
	}

	// Updating the AuthenticatorSelection options.
	// See the struct declarations for values
	//authSelect := protocol.AuthenticatorSelection{
	//	AuthenticatorAttachment: protocol.Platform,                // 使用平台验证器
	//	RequireResidentKey:      protocol.ResidentKeyUnrequired(), // 不需要常驻密钥
	//	UserVerification:        protocol.VerificationRequired,    // 需要用户验证
	//}

	// Updating the ConveyancePreference options.
	// See the struct declarations for values
	//conveyancePref := protocol.PreferNoAttestation // 证明传输偏好 不需要用户同意

	// Handle next steps

	//options, sessionData, err := cfg.WAWeb.BeginRegistration(&user, webauthn.WithAuthenticatorSelection(authSelect), webauthn.WithConveyancePreference(conveyancePref))
	options, sessionData, err := cfg.WAWeb.BeginRegistration(&user)
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

	success(c, returnNoData(http.StatusOK, "注册成功")) // Handle next steps

}
