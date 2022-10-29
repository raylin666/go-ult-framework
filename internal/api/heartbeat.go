package api

import (
	"ult/internal/service"
	"ult/pkg/http"
	"ult/pkg/logger"
)

var _ HeartbeatInterface = (*HeartbeatHandler)(nil)

type HeartbeatInterface interface {
	PONE() http.HandlerFunc
}

type HeartbeatHandler struct {
	logger  *logger.Logger
	service *service.HeartbeatService
}

func NewHeartbeatHandler(logger *logger.Logger, service *service.HeartbeatService) HeartbeatInterface {
	return &HeartbeatHandler{
		logger:  logger,
		service: service,
	}
}

func (h *HeartbeatHandler) PONE() http.HandlerFunc {
	return func(ctx http.Context) {
		ctx.WithPayload(h.service.PONE(ctx.RequestContext()))
	}
}
