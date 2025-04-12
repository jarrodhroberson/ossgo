/*
Package errors contains general purpose errors and error traits for you to build more specific errors in your own packages.
*/
package errors

import (
	"github.com/joomcode/errorx"
)

var MustNamespace = errorx.NewNamespace("Must")
var HttpNamespace = errorx.NewNamespace("HTTP")

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
var InvalidTrait = errorx.RegisterTrait("Invalid")
var NilTrait = errorx.RegisterTrait("Nil Reference")
var IterationTrait = errorx.RegisterTrait("iteration failed")
var TemporaryTrait = errorx.Temporary()
var PermanentTrait = errorx.RegisterTrait("Permanent")

// http error status traits
var HttpRedirectionTrait = errorx.RegisterTrait("Redirection")
var HttpClientTrait = errorx.RegisterTrait("Client Error")
var HttpServerTrait = errorx.RegisterTrait("Server Error")

var MustNeverError = errorx.NewType(MustNamespace, "Must Never Fail", MustNeverErrorTrait)

// security (authentication, authorization)
var Unauthorized = MustNeverError.NewSubtype("Unauthorized")

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
var NotFoundError = errorx.NewType(errorx.CommonErrors, "Not Found", errorx.NotFound())
var CookieNotFoundError = NotFoundError.NewSubtype("CookieNotFoundError")

// constraint/validation errors
var InvalidData = errorx.NewType(errorx.CommonErrors, "Invalid Data", InvalidTrait)
var InvalidState = errorx.NewType(errorx.CommonErrors, "Invalid State", InvalidTrait)
var MustNotBeNil = MustNeverError.NewSubtype("Must Not Be Nil", NilTrait)
var MinSizeExceededError = MustNeverError.NewSubtype("Min Required Size", InvalidSizeTrait)
var MaxSizeExceededError = MustNeverError.NewSubtype("Max Size Exceeded", InvalidSizeTrait)
var ExpiredError = errorx.TimeoutElapsed.NewSubtype("Expired", errorx.Timeout())
var DisabledError = errorx.IllegalState.NewSubtype("Disabled", errorx.Temporary())
var InvalidJsonPayloadReceived = InvalidData.NewSubtype("invalid json payload received.")
var CanNotBindQueryParameter = UnMarshalError.NewSubtype("can not bind query parameter.")
var IterationError = MustNeverError.NewSubtype("iteration error", IterationTrait)

// standard http errors
var HttpRedirectionStatus = errorx.NewType(HttpNamespace, "Redirect Error", HttpRedirectionTrait)
var StatusMultipleChoices = HttpRedirectionStatus.NewSubtype( "Status Multiple Choices")
var StatusMovedPermanently = HttpRedirectionStatus.NewSubtype( "Status Moved Permanently", PermanentTrait)
var StatusFound = HttpRedirectionStatus.NewSubtype( "Status Found", TemporaryTrait)
var StatusSeeOther = HttpRedirectionStatus.NewSubtype( "Status See Other")
var StatusNotModified = HttpRedirectionStatus.NewSubtype( "Status Not Modified")
var StatusTemporaryRedirect = HttpRedirectionStatus.NewSubtype( "Status Temporary Redirect", TemporaryTrait)
var StatusPermanentRedirect = HttpRedirectionStatus.NewSubtype( "Status Permanent Redirect", PermanentTrait)

var HttpClientErrorStatus = errorx.NewType(HttpNamespace, "Client Error", HttpClientTrait)
var StatusBadRequest = HttpClientErrorStatus.NewSubtype( "Status Bad Request")
var StatusUnauthorized = HttpClientErrorStatus.NewSubtype( "Status Unauthorized")
var StatusPaymentRequired = HttpClientErrorStatus.NewSubtype( "Status Payment Required")
var StatusForbidden = HttpClientErrorStatus.NewSubtype( "Status Forbidden")
var StatusNotFound = HttpClientErrorStatus.NewSubtype( "Status Not Found", errorx.NotFound())
var StatusMethodNotAllowed = HttpClientErrorStatus.NewSubtype( "Status Method Not Allowed")
var StatusNotAcceptable = HttpClientErrorStatus.NewSubtype( "Status Not Acceptable")
var StatusProxyAuthRequired = HttpClientErrorStatus.NewSubtype( "Status Proxy Auth Required")
var StatusRequestTimeout = HttpClientErrorStatus.NewSubtype( "Status Request Timeout", errorx.Timeout())
var StatusConflict = HttpClientErrorStatus.NewSubtype( "Status Conflict")
var StatusGone = HttpClientErrorStatus.NewSubtype( "Status Gone", PermanentTrait)
var StatusLengthRequired = HttpClientErrorStatus.NewSubtype( "Status Length Required", errorx.NotFound())
var StatusPreconditionFailed = HttpClientErrorStatus.NewSubtype( "Status Precondition Failed")
var StatusRequestEntityTooLarge = HttpClientErrorStatus.NewSubtype( "Status Request Entity Too Large", InvalidSizeTrait)
var StatusRequestURITooLong = HttpClientErrorStatus.NewSubtype( "Status Request URI Too Long", InvalidSizeTrait)
var StatusUnsupportedMediaType = HttpClientErrorStatus.NewSubtype( "Status Unsupported Media Type")
var StatusRequestedRangeNotSatisfiable = HttpClientErrorStatus.NewSubtype( "Status Requested Range Not Satisfiable")
var StatusExpectationFailed = HttpClientErrorStatus.NewSubtype( "Status Expectation Failed")
var StatusMisdirectedRequest = HttpClientErrorStatus.NewSubtype( "Status Misdirected Request")
var StatusUnprocessableEntity = HttpClientErrorStatus.NewSubtype( "Status Unprocessable Entity")
var StatusLocked = HttpClientErrorStatus.NewSubtype( "Status Locked")
var StatusFailedDependency = HttpClientErrorStatus.NewSubtype( "Status Failed Dependency")
var StatusTooEarly = HttpClientErrorStatus.NewSubtype( "Status Too Early")
var StatusUpgradeRequired = HttpClientErrorStatus.NewSubtype( "Status Upgrade Required")
var StatusPreconditionRequired = HttpClientErrorStatus.NewSubtype( "Status Precondition Required")
var StatusTooManyRequests = HttpClientErrorStatus.NewSubtype( "Status Too Many Requests")
var StatusRequestHeaderFieldsTooLarge = HttpClientErrorStatus.NewSubtype( "Status Request Header Fields Too Large", InvalidSizeTrait)
var StatusUnavailableForLegalReasons = HttpClientErrorStatus.NewSubtype( "Status Unavailable For Legal Reasons")

var HttpServerErrorStatus = errorx.NewType(HttpNamespace, "Server Error", HttpServerTrait)
var StatusInternalServerError = HttpServerErrorStatus.NewSubtype( "Status Internal Server Error")
var StatusNotImplemented = HttpServerErrorStatus.NewSubtype( "Status Not Implemented")
var StatusBadGateway = HttpServerErrorStatus.NewSubtype( "Status Bad Gateway",HttpServerTrait)
var StatusServiceUnavailable = HttpServerErrorStatus.NewSubtype( "Status Service Unavailable")
var StatusGatewayTimeout = HttpServerErrorStatus.NewSubtype( "Status Gateway Timeout")
var StatusHTTPVersionNotSupported = HttpServerErrorStatus.NewSubtype( "Status HTTP Version Not Supported")
var StatusVariantAlsoNegotiates = HttpServerErrorStatus.NewSubtype( "Status Variant Also Negotiates")
var StatusInsufficientStorage = HttpServerErrorStatus.NewSubtype( "Status Insufficient Storage", InvalidSizeTrait)
var StatusLoopDetected = HttpServerErrorStatus.NewSubtype( "Status Loop Detected")
var StatusNotExtended = HttpServerErrorStatus.NewSubtype( "Status Not Extended")
var StatusNetworkAuthenticationRequired = HttpServerErrorStatus.NewSubtype( "Status Network Authentication Required")
