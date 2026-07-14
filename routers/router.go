package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/egresados_crud/controllers"
)

func init() {
	// Los catálogos viven en el servicio institucional de parámetros (C-1); este
	// CRUD solo expone las entidades propias del módulo. Rutas registradas vía
	// NSInclude a partir de las anotaciones @router de cada controller.
	ns := web.NewNamespace("/v1",
		web.NSNamespace("/usuario",
			web.NSInclude(&controllers.UsuarioController{}),
		),
		web.NSNamespace("/egresado",
			web.NSInclude(&controllers.EgresadoController{}),
		),
		web.NSNamespace("/empresa",
			web.NSInclude(&controllers.EmpresaController{}),
		),
		web.NSNamespace("/usuario-empresa",
			web.NSInclude(&controllers.UsuarioEmpresaController{}),
		),
		web.NSNamespace("/beneficio",
			web.NSInclude(&controllers.BeneficioController{}),
		),
		web.NSNamespace("/solicitud-beneficio",
			web.NSInclude(&controllers.SolicitudBeneficioController{}),
		),
		web.NSNamespace("/historial-solicitud",
			web.NSInclude(&controllers.HistorialSolicitudController{}),
		),
		web.NSNamespace("/mensaje-solicitud",
			web.NSInclude(&controllers.MensajeSolicitudController{}),
		),
		web.NSNamespace("/documento-requerido-beneficio",
			web.NSInclude(&controllers.DocumentoRequeridoBeneficioController{}),
		),
		web.NSNamespace("/documento-solicitud",
			web.NSInclude(&controllers.DocumentoSolicitudController{}),
		),
		web.NSNamespace("/bitacora-acceso-pii",
			web.NSInclude(&controllers.BitacoraAccesoPiiController{}),
		),
	)
	web.AddNamespace(ns)
}
