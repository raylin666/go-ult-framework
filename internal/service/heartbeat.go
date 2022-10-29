package service

import (
	"context"
	"ult/pkg/global"
	"ult/pkg/logger"
)

type HeartbeatService struct {
	logger   *logger.Logger
	dataRepo global.DataRepo
}

func NewHeartbeatService(logger *logger.Logger, dataRepo global.DataRepo) *HeartbeatService {
	return &HeartbeatService{
		logger:   logger,
		dataRepo: dataRepo,
	}
}

func (h *HeartbeatService) PONE(ctx context.Context) string {
	return "PONE"
}
