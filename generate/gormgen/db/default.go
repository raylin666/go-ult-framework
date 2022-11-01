package db

import (
	"gorm.io/gen"
	"ult/pkg/db"
)

func NewGeneratorDefaultDb(dbInterface db.Db, outPath string) {
	g := gen.NewGenerator(gen.Config{
		// 生成目录存放位置
		OutPath: outPath,
		// WithContext 模式
		Mode: gen.WithDefaultQuery,
	})

	g.UseDB(dbInterface.Get().DB())

	// apply basic crud api on structs or table models which is specified by table name with function
	// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Execute.
	g.ApplyBasic(
	)

	// apply diy interfaces on structs or table models
	// g.ApplyInterface()

	g.Execute()
}
