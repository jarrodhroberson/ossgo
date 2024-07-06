package firestore

import (
	"github.com/joomcode/errorx"
)

var DocumentNotFound = errorx.IllegalArgument.NewSubtype("DocumentNotFound", errorx.NotFound())
var BulkWriterError = errorx.IllegalState.NewSubtype("Bulk Writer Error", errorx.Temporary())
