/*
Package errors contains general purpose errors and error traits for you to build more specific errors in your own packages.
*/
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
var MutuallyExclusiveTrait = errorx.RegisterTrait("Mutually Exclusive")
var InvalidSizeTrait = errorx.RegisterTrait("Invalid Size")

var MustNeverError = errorx.NewType(MustNamespace, "Must Never Fail", MustNeverErrorTrait)

// security (authentication, authorization)
var Unauthorized = MustNeverError.NewSubtype("Unauthorized")
var Invalid = errorx.AssertionFailed.NewSubtype("Invalid")

// create, read, write errors
var StructNotInitialized = MustNeverError.NewSubtype("Struct Not Initialized", UnableToCreateTrait)
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
var CookieNotFoundError = NotFoundError.NewSubtype("CookieNotFoundError")

// constraint/validation errors
var MinSizeExceededError = MustNeverError.NewSubtype("Min Required Size", InvalidSizeTrait)
var MaxSizeExceededError = MustNeverError.NewSubtype("Max Size Exceeded", InvalidSizeTrait)
var ExpiredError = errorx.TimeoutElapsed.NewSubtype("Expired", errorx.Timeout())
var DisabledError = errorx.IllegalState.NewSubtype("Disabled", errorx.Temporary())
var InvalidJsonPayloadReceived = Invalid.NewSubtype("invalid json payload received.")
var CanNotBindQueryParameter = UnMarshalError.NewSubtype("can not bind query parameter.")