package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// SolicitudBeneficio solicitud de un egresado sobre un beneficio.
type SolicitudBeneficio struct {
	Id                    int              `orm:"column(id);auto;pk" json:"id"`
	Radicado              string           `orm:"column(radicado);size(20);unique" json:"radicado"`
	Egresado              *Egresado        `orm:"column(egresado_id);rel(fk)" json:"egresado"`
	Beneficio             *Beneficio       `orm:"column(beneficio_id);rel(fk)" json:"beneficio"`
	EstadoSolicitud       *EstadoSolicitud `orm:"column(estado_solicitud_id);rel(fk)" json:"estado_solicitud"`
	DatosComplementarios  string           `orm:"column(datos_complementarios);type(jsonb);null" json:"datos_complementarios,omitempty"`
	FechaSolicitud        time.Time        `orm:"column(fecha_solicitud);auto_now_add;type(datetime)" json:"fecha_solicitud"`
	Activo                bool             `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion         time.Time        `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion     time.Time        `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (s *SolicitudBeneficio) TableName() string { return "solicitud_beneficio" }

func init() { orm.RegisterModel(new(SolicitudBeneficio)) }

func AddSolicitudBeneficio(m *SolicitudBeneficio) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetSolicitudBeneficioById(id int) (v *SolicitudBeneficio, err error) {
	o := orm.NewOrm()
	v = &SolicitudBeneficio{Id: id}
	if err = o.Read(v); err == nil {
		o.LoadRelated(v, "Egresado")
		o.LoadRelated(v, "Beneficio")
		o.LoadRelated(v, "EstadoSolicitud")
		return v, nil
	}
	return nil, err
}

func GetAllSolicitudBeneficio() (ml []SolicitudBeneficio, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(SolicitudBeneficio)).Filter("Activo", true).RelatedSel().All(&ml)
	return
}

func UpdateSolicitudBeneficioById(m *SolicitudBeneficio) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteSolicitudBeneficio(id int) (err error) {
	o := orm.NewOrm()
	v := SolicitudBeneficio{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
