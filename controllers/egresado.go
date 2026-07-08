package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/egresados_crud/models"
)

type EgresadoController struct{ web.Controller }

// @Title GetAll
// @Description Lista egresado según el contrato estándar de listado SGA (ver README: query, fields, sortby, order, limit, offset)
// @Param   query    query   string  false   "filtros k:v separados por coma (dot-notation para relaciones, ej. Usuario.Activo:true)"
// @Param   fields   query   string  false   "campos Go a devolver, separados por coma"
// @Param   sortby   query   string  false   "campo(s) de orden, separados por coma"
// @Param   order    query   string  false   "asc|desc, uno por sortby o único para todos"
// @Param   limit    query   int     false   "máximo de resultados (default 10, 0 = sin límite)"
// @Param   offset   query   int     false   "desplazamiento (default 0)"
// @Success 200 {array} models.Egresado
// @Failure 400 parámetros de query inválidos
// @Failure 404 error de consulta en la base de datos
// @router / [get]
func (c *EgresadoController) GetAll() {
	query, fields, sortby, order, offset, limit, err := parseGetAllParams(&c.Controller)
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	l, err := models.GetAllEgresado(query, fields, sortby, order, offset, limit)
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
// @Description Obtiene un egresado por id
// @Param   id    path    int    true    "id del egresado"
// @Success 200 {object} models.Egresado
// @Failure 400 id inválido
// @Failure 404 no encontrado
// @router /:id [get]
func (c *EgresadoController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	result, err := models.GetEgresadoById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(404); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = result
	}
	c.ServeJSON()
}

// @Title Post
// @Description Crea un egresado (subtipo exclusivo de usuario, C-7)
// @Param   body    body    models.Egresado    true    "objeto egresado a crear"
// @Success 201 {string} id "id numérico del registro creado"
// @Failure 400 error de parseo del body
// @Failure 500 error interno (ej. FK inválida)
// @router / [post]
func (c *EgresadoController) Post() {
	var v models.Egresado
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	if id, err := models.AddEgresado(&v); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Ctx.Output.SetStatus(201); c.Data["json"] = map[string]int64{"id": id}
	}
	c.ServeJSON()
}

// @Title Put
// @Description Reemplaza un egresado completo (Update sin lista de columnas: el caller debe enviar el objeto entero)
// @Param   id      path    int               true    "id del egresado"
// @Param   body    body    models.Egresado   true    "objeto completo a reemplazar"
// @Success 200 {string} OK
// @Failure 400 id o body inválido
// @Failure 500 error interno
// @router /:id [put]
func (c *EgresadoController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	var v models.Egresado
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	v.Id = id
	if err := models.UpdateEgresadoById(&v); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}

// @Title Delete
// @Description Borrado lógico de un egresado (activo=false)
// @Param   id    path    int    true    "id del egresado"
// @Success 200 {string} OK
// @Failure 400 id inválido
// @Failure 500 error interno
// @router /:id [delete]
func (c *EgresadoController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	if err := models.DeleteEgresado(id); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}
