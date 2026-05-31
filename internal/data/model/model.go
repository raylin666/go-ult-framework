package model

import (
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int            `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CreatedAt int64          `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int64          `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}