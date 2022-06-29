package gen

import (
	"github.com/oaago/cloud/mysql"
	"gorm.io/gen"
)

func GenStruct(dbName string, tables []string) {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/model/" + dbName + "_query/",
		ModelPkgPath: "./internal/model/" + dbName + "_model/",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery,
		//if you want the nullable field generation property to be pointer type, set FieldNullable true
		//FieldNullable: true,
		//if you want to assign field which has default value in `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		//FieldCoverable: true,
		//if you want to generate index tags from database, set FieldWithIndexTag true
		FieldWithIndexTag: true,
		//if you want to generate type tags from database, set FieldWithTypeTag true
		FieldWithTypeTag: true,
		//if you need unit tests for query code, set WithUnitTest true
		/* WithUnitTest: true, */
	})

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary or it will panic
	db := mysql.GetDBByName(dbName)
	g.UseDB(db)

	// apply basic crud api on structs or table models which is specified by table name with function
	// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Excute.
	for _, table := range tables {
		g.ApplyBasic(g.GenerateModel(table))
	}
	g.Execute()
}
