package api

import (
	"ult/internal/app"
	"ult/internal/service"
	"ult/pkg/http"
)

var _ HeartbeatInterface = (*HeartbeatHandler)(nil)

type HeartbeatInterface interface {
	PONE() http.HandlerFunc
}

type HeartbeatHandler struct {
	service *service.HeartbeatService
	tools   *app.Tools
}

func NewHeartbeatHandler(service *service.HeartbeatService, tools *app.Tools) HeartbeatInterface {
	return &HeartbeatHandler{
		service: service,
		tools:   tools,
	}
}

func (h *HeartbeatHandler) PONE() http.HandlerFunc {
	return func(ctx http.Context) {
		var resp = h.service.PONE(ctx.RequestContext())
		ctx.WithPayload(resp)
	}
}
