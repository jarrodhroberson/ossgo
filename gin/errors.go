package gin

import (
	"github.com/joomcode/errorx"
)

var COOKIE_NOT_FOUND = errorx.DataUnavailable.NewSubtype("COOKIE_NOT_FOUND", errorx.NotFound())
