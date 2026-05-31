package repo

import (
	"context"

	"ult/internal/data"
	"ult/internal/data/dbquery"
	"ult/internal/data/model"

	"gorm.io/gorm"
)

var _ TestRepo = (*testRepo)(nil)

type TestRepo interface {
	GetByID(ctx context.Context, id int) (*model.Test, error)
	List(ctx context.Context) ([]*model.Test, error)
	Create(ctx context.Context, test *model.Test) error
	Update(ctx context.Context, test *model.Test) error
	Delete(ctx context.Context, id int) error
}

type testRepo struct {
	data data.Data
}

func NewTestRepo(data data.Data) TestRepo {
	return &testRepo{data: data}
}

func (t *testRepo) db(ctx context.Context) *gorm.DB {
	return t.data.WithContext(ctx).GormDB()
}

func (t *testRepo) query(ctx context.Context) dbquery.ITestDo {
	return dbquery.Use(t.db(ctx)).Test.WithContext(ctx)
}

func (t *testRepo) GetByID(ctx context.Context, id int) (*model.Test, error) {
	return t.query(ctx).Where(dbquery.Test.ID.Eq(id)).First()
}

func (t *testRepo) List(ctx context.Context) ([]*model.Test, error) {
	return t.query(ctx).Find()
}

func (t *testRepo) Create(ctx context.Context, test *model.Test) error {
	return t.query(ctx).Create(test)
}

func (t *testRepo) Update(ctx context.Context, test *model.Test) error {
	return t.query(ctx).Save(test)
}

func (t *testRepo) Delete(ctx context.Context, id int) error {
	_, err := t.query(ctx).Where(dbquery.Test.ID.Eq(id)).Delete()
	return err
}
