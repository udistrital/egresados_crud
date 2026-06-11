package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/sga_crud_beneficios_egresados/models"
)

// BitacoraAccesoPiiController — log inmutable, solo permite GET y POST.
type BitacoraAccesoPiiController struct{ web.Controller }

func (c *BitacoraAccesoPiiController) GetAll() {
	query, fields, sortby, order, offset, limit, err := parseGetAllParams(&c.Controller)
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	l, err := models.GetAllBitacoraAccesoPii(query, fields, sortby, order, offset, limit)
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

func (c *BitacoraAccesoPiiController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = "id inválido"; c.ServeJSON(); return
	}
	result, err := models.GetBitacoraAccesoPiiById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(404); c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = result
	}
	c.ServeJSON()
}

func (c *BitacoraAccesoPiiController) Post() {
	var v models.BitacoraAccesoPii
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Ctx.Output.SetStatus(400); c.Data["json"] = err.Error(); c.ServeJSON(); return
	}
	if id, err := models.AddBitacoraAccesoPii(&v); err != nil {
		c.Ctx.Output.SetStatus(500); c.Data["json"] = err.Error()
	} else {
		c.Ctx.Output.SetStatus(201); c.Data["json"] = map[string]int64{"id": id}
	}
	c.ServeJSON()
}
