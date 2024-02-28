package slices

import (
	"github.com/joomcode/errorx"
)

var struct_not_found = errorx.NewType(errorx.NewNamespace("SERVER"), "STRUCT NOT FOUND", errorx.NotFound())
