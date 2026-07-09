package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/egresados_crud/models"
)

type BeneficioController struct{ web.Controller }

// @Title GetAll
// @Description Lista beneficio según el contrato estándar de listado SGA (ver README: query, fields, sortby, order, limit, offset)
// @Param   query    query   string  false   "filtros k:v separados por coma (dot-notation para relaciones, ej. Empresa.Id:1,Activo:true)"
// @Param   fields   query   string  false   "campos Go a devolver, separados por coma"
// @Param   sortby   query   string  false   "campo(s) de orden, separados por coma"
// @Param   order    query   string  false   "asc|desc, uno por sortby o único para todos"
// @Param   limit    query   int     false   "máximo de resultados (default 10, 0 = sin límite)"
// @Param   offset   query   int     false   "desplazamiento (default 0)"
// @Success 200 {array} models.Beneficio
// @Failure 400 parámetros de query inválidos
// @Failure 404 error de consulta en la base de datos
// @router / [get]
func (c *BeneficioController) GetAll() {
	query, fields, sortby, order, offset, limit, err := parseGetAllParams(&c.Controller)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	l, err := models.GetAllBeneficio(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = err.Error()
	} else {
		if l == nil {
			// lista vacía → [{}]: idioma estándar de los *_crud del SGA
			l = append(l, map[string]interface{}{})
		}
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// @Title GetOne
// @Description Obtiene un beneficio por id
// @Param   id    path    int    true    "id del beneficio"
// @Success 200 {object} models.Beneficio
// @Failure 400 id inválido
// @Failure 404 no encontrado
// @router /:id [get]
func (c *BeneficioController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "id inválido"
		c.ServeJSON()
		return
	}
	result, err := models.GetBeneficioById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = result
	}
	c.ServeJSON()
}

// @Title Post
// @Description Publica un beneficio
// @Param   body    body    models.Beneficio    true    "objeto beneficio a crear"
// @Success 201 {string} id "id numérico del registro creado"
// @Failure 400 error de parseo del body
// @Failure 500 error interno (ej. FK inválida)
// @router / [post]
func (c *BeneficioController) Post() {
	var v models.Beneficio
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	if id, err := models.AddBeneficio(&v); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = map[string]int64{"id": id}
	}
	c.ServeJSON()
}

// @Title Put
// @Description Reemplaza un beneficio completo (Update sin lista de columnas: el caller debe enviar el objeto entero)
// @Param   id      path    int                true    "id del beneficio"
// @Param   body    body    models.Beneficio   true    "objeto completo a reemplazar"
// @Success 200 {string} OK
// @Failure 400 id o body inválido
// @Failure 500 error interno
// @router /:id [put]
func (c *BeneficioController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "id inválido"
		c.ServeJSON()
		return
	}
	var v models.Beneficio
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	v.Id = id
	if err := models.UpdateBeneficioById(&v); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}

// @Title Delete
// @Description Borrado lógico de un beneficio (activo=false)
// @Param   id    path    int    true    "id del beneficio"
// @Success 200 {string} OK
// @Failure 400 id inválido
// @Failure 500 error interno
// @router /:id [delete]
func (c *BeneficioController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "id inválido"
		c.ServeJSON()
		return
	}
	if err := models.DeleteBeneficio(id); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}

// @Title DescontarCupo
// @Description Descuento atómico de un cupo del beneficio (RN-002b): UPDATE ... WHERE cupos_disponibles > 0, sin race conditions
// @Param   id    path    int    true    "id del beneficio"
// @Success 200 {string} descontado "true si se descontó el cupo"
// @Failure 400 id inválido
// @Failure 500 sin cupos disponibles o error interno
// @router /:id/cupo/descontar [post]
func (c *BeneficioController) DescontarCupo() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "id inválido"
		c.ServeJSON()
		return
	}
	descontado, err := models.DescontarCupo(id)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = map[string]bool{"descontado": descontado}
	}
	c.ServeJSON()
}

// @Title DevolverCupo
// @Description Devolución atómica de un cupo del beneficio (RN-002c): UPDATE ... WHERE cupos_disponibles < cupos_total, sin race conditions
// @Param   id    path    int    true    "id del beneficio"
// @Success 200 {string} devuelto "true si se devolvió el cupo"
// @Failure 400 id inválido
// @Failure 500 sin cupos que devolver o error interno
// @router /:id/cupo/devolver [post]
func (c *BeneficioController) DevolverCupo() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "id inválido"
		c.ServeJSON()
		return
	}
	devuelto, err := models.DevolverCupo(id)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = map[string]bool{"devuelto": devuelto}
	}
	c.ServeJSON()
}
