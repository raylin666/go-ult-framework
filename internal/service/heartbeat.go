// Package service 提供业务逻辑层实现。
// 服务层负责处理业务逻辑，调用数据层获取数据，为 API 层提供服务接口。
package service

import (
	"context"
	"time"
	"ult/internal/app"
	"ult/internal/data"
	"ult/internal/data/repo"
)

// HeartbeatStatusResponse 健康检查状态结构体。
// 包含整体状态、时间戳和各组件状态。
type HeartbeatStatusResponse struct {
	Status     string                        `json:"status"`               // 整体状态: healthy/degraded/unhealthy
	Timestamp  string                        `json:"timestamp"`            // 检查时间戳
	Components map[string]HeartbeatComponent `json:"components,omitempty"` // 各组件状态
}

// HeartbeatComponent 健康检查组件状态结构体。
type HeartbeatComponent struct {
	Status  string `json:"status"`            // 组件状态: healthy/unhealthy
	Message string `json:"message,omitempty"` // 状态消息（错误信息）
}

// HeartbeatService 健康检查服务。
// 提供数据库和 Redis 连接状态检查功能。
type HeartbeatService struct {
	dataRepo *data.DataRepo // 数据仓库
	testRepo repo.TestRepo  // 测试数据仓库
	tools    *app.Tools     // 应用工具包
}

// NewHeartbeatService 创建新的健康检查服务实例。
//
// 参数:
//   - dataRepo: 数据仓库
//   - testRepo: 测试数据仓库
//   - tools: 应用工具包
//
// 返回:
//   - *HeartbeatService: 健康检查服务实例
func NewHeartbeatService(dataRepo *data.DataRepo, testRepo repo.TestRepo, tools *app.Tools) *HeartbeatService {
	return &HeartbeatService{
		dataRepo: dataRepo,
		testRepo: testRepo,
		tools:    tools,
	}
}

// State 获取系统健康状态。
// 检查数据库和 Redis 连接状态，返回整体健康状态。
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - HeartbeatStatusResponse: 健康状态
func (h *HeartbeatService) State(ctx context.Context) HeartbeatStatusResponse {
	now := time.Now().Format(time.RFC3339)
	status := HeartbeatStatusResponse{
		Status:     "healthy",
		Timestamp:  now,
		Components: make(map[string]HeartbeatComponent),
	}

	// 检查数据库连接状态
	if h.dataRepo != nil && h.dataRepo.DbRepo != nil {
		if h.dataRepo.DbRepo.Count() > 0 {
			for name, db := range h.dataRepo.DbRepo.All() {
				if db == nil {
					status.Components["db_"+name] = HeartbeatComponent{
						Status:  "unhealthy",
						Message: "db connection is nil",
					}
					status.Status = "degraded"
					continue
				}

				if err := db.Ping(); err != nil {
					status.Components["db_"+name] = HeartbeatComponent{
						Status:  "unhealthy",
						Message: err.Error(),
					}
					status.Status = "degraded"
				} else {
					status.Components["db_"+name] = HeartbeatComponent{
						Status: "healthy",
					}
				}
			}
		}
	}

	// 检查 Redis 连接状态
	if h.dataRepo != nil && h.dataRepo.RedisRepo != nil {
		if h.dataRepo.RedisRepo.Count() > 0 {
			for name, redis := range h.dataRepo.RedisRepo.All() {
				if redis == nil {
					status.Components["redis_"+name] = HeartbeatComponent{
						Status:  "unhealthy",
						Message: "redis connection is nil",
					}
					status.Status = "degraded"
					continue
				}

				if err := redis.Ping(ctx); err != nil {
					status.Components["redis_"+name] = HeartbeatComponent{
						Status:  "unhealthy",
						Message: err.Error(),
					}
					status.Status = "degraded"
				} else {
					status.Components["redis_"+name] = HeartbeatComponent{
						Status: "healthy",
					}
				}
			}
		}
	}

	return status
}
