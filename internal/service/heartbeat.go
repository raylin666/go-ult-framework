package service

import (
	"context"
	"ult/pkg/logger"
)

type HeartbeatService struct {
	logger   *logger.Logger
}

func NewHeartbeatService(logger *logger.Logger) *HeartbeatService {
	return &HeartbeatService{
		logger:   logger,
	}
}

func (h *HeartbeatService) PONE(ctx context.Context) string {
	return "PONE"
}
