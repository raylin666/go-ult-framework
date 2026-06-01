// Package repo 提供数据仓库实现。
// 数据仓库层封装数据库操作，为服务层提供数据访问接口。
package repo

import (
	"context"

	"ult/internal/data"
	"ult/internal/data/dbquery"
	"ult/internal/data/model"

	"gorm.io/gorm"
)

// TestRepo 接口验证。
var _ TestRepo = (*testRepo)(nil)

// TestRepo 测试数据仓库接口。
// 定义测试数据的 CRUD 操作方法。
type TestRepo interface {
	GetByID(ctx context.Context, id int) (*model.Test, error) // 根据 ID 获取测试数据
	List(ctx context.Context) ([]*model.Test, error)          // 获取测试数据列表
	Create(ctx context.Context, test *model.Test) error       // 创建测试数据
	Update(ctx context.Context, test *model.Test) error       // 更新测试数据
	Delete(ctx context.Context, id int) error                 // 删除测试数据
}

// testRepo 测试数据仓库实现。
type testRepo struct {
	data data.Data // 数据实例
}

// NewTestRepo 创建新的测试数据仓库实例。
//
// 参数:
//   - data: 数据实例
//
// 返回:
//   - TestRepo: 测试数据仓库接口
func NewTestRepo(data data.Data) TestRepo {
	return &testRepo{data: data}
}

// db 获取带上下文的 GORM DB 实例。
func (t *testRepo) db(ctx context.Context) *gorm.DB {
	return t.data.WithContext(ctx).GormDB()
}

// query 获取带上下文的测试模型查询器。
func (t *testRepo) query(ctx context.Context) dbquery.ITestDo {
	return dbquery.Use(t.db(ctx)).Test.WithContext(ctx)
}

// GetByID 根据 ID 获取测试数据。
//
// 参数:
//   - ctx: 上下文
//   - id: 测试数据 ID
//
// 返回:
//   - *model.Test: 测试数据
//   - error: 查询错误
func (t *testRepo) GetByID(ctx context.Context, id int) (*model.Test, error) {
	return t.query(ctx).Where(dbquery.Test.ID.Eq(id)).First()
}

// List 获取所有测试数据列表。
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - []*model.Test: 测试数据列表
//   - error: 查询错误
func (t *testRepo) List(ctx context.Context) ([]*model.Test, error) {
	return t.query(ctx).Find()
}

// Create 创建新的测试数据。
//
// 参数:
//   - ctx: 上下文
//   - test: 测试数据
//
// 返回:
//   - error: 创建错误
func (t *testRepo) Create(ctx context.Context, test *model.Test) error {
	return t.query(ctx).Create(test)
}

// Update 更新测试数据。
//
// 参数:
//   - ctx: 上下文
//   - test: 测试数据
//
// 返回:
//   - error: 更新错误
func (t *testRepo) Update(ctx context.Context, test *model.Test) error {
	return t.query(ctx).Save(test)
}

// Delete 根据 ID 删除测试数据。
//
// 参数:
//   - ctx: 上下文
//   - id: 测试数据 ID
//
// 返回:
//   - error: 删除错误
func (t *testRepo) Delete(ctx context.Context, id int) error {
	_, err := t.query(ctx).Where(dbquery.Test.ID.Eq(id)).Delete()
	return err
}
