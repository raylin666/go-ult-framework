// Package http 提供 HTTP 服务器实现，基于 Gin 框架封装。
package http

// SuccessResponse 成功响应结构体。
type SuccessResponse struct {
	TraceId string      `json:"trace_id"` // 链路追踪 ID
	Data    interface{} `json:"data"`     // 响应数据
}

// ErrorResponse 错误响应结构体。
type ErrorResponse struct {
	TraceId string `json:"trace_id"` // 链路追踪 ID
	Code    int    `json:"code"`     // 业务错误码
	Message string `json:"message"`  // 错误消息
	Desc    string `json:"desc"`     // 错误描述
}

// NewSuccessResponse 创建成功响应。
//
// 参数:
//   - traceId: 链路追踪 ID
//   - data: 响应数据
//
// 返回:
//   - SuccessResponse: 成功响应结构体
func NewSuccessResponse(traceId string, data interface{}) SuccessResponse {
	return SuccessResponse{traceId, data}
}

// NewErrorResponse 创建错误响应。
//
// 参数:
//   - traceId: 链路追踪 ID
//   - code: 业务错误码
//   - message: 错误消息
//   - desc: 错误描述
//
// 返回:
//   - ErrorResponse: 错误响应结构体
func NewErrorResponse(traceId string, code int, message, desc string) ErrorResponse {
	return ErrorResponse{traceId, code, message, desc}
}
