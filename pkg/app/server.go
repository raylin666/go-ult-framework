package app

import "github.com/raylin666/go-utils/v2/server"

type ServerAgreement struct {
	Network string
	Addr    string
	Target  string
}

type Server interface {
	server.Server

	ServerType() string
	StartBefore()
	StartAfter()
	CancelBefore()
	CancelAfter()
}