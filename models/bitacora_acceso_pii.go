package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// BitacoraAccesoPii bitácora inmutable de accesos a datos personales (Ley 1581 / RNF-002a).
// Retención mínima 6 meses. No tiene borrado lógico por ser log inmutable.
type BitacoraAccesoPii struct {
	Id          int       `orm:"column(id);auto;pk" json:"id"`
	Usuario     *Usuario  `orm:"column(usuario_id);rel(fk)" json:"usuario"`
	RecursoTipo string    `orm:"column(recurso_tipo);size(50)" json:"recurso_tipo"`
	RecursoId   int       `orm:"column(recurso_id);null" json:"recurso_id,omitempty"`
	Accion      string    `orm:"column(accion);size(50)" json:"accion"`
	DireccionIp string    `orm:"column(direccion_ip);size(45);null" json:"direccion_ip,omitempty"`
	UserAgent   string    `orm:"column(user_agent);size(500);null" json:"user_agent,omitempty"`
	Detalle     string    `orm:"column(detalle);type(jsonb);null" json:"detalle,omitempty"`
	FechaEvento time.Time `orm:"column(fecha_evento);auto_now_add;type(datetime)" json:"fecha_evento"`
}

func (b *BitacoraAccesoPii) TableName() string { return "bitacora_acceso_pii" }

func init() { orm.RegisterModel(new(BitacoraAccesoPii)) }

func AddBitacoraAccesoPii(m *BitacoraAccesoPii) (id int64, err error) {
	// log inmutable: el modelo no mapea 'activo' (la columna tiene DEFAULT TRUE en la BD)
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetBitacoraAccesoPiiById(id int) (v *BitacoraAccesoPii, err error) {
	o := orm.NewOrm()
	v = &BitacoraAccesoPii{Id: id}
	if err = o.Read(v); err == nil {
		o.LoadRelated(v, "Usuario")
		return v, nil
	}
	return nil, err
}

func GetAllBitacoraAccesoPii(query map[string]string, fields []string, sortby []string,
	order []string, offset int64, limit int64) (ml []interface{}, err error) {
	qs := orm.NewOrm().QueryTable(new(BitacoraAccesoPii)).RelatedSel()
	var l []BitacoraAccesoPii
	return getAllQuery(qs, query, fields, sortby, order, offset, limit, &l)
}
