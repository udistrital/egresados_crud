package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Usuario identidad local creada por JIT provisioning al autenticarse contra SGA (egresados) o Ágora (empresas).
type Usuario struct {
	Id                int          `orm:"column(id);auto;pk" json:"id"`
	Documento         string       `orm:"column(documento);size(20);unique" json:"documento"`
	Nombre            string       `orm:"column(nombre);size(200)" json:"nombre"`
	Correo            string       `orm:"column(correo);size(150)" json:"correo"`
	TipoUsuario       *TipoUsuario `orm:"column(tipo_usuario_id);rel(fk)" json:"tipo_usuario"`
	IdExterno         string       `orm:"column(id_externo);size(50);null" json:"id_externo,omitempty"`
	SistemaOrigen     string       `orm:"column(sistema_origen);size(20)" json:"sistema_origen"`
	UltimoAcceso      time.Time    `orm:"column(ultimo_acceso);null;type(datetime)" json:"ultimo_acceso,omitempty"`
	Activo            bool         `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion     time.Time    `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion time.Time    `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (u *Usuario) TableName() string { return "usuario" }

func init() { orm.RegisterModel(new(Usuario)) }

func AddUsuario(m *Usuario) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetUsuarioById(id int) (v *Usuario, err error) {
	o := orm.NewOrm()
	v = &Usuario{Id: id}
	if err = o.Read(v); err == nil {
		o.LoadRelated(v, "TipoUsuario")
		return v, nil
	}
	return nil, err
}

func GetUsuarioByDocumento(documento string) (v *Usuario, err error) {
	o := orm.NewOrm()
	v = &Usuario{Documento: documento}
	if err = o.Read(v, "Documento"); err == nil {
		o.LoadRelated(v, "TipoUsuario")
		return v, nil
	}
	return nil, err
}

func GetAllUsuario() (ml []Usuario, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(Usuario)).Filter("Activo", true).RelatedSel().All(&ml)
	return
}

func UpdateUsuarioById(m *Usuario) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteUsuario(id int) (err error) {
	o := orm.NewOrm()
	v := Usuario{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
