package http

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joomcode/errorx"

	"github.com/jarrodhroberson/ossgo/functions/must"

	"github.com/gin-gonic/gin"
	"github.com/jellydator/ttlcache/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

// RecoveryMiddleware catches panics and prevents the server from crashing
func RecoveryMiddleware() func(ctx *gin.Context, err any) {
	return func(c *gin.Context, recovered any) {
		if err, ok := errorx.ErrorFromPanic(recovered); ok {
			err = errorx.EnsureStackTrace(err)
			errx := errorx.Cast(err)
			errx = errorx.EnhanceStackTrace(err, "Global Panic Recovery")
			log.Error().Stack().Err(err).Msgf("Panic recovered: %s", errx)
			if log.Logger.GetLevel() == zerolog.DebugLevel {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("Internal Server Error %s", errx),
				})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error",
				})
				return
			}
		}
	}
}

// StackTraceLoggingErrorHandler is a middleware that logs stack traces for errors.
func StackTraceLoggingErrorHandler(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Errors != nil {
			for _, ginErr := range c.Errors {
				log.Error().Stack().Err(ginErr).Msgf(ginErr.Error())
			}
		}
		c.Next()
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

func EnableGzipCompression(paths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(paths) == 0 {
			c.Writer.Header().Set("Content-Encoding", "gzip")
		} else {
			path := c.Request.URL.Path
			for _, p := range paths {
				if strings.HasPrefix(path, p) {
					c.Writer.Header().Set("Content-Encoding", "gzip")
					break
				}
			}
		}
		c.Next()
	}
}

func DisableGzipCompression(paths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(paths) == 0 {
			c.Writer.Header().Del("Content-Encoding")
		} else {
			path := c.Request.URL.Path
			for _, p := range paths {
				if strings.HasPrefix(path, p) {
					c.Writer.Header().Del("Content-Encoding")
					break
				}
			}
		}
		c.Next()
	}
}

// GoogleAppEngineHttpTaskAuth is a Gin middleware function that ensures API requests contain
// the required headers indicating that the request is from a Google Cloud Task.
// It validates the presence of specific headers and aborts the request with an
// HTTP 401 Unauthorized status if any required headers are missing.
//
// The middleware checks for the following headers:
// - X-AppEngine-QueueName
// - X-AppEngine-TaskName
// - X-AppEngine-TaskRetryCount
// - X-AppEngine-TaskExecutionCount
// - X-AppEngine-TaskETA
//
// The following headers are optional because they provide additional metadata
// about the task execution and are not critical for task authentication or processing:
//
// - X-AppEngine-TaskPreviousResponse
// - X-AppEngine-TaskRetryReason
// - X-AppEngine-FailFast
//
// Specification for [Reading App Engine task request headers]
//
// For each of these headers, the middleware checks their presence and extracts their values,
// setting them into the Gin context for further processing.
//
// Returns a gin.HandlerFunc for use in your Gin HTTP handler chain.
//
// [Reading App Engine task request headers]: https://cloud.google.com/tasks/docs/creating-appengine-handlers#reading-headers
func GoogleAppEngineHttpTaskAuth() gin.HandlerFunc {
	requiredAppEngineTaskHeaders := []string{"X-AppEngine-QueueName", "X-AppEngine-TaskName", "X-AppEngine-TaskRetryCount", "X-AppEngine-TaskExecutionCount", "X-AppEngine-TaskETA"}
	optionalAppEngineTaskHeaders := []string{"X-AppEngine-TaskPreviousResponse", "X-AppEngine-TaskRetryReason", "X-AppEngine-FailFast"}
	return func(c *gin.Context) {
		for _, header := range requiredAppEngineTaskHeaders {
			if _, ok := c.Request.Header[header]; !ok {
				log.Error().Msgf("%s header not found", header)
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.Set(header, c.Request.Header.Get(header))
		}
		for _, header := range optionalAppEngineTaskHeaders {
			if _, ok := c.Request.Header[header]; ok {
				c.Set(header, c.Request.Header.Get(header))
			}
		}
		c.Next()
	}
}

func GoogleHttpTaskAuth() gin.HandlerFunc {
	requiredCloudTaskHeaders := []string{"X-CloudTasks-TaskName", "X-CloudTasks-TaskRetryCount",
		"X-CloudTasks-TaskRetryCount", "X-CloudTasks-TaskExecutionCount", "X-CloudTasks-TaskETA"}
	return func(c *gin.Context) {
		for _, header := range requiredCloudTaskHeaders {
			if _, ok := c.Request.Header[header]; !ok {
				log.Error().Msgf("%s header not found in %s", header, must.MarshalJson(c.Request.Header))
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.Set(header, c.Request.Header.Get(header))
		}
		c.Next()
	}
}
