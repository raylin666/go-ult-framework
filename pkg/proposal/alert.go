// Package proposal 提供告警消息定义和处理。
// 定义告警消息结构和通知处理函数类型。
package proposal

import (
	"encoding/json"
	"time"
)

// AlertMessage 告警消息结构体。
// 包含项目信息、请求信息、错误信息和时间戳等。
type AlertMessage struct {
	ProjectName  string      `json:"project_name"`  // 项目名，用于区分不同项目告警信息
	Environment  string      `json:"environment"`   // 运行环境
	TraceID      string      `json:"trace_id"`      // 唯一 ID，用于链路追踪关联
	HOST         string      `json:"host"`          // 请求 HOST
	URI          string      `json:"uri"`           // 请求 URI
	Method       string      `json:"method"`        // 请求 Method
	ErrorMessage interface{} `json:"error_message"` // 错误信息
	ErrorStack   string      `json:"error_stack"`   // 堆栈信息
	Timestamp    time.Time   `json:"timestamp"`     // 时间戳
}

// Marshal 序列化告警消息为 JSON。
//
// 返回:
//   - jsonRaw: JSON 字节数组
func (a *AlertMessage) Marshal() (jsonRaw []byte) {
	jsonRaw, _ = json.Marshal(a)
	return
}

// NotifyHandler 告警通知处理函数类型。
// 用于发送告警消息（如邮件、钉钉等）。
type NotifyHandler func(msg *AlertMessage)
