package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// EstadoSolicitud catálogo de estados de una solicitud de beneficio.
type EstadoSolicitud struct {
	Id                int       `orm:"column(id);auto;pk" json:"id"`
	Nombre            string    `orm:"column(nombre);size(100)" json:"nombre"`
	Descripcion       string    `orm:"column(descripcion);size(500);null" json:"descripcion,omitempty"`
	CodigoAbreviacion string    `orm:"column(codigo_abreviacion);size(50);unique;null" json:"codigo_abreviacion,omitempty"`
	NumeroOrden       float64   `orm:"column(numero_orden);digits(5);decimals(2);null" json:"numero_orden,omitempty"`
	Activo            bool      `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion     time.Time `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion time.Time `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (e *EstadoSolicitud) TableName() string { return "estado_solicitud" }

func init() { orm.RegisterModel(new(EstadoSolicitud)) }

func AddEstadoSolicitud(m *EstadoSolicitud) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetEstadoSolicitudById(id int) (v *EstadoSolicitud, err error) {
	o := orm.NewOrm()
	v = &EstadoSolicitud{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetAllEstadoSolicitud() (ml []EstadoSolicitud, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(EstadoSolicitud)).Filter("Activo", true).OrderBy("NumeroOrden").All(&ml)
	return
}

func UpdateEstadoSolicitudById(m *EstadoSolicitud) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteEstadoSolicitud(id int) (err error) {
	o := orm.NewOrm()
	v := EstadoSolicitud{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
