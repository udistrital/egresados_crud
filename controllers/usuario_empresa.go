package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/egresados_crud/models"
)

type UsuarioEmpresaController struct{ web.Controller }

// @Title GetAll
// @Description Lista usuario_empresa según el contrato estándar de listado SGA (ver README: query, fields, sortby, order, limit, offset)
// @Param   query    query   string  false   "filtros k:v separados por coma (dot-notation para relaciones)"
// @Param   fields   query   string  false   "campos Go a devolver, separados por coma"
// @Param   sortby   query   string  false   "campo(s) de orden, separados por coma"
// @Param   order    query   string  false   "asc|desc, uno por sortby o único para todos"
// @Param   limit    query   int     false   "máximo de resultados (default 10, 0 = sin límite)"
// @Param   offset   query   int     false   "desplazamiento (default 0)"
// @Success 200 {array} models.UsuarioEmpresa
// @Failure 400 parámetros de query inválidos
// @Failure 404 error de consulta en la base de datos
// @router / [get]
func (c *UsuarioEmpresaController) GetAll() {
	query, fields, sortby, order, offset, limit, err := parseGetAllParams(&c.Controller)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	l, err := models.GetAllUsuarioEmpresa(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = err.Error()
	} else {
		if l == nil {
			l = append(l, map[string]interface{}{})
		}
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// @Title GetOne
// @Description Obtiene un vínculo usuario_empresa por id
// @Param   id    path    int    true    "id del vínculo usuario_empresa"
// @Success 200 {object} models.UsuarioEmpresa
// @Failure 400 id inválido
// @Failure 404 no encontrado
// @router /:id [get]
func (c *UsuarioEmpresaController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "id inválido"
		c.ServeJSON()
		return
	}
	result, err := models.GetUsuarioEmpresaById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = result
	}
	c.ServeJSON()
}

// @Title Post
// @Description Crea el vínculo usuario(EMP) ↔ empresa (C-7: tipo_usuario fijado a 'EMP')
// @Param   body    body    models.UsuarioEmpresa    true    "objeto usuario_empresa a crear"
// @Success 201 {string} id "id numérico del registro creado"
// @Failure 400 error de parseo del body
// @Failure 500 error interno (ej. FK inválida)
// @router / [post]
func (c *UsuarioEmpresaController) Post() {
	var v models.UsuarioEmpresa
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	if id, err := models.AddUsuarioEmpresa(&v); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = map[string]int64{"id": id}
	}
	c.ServeJSON()
}

// @Title Put
// @Description Reemplaza un usuario_empresa completo (Update sin lista de columnas: el caller debe enviar el objeto entero)
// @Param   id      path    int                     true    "id del vínculo usuario_empresa"
// @Param   body    body    models.UsuarioEmpresa   true    "objeto completo a reemplazar"
// @Success 200 {string} OK
// @Failure 400 id o body inválido
// @Failure 500 error interno
// @router /:id [put]
func (c *UsuarioEmpresaController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "id inválido"
		c.ServeJSON()
		return
	}
	var v models.UsuarioEmpresa
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	v.Id = id
	if err := models.UpdateUsuarioEmpresaById(&v); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}

// @Title Delete
// @Description Borrado lógico de un usuario_empresa (activo=false)
// @Param   id    path    int    true    "id del vínculo usuario_empresa"
// @Success 200 {string} OK
// @Failure 400 id inválido
// @Failure 500 error interno
// @router /:id [delete]
func (c *UsuarioEmpresaController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "id inválido"
		c.ServeJSON()
		return
	}
	if err := models.DeleteUsuarioEmpresa(id); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}
