package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/egresados_crud/controllers"
)

func init() {
	// Los catálogos (tipo_usuario, estados, categorías, sectores, parámetros de
	// sistema) viven en el servicio institucional de parámetros (C-1);
	// este CRUD solo expone las entidades propias del módulo.
	//
	// Rutas registradas vía NSInclude a partir de las anotaciones @router de cada
	// controller (routers/commentsRouter.go, generado con `bee generate routers`).
	// C-5: el radicado se genera con la SEQUENCE nativa seq_radicado_beneficio vía
	// fn_siguiente_radicado() (DEFAULT de solicitud_beneficio.radicado). Ya no hay
	// tabla/controlador secuencia_radicado.
	// historial_solicitud: única fuente de estado de las solicitudes (C-4b)
	// documento_requerido_beneficio: qué documentos exige la empresa al publicar un beneficio (RF-005)
	// documento_solicitud: PDFs subidos por el egresado para cumplir los documentos requeridos de su solicitud
	// bitacora_acceso_pii: solo lectura (log inmutable, no DELETE/PUT)
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
