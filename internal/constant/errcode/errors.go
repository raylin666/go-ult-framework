package errcode

import "ult/pkg/code"

var (
	ErrorServerError          = NewError(code.ServerError)
	ErrorAuthorizationError   = NewError(code.AuthorizationError)
	ErrorParamBindError       = NewError(code.ParamBindError)
	ErrorRequestError         = NewError(code.RequestError)
	ErrorParamValidateError   = NewError(code.ParamValidateError)
	ErrorUnknownError         = NewError(code.UnknownError)
	ErrorDataNotExistError    = NewError(code.DataNotExistError)
	ErrorDataExistError       = NewError(code.DataExistError)
	ErrorRequestNotFoundError = NewError(code.RequestNotFoundError)
	ErrorDataDeleteError      = NewError(code.DataDeleteError)
	ErrorDataSelectError      = NewError(code.DataSelectError)
	ErrorDataCreateError      = NewError(code.DataCreateError)
	ErrorDataUpdateError      = NewError(code.DataUpdateError)
)
