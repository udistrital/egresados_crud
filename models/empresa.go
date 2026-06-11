package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Empresa empresa aliada. Espejo local de Ágora + estado de ciclo de vida del módulo.
// SectorEconomicoId y EstadoEmpresaId referencian parametro.parametro
// (tipos SECTOR_ECONOMICO y ESTADO_EMPRESA); sin FK local (C-1).
type Empresa struct {
	Id                int       `orm:"column(id);auto;pk" json:"id"`
	Nit               string    `orm:"column(nit);size(20);unique" json:"nit"`
	RazonSocial       string    `orm:"column(razon_social);size(200)" json:"razon_social"`
	AgoraIdExterno    string    `orm:"column(agora_id_externo);size(50);null" json:"agora_id_externo,omitempty"`
	SectorEconomicoId *int      `orm:"column(sector_economico_id);null" json:"sector_economico_id,omitempty"`
	EstadoEmpresaId   int       `orm:"column(estado_empresa_id)" json:"estado_empresa_id"`
	SitioWeb          string    `orm:"column(sitio_web);size(255);null" json:"sitio_web,omitempty"`
	CorreoContacto    string    `orm:"column(correo_contacto);size(150);null" json:"correo_contacto,omitempty"`
	TelefonoContacto  string    `orm:"column(telefono_contacto);size(20);null" json:"telefono_contacto,omitempty"`
	Direccion         string    `orm:"column(direccion);size(255);null" json:"direccion,omitempty"`
	FechaAprobacion   time.Time `orm:"column(fecha_aprobacion);null;type(datetime)" json:"fecha_aprobacion,omitempty"`
	UsuarioAprobador  *Usuario  `orm:"column(usuario_aprobador_id);rel(fk);null" json:"usuario_aprobador,omitempty"`
	Activo            bool      `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion     time.Time `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion time.Time `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (e *Empresa) TableName() string { return "empresa" }

func init() { orm.RegisterModel(new(Empresa)) }

func AddEmpresa(m *Empresa) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetEmpresaById(id int) (v *Empresa, err error) {
	o := orm.NewOrm()
	v = &Empresa{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetAllEmpresa(query map[string]string, fields []string, sortby []string,
	order []string, offset int64, limit int64) (ml []interface{}, err error) {
	qs := orm.NewOrm().QueryTable(new(Empresa)).RelatedSel()
	var l []Empresa
	return getAllQuery(qs, query, fields, sortby, order, offset, limit, &l)
}

func UpdateEmpresaById(m *Empresa) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteEmpresa(id int) (err error) {
	o := orm.NewOrm()
	v := Empresa{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
