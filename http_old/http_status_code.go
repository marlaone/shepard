package http

import "strconv"

type StatusCode int

const (
	StatusCodeOk                           StatusCode = 200
	StatusCodeCreated                      StatusCode = 201
	StatusCodeAccepted                     StatusCode = 202
	StatusCodeNonAuthoritativeInformation  StatusCode = 203
	StatusCodeNoContent                    StatusCode = 204
	StatusCodeResetContent                 StatusCode = 205
	StatusCodePartialContent               StatusCode = 206
	StatusCodeMultipleChoices              StatusCode = 300
	StatusCodeMovedPermanently             StatusCode = 301
	StatusCodeFound                        StatusCode = 302
	StatusCodeSeeOther                     StatusCode = 303
	StatusCodeNotModified                  StatusCode = 304
	StatusCodeUseProxy                     StatusCode = 305
	StatusCodeTemporaryRedirect            StatusCode = 307
	StatusCodeBadRequest                   StatusCode = 400
	StatusCodeUnauthorized                 StatusCode = 401
	StatusCodePaymentRequired              StatusCode = 402
	StatusCodeForbidden                    StatusCode = 403
	StatusCodeNotFound                     StatusCode = 404
	StatusCodeMethodNotAllowed             StatusCode = 405
	StatusCodeNotAcceptable                StatusCode = 406
	StatusCodeProxyAuthenticationRequired  StatusCode = 407
	StatusCodeRequestTimeout               StatusCode = 408
	StatusCodeConflict                     StatusCode = 409
	StatusCodeGone                         StatusCode = 410
	StatusCodeLengthRequired               StatusCode = 411
	StatusCodePreconditionFailed           StatusCode = 412
	StatusCodeRequestEntityTooLarge        StatusCode = 413
	StatusCodeRequestURITooLong            StatusCode = 414
	StatusCodeUnsupportedMediaType         StatusCode = 415
	StatusCodeRequestedRangeNotSatisfiable StatusCode = 416
	StatusCodeExpectationFailed            StatusCode = 417
	StatusCodeInternalServerError          StatusCode = 500
	StatusCodeNotImplemented               StatusCode = 501
	StatusCodeBadGateway                   StatusCode = 502
	StatusCodeServiceUnavailable           StatusCode = 503
	StatusCodeGatewayTimeout               StatusCode = 504
	StatusCodeHTTPVersionNotSupported      StatusCode = 505
)

func (s StatusCode) Default() StatusCode {
	return StatusCodeOk
}

func (s StatusCode) String() string {
	return strconv.Itoa(int(s))
}
