package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/sga_crud_beneficios_egresados/controllers"
)

func init() {
	// Los catálogos (tipo_usuario, estados, categorías, sectores, parámetros de
	// sistema) viven en el servicio institucional de parámetros (C-1);
	// este CRUD solo expone las entidades propias del módulo.

	// ── Entidades ─────────────────────────────────────────────────────────────
	web.Router("/v1/usuario", &controllers.UsuarioController{}, "get:GetAll;post:Post")
	web.Router("/v1/usuario/:id", &controllers.UsuarioController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/egresado", &controllers.EgresadoController{}, "get:GetAll;post:Post")
	web.Router("/v1/egresado/:id", &controllers.EgresadoController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/empresa", &controllers.EmpresaController{}, "get:GetAll;post:Post")
	web.Router("/v1/empresa/:id", &controllers.EmpresaController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/usuario_empresa", &controllers.UsuarioEmpresaController{}, "get:GetAll;post:Post")
	web.Router("/v1/usuario_empresa/:id", &controllers.UsuarioEmpresaController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/beneficio", &controllers.BeneficioController{}, "get:GetAll;post:Post")
	web.Router("/v1/beneficio/:id", &controllers.BeneficioController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/secuencia_radicado", &controllers.SecuenciaRadicadoController{}, "get:GetAll;post:Post")
	web.Router("/v1/secuencia_radicado/siguiente/:anio", &controllers.SecuenciaRadicadoController{}, "post:Siguiente")
	web.Router("/v1/secuencia_radicado/:id", &controllers.SecuenciaRadicadoController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/solicitud_beneficio", &controllers.SolicitudBeneficioController{}, "get:GetAll;post:Post")
	web.Router("/v1/solicitud_beneficio/:id", &controllers.SolicitudBeneficioController{}, "get:GetOne;put:Put;delete:Delete")

	// historial_solicitud: única fuente de estado de las solicitudes (C-4b)
	web.Router("/v1/historial_solicitud", &controllers.HistorialSolicitudController{}, "get:GetAll;post:Post")
	web.Router("/v1/historial_solicitud/solicitud/:solicitud_id", &controllers.HistorialSolicitudController{}, "get:GetBySolicitud")
	web.Router("/v1/historial_solicitud/solicitud/:solicitud_id/vigente", &controllers.HistorialSolicitudController{}, "get:GetVigente")
	web.Router("/v1/historial_solicitud/:id", &controllers.HistorialSolicitudController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/mensaje_solicitud", &controllers.MensajeSolicitudController{}, "get:GetAll;post:Post")
	web.Router("/v1/mensaje_solicitud/:id", &controllers.MensajeSolicitudController{}, "get:GetOne;put:Put;delete:Delete")

	// bitacora_acceso_pii: solo lectura (log inmutable, no DELETE/PUT)
	web.Router("/v1/bitacora_acceso_pii", &controllers.BitacoraAccesoPiiController{}, "get:GetAll;post:Post")
	web.Router("/v1/bitacora_acceso_pii/:id", &controllers.BitacoraAccesoPiiController{}, "get:GetOne")
}
