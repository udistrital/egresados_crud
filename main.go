package main

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
	_ "github.com/udistrital/egresados_crud/models"
	_ "github.com/udistrital/egresados_crud/routers"
)

func init() {
	dbUser := web.AppConfig.DefaultString("PGuser", "postgres")
	dbPassword := web.AppConfig.DefaultString("PGpass", "")
	dbHost := web.AppConfig.DefaultString("PGhost", "127.0.0.1")
	dbPort := web.AppConfig.DefaultString("PGport", "5432")
	dbName := web.AppConfig.DefaultString("PGdb", "beneficios_egresados")
	dbSchema := web.AppConfig.DefaultString("PGschema", "beneficios_egresados")

	dataSource := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable search_path=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSchema,
	)

	orm.RegisterDriver("postgres", orm.DRPostgres)
	if err := orm.RegisterDataBase("default", "postgres", dataSource); err != nil {
		panic(fmt.Sprintf("error registrando base de datos: %v", err))
	}

	// EnableDocs (conf/app.conf) solo activa la bandera; Beego v2 no sirve
	// swagger/ automáticamente, hay que exponerla como estática (bee generate docs).
	if web.BConfig.WebConfig.EnableDocs {
		web.SetStaticPath("/swagger", "swagger")
	}
}

func main() {
	web.Run()
}
