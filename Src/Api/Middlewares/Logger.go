package Middlewares

import (
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func DefaultStructuredLogger(cfg *Config.Config) gin.HandlerFunc {
	logger := Log.NewLogger(cfg)
	return structuredLogger(logger)
}

func structuredLogger(logger Log.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.FullPath(), "swagger") {
			c.Next()
		} else {
			blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			start := time.Now() // start
			path := c.FullPath()
			raw := c.Request.URL.RawQuery

			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body.Close()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			c.Writer = blw
			c.Next()

			param := gin.LogFormatterParams{}
			param.TimeStamp = time.Now() // stop
			param.Latency = param.TimeStamp.Sub(start)
			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path

			keys := map[Log.Extrakey]interface{}{}
			keys[Log.Path] = param.Path
			keys[Log.ClientIp] = param.ClientIP
			keys[Log.Method] = param.Method
			keys[Log.Latency] = param.Latency
			keys[Log.StatusCode] = param.StatusCode
			keys[Log.ErrorMessage] = param.ErrorMessage
			keys[Log.BodySize] = param.BodySize
			keys[Log.RequestBody] = string(bodyBytes)
			keys[Log.ResponseBody] = blw.body.String()

			logger.Info(Log.RequestResponse, Log.Api, "", keys)
		}
	}
}
