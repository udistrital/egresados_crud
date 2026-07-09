package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"],
		beego.ControllerComments{
			Method:           "DescontarCupo",
			Router:           `/:id/cupo/descontar`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BeneficioController"],
		beego.ControllerComments{
			Method:           "DevolverCupo",
			Router:           `/:id/cupo/devolver`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BitacoraAccesoPiiController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BitacoraAccesoPiiController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BitacoraAccesoPiiController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BitacoraAccesoPiiController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BitacoraAccesoPiiController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:BitacoraAccesoPiiController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoRequeridoBeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoRequeridoBeneficioController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoRequeridoBeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoRequeridoBeneficioController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoRequeridoBeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoRequeridoBeneficioController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoRequeridoBeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoRequeridoBeneficioController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoRequeridoBeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoRequeridoBeneficioController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoSolicitudController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoSolicitudController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoSolicitudController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoSolicitudController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoSolicitudController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:DocumentoSolicitudController"],
		beego.ControllerComments{
			Method:           "GetBySolicitud",
			Router:           `/solicitud/:solicitud_id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EgresadoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EgresadoController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EgresadoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EgresadoController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EgresadoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EgresadoController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EgresadoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EgresadoController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EgresadoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EgresadoController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EmpresaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EmpresaController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EmpresaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EmpresaController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EmpresaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EmpresaController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EmpresaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EmpresaController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EmpresaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:EmpresaController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"],
		beego.ControllerComments{
			Method:           "GetBySolicitud",
			Router:           `/solicitud/:solicitud_id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:HistorialSolicitudController"],
		beego.ControllerComments{
			Method:           "GetVigente",
			Router:           `/solicitud/:solicitud_id/vigente`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:MensajeSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:MensajeSolicitudController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:MensajeSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:MensajeSolicitudController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:MensajeSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:MensajeSolicitudController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:MensajeSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:MensajeSolicitudController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:MensajeSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:MensajeSolicitudController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:SolicitudBeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:SolicitudBeneficioController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:SolicitudBeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:SolicitudBeneficioController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:SolicitudBeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:SolicitudBeneficioController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:SolicitudBeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:SolicitudBeneficioController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:SolicitudBeneficioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:SolicitudBeneficioController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioEmpresaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioEmpresaController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioEmpresaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioEmpresaController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioEmpresaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioEmpresaController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioEmpresaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioEmpresaController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioEmpresaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/egresados_crud/controllers:UsuarioEmpresaController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
