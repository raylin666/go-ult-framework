package global

import "github.com/raylin666/go-utils/server"

type ServerAgreement struct {
	Network string
	Addr    string
	Target  string
}

type Server interface {
	server.Server

	// ServerType 服务类型
	ServerType() string
	// StartBefore 服务启动之前操作
	StartBefore()
	// StartAfter 服务启动之后操作
	StartAfter()
	// CancelBefore 服务取消之前操作
	CancelBefore()
	// CancelAfter 服务取消之后操作
	CancelAfter()
}
