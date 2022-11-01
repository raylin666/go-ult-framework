package global

// SuccessResponse 成功响应数据格式
type SuccessResponse struct {
	TraceId string      `json:"trace_id"`
	Data    interface{} `json:"data"`
}

// ErrorResponse 失败响应数据格式
type ErrorResponse struct {
	TraceId string `json:"trace_id"`
	Code    int    `json:"code"`    // 业务码
	Message string `json:"message"` // 描述信息
	Desc    string `json:"desc"`    // 描述说明
}

func NewSuccessResponse(traceId string, data interface{}) SuccessResponse {
	return SuccessResponse{traceId, data}
}

func NewErrorResponse(traceId string, code int, message, desc string) ErrorResponse {
	return ErrorResponse{traceId, code, message, desc}
}
