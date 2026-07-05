package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// DocumentoRequeridoBeneficio documento que la empresa exige al egresado para
// postularse a un beneficio (definido al publicar, RF-005). El archivo en sí
// lo sube el egresado por solicitud (ver DocumentoSolicitud).
type DocumentoRequeridoBeneficio struct {
	Id                int        `orm:"column(id);auto;pk" json:"id"`
	Beneficio         *Beneficio `orm:"column(beneficio_id);rel(fk)" json:"beneficio"`
	Nombre            string     `orm:"column(nombre);size(200)" json:"nombre"`
	Descripcion       string     `orm:"column(descripcion);type(text)" json:"descripcion"`
	Activo            bool       `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion     time.Time  `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion time.Time  `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (d *DocumentoRequeridoBeneficio) TableName() string { return "documento_requerido_beneficio" }

func init() { orm.RegisterModel(new(DocumentoRequeridoBeneficio)) }

func AddDocumentoRequeridoBeneficio(m *DocumentoRequeridoBeneficio) (id int64, err error) {
	m.Activo = true // toda fila creada nace activa (el default(true) del ORM no aplica en INSERT)
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetDocumentoRequeridoBeneficioById(id int) (v *DocumentoRequeridoBeneficio, err error) {
	o := orm.NewOrm()
	v = &DocumentoRequeridoBeneficio{Id: id}
	if err = o.Read(v); err == nil {
		o.LoadRelated(v, "Beneficio")
		return v, nil
	}
	return nil, err
}

func GetAllDocumentoRequeridoBeneficio(query map[string]string, fields []string, sortby []string,
	order []string, offset int64, limit int64) (ml []interface{}, err error) {
	qs := orm.NewOrm().QueryTable(new(DocumentoRequeridoBeneficio)).RelatedSel()
	var l []DocumentoRequeridoBeneficio
	return getAllQuery(qs, query, fields, sortby, order, offset, limit, &l)
}

func UpdateDocumentoRequeridoBeneficioById(m *DocumentoRequeridoBeneficio) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteDocumentoRequeridoBeneficio(id int) (err error) {
	o := orm.NewOrm()
	v := DocumentoRequeridoBeneficio{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
