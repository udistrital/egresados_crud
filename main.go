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
	// database.BuildPostgresConnectionString lee PGuser/PGpass/PGhost/PGport/PGdb/
	// PGschema de conf/app.conf (y, si parameterStore está seteado, credenciales desde
	// AWS SSM en su lugar). Reemplaza el armado manual que teníamos antes.
	dataSource, err := database.BuildPostgresConnectionString()
	if err != nil {
		panic(fmt.Sprintf("error armando la cadena de conexión: %v", err))
	}

	orm.RegisterDriver("postgres", orm.DRPostgres)
	if err := orm.RegisterDataBase("default", "postgres", dataSource); err != nil {
		panic(fmt.Sprintf("error registrando base de datos: %v", err))
	}

	// orm.Debug (configuracion_ci_cd.md): loguea cada sentencia SQL ejecutada, solo
	// en dev. De paso alimenta el campo sql_statement del log de auditoria.InitMiddleware
	// (más abajo) — sin esto, ese campo queda siempre vacío.
	if web.BConfig.RunMode == web.DEV {
		orm.Debug = true
	}

	// EnableDocs (conf/app.conf) solo activa la bandera; Beego v2 no sirve
	// swagger/ automáticamente, hay que exponerla como estática (bee generate docs).
	// Se lee directo del config genérico (no de BConfig.WebConfig.EnableDocs, que
	// Beego v2 marcó deprecated: "Beego didn't use it anymore" — sigue funcionando
	// hoy, pero podrían quitarlo en una versión futura).
	if web.AppConfig.DefaultBool("EnableDocs", false) {
		web.SetStaticPath("/swagger", "swagger")
	}

	// Integraciones institucionales (configuracion_ci_cd.md). Habilitadas ahora que
	// utils_oas/v2 (v2.0.0-beta.1, 2026-07) migró a Beego v2 — ya no choca con nuestro
	// stack (antes: panic: flag redefined: graceful, por astaxie/beego transitivo).
	apistatus.Init()              // GET / — healthcheck institucional ({"status":"ok"})
	auditoria.InitMiddleware()    // log de auditoría por request + última sentencia SQL
	security.SetSecurityHeaders() // headers CSP/HSTS/X-Frame-Options/etc.
	xray.Init()                   // tracing AWS X-Ray (no-op si PARAMETER_STORE no está seteado)
	web.ErrorController(&customerrorv2.CustomErrorController{})
}

func main() {
	web.Run()
}
