package firestore

import (
	"github.com/joomcode/errorx"
)

var DocumentNotFound = errorx.IllegalArgument.NewSubtype("DocumentNotFound", errorx.NotFound())
