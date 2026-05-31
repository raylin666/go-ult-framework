package api

import (
	"ult/internal/app"
	"ult/internal/service"
	"ult/pkg/http"
)

var _ HeartbeatInterface = (*HeartbeatHandler)(nil)

type HeartbeatInterface interface {
	State() http.HandlerFunc
}

type HeartbeatHandler struct {
	service *service.HealtbeatService
	tools   *app.Tools
}

func NewHeartbeatHandler(service *service.HealtbeatService, tools *app.Tools) HeartbeatInterface {
	return &HeartbeatHandler{
		service: service,
		tools:   tools,
	}
}

func (h *HeartbeatHandler) State() http.HandlerFunc {
	return func(ctx http.Context) {
		var resp = h.service.State(ctx.RequestContext())
		ctx.WithPayload(resp)
	}
}
