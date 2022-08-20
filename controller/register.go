package controller

import (
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
	options, sessionData, err := cfg.WAWeb.BeginRegistration(&user)
	// handle errors if present
	if err != nil {
		fail(c, err)
		return
	}
	// store the sessionData values
	redis.GetRds().Set("session:"+c.Query("username"), sessionData, 0)
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
	err = redis.GetRds().Get("session:" + c.Query("username")).Scan(&sessionData)
	if err != nil {
		fail(c, err)
		return
	}
	// using gorilla/sessions it could look like this
	//sessionData := store.Get(r, "registration-session")
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(c.Request.Body)
	if err != nil {
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
	user.AddCredential(*credential)

	success(c, returnNoData(http.StatusOK, "注册成功")) // Handle next steps

}
