package global

// ResponseOK 成功响应数据格式
type ResponseOK struct {
	TraceId string      `json:"trace_id"`
	Data    interface{} `json:"data"`
}

// ResponseErr 失败响应数据格式
type ResponseErr struct {
	TraceId string `json:"trace_id"`
	Code    int    `json:"code"`    // 业务码
	Message string `json:"message"` // 描述信息
	Desc    string `json:"desc"`    // 描述说明
}
