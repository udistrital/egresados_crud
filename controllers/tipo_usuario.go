package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/sga_crud_beneficios_egresados/models"
)

type TipoUsuarioController struct{ web.Controller }

// GetAll @router /v1/tipo_usuario [get]
func (c *TipoUsuarioController) GetAll() {
	results, err := models.GetAllTipoUsuario()
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = results
	}
	c.ServeJSON()
}

// GetOne @router /v1/tipo_usuario/:id [get]
func (c *TipoUsuarioController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "id inválido"
		c.ServeJSON()
		return
	}
	result, err := models.GetTipoUsuarioById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = result
	}
	c.ServeJSON()
}

// Post @router /v1/tipo_usuario [post]
func (c *TipoUsuarioController) Post() {
	var v models.TipoUsuario
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	if id, err := models.AddTipoUsuario(&v); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = map[string]int64{"id": id}
	}
	c.ServeJSON()
}

// Put @router /v1/tipo_usuario/:id [put]
func (c *TipoUsuarioController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "id inválido"
		c.ServeJSON()
		return
	}
	var v models.TipoUsuario
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	v.Id = id
	if err := models.UpdateTipoUsuarioById(&v); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}

// Delete @router /v1/tipo_usuario/:id [delete]
func (c *TipoUsuarioController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "id inválido"
		c.ServeJSON()
		return
	}
	if err := models.DeleteTipoUsuario(id); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}
