package errors

import (
	"github.com/joomcode/errorx"
)

var MustNamespace = errorx.NewNamespace("Must")

var MustNeverErrorTrait = errorx.RegisterTrait("Must Never Error")
var UnableToParseTrait = errorx.RegisterTrait("Unable to Parse")
var UnableToMarshalTrait = errorx.RegisterTrait("Unable to Marshal")
var UnableToUnmarshalTrait = errorx.RegisterTrait("Unable to Marshal")
var UnableToCreateTrait = errorx.RegisterTrait("Unable to Create")

var MustNeverError = errorx.NewType(MustNamespace, "Must Never Fail", MustNeverErrorTrait)

var ParseError = MustNeverError.NewSubtype("Unable to Parse", UnableToParseTrait)
var NotFoundError = MustNeverError.NewSubtype("Not Found", errorx.NotFound())
var NotCreatedError = MustNeverError.NewSubtype("Not Created", UnableToCreateTrait)
var MarshalError = MustNeverError.NewSubtype("Unable To Marshal", UnableToMarshalTrait)
var UnMarshalError = MustNeverError.NewSubtype("Unable To Marshal", UnableToUnmarshalTrait)
var StructNotFound = NotFoundError.NewSubtype("Struct not found", errorx.NotFound())

var DuplicateFound = errorx.IllegalState.NewSubtype("Duplicate Found")
