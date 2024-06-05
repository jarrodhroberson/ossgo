package gin

import (
	g "github.com/gin-gonic/gin"
	errs "github.com/jarrodhroberson/ossgo/errors"
)

func MustGetValue[T any](c *g.Context, key string) T {
	value, exists := c.Get("account_id")
	if !exists {
		panic(errs.NotFoundError.New("%s not found in context", key))
	}
	return value.(T)
}
