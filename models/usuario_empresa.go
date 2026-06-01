package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// UsuarioEmpresa relación N:M entre usuarios tipo EMPRESA y empresas.
type UsuarioEmpresa struct {
	Id                int       `orm:"column(id);auto;pk" json:"id"`
	Usuario           *Usuario  `orm:"column(usuario_id);rel(fk)" json:"usuario"`
	Empresa           *Empresa  `orm:"column(empresa_id);rel(fk)" json:"empresa"`
	Cargo             string    `orm:"column(cargo);size(100);null" json:"cargo,omitempty"`
	EsPrincipal       bool      `orm:"column(es_principal);default(false)" json:"es_principal"`
	Activo            bool      `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion     time.Time `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion time.Time `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (u *UsuarioEmpresa) TableName() string { return "usuario_empresa" }

func init() { orm.RegisterModel(new(UsuarioEmpresa)) }

func AddUsuarioEmpresa(m *UsuarioEmpresa) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetUsuarioEmpresaById(id int) (v *UsuarioEmpresa, err error) {
	o := orm.NewOrm()
	v = &UsuarioEmpresa{Id: id}
	if err = o.Read(v); err == nil {
		o.LoadRelated(v, "Usuario")
		o.LoadRelated(v, "Empresa")
		return v, nil
	}
	return nil, err
}

func GetAllUsuarioEmpresa() (ml []UsuarioEmpresa, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(UsuarioEmpresa)).Filter("Activo", true).RelatedSel().All(&ml)
	return
}

func UpdateUsuarioEmpresaById(m *UsuarioEmpresa) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteUsuarioEmpresa(id int) (err error) {
	o := orm.NewOrm()
	v := UsuarioEmpresa{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
