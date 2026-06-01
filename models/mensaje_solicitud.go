package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// MensajeSolicitud intercambio empresa ↔ egresado cuando la solicitud está en REQUIERE_INFO (RF-007).
type MensajeSolicitud struct {
	Id                  int                 `orm:"column(id);auto;pk" json:"id"`
	SolicitudBeneficio  *SolicitudBeneficio `orm:"column(solicitud_beneficio_id);rel(fk)" json:"solicitud_beneficio"`
	Usuario             *Usuario            `orm:"column(usuario_id);rel(fk)" json:"usuario"`
	Mensaje             string              `orm:"column(mensaje);type(text)" json:"mensaje"`
	FechaEnvio          time.Time           `orm:"column(fecha_envio);auto_now_add;type(datetime)" json:"fecha_envio"`
	Activo              bool                `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion       time.Time           `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion   time.Time           `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (m *MensajeSolicitud) TableName() string { return "mensaje_solicitud" }

func init() { orm.RegisterModel(new(MensajeSolicitud)) }

func AddMensajeSolicitud(m *MensajeSolicitud) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetMensajeSolicitudById(id int) (v *MensajeSolicitud, err error) {
	o := orm.NewOrm()
	v = &MensajeSolicitud{Id: id}
	if err = o.Read(v); err == nil {
		o.LoadRelated(v, "SolicitudBeneficio")
		o.LoadRelated(v, "Usuario")
		return v, nil
	}
	return nil, err
}

func GetAllMensajeSolicitud() (ml []MensajeSolicitud, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(MensajeSolicitud)).Filter("Activo", true).RelatedSel().All(&ml)
	return
}

func UpdateMensajeSolicitudById(m *MensajeSolicitud) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteMensajeSolicitud(id int) (err error) {
	o := orm.NewOrm()
	v := MensajeSolicitud{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
