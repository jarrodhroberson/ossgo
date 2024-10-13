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
var UnableToDeleteTrait = errorx.RegisterTrait("Unable to Delete")
var UnableToWriteTrait = errorx.RegisterTrait("Unable to Write")
var UnableToReadTrait = errorx.RegisterTrait("Unable to Read")
var MultipleErrorTrait = errorx.RegisterTrait("Multiple Errors")

var MustNeverError = errorx.NewType(MustNamespace, "Must Never Fail", MustNeverErrorTrait)

// create, read, write errors
var NotCreatedError = MustNeverError.NewSubtype("Not Created", UnableToCreateTrait)
var NotDeletedError = MustNeverError.NewSubtype("Not Deleted", UnableToDeleteTrait)
var NotWrittenError = MustNeverError.NewSubtype("Not Written", UnableToWriteTrait)
var NotUpdatedError = MustNeverError.NewSubtype("Not Updated", UnableToWriteTrait)
var NotReadError = MustNeverError.NewSubtype("Not Read", UnableToReadTrait)
var DuplicateExistsError = errorx.IllegalState.NewSubtype("Duplicate Exists", errorx.Duplicate())

// marshalling errors
var ParseError = MustNeverError.NewSubtype("Unable to Parse", UnableToParseTrait)
var MarshalError = MustNeverError.NewSubtype("Unable To Marshal", UnableToMarshalTrait)
var UnMarshalError = MustNeverError.NewSubtype("Unable To UnMarshal", UnableToUnmarshalTrait)

// searching errors
var NotFoundError = MustNeverError.NewSubtype("Not Found", errorx.NotFound())
