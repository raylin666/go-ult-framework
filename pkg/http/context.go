// Package http 提供 HTTP 服务器实现，基于 Gin 框架封装。
package http

import (
	"bytes"
	stdCtx "context"
	"io"
	"net/http"
	"net/url"
	"sync"
	"ult/errcode"
	"ult/pkg/types"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/raylin666/go-utils/v2/validator"
)

// 内部上下文键，用于存储请求数据。
const (
	_BodyName_       = "_body_"        // 存储原始请求体的键
	_PayloadName_    = "_payload_"     // 存储响应数据的键
	_AbortErrorName_ = "_abort_error_" // 存储中止错误的键
	_ValidatorName_  = "_validator_"   // 存储验证器实例的键
)

// HandlerFunc 定义本包使用的处理函数类型。
// 接受自定义 Context 接口而非直接使用 Gin.Context。
type HandlerFunc func(ctx Context)

// Context 接口验证。
var _ Context = (*context)(nil)

// Context 定义 HTTP 请求上下文接口。
// 封装 Gin.Context 并提供额外功能：
// - 请求绑定（query、form、JSON、URI）
// - 响应处理（payload、error）
// - 请求信息访问
// - 分布式追踪的 TraceID 管理
type Context interface {
	init()

	// ShouldBindQuery 反序列化查询字符串参数。
	// 在结构体字段中使用 `form:"xxx"` 标签（不是 `query`）。
	ShouldBindQuery(obj interface{}) error

	// ShouldBindPostForm 反序列化 POST 表单数据（忽略查询字符串）。
	// 在结构体字段中使用 `form:"xxx"` 标签。
	ShouldBindPostForm(obj interface{}) error

	// ShouldBindForm 反序列化查询字符串和 POST 表单。
	// 当同一字段同时存在时，POST 表单优先。
	// 在结构体字段中使用 `form:"xxx"` 标签。
	ShouldBindForm(obj interface{}) error

	// ShouldBindJSON 反序列化 JSON 请求体。
	// 在结构体字段中使用 `json:"xxx"` 标签。
	ShouldBindJSON(obj interface{}) error

	// ShouldBindURI 反序列化路径参数（如 /user/:name）。
	// 在结构体字段中使用 `uri:"xxx"` 标签。
	ShouldBindURI(obj interface{}) error

	// Redirect 执行 HTTP 重定向。
	Redirect(code int, location string)

	// Param 通过键名返回路径参数值。
	Param(key string) string

	// TraceID 返回分布式追踪的 TraceID。
	// 如果不存在，则生成新的 UUID 并存储。
	TraceID() string

	// Validator 使用配置的验证器验证请求结构体。
	// 如果验证失败（发生错误）则返回 true。
	Validator(req interface{}) bool
	WithValidator(v validator.Validator)

	// WithAbortError 设置错误以中止请求。
	// 该错误将用于响应处理。
	WithAbortError(err errcode.BusinessError)
	getAbortError() errcode.BusinessError

	// WithPayload 设置成功响应的数据负载。
	WithPayload(payload interface{})
	getPayload() interface{}

	// Header 返回请求头的克隆副本。
	Header() http.Header
	// GetHeader 通过键名返回特定的请求头值。
	GetHeader(key string) string
	// SetHeader 设置响应头。
	SetHeader(key, value string)

	// RequestInputParams 返回所有查询和表单参数。
	RequestInputParams() url.Values
	// RequestPostFormParams 仅返回 POST 表单参数。
	RequestPostFormParams() url.Values
	// Request 返回底层的 http.Request 对象。
	Request() *http.Request
	// RawData 返回原始请求体字节。
	RawData() []byte
	// Method 返回 HTTP 方法。
	Method() string
	// Host 返回请求主机。
	Host() string
	// Path 返回不带查询字符串的请求路径。
	Path() string
	// URI 返回未转义的请求 URI。
	URI() string
	// RequestContext 返回带有 TraceID 的请求上下文。
	// 当客户端关闭连接时，该上下文会被取消。
	RequestContext() stdCtx.Context

	// ResponseWriter 返回 Gin 响应写入器。
	ResponseWriter() gin.ResponseWriter
}

// context 是 Context 接口的内部实现。
// 封装 gin.Context 并存储额外的请求/响应数据。
type context struct {
	ctx *gin.Context
}

// init 初始化上下文，读取并存储原始请求体。
// 这允许在请求处理过程中多次读取请求体。
func (c *context) init() {
	body, err := c.ctx.GetRawData()
	if err != nil {
		c.ctx.AbortWithStatus(http.StatusInternalServerError)
		c.ctx.Set(_AbortErrorName_, errcode.New(errcode.ServerError).WithDesc(err.Error()))
		return
	}

	c.ctx.Set(_BodyName_, body)
	c.ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
}

// ShouldBindQuery 将查询字符串参数绑定到结构体。
func (c *context) ShouldBindQuery(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Query)
}

// ShouldBindPostForm 将 POST 表单数据绑定到结构体。
func (c *context) ShouldBindPostForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.FormPost)
}

// ShouldBindForm 将查询字符串和 POST 表单绑定到结构体。
func (c *context) ShouldBindForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Form)
}

// ShouldBindJSON 将 JSON 请求体绑定到结构体。
func (c *context) ShouldBindJSON(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.JSON)
}

// ShouldBindURI 将 URI 路径参数绑定到结构体。
func (c *context) ShouldBindURI(obj interface{}) error {
	return c.ctx.ShouldBindUri(obj)
}

// Param 返回 URI 路径参数的值。
func (c *context) Param(key string) string {
	return c.ctx.Param(key)
}

// Redirect 执行 HTTP 重定向到指定位置。
func (c *context) Redirect(code int, location string) {
	c.ctx.Redirect(code, location)
}

// TraceID 返回分布式追踪的 TraceID。
// 首先检查上下文或请求头中是否存在 TraceID。
// 如果未找到，则生成新的 UUID 并存储。
func (c *context) TraceID() string {
	traceId, ok := c.ctx.Get(types.TraceIdName)
	if ok {
		if tid, ok := traceId.(string); ok {
			return tid
		}
	}

	var headerTraceId = c.GetHeader(types.TraceIdName)
	if len(headerTraceId) <= 0 {
		headerTraceId = uuid.New().String()
	}

	c.ctx.Set(types.TraceIdName, headerTraceId)
	return headerTraceId
}

// Validator 执行请求验证：绑定和校验。
// 如果验证失败返回 true，成功返回 false。
//
// 该方法：
// 1. 将表单数据绑定到请求结构体
// 2. 使用配置的验证器验证结构体
// 3. 如果验证失败则设置适当的错误
func (c *context) Validator(req interface{}) (isErr bool) {
	if err := c.ShouldBindForm(req); err != nil {
		c.WithAbortError(errcode.New(errcode.ParamValidateError).WithStackError(err))
		return true
	}

	validate, ok := c.ctx.Get(_ValidatorName_)
	if !ok {
		return true
	}

	validatorInst, ok := validate.(validator.Validator)
	if !ok {
		c.WithAbortError(errcode.New(errcode.ServerError).WithDesc("validator not found"))
		return true
	}
	if errStr := validatorInst.Validate(req); errStr != nil {
		c.WithAbortError(errcode.New(errcode.ParamValidateError).WithDesc(errStr.Error()))
		return true
	}

	return false
}

// WithValidator 设置请求验证的验证器实例。
func (c *context) WithValidator(v validator.Validator) {
	c.ctx.Set(_ValidatorName_, v)
}

// WithAbortError 设置业务错误并中止请求。
// HTTP 状态码由错误的 HTTPCode 方法决定。
func (c *context) WithAbortError(err errcode.BusinessError) {
	if err != nil {
		httpCode := err.HTTPCode()
		if httpCode == 0 {
			httpCode = http.StatusInternalServerError
		}

		c.ctx.AbortWithStatus(httpCode)
		c.ctx.Set(_AbortErrorName_, err)
	}
}

func (c *context) getAbortError() errcode.BusinessError {
	err, ok := c.ctx.Get(_AbortErrorName_)
	if !ok {
		return nil
	}
	businessErr, ok := err.(errcode.BusinessError)
	if !ok {
		return nil
	}
	return businessErr
}

// WithPayload 设置成功响应的数据负载。
func (c *context) WithPayload(payload interface{}) {
	c.ctx.Set(_PayloadName_, payload)
}

// getPayload 返回存储的响应数据负载。
func (c *context) getPayload() interface{} {
	if payload, ok := c.ctx.Get(_PayloadName_); ok != false {
		return payload
	}
	return nil
}

// Header 返回请求头的克隆副本。
// 该克隆副本可以安全修改而不影响原始请求头。
func (c *context) Header() http.Header {
	header := c.ctx.Request.Header
	clone := make(http.Header, len(header))
	for k, v := range header {
		value := make([]string, len(v))
		copy(value, v)

		clone[k] = value
	}

	return clone
}

// GetHeader 通过键名返回特定的请求头值。
func (c *context) GetHeader(key string) string {
	return c.ctx.GetHeader(key)
}

// SetHeader 设置响应头。
func (c *context) SetHeader(key, value string) {
	c.ctx.Header(key, value)
}

// RequestInputParams 返回所有查询和表单参数的组合。
func (c *context) RequestInputParams() url.Values {
	_ = c.ctx.Request.ParseForm()
	return c.ctx.Request.Form
}

// RequestPostFormParams 仅返回 POST 表单参数。
func (c *context) RequestPostFormParams() url.Values {
	_ = c.ctx.Request.ParseForm()
	return c.ctx.Request.PostForm
}

// Request 返回底层的 http.Request 对象。
func (c *context) Request() *http.Request {
	return c.ctx.Request
}

// RawData 返回存储的原始请求体字节。
func (c *context) RawData() []byte {
	body, ok := c.ctx.Get(_BodyName_)
	if !ok {
		return nil
	}

	bodyBytes, ok := body.([]byte)
	if !ok {
		return nil
	}
	return bodyBytes
}

// Method 返回 HTTP 方法（GET、POST 等）。
func (c *context) Method() string {
	return c.ctx.Request.Method
}

// Host 返回请求主机。
func (c *context) Host() string {
	return c.ctx.Request.Host
}

// Path 返回不带查询字符串的请求路径。
func (c *context) Path() string {
	return c.ctx.Request.URL.Path
}

// URI 返回带查询字符串的未转义请求 URI。
func (c *context) URI() string {
	uri, _ := url.QueryUnescape(c.ctx.Request.URL.RequestURI())
	return uri
}

// RequestContext 返回带有 TraceID 的请求上下文，用于分布式追踪。
// 当客户端关闭连接时，该上下文会被取消。
func (c *context) RequestContext() stdCtx.Context {
	var reqContext = new(types.RequestContext)
	reqContext.WithTraceID(c.TraceID())
	return types.NewRequestContext(c.ctx.Request.Context(), reqContext)
}

// ResponseWriter 返回 Gin 响应写入器，用于直接操作响应。
func (c *context) ResponseWriter() gin.ResponseWriter {
	return c.ctx.Writer
}

// contextPool 是用于复用上下文对象的 sync.Pool。
// 通过减少内存分配来提高性能。
var contextPool = &sync.Pool{
	New: func() interface{} {
		return new(context)
	},
}

// newContext 从池中创建或获取上下文对象。
// 用额外功能封装 Gin 上下文。
func newContext(ctx *gin.Context) Context {
	context := contextPool.Get().(*context)
	context.ctx = ctx
	return context
}

// recoveryContext 使用后将上下文归还到池中。
// 清除 Gin 上下文引用以防止内存泄漏。
func recoveryContext(ctx Context) {
	c, ok := ctx.(*context)
	if !ok {
		return
	}
	c.ctx = nil
	contextPool.Put(c)
}
