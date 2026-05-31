package errcode

const (
	ZhCN = "zh-cn"
	EnUS = "en-us"
)

const (
	ServerError           = 100001

	AuthorizationError    = 200001
	ParamBindError        = 200002
	RequestError          = 200003
	ParamValidateError    = 200004
	UnknownError          = 200005
	DataNotExistError     = 200006
	DataExistError        = 200007
	RequestNotFoundError  = 200008
	DataDeleteError       = 200009
	ResourceNotExistError = 200010
	DataSelectError       = 200011
	DataCreateError       = 200012
	DataUpdateError       = 200013
)