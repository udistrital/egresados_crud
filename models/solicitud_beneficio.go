package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// SolicitudBeneficio solicitud de un egresado sobre un beneficio.
// No lleva estado propio (C-4b): el estado vigente es el último registro
// en historial_solicitud (ver GetEstadoVigenteBySolicitud).
type SolicitudBeneficio struct {
	Id                   int        `orm:"column(id);auto;pk" json:"id"`
	Radicado             string     `orm:"column(radicado);size(20);unique" json:"radicado"`
	Egresado             *Egresado  `orm:"column(egresado_id);rel(fk)" json:"egresado"`
	Beneficio            *Beneficio `orm:"column(beneficio_id);rel(fk)" json:"beneficio"`
	DatosComplementarios string     `orm:"column(datos_complementarios);type(jsonb);null" json:"datos_complementarios,omitempty"`
	FechaSolicitud       time.Time  `orm:"column(fecha_solicitud);auto_now_add;type(datetime)" json:"fecha_solicitud"`
	Activo               bool       `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion        time.Time  `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion    time.Time  `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (s *SolicitudBeneficio) TableName() string { return "solicitud_beneficio" }

func init() { orm.RegisterModel(new(SolicitudBeneficio)) }

func AddSolicitudBeneficio(m *SolicitudBeneficio) (id int64, err error) {
	m.Activo = true // toda fila creada nace activa (el default(true) del ORM no aplica en INSERT)
	o := orm.NewOrm()
	// C-5: el radicado BNF-YYYY-NNNNNN lo genera la secuencia nativa de PostgreSQL.
	// Si el caller no lo envía, se resuelve con fn_siguiente_radicado() (nextval atómico)
	// antes de insertar, de modo que el POST pueda devolverlo en la respuesta.
	if m.Radicado == "" {
		if err = o.Raw("SELECT fn_siguiente_radicado()").QueryRow(&m.Radicado); err != nil {
			return 0, err
		}
	}
	id, err = o.Insert(m)
	return
}

func GetSolicitudBeneficioById(id int) (v *SolicitudBeneficio, err error) {
	o := orm.NewOrm()
	v = &SolicitudBeneficio{Id: id}
	if err = o.Read(v); err == nil {
		o.LoadRelated(v, "Egresado")
		o.LoadRelated(v, "Beneficio")
		return v, nil
	}
	return nil, err
}

func GetAllSolicitudBeneficio(query map[string]string, fields []string, sortby []string,
	order []string, offset int64, limit int64) (ml []interface{}, err error) {
	qs := orm.NewOrm().QueryTable(new(SolicitudBeneficio)).RelatedSel()
	var l []SolicitudBeneficio
	return getAllQuery(qs, query, fields, sortby, order, offset, limit, &l)
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
