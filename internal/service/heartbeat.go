package service

import (
	"context"

	"ult/internal/app"
	"ult/internal/data/repo"
)

type HeartbeatService struct {
	testRepo repo.TestRepo
	tools    *app.Tools
}

func NewHeartbeatService(testRepo repo.TestRepo, tools *app.Tools) *HeartbeatService {
	return &HeartbeatService{
		testRepo: testRepo,
		tools:    tools,
	}
}

func (h *HeartbeatService) PONE(ctx context.Context) string {
	return "PONE"
}
