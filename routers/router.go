package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/sga_crud_beneficios_egresados/controllers"
)

func init() {
	// ── Catálogos ─────────────────────────────────────────────────────────────
	web.Router("/v1/tipo_usuario", &controllers.TipoUsuarioController{}, "get:GetAll;post:Post")
	web.Router("/v1/tipo_usuario/:id", &controllers.TipoUsuarioController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/estado_empresa", &controllers.EstadoEmpresaController{}, "get:GetAll;post:Post")
	web.Router("/v1/estado_empresa/:id", &controllers.EstadoEmpresaController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/estado_beneficio", &controllers.EstadoBeneficioController{}, "get:GetAll;post:Post")
	web.Router("/v1/estado_beneficio/:id", &controllers.EstadoBeneficioController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/estado_solicitud", &controllers.EstadoSolicitudController{}, "get:GetAll;post:Post")
	web.Router("/v1/estado_solicitud/:id", &controllers.EstadoSolicitudController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/categoria_beneficio", &controllers.CategoriaBeneficioController{}, "get:GetAll;post:Post")
	web.Router("/v1/categoria_beneficio/:id", &controllers.CategoriaBeneficioController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/sector_economico", &controllers.SectorEconomicoController{}, "get:GetAll;post:Post")
	web.Router("/v1/sector_economico/:id", &controllers.SectorEconomicoController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/parametro_sistema", &controllers.ParametroSistemaController{}, "get:GetAll;post:Post")
	web.Router("/v1/parametro_sistema/:id", &controllers.ParametroSistemaController{}, "get:GetOne;put:Put;delete:Delete")

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
	web.Router("/v1/secuencia_radicado/:id", &controllers.SecuenciaRadicadoController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/solicitud_beneficio", &controllers.SolicitudBeneficioController{}, "get:GetAll;post:Post")
	web.Router("/v1/solicitud_beneficio/:id", &controllers.SolicitudBeneficioController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/historial_estado_solicitud", &controllers.HistorialEstadoSolicitudController{}, "get:GetAll;post:Post")
	web.Router("/v1/historial_estado_solicitud/:id", &controllers.HistorialEstadoSolicitudController{}, "get:GetOne;put:Put;delete:Delete")

	web.Router("/v1/mensaje_solicitud", &controllers.MensajeSolicitudController{}, "get:GetAll;post:Post")
	web.Router("/v1/mensaje_solicitud/:id", &controllers.MensajeSolicitudController{}, "get:GetOne;put:Put;delete:Delete")

	// bitacora_acceso_pii: solo lectura (log inmutable, no DELETE/PUT)
	web.Router("/v1/bitacora_acceso_pii", &controllers.BitacoraAccesoPiiController{}, "get:GetAll;post:Post")
	web.Router("/v1/bitacora_acceso_pii/:id", &controllers.BitacoraAccesoPiiController{}, "get:GetOne")
}
