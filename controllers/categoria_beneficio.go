package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/sga_crud_beneficios_egresados/models"
)

type CategoriaBeneficioController struct{ web.Controller }

func (c *CategoriaBeneficioController) GetAll() {
	results, err := models.GetAllCategoriaBeneficio()
	if err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = results
	}
	c.ServeJSON()
}

func (c *CategoriaBeneficioController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	result, err := models.GetCategoriaBeneficioById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(404); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = result
	}
	c.ServeJSON()
}

func (c *CategoriaBeneficioController) Post() {
	var v models.CategoriaBeneficio
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	if id, err := models.AddCategoriaBeneficio(&v); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Ctx.Output.SetStatus(201); c.Data["json"] = map[string]int64{"id": id}
	}
	c.ServeJSON()
}

func (c *CategoriaBeneficioController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	var v models.CategoriaBeneficio
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	v.Id = id
	if err := models.UpdateCategoriaBeneficioById(&v); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}

func (c *CategoriaBeneficioController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	if err := models.DeleteCategoriaBeneficio(id); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "OK"
	}
	c.ServeJSON()
}
