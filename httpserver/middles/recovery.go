package middles

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/yunxiaoyang01/open_sdk/logger"
)

func Recovery() Middle {
	return RecoveryWithLogger(logger.StandardLogger())
}

func RecoveryWithLogger(l logger.Logger) Middle {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l.Errorf(c.Request.Context(), "panic recovered: err = %v, stack = %s", err, debug.Stack())
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
