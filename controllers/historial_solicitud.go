package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/sga_crud_beneficios_egresados/models"
)

type HistorialSolicitudController struct{ web.Controller }

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

// GetBySolicitud GET /v1/historial_solicitud/solicitud/:solicitud_id
// Bitácora completa de una solicitud (más reciente primero).
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

// GetVigente GET /v1/historial_solicitud/solicitud/:solicitud_id/vigente
// Último registro de historial = estado vigente de la solicitud (C-4b).
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
