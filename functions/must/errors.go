package must

import (
	"github.com/joomcode/errorx"
)

var MustNeverErrorTrait = errorx.RegisterTrait("Must Never Error")
var UnableToParseTrait = errorx.RegisterTrait("Unable to Parse")
var UnableToMarshalTrait = errorx.RegisterTrait("Unable to Marshal")
var UnableToUnmarshalTrait = errorx.RegisterTrait("Unable to Marshal")

var mustNeverError = errorx.NewType(errorx.NewNamespace("Must"), "Must Never Fail", MustNeverErrorTrait)

var parseError = mustNeverError.NewSubtype("Unable to Parse", UnableToParseTrait)
var notFoundError = mustNeverError.NewSubtype("Not Found", errorx.NotFound())
var marshalError = mustNeverError.NewSubtype("Unable To Marshal", UnableToMarshalTrait)
var unMarshalError = mustNeverError.NewSubtype("Unable To Marshal", UnableToUnmarshalTrait)
