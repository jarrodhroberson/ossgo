package http

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jellydator/ttlcache/v3"
	"github.com/rs/zerolog"
)

// AntiHackingMiddleware is a middleware that blocks requests from blacklisted IPs and paths.
func AntiHackingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIp := c.ClientIP()
		if blackHole.Get(clientIp) != nil {
			hardCloseConnection(c)
			c.Abort()
			return
		}
		path := c.Request.URL.Path
		if item := banned.Get(path); item != nil {
			count := item.Value()
			count.Add(1)
			blackHole.Set(clientIp, time.Now().UTC(), ttlcache.DefaultTTL)
			hardCloseConnection(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

// StackTraceLoggingErrorHandler is a middleware that logs stack traces for errors.
func StackTraceLoggingErrorHandler(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, ginErr := range c.Errors {
			log.Error().Err(ginErr).Stack().Msgf(ginErr.Error())
		}
	}
}

// RequestMetricsLogger is a middleware that logs request metrics.
func RequestMetricsLogger(serName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		// after request
		// latency := time.Since(t)
		// clientIP := c.ClientIP()
		// method := c.Request.Method
		// statusCode := c.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}
		msg := c.Errors.String()
		if msg == "" {
			msg = "Request"
		}
		cData := &ginHands{
			LogLevel:   os.Getenv("GLOBAL_LOG_LEVEL"),
			SerName:    serName,
			Path:       path,
			Latency:    time.Since(t),
			Method:     c.Request.Method,
			StatusCode: c.Writer.Status(),
			ClientIP:   c.Request.RemoteAddr,
			MsgStr:     msg,
		}

		logSwitch(cData)
	}
}
