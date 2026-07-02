package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// UsuarioEmpresa relación N:M entre usuarios tipo EMP y empresas.
// Subtipo EXCLUSIVO de usuario (C-7): TipoUsuario fijado a 'EMP' participa en la FK compuesta
// (usuario_id, tipo_usuario) -> usuario(id, tipo_usuario) declarada a nivel DDL.
type UsuarioEmpresa struct {
	Id                int       `orm:"column(id);auto;pk" json:"id"`
	Usuario           *Usuario  `orm:"column(usuario_id);rel(fk)" json:"usuario"`
	TipoUsuario       string    `orm:"column(tipo_usuario);size(3);default(EMP)" json:"tipo_usuario"`
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
	if m.TipoUsuario == "" {
		m.TipoUsuario = "EMP" // discriminador fijo del subtipo (C-7)
	}
	m.Activo = true // toda fila creada nace activa (el default(true) del ORM no aplica en INSERT)
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

func GetAllUsuarioEmpresa(query map[string]string, fields []string, sortby []string,
	order []string, offset int64, limit int64) (ml []interface{}, err error) {
	qs := orm.NewOrm().QueryTable(new(UsuarioEmpresa)).RelatedSel()
	var l []UsuarioEmpresa
	return getAllQuery(qs, query, fields, sortby, order, offset, limit, &l)
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
