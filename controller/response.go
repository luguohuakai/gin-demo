package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var version = "1.0.0"

type result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int         `json:"total,omitempty"`
	Version string      `json:"version"`
	Error   string      `json:"error,omitempty"`
}

func returnData(data interface{}, code int, messages ...string) *result {
	var message string
	if len(messages) < 1 {
		message = http.StatusText(code)
		if message == "" {
			message = "unknown code"
		}
	} else {
		message = messages[0]
	}
	return &result{
		Code:    code,
		Message: message,
		Data:    data,
		Total:   1,
		Version: version,
	}
}

func returnNoData(code int, messages ...string) *result {
	var message string
	if len(messages) < 1 {
		message = http.StatusText(code)
		if message == "" {
			message = "unknown code"
		}
	} else {
		message = messages[0]
	}
	return &result{
		Code:    code,
		Message: message,
		Version: version,
	}
}

func success(c *gin.Context, data ...interface{}) {
	if len(data) < 1 {
		c.JSON(http.StatusOK, returnNoData(http.StatusOK))
	} else {
		if dt, ok := data[0].(*result); ok {
			if dt.Message == "" {
				dt.Message = http.StatusText(dt.Code)
				if dt.Message == "" {
					dt.Message = "unknown code"
				}
			}
			c.JSON(http.StatusOK, dt)
		} else {
			if dt, ok := data[0].(result); ok {
				if dt.Message == "" {
					dt.Message = http.StatusText(dt.Code)
					if dt.Message == "" {
						dt.Message = "unknown code"
					}
				}
				c.JSON(http.StatusOK, dt)
			} else {
				c.JSON(http.StatusOK, returnData(data[0], http.StatusOK))
			}
		}
	}
}

func list(c *gin.Context, data interface{}, total int) {
	if dt, ok := data.(*result); ok {
		dt.Total = total
		c.JSON(http.StatusOK, dt)
	} else {
		res := returnData(data, http.StatusOK)
		res.Total = total
		c.JSON(http.StatusOK, res)
	}
}

func fail(c *gin.Context, err error, data ...interface{}) {
	if err != nil {
		if len(data) < 1 {
			res := returnNoData(http.StatusBadRequest, err.Error())
			res.Error = err.Error()
			c.JSON(http.StatusOK, res)
		} else {
			if dt, ok := data[0].(*result); ok {
				if dt.Message == "" {
					dt.Message = err.Error()
				}
				dt.Error = err.Error()
				c.JSON(http.StatusOK, dt)
			} else {
				if dt, ok := data[0].(result); ok {
					if dt.Message == "" {
						dt.Message = err.Error()
					}
					dt.Error = err.Error()
					c.JSON(http.StatusOK, dt)
				} else {
					res := returnData(data[0], http.StatusBadRequest, err.Error())
					res.Error = err.Error()
					c.JSON(http.StatusOK, res)
				}
			}
		}
	} else {
		if len(data) < 1 {
			c.JSON(http.StatusOK, returnNoData(http.StatusBadRequest))
		} else {
			if dt, ok := data[0].(*result); ok {
				c.JSON(http.StatusOK, dt)
			} else {
				if dt, ok := data[0].(result); ok {
					c.JSON(http.StatusOK, dt)
				} else {
					c.JSON(http.StatusOK, returnData(data[0], http.StatusBadRequest))
				}
			}
		}
	}
}
