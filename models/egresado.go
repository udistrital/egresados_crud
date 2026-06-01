package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Egresado perfil de egresado (1:1 con usuario tipo EGRESADO). Espejo del SGA + datos propios del módulo.
type Egresado struct {
	Id                   int       `orm:"column(id);auto;pk" json:"id"`
	Usuario              *Usuario  `orm:"column(usuario_id);rel(fk);unique" json:"usuario"`
	CodigoInstitucional  string    `orm:"column(codigo_institucional);size(20);unique" json:"codigo_institucional"`
	ProgramaAcademico    string    `orm:"column(programa_academico);size(150);null" json:"programa_academico,omitempty"`
	Facultad             string    `orm:"column(facultad);size(150);null" json:"facultad,omitempty"`
	FechaGrado           time.Time `orm:"column(fecha_grado);null;type(date)" json:"fecha_grado,omitempty"`
	TelefonoContacto     string    `orm:"column(telefono_contacto);size(20);null" json:"telefono_contacto,omitempty"`
	Activo               bool      `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion        time.Time `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion    time.Time `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (e *Egresado) TableName() string { return "egresado" }

func init() { orm.RegisterModel(new(Egresado)) }

func AddEgresado(m *Egresado) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetEgresadoById(id int) (v *Egresado, err error) {
	o := orm.NewOrm()
	v = &Egresado{Id: id}
	if err = o.Read(v); err == nil {
		o.LoadRelated(v, "Usuario")
		return v, nil
	}
	return nil, err
}

func GetAllEgresado() (ml []Egresado, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(Egresado)).Filter("Activo", true).RelatedSel().All(&ml)
	return
}

func UpdateEgresadoById(m *Egresado) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteEgresado(id int) (err error) {
	o := orm.NewOrm()
	v := Egresado{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
