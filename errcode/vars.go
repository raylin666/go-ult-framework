package errcode

var (
	ErrServerError          = New(ServerError)
	ErrAuthorizationError   = New(AuthorizationError)
	ErrParamBindError       = New(ParamBindError)
	ErrRequestError         = New(RequestError)
	ErrParamValidateError   = New(ParamValidateError)
	ErrUnknownError         = New(UnknownError)
	ErrDataNotExistError    = New(DataNotExistError)
	ErrDataExistError       = New(DataExistError)
	ErrRequestNotFoundError = New(RequestNotFoundError)
	ErrDataDeleteError      = New(DataDeleteError)
	ErrDataSelectError      = New(DataSelectError)
	ErrDataCreateError      = New(DataCreateError)
	ErrDataUpdateError      = New(DataUpdateError)
)