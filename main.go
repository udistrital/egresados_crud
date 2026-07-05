package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
	_ "github.com/udistrital/sga_crud_beneficios_egresados/models"
	_ "github.com/udistrital/sga_crud_beneficios_egresados/routers"
)

func init() {
	dbUser := getEnv("BENEFICIOS_EGRESADOS_CRUD_DB_USER", "postgres")
	dbPassword := getEnv("BENEFICIOS_EGRESADOS_CRUD_DB_PASSWORD", "1234")
	dbHost := getEnv("BENEFICIOS_EGRESADOS_CRUD_DB_HOST", "127.0.0.1")
	dbPort := getEnv("BENEFICIOS_EGRESADOS_CRUD_DB_PORT", "5432")
	dbName := getEnv("BENEFICIOS_EGRESADOS_CRUD_DB_NAME", "beneficios_egresados")
	dbSchema := getEnv("BENEFICIOS_EGRESADOS_CRUD_DB_SCHEMA", "beneficios_egresados")

	dataSource := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable search_path=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSchema,
	)

	orm.RegisterDriver("postgres", orm.DRPostgres)
	if err := orm.RegisterDataBase("default", "postgres", dataSource); err != nil {
		panic(fmt.Sprintf("error registrando base de datos: %v", err))
	}

	if port := os.Getenv("BENEFICIOS_EGRESADOS_CRUD_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			web.BConfig.Listen.HTTPPort = p
		}
	}
	if runmode := os.Getenv("BENEFICIOS_EGRESADOS_CRUD_RUNMODE"); runmode != "" {
		web.BConfig.RunMode = runmode
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	web.Run()
}
