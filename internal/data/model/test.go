// Package model 提供数据模型定义。
package model

// Test 测试模型结构体。
// 对应数据库中的 test 表。
type Test struct {
	BaseModel

	Name string `gorm:"column:name;type:string;size:30;unique:uk_name;comment:测试名称" json:"name"` // 测试名称
}
