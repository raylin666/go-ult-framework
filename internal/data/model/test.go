package model

type Test struct {
	BaseModel

	Name string `gorm:"column:name;type:string;size:30;unique:uk_name;comment:测试名称" json:"name"`
}