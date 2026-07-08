package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/egresados_crud/models"
)

type HistorialSolicitudController struct{ web.Controller }

// @Title GetAll
// @Description Lista historial_solicitud según el contrato estándar de listado SGA (ver README: query, fields, sortby, order, limit, offset)
// @Param   query    query   string  false   "filtros k:v separados por coma (dot-notation para relaciones)"
// @Param   fields   query   string  false   "campos Go a devolver, separados por coma"
// @Param   sortby   query   string  false   "campo(s) de orden, separados por coma"
// @Param   order    query   string  false   "asc|desc, uno por sortby o único para todos"
// @Param   limit    query   int     false   "máximo de resultados (default 10, 0 = sin límite)"
// @Param   offset   query   int     false   "desplazamiento (default 0)"
// @Success 200 {array} models.HistorialSolicitud
// @Failure 400 parámetros de query inválidos
// @Failure 404 error de consulta en la base de datos
// @router / [get]
func (c *HistorialSolicitudController) GetAll() {
	query, fields, sortby, order, offset, limit, err := parseGetAllParams(&c.Controller)
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	l, err := models.GetAllHistorialSolicitud(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Ctx.Output.SetStatus(404); c.Data["json"] = err.Error()
	} else {
		if l == nil {
			l = append(l, map[string]interface{}{})
		}
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// @Title GetOne
// @Description Obtiene un registro de historial_solicitud por id
// @Param   id    path    int    true    "id del registro de historial"
// @Success 200 {object} models.HistorialSolicitud
// @Failure 400 id inválido
// @Failure 404 no encontrado
// @router /:id [get]
func (c *HistorialSolicitudController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	result, err := models.GetHistorialSolicitudById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(404); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = result
	}
	c.ServeJSON()
}

// @Title GetBySolicitud
// @Description Bitácora completa de una solicitud (más reciente primero)
// @Param   solicitud_id    path    int    true    "id de la solicitud_beneficio"
// @Success 200 {array} models.HistorialSolicitud
// @Failure 400 solicitud_id inválido
// @Failure 500 error interno
// @router /solicitud/:solicitud_id [get]
func (c *HistorialSolicitudController) GetBySolicitud() {
	solicitudId, err := strconv.Atoi(c.Ctx.Input.Param(":solicitud_id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "solicitud_id inválido"; c.ServeJSON(); return
	}
	results, err := models.GetHistorialSolicitudBySolicitud(solicitudId)
	if err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = results
	}
	c.ServeJSON()
}

// @Title GetVigente
// @Description Último registro de historial = estado vigente de la solicitud (C-4b: el historial es la única fuente de estado)
// @Param   solicitud_id    path    int    true    "id de la solicitud_beneficio"
// @Success 200 {object} models.HistorialSolicitud
// @Failure 400 solicitud_id inválido
// @Failure 404 la solicitud no tiene historial
// @router /solicitud/:solicitud_id/vigente [get]
func (c *HistorialSolicitudController) GetVigente() {
	solicitudId, err := strconv.Atoi(c.Ctx.Input.Param(":solicitud_id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "solicitud_id inválido"; c.ServeJSON(); return
	}
	result, err := models.GetEstadoVigenteBySolicitud(solicitudId)
	if err != nil {
		c.Ctx.Output.SetStatus(404); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = result
	}
	c.ServeJSON()
}

// @Title Post
// @Description Inserta un cambio de estado en el historial (los cambios de estado son INSERT, no UPDATE, C-4b)
// @Param   body    body    models.HistorialSolicitud    true    "objeto historial_solicitud a crear"
// @Success 201 {string} id "id numérico del registro creado"
// @Failure 400 error de parseo del body
// @Failure 500 error interno (ej. FK inválida)
// @router / [post]
func (c *HistorialSolicitudController) Post() {
	var v models.HistorialSolicitud
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	if id, err := models.AddHistorialSolicitud(&v); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Ctx.Output.SetStatus(201); c.Data["json"] = map[string]int64{"id": id}
	}
	c.ServeJSON()
}

// @Title Put
// @Description Reemplaza un registro de historial_solicitud completo (Update sin lista de columnas: el caller debe enviar el objeto entero)
// @Param   id      path    int                          true    "id del registro de historial"
// @Param   body    body    models.HistorialSolicitud    true    "objeto completo a reemplazar"
// @Success 200 {string} OK
// @Failure 400 id o body inválido
// @Failure 500 error interno
// @router /:id [put]
func (c *HistorialSolicitudController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	var v models.HistorialSolicitud
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	v.Id = id
	if err := models.UpdateHistorialSolicitudById(&v); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}

// @Title Delete
// @Description Borrado lógico de un registro de historial_solicitud (activo=false)
// @Param   id    path    int    true    "id del registro de historial"
// @Success 200 {string} OK
// @Failure 400 id inválido
// @Failure 500 error interno
// @router /:id [delete]
func (c *HistorialSolicitudController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	if err := models.DeleteHistorialSolicitud(id); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}
