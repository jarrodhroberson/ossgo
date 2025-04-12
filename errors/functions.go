package errors

import (
	"github.com/joomcode/errorx"
)

var statusCodeToError = make(map[string]*errorx.Type)

func init() {
	statusCodeToError["300"] = StatusMultipleChoices // Multiple Choices
	statusCodeToError["301"] = StatusMovedPermanently
	statusCodeToError["302"] = StatusFound
	statusCodeToError["303"] = StatusSeeOther
	statusCodeToError["304"] = StatusNotModified
	statusCodeToError["307"] = StatusTemporaryRedirect
	statusCodeToError["308"] = StatusPermanentRedirect

	statusCodeToError["400"] = StatusBadRequest
	statusCodeToError["401"] = StatusUnauthorized
	statusCodeToError["402"] = StatusPaymentRequired
	statusCodeToError["403"] = StatusForbidden
	statusCodeToError["404"] = StatusNotFound
	statusCodeToError["405"] = StatusMethodNotAllowed
	statusCodeToError["406"] = StatusNotAcceptable
	statusCodeToError["407"] = StatusProxyAuthRequired
	statusCodeToError["408"] = StatusRequestTimeout
	statusCodeToError["409"] = StatusConflict
	statusCodeToError["410"] = StatusGone
	statusCodeToError["411"] = StatusLengthRequired
	statusCodeToError["412"] = StatusPreconditionFailed
	statusCodeToError["413"] = StatusRequestEntityTooLarge
	statusCodeToError["414"] = StatusRequestURITooLong
	statusCodeToError["415"] = StatusUnsupportedMediaType
	statusCodeToError["416"] = StatusRequestedRangeNotSatisfiable
	statusCodeToError["417"] = StatusExpectationFailed
	statusCodeToError["421"] = StatusMisdirectedRequest
	statusCodeToError["422"] = StatusUnprocessableEntity
	statusCodeToError["423"] = StatusLocked
	statusCodeToError["424"] = StatusFailedDependency
	statusCodeToError["425"] = StatusTooEarly
	statusCodeToError["426"] = StatusUpgradeRequired
	statusCodeToError["428"] = StatusPreconditionRequired
	statusCodeToError["429"] = StatusTooManyRequests
	statusCodeToError["431"] = StatusRequestHeaderFieldsTooLarge
	statusCodeToError["451"] = StatusUnavailableForLegalReasons

	statusCodeToError["500"] = StatusInternalServerError
	statusCodeToError["501"] = StatusNotImplemented
	statusCodeToError["502"] = StatusBadGateway
	statusCodeToError["503"] = StatusServiceUnavailable
	statusCodeToError["504"] = StatusGatewayTimeout
	statusCodeToError["505"] = StatusHTTPVersionNotSupported
	statusCodeToError["506"] = StatusVariantAlsoNegotiates
	statusCodeToError["507"] = StatusInsufficientStorage
	statusCodeToError["508"] = StatusLoopDetected
	statusCodeToError["510"] = StatusNotExtended
	statusCodeToError["511"] = StatusNetworkAuthenticationRequired
}

func FromHttpStatusCode(code string) *errorx.Type {
	return statusCodeToError[code]
}
