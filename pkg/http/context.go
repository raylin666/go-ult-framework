package http

import (
	"bytes"
	stdCtx "context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/raylin666/go-utils/validator"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"ult/pkg/code"
	"ult/pkg/errors"
)

const (
	headerXTraceIdName = "X-Trace-Id"

	_BodyName_       = "_body_"
	_PayloadName_    = "_payload_"
	_AbortErrorName_ = "_abort_error_"
	_ValidatorName_  = "_validator_"
	_TraceIdName_    = "_trace_id_"
)

type HandlerFunc func(ctx Context)

var _ Context = (*context)(nil)

type Context interface {
	init()

	// ShouldBindQuery 反序列化 querystring
	// tag: `form:"xxx"` (注：不要写成 query)
	ShouldBindQuery(obj interface{}) error

	// ShouldBindPostForm 反序列化 postform (querystring 会被忽略)
	// tag: `form:"xxx"`
	ShouldBindPostForm(obj interface{}) error

	// ShouldBindForm 同时反序列化 querystring 和 postform;
	// 当 querystring 和 postform 存在相同字段时，postform 优先使用。
	// tag: `form:"xxx"`
	ShouldBindForm(obj interface{}) error

	// ShouldBindJSON 反序列化 postjson
	// tag: `json:"xxx"`
	ShouldBindJSON(obj interface{}) error

	// ShouldBindURI 反序列化 path 参数(如路由路径为 /user/:name)
	// tag: `uri:"xxx"`
	ShouldBindURI(obj interface{}) error

	// Redirect 重定向
	Redirect(code int, location string)

	// Param 获取路径参数
	Param(key string) string

	// TraceID 获取链路追踪ID
	TraceID() string

	// Validator 数据验证器
	Validator(req interface{}) bool
	WithValidator(v validator.Validator)

	// WithAbortError 错误返回
	WithAbortError(err errors.BusinessError)
	getAbortError() errors.BusinessError

	// WithPayload 正确返回
	WithPayload(payload interface{})
	getPayload() interface{}

	// Header 获取 Header 对象
	Header() http.Header
	// GetHeader 获取 Header
	GetHeader(key string) string
	// SetHeader 设置 Header
	SetHeader(key, value string)

	// RequestInputParams 获取所有参数
	RequestInputParams() url.Values
	// RequestPostFormParams  获取 PostForm 参数
	RequestPostFormParams() url.Values
	// Request 获取 Request 对象
	Request() *http.Request
	// RawData 获取 Request.Body
	RawData() []byte
	// Method 获取 Request.Method
	Method() string
	// Host 获取 Request.Host
	Host() string
	// Path 获取 请求的路径 Request.URL.Path (不附带 querystring)
	Path() string
	// URI 获取 unescape 后的 Request.URL.RequestURI()
	URI() string
	// RequestContext 获取请求的 Context (当 client 关闭后，会自动 canceled)
	RequestContext() stdCtx.Context

	// ResponseWriter 获取 ResponseWriter 对象
	ResponseWriter() gin.ResponseWriter
}

type context struct {
	ctx *gin.Context
}

func (c *context) init() {
	body, err := c.ctx.GetRawData()
	if err != nil {
		panic(err)
	}

	c.ctx.Set(_BodyName_, body)
	c.ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // re-construct req body
}

// ShouldBindQuery 反序列化querystring
// tag: `form:"xxx"` (注：不要写成query)
func (c *context) ShouldBindQuery(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Query)
}

// ShouldBindPostForm 反序列化 postform (querystring 会被忽略)
// tag: `form:"xxx"`
func (c *context) ShouldBindPostForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.FormPost)
}

// ShouldBindForm 同时反序列化querystring和postform;
// 当querystring和postform存在相同字段时，postform优先使用。
// tag: `form:"xxx"`
func (c *context) ShouldBindForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Form)
}

// ShouldBindJSON 反序列化postjson
// tag: `json:"xxx"`
func (c *context) ShouldBindJSON(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.JSON)
}

// ShouldBindURI 反序列化path参数(如路由路径为 /user/:name)
// tag: `uri:"xxx"`
func (c *context) ShouldBindURI(obj interface{}) error {
	return c.ctx.ShouldBindUri(obj)
}

// Param 获取路径参数
func (c *context) Param(key string) string {
	return c.ctx.Param(key)
}

// Redirect 重定向
func (c *context) Redirect(code int, location string) {
	c.ctx.Redirect(code, location)
}

// TraceID 获取链路追踪ID
func (c *context) TraceID() string {
	traceId, ok := c.ctx.Get(_TraceIdName_)
	if ok {
		return traceId.(string)
	}

	var headerTraceId = c.GetHeader(headerXTraceIdName)
	if len(headerTraceId) <= 0 {
		headerTraceId = uuid.New().String()
	}

	c.ctx.Set(_TraceIdName_, headerTraceId)
	return headerTraceId
}

func (c *context) Validator(req interface{}) (isErr bool) {
	// 参数数据绑定
	if err := c.ShouldBindForm(req); err != nil {
		c.WithAbortError(errors.NewError(
			http.StatusBadRequest,
			code.ParamBindError,
			code.Get().GetText(code.ParamBindError)).WithStackError(err))

		return true
	}

	// 参数数据验证
	validate, ok := c.ctx.Get(_ValidatorName_)
	if !ok {
		return true
	}

	if errStr := validate.(validator.Validator).Validate(req); errStr != "" {
		c.WithAbortError(errors.NewError(
			http.StatusUnprocessableEntity,
			code.ParamValidateError,
			errStr))

		return true
	}

	return false
}

func (c *context) WithValidator(v validator.Validator) {
	c.ctx.Set(_ValidatorName_, v)
}

func (c *context) WithAbortError(err errors.BusinessError) {
	if err != nil {
		httpCode := err.HTTPCode()
		if httpCode == 0 {
			httpCode = http.StatusInternalServerError
		}

		c.ctx.AbortWithStatus(httpCode)
		c.ctx.Set(_AbortErrorName_, err)
	}
}

func (c *context) getAbortError() errors.BusinessError {
	err, _ := c.ctx.Get(_AbortErrorName_)
	return err.(errors.BusinessError)
}

func (c *context) WithPayload(payload interface{}) {
	c.ctx.Set(_PayloadName_, payload)
}

func (c *context) getPayload() interface{} {
	if payload, ok := c.ctx.Get(_PayloadName_); ok != false {
		return payload
	}
	return nil
}

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

func (c *context) GetHeader(key string) string {
	return c.ctx.GetHeader(key)
}

func (c *context) SetHeader(key, value string) {
	c.ctx.Header(key, value)
}

// RequestInputParams 获取所有参数
func (c *context) RequestInputParams() url.Values {
	_ = c.ctx.Request.ParseForm()
	return c.ctx.Request.Form
}

// RequestPostFormParams 获取 PostForm 参数
func (c *context) RequestPostFormParams() url.Values {
	_ = c.ctx.Request.ParseForm()
	return c.ctx.Request.PostForm
}

// Request 获取请求 Request
func (c *context) Request() *http.Request {
	return c.ctx.Request
}

func (c *context) RawData() []byte {
	body, ok := c.ctx.Get(_BodyName_)
	if !ok {
		return nil
	}

	return body.([]byte)
}

// Method 获取请求 Method
func (c *context) Method() string {
	return c.ctx.Request.Method
}

// Host 获取请求 Host
func (c *context) Host() string {
	return c.ctx.Request.Host
}

// Path 请求的路径(不附带querystring)
func (c *context) Path() string {
	return c.ctx.Request.URL.Path
}

// URI Unescape 后的 URI
func (c *context) URI() string {
	uri, _ := url.QueryUnescape(c.ctx.Request.URL.RequestURI())
	return uri
}

// RequestContext 获取请求的 Context (当client关闭后，会自动canceled)
func (c *context) RequestContext() stdCtx.Context {
	return c.ctx.Request.Context()
}

// ResponseWriter 获取 ResponseWriter
func (c *context) ResponseWriter() gin.ResponseWriter {
	return c.ctx.Writer
}

var contextPool = &sync.Pool{
	New: func() interface{} {
		return new(context)
	},
}

// newContext 创建 Context 上下文切换
func newContext(ctx *gin.Context) Context {
	context := contextPool.Get().(*context)
	context.ctx = ctx
	return context
}

// recoveryContext 回收 Context 上下文切换
func recoveryContext(ctx Context) {
	c := ctx.(*context)
	c.ctx = nil
	contextPool.Put(c)
}
