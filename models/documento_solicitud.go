package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// DocumentoSolicitud PDF subido por el egresado para cumplir un
// DocumentoRequeridoBeneficio dentro de una solicitud puntual. El binario vive
// en el servicio institucional gestor_documental_mid (Nuxeo);
// EnlaceGestorDocumental es el uid ("Enlace") que ese servicio devuelve al
// subir — referencia lógica, sin FK (el servicio es externo a este esquema).
type DocumentoSolicitud struct {
	Id                   int                          `orm:"column(id);auto;pk" json:"id"`
	SolicitudBeneficio   *SolicitudBeneficio           `orm:"column(solicitud_beneficio_id);rel(fk)" json:"solicitud_beneficio"`
	DocumentoRequerido   *DocumentoRequeridoBeneficio  `orm:"column(documento_requerido_id);rel(fk)" json:"documento_requerido"`
	NombreArchivo        string                       `orm:"column(nombre_archivo);size(300)" json:"nombre_archivo"`
	EnlaceGestorDocumental string                     `orm:"column(enlace_gestor_documental);size(100)" json:"enlace_gestor_documental"`
	ComentarioEmpresa    string                       `orm:"column(comentario_empresa);type(text);null" json:"comentario_empresa,omitempty"`
	FechaComentario      time.Time                    `orm:"column(fecha_comentario);null;type(datetime)" json:"fecha_comentario,omitempty"`
	Activo               bool                         `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion        time.Time                    `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion    time.Time                    `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (d *DocumentoSolicitud) TableName() string { return "documento_solicitud" }

func init() { orm.RegisterModel(new(DocumentoSolicitud)) }

func AddDocumentoSolicitud(m *DocumentoSolicitud) (id int64, err error) {
	m.Activo = true // toda fila creada nace activa (el default(true) del ORM no aplica en INSERT)
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetDocumentoSolicitudById(id int) (v *DocumentoSolicitud, err error) {
	o := orm.NewOrm()
	v = &DocumentoSolicitud{Id: id}
	if err = o.Read(v); err == nil {
		o.LoadRelated(v, "SolicitudBeneficio")
		o.LoadRelated(v, "DocumentoRequerido")
		return v, nil
	}
	return nil, err
}

func GetAllDocumentoSolicitud(query map[string]string, fields []string, sortby []string,
	order []string, offset int64, limit int64) (ml []interface{}, err error) {
	qs := orm.NewOrm().QueryTable(new(DocumentoSolicitud)).RelatedSel()
	var l []DocumentoSolicitud
	return getAllQuery(qs, query, fields, sortby, order, offset, limit, &l)
}

// GetDocumentoSolicitudBySolicitud lista los documentos subidos (activos) de una
// solicitud. Mismo patrón que HistorialSolicitud.GetHistorialSolicitudBySolicitud.
func GetDocumentoSolicitudBySolicitud(solicitudId int) (ml []DocumentoSolicitud, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(DocumentoSolicitud)).RelatedSel().
		Filter("SolicitudBeneficio__Id", solicitudId).
		Filter("Activo", true).
		All(&ml)
	return
}

func UpdateDocumentoSolicitudById(m *DocumentoSolicitud) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteDocumentoSolicitud(id int) (err error) {
	o := orm.NewOrm()
	v := DocumentoSolicitud{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
