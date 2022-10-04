package logging

import (
	"time"

	"github.com/gin-gonic/gin"
)

func GinLogging(logger ContextLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger = logger.With(
			"method", c.Request.Method,
			"path", c.Request.RequestURI,
			"status", c.Writer.Status(),
			"duration", time.Since(start),
			"referrer", c.Request.Referer(),
		)

		if c.Writer.Status() >= 500 {
			logger.Errorw("⚠️ Error", "err", c.Errors.String())
		} else {
			logger.Infow("✅ Served")
		}
	}
}
