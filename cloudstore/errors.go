package cloudstore

import (
	errs "github.com/jarrodhroberson/ossgo/errors"
)

var ObjectWriteError = errs.NotWrittenError.NewSubtype("Object Write Error")
var ObjectCreateError = errs.NotCreatedError.NewSubtype("Object Write Error")
