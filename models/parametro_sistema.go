package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// ParametroSistema parámetros configurables del sistema (ej. límite de solicitudes activas).
type ParametroSistema struct {
	Id                int       `orm:"column(id);auto;pk" json:"id"`
	Clave             string    `orm:"column(clave);size(100);unique" json:"clave"`
	Valor             string    `orm:"column(valor);size(500)" json:"valor"`
	TipoDato          string    `orm:"column(tipo_dato);size(20)" json:"tipo_dato"`
	Descripcion       string    `orm:"column(descripcion);size(500);null" json:"descripcion,omitempty"`
	Activo            bool      `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion     time.Time `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion time.Time `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (p *ParametroSistema) TableName() string { return "parametro_sistema" }

func init() { orm.RegisterModel(new(ParametroSistema)) }

func AddParametroSistema(m *ParametroSistema) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetParametroSistemaById(id int) (v *ParametroSistema, err error) {
	o := orm.NewOrm()
	v = &ParametroSistema{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetParametroSistemaByClave(clave string) (v *ParametroSistema, err error) {
	o := orm.NewOrm()
	v = &ParametroSistema{Clave: clave}
	if err = o.Read(v, "Clave"); err == nil {
		return v, nil
	}
	return nil, err
}

func GetAllParametroSistema() (ml []ParametroSistema, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(ParametroSistema)).Filter("Activo", true).All(&ml)
	return
}

func UpdateParametroSistemaById(m *ParametroSistema) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteParametroSistema(id int) (err error) {
	o := orm.NewOrm()
	v := ParametroSistema{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
