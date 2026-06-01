package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// HistorialEstadoSolicitud bitácora de transiciones de estado de cada solicitud (RN-004).
type HistorialEstadoSolicitud struct {
	Id                   int                 `orm:"column(id);auto;pk" json:"id"`
	SolicitudBeneficio   *SolicitudBeneficio `orm:"column(solicitud_beneficio_id);rel(fk)" json:"solicitud_beneficio"`
	EstadoAnterior       *EstadoSolicitud    `orm:"column(estado_anterior_id);rel(fk);null" json:"estado_anterior,omitempty"`
	EstadoNuevo          *EstadoSolicitud    `orm:"column(estado_nuevo_id);rel(fk)" json:"estado_nuevo"`
	Usuario              *Usuario            `orm:"column(usuario_id);rel(fk)" json:"usuario"`
	Justificacion        string              `orm:"column(justificacion);type(text);null" json:"justificacion,omitempty"`
	FechaCambio          time.Time           `orm:"column(fecha_cambio);auto_now_add;type(datetime)" json:"fecha_cambio"`
	Activo               bool                `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion        time.Time           `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion    time.Time           `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (h *HistorialEstadoSolicitud) TableName() string { return "historial_estado_solicitud" }

func init() { orm.RegisterModel(new(HistorialEstadoSolicitud)) }

func AddHistorialEstadoSolicitud(m *HistorialEstadoSolicitud) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetHistorialEstadoSolicitudById(id int) (v *HistorialEstadoSolicitud, err error) {
	o := orm.NewOrm()
	v = &HistorialEstadoSolicitud{Id: id}
	if err = o.Read(v); err == nil {
		o.LoadRelated(v, "SolicitudBeneficio")
		o.LoadRelated(v, "EstadoAnterior")
		o.LoadRelated(v, "EstadoNuevo")
		o.LoadRelated(v, "Usuario")
		return v, nil
	}
	return nil, err
}

func GetAllHistorialEstadoSolicitud() (ml []HistorialEstadoSolicitud, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(HistorialEstadoSolicitud)).Filter("Activo", true).RelatedSel().All(&ml)
	return
}

func UpdateHistorialEstadoSolicitudById(m *HistorialEstadoSolicitud) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteHistorialEstadoSolicitud(id int) (err error) {
	o := orm.NewOrm()
	v := HistorialEstadoSolicitud{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
