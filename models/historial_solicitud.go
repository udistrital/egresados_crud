package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// HistorialSolicitud bitácora de transiciones de estado de cada solicitud (RN-004).
// Única fuente de estado (C-4b): el estado vigente de una solicitud es el registro
// con mayor fecha_cambio. EstadoAnteriorId/EstadoNuevoId referencian
// parametro.parametro (tipo ESTADO_SOLICITUD); sin FK local (C-1).
type HistorialSolicitud struct {
	Id                       int                 `orm:"column(id);auto;pk" json:"id"`
	SolicitudBeneficio       *SolicitudBeneficio `orm:"column(solicitud_beneficio_id);rel(fk)" json:"solicitud_beneficio"`
	EstadoAnteriorId         *int                `orm:"column(estado_anterior_id);null" json:"estado_anterior_id,omitempty"`
	EstadoNuevoId            int                 `orm:"column(estado_nuevo_id)" json:"estado_nuevo_id"`
	Usuario                  *Usuario            `orm:"column(usuario_id);rel(fk)" json:"usuario"`
	Justificacion            string              `orm:"column(justificacion);type(text);null" json:"justificacion,omitempty"`
	// Comprobante OPCIONAL que la empresa adjunta al aprobar (solo en la transición a APROBADA).
	NombreArchivoComprobante string    `orm:"column(nombre_archivo_comprobante);size(300);null" json:"nombre_archivo_comprobante,omitempty"`
	EnlaceComprobante        string    `orm:"column(enlace_comprobante);size(100);null" json:"enlace_comprobante,omitempty"`
	FechaCambio              time.Time `orm:"column(fecha_cambio);auto_now_add;type(datetime)" json:"fecha_cambio"`
	Activo                   bool      `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion            time.Time `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion        time.Time `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (h *HistorialSolicitud) TableName() string { return "historial_solicitud" }

func init() { orm.RegisterModel(new(HistorialSolicitud)) }

func AddHistorialSolicitud(m *HistorialSolicitud) (id int64, err error) {
	m.Activo = true // toda fila creada nace activa (el default(true) del ORM no aplica en INSERT)
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetHistorialSolicitudById(id int) (v *HistorialSolicitud, err error) {
	o := orm.NewOrm()
	v = &HistorialSolicitud{Id: id}
	if err = o.Read(v); err == nil {
		o.LoadRelated(v, "SolicitudBeneficio")
		o.LoadRelated(v, "Usuario")
		return v, nil
	}
	return nil, err
}

func GetAllHistorialSolicitud(query map[string]string, fields []string, sortby []string,
	order []string, offset int64, limit int64) (ml []interface{}, err error) {
	qs := orm.NewOrm().QueryTable(new(HistorialSolicitud)).RelatedSel()
	var l []HistorialSolicitud
	return getAllQuery(qs, query, fields, sortby, order, offset, limit, &l)
}

// GetHistorialSolicitudBySolicitud retorna la bitácora completa de una solicitud,
// ordenada de más reciente a más antigua.
func GetHistorialSolicitudBySolicitud(solicitudId int) (ml []HistorialSolicitud, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(HistorialSolicitud)).
		Filter("SolicitudBeneficio__Id", solicitudId).
		Filter("Activo", true).
		OrderBy("-FechaCambio").
		All(&ml)
	return
}

// GetEstadoVigenteBySolicitud retorna el último registro de historial de la solicitud,
// que define su estado vigente (C-4b).
func GetEstadoVigenteBySolicitud(solicitudId int) (v *HistorialSolicitud, err error) {
	o := orm.NewOrm()
	v = &HistorialSolicitud{}
	err = o.QueryTable(new(HistorialSolicitud)).
		Filter("SolicitudBeneficio__Id", solicitudId).
		Filter("Activo", true).
		OrderBy("-FechaCambio").
		Limit(1).
		One(v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func UpdateHistorialSolicitudById(m *HistorialSolicitud) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteHistorialSolicitud(id int) (err error) {
	o := orm.NewOrm()
	v := HistorialSolicitud{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
