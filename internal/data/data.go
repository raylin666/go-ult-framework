// Package data 提供数据层实现。
// 封装数据库和 Redis 连接，提供统一的数据访问接口。
package data

import (
	"context"
	"fmt"

	"ult/config"
	"ult/internal/app"
	"ult/internal/data/dbquery"
	"ult/pkg/repositories"

	"github.com/google/wire"
	"github.com/raylin666/go-utils/v2/cache/redis"
	"gorm.io/gorm"
)

// ProviderSet Wire 依赖注入提供者集合。
var ProviderSet = wire.NewSet(NewData, NewDataRepo)

var _ Data = (*dataImpl)(nil)

// Data 数据接口，定义数据访问操作。
type Data interface {
	DB() *dbquery.Query                   // 获取 GORM 查询器
	Redis() redis.Client                  // 获取 Redis 客户端
	GormDB() *gorm.DB                     // 获取 GORM DB 实例
	WithContext(ctx context.Context) Data // 创建带上下文的数据实例
	DataRepo() *DataRepo                  // 获取数据仓库
}

// dataImpl 数据实例实现。
type dataImpl struct {
	db       *dbquery.Query  // GORM 查询器
	gorm     *gorm.DB        // GORM DB 实例
	redis    redis.Client    // Redis 客户端
	ctx      context.Context // 上下文
	dataRepo *DataRepo       // 数据仓库
}

// DataRepo 数据仓库结构体，包含数据库和 Redis 仓库。
type DataRepo struct {
	DbRepo    repositories.DbRepo    // 数据库仓库
	RedisRepo repositories.RedisRepo // Redis 仓库
}

// NewData 创建新的数据实例。
// 初始化数据库和 Redis 连接，返回清理函数。
//
// 参数:
//   - c: 应用配置
//   - tools: 应用工具包
//
// 返回:
//   - Data: 数据实例
//   - func(): 清理函数
func NewData(c *config.Config, tools *app.Tools) (Data, func()) {
	var (
		ctx  = context.TODO()
		data = new(dataImpl)
		repo = repositories.NewDataRepo(tools.Logger(), c)
	)

	data.dataRepo = &DataRepo{
		DbRepo:    repo.DbRepo(),
		RedisRepo: repo.RedisRepo(),
	}

	data.gorm = newDefaultDb(data.dataRepo.DbRepo)
	if data.gorm != nil {
		data.db = dbquery.Use(data.gorm)
	}

	data.redis = newDefaultRedis(data.dataRepo.RedisRepo)

	data.ctx = context.Background()

	cleanup := func() {
		if data.dataRepo.DbRepo != nil {
			for dbName, dbRepo := range data.dataRepo.DbRepo.All() {
				if dbRepo != nil {
					_ = dbRepo.Close()
				}
				tools.Logger().UseApp(ctx).Info(fmt.Sprintf("关闭数据库连接: %s", dbName))
			}
		}
		if data.dataRepo.RedisRepo != nil {
			for redisName, redisRepo := range data.dataRepo.RedisRepo.All() {
				if redisRepo != nil {
					_ = redisRepo.Close()
				}
				tools.Logger().UseApp(ctx).Info(fmt.Sprintf("关闭 Redis 连接: %s", redisName))
			}
		}
	}

	return data, cleanup
}

// NewDataRepo 从数据实例创建数据仓库。
//
// 参数:
//   - data: 数据实例
//
// 返回:
//   - *DataRepo: 数据仓库
func NewDataRepo(data Data) *DataRepo {
	return data.DataRepo()
}

// DB 获取 GORM 查询器。
//
// 返回:
//   - *dbquery.Query: GORM 查询器
func (d *dataImpl) DB() *dbquery.Query {
	return d.db
}

// Redis 获取 Redis 客户端。
//
// 返回:
//   - redis.Client: Redis 客户端
func (d *dataImpl) Redis() redis.Client {
	return d.redis
}

// GormDB 获取 GORM DB 实例。
//
// 返回:
//   - *gorm.DB: GORM DB 实例
func (d *dataImpl) GormDB() *gorm.DB {
	return d.gorm
}

// DataRepo 获取数据仓库。
//
// 返回:
//   - *DataRepo: 数据仓库
func (d *dataImpl) DataRepo() *DataRepo {
	return d.dataRepo
}

// WithContext 创建带上下文的数据实例。
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - Data: 带上下文的数据实例
func (d *dataImpl) WithContext(ctx context.Context) Data {
	var db *dbquery.Query
	var gormDB *gorm.DB
	if d.gorm != nil {
		gormDB = d.gorm.WithContext(ctx)
		db = dbquery.Use(gormDB)
	}
	return &dataImpl{
		db:       db,
		gorm:     gormDB,
		redis:    d.redis,
		ctx:      ctx,
		dataRepo: d.dataRepo,
	}
}

// newDefaultDb 获取默认数据库连接的 GORM DB 实例。
//
// 参数:
//   - dbRepo: 数据库仓库
//
// 返回:
//   - *gorm.DB: GORM DB 实例
func newDefaultDb(dbRepo repositories.DbRepo) *gorm.DB {
	dbConn := dbRepo.DB(repositories.DbConnectionDefaultName)
	if dbConn == nil {
		return nil
	}
	return dbConn.Get().DB()
}

// newDefaultRedis 获取默认 Redis 连接的客户端实例。
//
// 参数:
//   - redisRepo: Redis 仓库
//
// 返回:
//   - redis.Client: Redis 客户端
func newDefaultRedis(redisRepo repositories.RedisRepo) redis.Client {
	redisConn := redisRepo.Redis(repositories.RedisConnectionDefaultName)
	if redisConn == nil {
		return nil
	}
	return redisConn.Get()
}
