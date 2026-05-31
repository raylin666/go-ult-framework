package errcode

import "net/http"

var httpCode = map[int]int{
	ServerError:           http.StatusInternalServerError,
	AuthorizationError:    http.StatusUnauthorized,
	ParamBindError:        http.StatusBadRequest,
	RequestError:          http.StatusBadRequest,
	ParamValidateError:    http.StatusUnprocessableEntity,
	UnknownError:          http.StatusBadRequest,
	DataNotExistError:     http.StatusBadRequest,
	DataExistError:        http.StatusBadRequest,
	RequestNotFoundError:  http.StatusBadRequest,
	DataDeleteError:       http.StatusBadRequest,
	ResourceNotExistError: http.StatusBadRequest,
	DataSelectError:       http.StatusBadRequest,
	DataCreateError:       http.StatusBadRequest,
	DataUpdateError:       http.StatusBadRequest,
}