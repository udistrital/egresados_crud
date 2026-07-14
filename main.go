package main

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
	_ "github.com/udistrital/egresados_crud/models"
	_ "github.com/udistrital/egresados_crud/routers"
	apistatus "github.com/udistrital/utils_oas/v2/apiStatusLib"
	"github.com/udistrital/utils_oas/v2/auditoria"
	customerrorv2 "github.com/udistrital/utils_oas/v2/customerror"
	"github.com/udistrital/utils_oas/v2/database"
	"github.com/udistrital/utils_oas/v2/security"
	"github.com/udistrital/utils_oas/v2/xray"
)

func init() {
	dataSource, err := database.BuildPostgresConnectionString()
	if err != nil {
		panic(fmt.Sprintf("error armando la cadena de conexión: %v", err))
	}

	orm.RegisterDriver("postgres", orm.DRPostgres)
	if err := orm.RegisterDataBase("default", "postgres", dataSource); err != nil {
		panic(fmt.Sprintf("error registrando base de datos: %v", err))
	}

	if web.BConfig.RunMode == web.DEV {
		orm.Debug = true
	}

	apistatus.Init()
	auditoria.InitMiddleware()
	security.SetSecurityHeaders()
	xray.Init()
	web.ErrorController(&customerrorv2.CustomErrorController{})
}

func main() {
	web.Run()
}
