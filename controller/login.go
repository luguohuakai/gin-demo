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

func LoginBegin(c *gin.Context) {
	user, err := model.GetUser(c.Query("username")) // Find the user
	if err != nil {
		fail(c, err)
		return
	}
	options, sessionData, err := cfg.WAWeb.BeginLogin(&user)
	// handle errors if present
	if err != nil {
		fmt.Println("======")
		fmt.Println(err.Error())
		fail(c, err)
		return
	}
	// store the sessionData values
	marshal, err := json.Marshal(sessionData)
	if err != nil {
		fmt.Println("++++++")
		fmt.Println(err.Error())
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
	fmt.Println(credential)
	// If login was successful, handle next steps
	success(c, returnNoData(http.StatusOK, "登录成功"))
}
