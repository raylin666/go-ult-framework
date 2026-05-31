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

var ProviderSet = wire.NewSet(NewData)

var _ Data = (*dataImpl)(nil)

type Data interface {
	DB() *dbquery.Query
	Redis() redis.Client
	GormDB() *gorm.DB
	WithContext(ctx context.Context) Data
}

type dataImpl struct {
	db    *dbquery.Query
	gorm  *gorm.DB
	redis redis.Client
	ctx   context.Context
}

func NewData(c *config.Config, tools *app.Tools) (Data, func()) {
	var (
		ctx  = context.TODO()
		data = new(dataImpl)
		repo = repositories.NewDataRepo(tools.Logger(), c)
	)

	data.gorm = newDefaultDb(repo.DbRepo(), ctx)
	if data.gorm != nil {
		data.db = dbquery.Use(data.gorm)
	}
	data.redis = newDefaultRedis(repo.RedisRepo(), ctx)
	data.ctx = context.Background()

	cleanup := func() {
		for dbName, dbRepo := range repo.DbRepo().All() {
			_ = dbRepo.Close()
			tools.Logger().UseApp(ctx).Info(fmt.Sprintf("closing db: %s", dbName))
		}
		for redisName, redisRepo := range repo.RedisRepo().All() {
			_ = redisRepo.Close()
			tools.Logger().UseApp(ctx).Info(fmt.Sprintf("closing redis: %s", redisName))
		}
	}

	return data, cleanup
}

func (d *dataImpl) DB() *dbquery.Query {
	return d.db
}

func (d *dataImpl) Redis() redis.Client {
	return d.redis
}

func (d *dataImpl) GormDB() *gorm.DB {
	return d.gorm
}

func (d *dataImpl) WithContext(ctx context.Context) Data {
	var db *dbquery.Query
	var gormDB *gorm.DB
	if d.gorm != nil {
		gormDB = d.gorm.WithContext(ctx)
		db = dbquery.Use(gormDB)
	}
	return &dataImpl{
		db:    db,
		gorm:  gormDB,
		redis: d.redis,
		ctx:   ctx,
	}
}

func newDefaultDb(dbRepo repositories.DbRepo, ctx context.Context) *gorm.DB {
	dbConn := dbRepo.DB(repositories.DbConnectionDefaultName)
	if dbConn == nil {
		return nil
	}
	return dbConn.Get().DB()
}

func newDefaultRedis(redisRepo repositories.RedisRepo, ctx context.Context) redis.Client {
	redisConn := redisRepo.Redis(repositories.RedisConnectionDefaultName)
	if redisConn == nil {
		return nil
	}
	return redisConn.Get()
}