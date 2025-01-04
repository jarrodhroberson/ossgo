package gin

import (
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
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
