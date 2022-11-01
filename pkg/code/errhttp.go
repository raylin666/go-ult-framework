package code

/**
HTTP 状态码设定
*/

var httpCode = map[int]int{
	/* 系统相关 */
	ServerError:           500,
	AuthorizationError:    401,
	ParamBindError:        400,
	RequestError:          400,
	ParamValidateError:    422,
	UnknownError:          400,
	DataNotExistError:     400,
	DataExistError:        400,
	RequestNotFoundError:  400,
	DataDeleteError:       400,
	ResourceNotExistError: 400,
	DataSelectError:       400,
	DataCreateError:       400,
	DataUpdateError:       400,
}
