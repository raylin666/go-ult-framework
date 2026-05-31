// Package model 提供数据模型定义。
// 定义数据库表对应的 Go 结构体模型。
package model

import (
	"gorm.io/gorm"
)

// BaseModel 基础模型结构体。
// 包含通用的 ID、创建时间、更新时间和软删除字段。
type BaseModel struct {
	ID        int            `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`         // 主键 ID
	CreatedAt int64          `gorm:"column:created_at" json:"created_at"`                       // 创建时间
	UpdatedAt int64          `gorm:"column:updated_at" json:"updated_at"`                       // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`                       // 软删除时间
}