package log

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	green  = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white  = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red    = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	reset  = string([]byte{27, 91, 48, 109})
)

// Instances a Logger middleware that will write the logs to gin.DefaultWriter
// By default gin.DefaultWriter = os.Stdout
func Logger() gin.HandlerFunc {
	return LoggerWithWriter(gin.DefaultWriter)
}

// Instance a Logger middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(out io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		statusColor := colorForStatus(statusCode)

		if len(c.Request.Form) > 0 {
			fmt.Fprintf(out, "[GIN] %v |%s %3d %s| %11v | %s |%-7s %s\n      Form: %v \n",
				end.Format("2006/01/02 15:04:05"),
				statusColor, statusCode, reset,
				latency,
				clientIP,
				method,
				path,
				c.Request.Form,
			)
		} else {
			fmt.Fprintf(out, "[GIN] %v |%s %3d %s| %11v | %s |%-7s %s\n",
				end.Format("2006/01/02 15:04:05"),
				statusColor, statusCode, reset,
				latency,
				clientIP,
				method,
				path,
			)
		}

	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}
