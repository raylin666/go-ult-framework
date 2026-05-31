package db

import (
	"ult/internal/data/model"
	"ult/pkg/db"

	"gorm.io/gen"
)

func NewGeneratorDefaultDb(dbInterface db.Db, outPath string) {
	g := gen.NewGenerator(gen.Config{
		OutPath: outPath,
		Mode:    gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
	})

	g.UseDB(dbInterface.Get().DB())

	testModel := model.Test{}

	// apply basic crud api on structs or table models which is specified by table name with function
	// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Execute.

	g.ApplyBasic(
		testModel,
	)
	
	// apply diy interfaces on structs or table models
	// g.ApplyInterface()


	g.Execute()
}
