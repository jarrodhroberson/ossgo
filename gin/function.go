package gin

import (
	"fmt"
	"strings"

	"firebase.google.com/go/auth"

	g "github.com/gin-gonic/gin"
	errs "github.com/jarrodhroberson/ossgo/errors"
)

const CURRENT_USER_ID_TOKEN = "CURRENT_USER_ID_TOKEN"

// GetValue retrieves a value from the Gin context with the specified key.
// It returns the value as type T and an error if the key is not found.
// If the key is not found, it returns a NotFoundError.
// It panics if the value is not of type T.
func GetValue[T any](c *g.Context, key string) (T, error) {
	value, exists := c.Get(key)
	if !exists {
		return *new(T), errs.NotFoundError.New("%s not found in context", key)
	}
	return value.(T), nil
}

// MustGetValue retrieves a value from the Gin context with the specified key.
// It returns the value as type T.
// If the key is not found, it panics with a NotFoundError.
// It panics if the value is not of type T.
// Deprecated: use must.Must(GetValue[T](c, key)) instead
func MustGetValue[T any](c *g.Context, key string) T {
	value, exists := c.Get(key)
	if !exists {
		panic(errs.NotFoundError.New("%s not found in context", key))
	}
	return value.(T)
}

// MustGetCurrentUserIdToken retrieves the current user's ID token from the Gin context.
// It returns the ID token as a *auth.Token.
// If the ID token is not found, it returns nil.
// Deprecated: use must.Must(GetCurrentUserIdToken(c)) instead
func MustGetCurrentUserIdToken(c *g.Context) *auth.Token {
	idToken, err := GetCurrentUserIdToken(c)
	if err != nil {
		//TODO: put some logging here this should never happen
		return nil
	}
	return idToken
}

// GetCurrentUserIdToken retrieves the current user's ID token from the Gin context.
// It returns the ID token as a *auth.Token and an error if the ID token is not found.
// If the ID token is not found, it returns a COOKIE_NOT_FOUND error.
func GetCurrentUserIdToken(c *g.Context) (*auth.Token, error) {
	idToken, ok := c.Get(CURRENT_USER_ID_TOKEN)
	if !ok {
		return nil, COOKIE_NOT_FOUND.Wrap(fmt.Errorf("%s not found in gin.Context", CURRENT_USER_ID_TOKEN), "")
	}
	return idToken.(*auth.Token), nil
}

// HasHeader checks if a header with the specified name exists in the Gin context.
// It returns true if the header exists, false otherwise.
func HasHeader(c *g.Context, name string) bool {
	return c.GetHeader(name) != ""
}

// HasCookie checks if a cookie with the specified name exists in the Gin context.
// It returns true if the cookie exists, false otherwise.
func HasCookie(c *g.Context, name string) bool {
	_, err := c.Cookie(name)
	return err == nil
}

// NewAcceptHandlerRegistry creates a new AcceptHandlerRegistry for a specific content type.
func NewAcceptHandlerRegistry(contentType string) *AcceptHandlerRegistry {
	return &AcceptHandlerRegistry{
		contentType:            strings.TrimSpace(strings.ToLower(contentType)), // Store the base content type
		acceptHeaderHandlerMap: make(map[string]g.HandlerFunc),                  // Initialize the map
	}
}
