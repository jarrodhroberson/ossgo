package gin

import (
	"fmt"

	"firebase.google.com/go/auth"
	g "github.com/gin-gonic/gin"
	errs "github.com/jarrodhroberson/ossgo/errors"
)

const CURRENT_USER_ID_TOKEN = "CURRENT_USER_ID_TOKEN"

func GetValue[T any](c *g.Context, key string) (T, error) {
	value, exists := c.Get(key)
	if !exists {
		return *new(T), errs.NotFoundError.New("%s not found in context", key)
	}
	return value.(T), nil
}

func MustGetValue[T any](c *g.Context, key string) T {
	value, exists := c.Get(key)
	if !exists {
		panic(errs.NotFoundError.New("%s not found in context", key))
	}
	return value.(T)
}

func MustGetCurrentUserIdToken(c *g.Context) *auth.Token {
	idToken, err := GetCurrentUserIdToken(c)
	if err != nil {
		//TODO: put some logging here this should never happen
		return nil
	}
	return idToken
}

func GetCurrentUserIdToken(c *g.Context) (*auth.Token, error) {
	idToken, ok := c.Get(CURRENT_USER_ID_TOKEN)
	if !ok {
		return nil, COOKIE_NOT_FOUND.Wrap(fmt.Errorf("%s not found in gin.Context", CURRENT_USER_ID_TOKEN), "")
	}
	return idToken.(*auth.Token), nil
}
