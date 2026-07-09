package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Beneficio publicado por una empresa aliada para egresados.
// CategoriaBeneficioId y EstadoBeneficioId referencian parametro.parametro
// (tipos CATEGORIA_BENEFICIO y ESTADO_BENEFICIO); sin FK local (C-1).
type Beneficio struct {
	Id                   int       `orm:"column(id);auto;pk" json:"id"`
	Empresa              *Empresa  `orm:"column(empresa_id);rel(fk)" json:"empresa"`
	CategoriaBeneficioId int       `orm:"column(categoria_beneficio_id)" json:"categoria_beneficio_id"`
	EstadoBeneficioId    int       `orm:"column(estado_beneficio_id)" json:"estado_beneficio_id"`
	Titulo               string    `orm:"column(titulo);size(200)" json:"titulo"`
	Descripcion          string    `orm:"column(descripcion);type(text)" json:"descripcion"`
	Condiciones          string    `orm:"column(condiciones);type(text)" json:"condiciones"`
	FechaInicio          time.Time `orm:"column(fecha_inicio);type(date)" json:"fecha_inicio"`
	FechaFin             time.Time `orm:"column(fecha_fin);type(date)" json:"fecha_fin"`
	CuposTotal           int       `orm:"column(cupos_total)" json:"cupos_total"`
	CuposDisponibles     int       `orm:"column(cupos_disponibles)" json:"cupos_disponibles"`
	ImagenUrl            string    `orm:"column(imagen_url);size(500);null" json:"imagen_url,omitempty"`
	FechaPublicacion     time.Time `orm:"column(fecha_publicacion);null;type(datetime)" json:"fecha_publicacion,omitempty"`
	UsuarioCreador       *Usuario  `orm:"column(usuario_creador_id);rel(fk)" json:"usuario_creador"`
	Activo               bool      `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion        time.Time `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion    time.Time `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (b *Beneficio) TableName() string { return "beneficio" }

func init() { orm.RegisterModel(new(Beneficio)) }

func AddBeneficio(m *Beneficio) (id int64, err error) {
	m.Activo = true // toda fila creada nace activa (el default(true) del ORM no aplica en INSERT)
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetBeneficioById(id int) (v *Beneficio, err error) {
	o := orm.NewOrm()
	v = &Beneficio{Id: id}
	if err = o.Read(v); err == nil {
		o.LoadRelated(v, "Empresa")
		o.LoadRelated(v, "UsuarioCreador")
		return v, nil
	}
	return nil, err
}

func GetAllBeneficio(query map[string]string, fields []string, sortby []string,
	order []string, offset int64, limit int64) (ml []interface{}, err error) {
	qs := orm.NewOrm().QueryTable(new(Beneficio)).RelatedSel()
	var l []Beneficio
	return getAllQuery(qs, query, fields, sortby, order, offset, limit, &l)
}

func UpdateBeneficioById(m *Beneficio) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteBeneficio(id int) (err error) {
	o := orm.NewOrm()
	v := Beneficio{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}

// DescontarCupo resta 1 a cupos_disponibles de forma ATÓMICA (RN-002b). El guard
// cupos_disponibles > 0 en el propio UPDATE evita la condición de carrera de dos
// solicitudes concurrentes sobre el último cupo. Devuelve true si descontó.
func DescontarCupo(id int) (descontado bool, err error) {
	res, err := orm.NewOrm().Raw(
		"UPDATE beneficio SET cupos_disponibles = cupos_disponibles - 1, fecha_modificacion = NOW() "+
			"WHERE id = ? AND activo = TRUE AND cupos_disponibles > 0", id).Exec()
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// DevolverCupo suma 1 a cupos_disponibles sin exceder cupos_total (RN-002c). También
// atómico. Devuelve true si devolvió (false si ya estaba en el tope, no es error).
func DevolverCupo(id int) (devuelto bool, err error) {
	res, err := orm.NewOrm().Raw(
		"UPDATE beneficio SET cupos_disponibles = cupos_disponibles + 1, fecha_modificacion = NOW() "+
			"WHERE id = ? AND cupos_disponibles < cupos_total", id).Exec()
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}
