package main

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
	_ "github.com/udistrital/egresados_crud/models"
	_ "github.com/udistrital/egresados_crud/routers"
	// Pendiente (ver nota "Integraciones institucionales" en init(), más abajo):
	// estos 5 imports rompen el build hoy por el choque Beego v1/v2. Descomentar
	// junto con el bloque de abajo cuando utils_oas soporte Beego v2.
	// apistatus "github.com/udistrital/utils_oas/apiStatusLib"
	// "github.com/udistrital/utils_oas/auditoria"
	// "github.com/udistrital/utils_oas/customerrorv2"
	// "github.com/udistrital/utils_oas/security"
	// "github.com/udistrital/utils_oas/xray"
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

	// ── Integraciones institucionales pendientes (configuracion_ci_cd.md) ──────────
	// Decisión con el Ingeniero (2026-07-09): este repo se queda en Beego v2 por ahora.
	// utils_oas se actualizará a v2 más adelante (de su lado o del nuestro); hasta que
	// eso pase, NO descomentar lo siguiente: apiStatusLib/auditoria/security/xray/
	// customerrorv2 importan github.com/astaxie/beego (v1), que registra la misma flag
	// "graceful" que beego/v2/server/web/grace — el binario hace panic al arrancar si
	// conviven los dos (confirmado: "panic: flag redefined: graceful").
	//
	// database.BuildPostgresConnectionString() (el helper que pide la plantilla para
	// APIs CRUD) NO se necesita agregar: arriba ya se arma la cadena de conexión con
	// las mismas claves institucionales PGuser/PGpass/PGhost/PGport/PGdb/PGschema, con
	// código v2-nativo (web.AppConfig), sin depender de utils_oas.
	//
	// Cuando utils_oas soporte Beego v2 (o se decida migrar este repo a v1): descomentar
	// esto + los 5 imports de arriba, correr
	//   go get github.com/udistrital/utils_oas@latest && go mod tidy
	// y validar que compile y arranque (incluye el ruteo, el ORM y el swagger, que hoy
	// ya están probados y funcionando sin estas piezas) antes de desplegar.
	//
	// apistatus.Init()
	// auditoria.InitMiddleware()
	// security.SetSecurityHeaders()
	// xray.Init()
	// TODO: revisar si customerrorv2.CustomErrorController tiene equivalente para
	// web.ErrorController (v2) o si customerrorv2 en sí ya es agnóstico de versión.
	// web.ErrorController(&customerrorv2.CustomErrorController{})
}

func main() {
	web.Run()
}
