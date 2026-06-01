// Package repositories 提供数据仓库抽象层。
// 管理数据库和 Redis 连接资源，支持多连接实例管理。
package repositories

import "ult/pkg/db"

const (
	DbConnectionDefaultName = "default" // 默认数据库连接名称
)

var _ DbRepo = (*dbRepo)(nil)

// DbRepo 数据库仓库接口，定义数据库连接管理操作。
type DbRepo interface {
	Count() int            // 获取连接数量
	Has(name string) bool  // 检查连接是否存在
	DB(name string) db.Db  // 获取指定名称的数据库连接
	All() map[string]db.Db // 获取所有数据库连接
}

// dbRepo 数据库仓库实例，管理多个数据库连接。
type dbRepo struct {
	resource map[string]db.Db // 数据库连接映射
}

// Count 获取数据库连接数量。
//
// 返回:
//   - int: 连接数量
func (repo *dbRepo) Count() int {
	return len(repo.resource)
}

// Has 检查指定名称的数据库连接是否存在。
//
// 参数:
//   - name: 连接名称
//
// 返回:
//   - bool: true 表示存在
func (repo *dbRepo) Has(name string) bool {
	if _, ok := repo.resource[name]; ok {
		return true
	}

	return false
}

// DB 获取指定名称的数据库连接。
//
// 参数:
//   - name: 连接名称
//
// 返回:
//   - db.Db: 数据库连接实例
func (repo *dbRepo) DB(name string) db.Db {
	return repo.resource[name]
}

// All 获取所有数据库连接。
//
// 返回:
//   - map[string]db.Db: 数据库连接映射
func (repo *dbRepo) All() map[string]db.Db {
	return repo.resource
}
