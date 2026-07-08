package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/egresados_crud/models"
)

type DocumentoRequeridoBeneficioController struct{ web.Controller }

// @Title GetAll
// @Description Lista documento_requerido_beneficio según el contrato estándar de listado SGA (ver README: query, fields, sortby, order, limit, offset)
// @Param   query    query   string  false   "filtros k:v separados por coma (dot-notation para relaciones)"
// @Param   fields   query   string  false   "campos Go a devolver, separados por coma"
// @Param   sortby   query   string  false   "campo(s) de orden, separados por coma"
// @Param   order    query   string  false   "asc|desc, uno por sortby o único para todos"
// @Param   limit    query   int     false   "máximo de resultados (default 10, 0 = sin límite)"
// @Param   offset   query   int     false   "desplazamiento (default 0)"
// @Success 200 {array} models.DocumentoRequeridoBeneficio
// @Failure 400 parámetros de query inválidos
// @Failure 404 error de consulta en la base de datos
// @router / [get]
func (c *DocumentoRequeridoBeneficioController) GetAll() {
	query, fields, sortby, order, offset, limit, err := parseGetAllParams(&c.Controller)
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	l, err := models.GetAllDocumentoRequeridoBeneficio(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Ctx.Output.SetStatus(404); c.Data["json"] = err.Error()
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
// @Description Obtiene un documento_requerido_beneficio por id
// @Param   id    path    int    true    "id del documento requerido"
// @Success 200 {object} models.DocumentoRequeridoBeneficio
// @Failure 400 id inválido
// @Failure 404 no encontrado
// @router /:id [get]
func (c *DocumentoRequeridoBeneficioController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	result, err := models.GetDocumentoRequeridoBeneficioById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(404); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = result
	}
	c.ServeJSON()
}

// @Title Post
// @Description Define un documento exigido por la empresa al publicar un beneficio (RF-005)
// @Param   body    body    models.DocumentoRequeridoBeneficio    true    "objeto documento_requerido_beneficio a crear"
// @Success 201 {string} id "id numérico del registro creado"
// @Failure 400 error de parseo del body
// @Failure 500 error interno (ej. FK inválida)
// @router / [post]
func (c *DocumentoRequeridoBeneficioController) Post() {
	var v models.DocumentoRequeridoBeneficio
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	if id, err := models.AddDocumentoRequeridoBeneficio(&v); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Ctx.Output.SetStatus(201); c.Data["json"] = map[string]int64{"id": id}
	}
	c.ServeJSON()
}

// @Title Put
// @Description Reemplaza un documento_requerido_beneficio completo (Update sin lista de columnas: el caller debe enviar el objeto entero)
// @Param   id      path    int                                   true    "id del documento requerido"
// @Param   body    body    models.DocumentoRequeridoBeneficio    true    "objeto completo a reemplazar"
// @Success 200 {string} OK
// @Failure 400 id o body inválido
// @Failure 500 error interno
// @router /:id [put]
func (c *DocumentoRequeridoBeneficioController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	var v models.DocumentoRequeridoBeneficio
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	v.Id = id
	if err := models.UpdateDocumentoRequeridoBeneficioById(&v); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}

// @Title Delete
// @Description Borrado lógico de un documento_requerido_beneficio (activo=false)
// @Param   id    path    int    true    "id del documento requerido"
// @Success 200 {string} OK
// @Failure 400 id inválido
// @Failure 500 error interno
// @router /:id [delete]
func (c *DocumentoRequeridoBeneficioController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	if err := models.DeleteDocumentoRequeridoBeneficio(id); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}
