package http

type SuccessResponse struct {
	TraceId string      `json:"trace_id"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	TraceId string `json:"trace_id"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Desc    string `json:"desc"`
}

func NewSuccessResponse(traceId string, data interface{}) SuccessResponse {
	return SuccessResponse{traceId, data}
}

func NewErrorResponse(traceId string, code int, message, desc string) ErrorResponse {
	return ErrorResponse{traceId, code, message, desc}
}