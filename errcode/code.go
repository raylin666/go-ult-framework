// Package errcode 提供统一错误码定义和管理。
package errcode

// 语言常量定义。
const (
	ZhCN = "zh-cn" // 中文
	EnUS = "en-us" // 英文
)

// 业务错误码常量定义。
// 100xxx: 服务端错误
// 200xxx: 客户端错误
const (
	ServerError           = 100001 // 内部服务器错误

	AuthorizationError    = 200001 // 签名信息错误
	ParamBindError        = 200002 // 参数信息错误
	RequestError          = 200003 // 请求错误
	ParamValidateError    = 200004 // 参数校验错误
	UnknownError          = 200005 // 未知错误
	DataNotExistError     = 200006 // 数据不存在
	DataExistError        = 200007 // 数据已存在
	RequestNotFoundError  = 200008 // 不存在的请求
	DataDeleteError       = 200009 // 数据删除错误
	ResourceNotExistError = 200010 // 资源不存在
	DataSelectError       = 200011 // 数据查询失败
	DataCreateError       = 200012 // 数据创建失败
	DataUpdateError       = 200013 // 数据更新失败
)