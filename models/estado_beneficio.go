package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// EstadoBeneficio catálogo de estados de un beneficio publicado.
type EstadoBeneficio struct {
	Id                int       `orm:"column(id);auto;pk" json:"id"`
	Nombre            string    `orm:"column(nombre);size(100)" json:"nombre"`
	Descripcion       string    `orm:"column(descripcion);size(500);null" json:"descripcion,omitempty"`
	CodigoAbreviacion string    `orm:"column(codigo_abreviacion);size(50);unique;null" json:"codigo_abreviacion,omitempty"`
	NumeroOrden       float64   `orm:"column(numero_orden);digits(5);decimals(2);null" json:"numero_orden,omitempty"`
	Activo            bool      `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion     time.Time `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion time.Time `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (e *EstadoBeneficio) TableName() string { return "estado_beneficio" }

func init() { orm.RegisterModel(new(EstadoBeneficio)) }

func AddEstadoBeneficio(m *EstadoBeneficio) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetEstadoBeneficioById(id int) (v *EstadoBeneficio, err error) {
	o := orm.NewOrm()
	v = &EstadoBeneficio{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetAllEstadoBeneficio() (ml []EstadoBeneficio, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(EstadoBeneficio)).Filter("Activo", true).OrderBy("NumeroOrden").All(&ml)
	return
}

func UpdateEstadoBeneficioById(m *EstadoBeneficio) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteEstadoBeneficio(id int) (err error) {
	o := orm.NewOrm()
	v := EstadoBeneficio{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
