package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Demo(c *gin.Context) {
	switch c.Query("type") {
	case "1":
		success(c)
	case "2":
		success(c, nil)
	case "3":
		success(c, "good")
	case "4":
		success(c, returnData("data content", http.StatusOK))
	case "5":
		success(c, &result{
			Code:    4012,
			Message: "没有认证啊",
			Total:   1,
			Version: version,
		})
	case "6":
		fail(c, nil)
	case "7":
		fail(c, errors.New("error found"), "content")
	case "8":
		fail(c, nil, "good")
	case "9":
		fail(c, nil, nil)
	case "10":
		list(c, "list data", 129)
	default:
		success(c)
	}
	return
}
