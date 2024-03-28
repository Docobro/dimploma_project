package db

import (
	"log"

	"github.com/docobro/dimploma_project/internal/db/driver"
	"github.com/docobro/dimploma_project/internal/model"
	"gorm.io/gen"
)

// Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func Generate() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/db/queries",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	d := driver.New()

	c, err := d.Connect(driver.DatabaseConnectionParams{
		Host:     "localhost",
		Port:     "8124",
		Database: "CryptoWallet",
		Username: "default",
		Password: "",
	})
	if err != nil {
		log.Fatalf("failed to connect to clickhouse.%v", err)
	}

	g.UseDB(c) // reuse your gorm db
	c.AutoMigrate(&model.User{})
	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(model.User{})

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	g.ApplyInterface(func(Querier) {}, model.User{})

	// Generate the code
	g.Execute()
}
