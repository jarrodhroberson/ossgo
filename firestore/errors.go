package firestore

import (
	"github.com/joomcode/errorx"
)

var BulkWriterError = errorx.IllegalState.NewSubtype("Bulk Writer Error")
