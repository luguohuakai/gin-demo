package controller

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
)

// Limit 定义限流方法
func Limit(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}
