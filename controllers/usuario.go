package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/sga_crud_beneficios_egresados/models"
)

type UsuarioController struct{ web.Controller }

func (c *UsuarioController) GetAll() {
	results, err := models.GetAllUsuario()
	if err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = results
	}
	c.ServeJSON()
}

func (c *UsuarioController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	result, err := models.GetUsuarioById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(404); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = result
	}
	c.ServeJSON()
}

func (c *UsuarioController) Post() {
	var v models.Usuario
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	if id, err := models.AddUsuario(&v); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Ctx.Output.SetStatus(201); c.Data["json"] = map[string]int64{"id": id}
	}
	c.ServeJSON()
}

func (c *UsuarioController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	var v models.Usuario
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	v.Id = id
	if err := models.UpdateUsuarioById(&v); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}

func (c *UsuarioController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	if err := models.DeleteUsuario(id); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}
