package service

import (
	"context"
	"time"
	"ult/internal/app"
	"ult/internal/data"
	"ult/internal/data/repo"
)

type HealtbeatStatus struct {
	Status     string   `json:"status"`
	Timestamp  string   `json:"timestamp"`
	Components map[string]HealtbeatComponent `json:"components,omitempty"`
}

type HealtbeatComponent struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type HealtbeatService struct {
	dataRepo *data.DataRepo
	testRepo repo.TestRepo
	tools    *app.Tools
}

func NewHeartbeatService(dataRepo *data.DataRepo, testRepo repo.TestRepo, tools *app.Tools) *HealtbeatService {
	return &HealtbeatService{
		dataRepo: dataRepo,
		testRepo: testRepo,
		tools:    tools,
	}
}

func (h *HealtbeatService) State(ctx context.Context) HealtbeatStatus {
	now := time.Now().Format(time.RFC3339)
	status := HealtbeatStatus{
		Status:     "healthy",
		Timestamp:  now,
		Components: make(map[string]HealtbeatComponent),
	}

	if h.dataRepo != nil && h.dataRepo.DbRepo != nil {
		if h.dataRepo.DbRepo.Count() > 0 {
			for name, db := range h.dataRepo.DbRepo.All() {
				if db == nil {
					status.Components["db_"+name] = HealtbeatComponent{
						Status:  "unhealthy",
						Message: "db connection is nil",
					}
					status.Status = "degraded"
					continue
				}

				if err := db.Ping(); err != nil {
					status.Components["db_"+name] = HealtbeatComponent{
						Status:  "unhealthy",
						Message: err.Error(),
					}
					status.Status = "degraded"
				} else {
					status.Components["db_"+name] = HealtbeatComponent{
						Status: "healthy",
					}
				}
			}
		}
	}

	if h.dataRepo != nil && h.dataRepo.RedisRepo != nil {
		if h.dataRepo.RedisRepo.Count() > 0 {
			for name, redis := range h.dataRepo.RedisRepo.All() {
				if redis == nil {
					status.Components["redis_"+name] = HealtbeatComponent{
						Status:  "unhealthy",
						Message: "redis connection is nil",
					}
					status.Status = "degraded"
					continue
				}

				if err := redis.Ping(ctx); err != nil {
					status.Components["redis_"+name] = HealtbeatComponent{
						Status:  "unhealthy",
						Message: err.Error(),
					}
					status.Status = "degraded"
				} else {
					status.Components["redis_"+name] = HealtbeatComponent{
						Status: "healthy",
					}
				}
			}
		}
	}

	return status
}
