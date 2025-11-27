package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware provides centralized error handling and logging.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check for errors after request processing
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				slog.Error("request error",
					"error", err.Error(),
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)
			}

			// If no response has been written yet, return a generic error
			if !c.Writer.Written() {
				c.JSON(
					http.StatusInternalServerError,
					ErrorResponse{Error: "internal server error"},
				)
			}
		}
	}
}
