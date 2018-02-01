package log

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a middleware that will write the logs to gin.DefaultWriter
// By default gin.DefaultWriter = os.Stdout
func Logger(fileName string, maxdays int, hideValues []string) gin.HandlerFunc {
	if fileName == "" {
		return LoggerWithWriter(gin.DefaultWriter, hideValues)
	}
	fw := new(fileLogWriter)
	fw.FileName = fileName
	fw.MaxDays = maxdays

	_, err := fw.createLogFile()
	if err != nil {
		panic(err)
	}
	return LoggerWithWriter(fw, hideValues)
}

// LoggerWithWriter is a middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(out io.Writer, hideValues []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		statusCode := c.Writer.Status()
		// if (strings.HasPrefix(path, "/system/") || strings.HasPrefix(path, "/assets/")) && statusCode < 400 {
		// 	return
		// }

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		formValues := c.Request.URL.Query()
		if formValues == nil {
			formValues = make(map[string][]string)
		}
		for k, v := range c.Request.Form {
			formValues[k] = v
		}

		if len(formValues) > 0 {
			for formValue, _ := range formValues {
				for _, hideValue := range hideValues {
					if formValue == hideValue || formValue == "password" {
						formValues[formValue] = []string{"***"}
					}
				}
			}
			fmt.Fprintf(out, "[GIN] %v | %3d | %11v | %s |%-7s %s\n      Params: %v \n",
				end.Format("2006/01/02 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
				formValues,
			)
		} else {
			fmt.Fprintf(out, "[GIN] %v | %3d | %11v | %s |%-7s %s\n",
				end.Format("2006/01/02 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
			)
		}
	}
}
