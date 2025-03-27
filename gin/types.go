package gin

import (
	"net/http"
	"strings"

	"firebase.google.com/go/auth"

	"github.com/gin-gonic/gin"
	errs "github.com/jarrodhroberson/ossgo/errors"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"
)

type Context gin.Context

func (c *Context) IsAuthenticated(t *auth.Token, err error) bool {
	return t != nil && err == nil
}

func (c *Context) CurrentIdToken() (*auth.Token, error) {
	gc := (*gin.Context)(c)
	return GetCurrentUserIdToken(gc)
}

func (c *Context) AccountId() (string, error) {
	gc := (*gin.Context)(c)
	return GetValue[string](gc, "account_id")
}

func (c *Context) YouTubeId() (string, error) {
	gc := (*gin.Context)(c)
	return GetValue[string](gc, "youtube_id")
}

// AcceptHandlerRegistry struct to manage the acceptHeaderHandlerMap.
type AcceptHandlerRegistry struct {
	contentType            string                     // The base content type, e.g., "application/vnd.example.resource+json"
	acceptHeaderHandlerMap map[string]gin.HandlerFunc // Key is the full Accept header, value is the handler
}

// RegisterAcceptHandler registers an Accept header and handler.
func (r *AcceptHandlerRegistry) RegisterAcceptHandler(acceptHeader string, handler gin.HandlerFunc) *AcceptHandlerRegistry {
	if acceptHeader == "" {
		err := errorx.EnsureStackTrace(errs.MustNeverError.WrapWithNoMessage(errorx.IllegalArgument.New("acceptHeader cannot be empty")))
		log.Error().Err(err).Msg(err.Error())
		errorx.Panic(err)
	}
	if handler == nil {
		err := errorx.EnsureStackTrace(errs.MustNeverError.WrapWithNoMessage(errorx.IllegalArgument.New("handler cannot be nil")))
		log.Error().Err(err).Msg(err.Error())
		errorx.Panic(err)
	}

	// Normalize the accept header
	normalizedHeader := strings.TrimSpace(strings.ToLower(acceptHeader))
	r.acceptHeaderHandlerMap[normalizedHeader] = handler

	return r // Return the registry for chaining
}

// HandlerFunc returns a gin.HandlerFunc that handles routing based on the Accept header.
func (r *AcceptHandlerRegistry) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptHeaderFromRequest := c.GetHeader("Accept")
		// Check if the acceptHeaderFromRequest *starts* with the content type.  This makes the logic more robust
		// if other parts of the Accept header are present (e.g., quality factors).
		if !strings.HasPrefix(strings.ToLower(acceptHeaderFromRequest), r.contentType) {
			c.Next() // Pass to the next handler if it doesn't match the content type.  Good for multiple registries.
			return
		}

		// Iterate through the map
		for key, h := range r.acceptHeaderHandlerMap {
			if strings.EqualFold(acceptHeaderFromRequest, key) {
				h(c) // Call the handler
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Unsupported content type or version"})
	}
}
