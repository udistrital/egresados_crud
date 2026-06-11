package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// SecuenciaRadicado contador para generar radicados BNF-YYYY-NNNNNN.
// Usar SELECT FOR UPDATE al asignar (RN-RADICADO).
type SecuenciaRadicado struct {
	Id                int       `orm:"column(id);auto;pk" json:"id"`
	Anio              int       `orm:"column(anio);unique" json:"anio"`
	UltimoNumero      int       `orm:"column(ultimo_numero);default(0)" json:"ultimo_numero"`
	Activo            bool      `orm:"column(activo);default(true)" json:"activo"`
	FechaCreacion     time.Time `orm:"column(fecha_creacion);auto_now_add;type(datetime)" json:"fecha_creacion"`
	FechaModificacion time.Time `orm:"column(fecha_modificacion);auto_now;type(datetime)" json:"fecha_modificacion"`
}

func (s *SecuenciaRadicado) TableName() string { return "secuencia_radicado" }

func init() { orm.RegisterModel(new(SecuenciaRadicado)) }

func AddSecuenciaRadicado(m *SecuenciaRadicado) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetSecuenciaRadicadoById(id int) (v *SecuenciaRadicado, err error) {
	o := orm.NewOrm()
	v = &SecuenciaRadicado{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetAllSecuenciaRadicado(query map[string]string, fields []string, sortby []string,
	order []string, offset int64, limit int64) (ml []interface{}, err error) {
	qs := orm.NewOrm().QueryTable(new(SecuenciaRadicado))
	var l []SecuenciaRadicado
	return getAllQuery(qs, query, fields, sortby, order, offset, limit, &l)
}

// SiguienteRadicado incrementa atómicamente el contador del año actual y retorna el número.
// Usa una transacción con FOR UPDATE para evitar race conditions (RN-RADICADO).
func SiguienteRadicado(anio int) (numero int, err error) {
	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var seq SecuenciaRadicado
	err = tx.Raw("SELECT * FROM secuencia_radicado WHERE anio = ? FOR UPDATE", anio).QueryRow(&seq)
	if err != nil {
		// Si no existe la fila para este año, crearla
		seq = SecuenciaRadicado{Anio: anio, UltimoNumero: 0, Activo: true}
		if _, err = tx.Insert(&seq); err != nil {
			return
		}
	}

	seq.UltimoNumero++
	seq.FechaModificacion = time.Now()
	if _, err = tx.Update(&seq, "UltimoNumero", "FechaModificacion"); err != nil {
		return
	}

	err = tx.Commit()
	numero = seq.UltimoNumero
	return
}

func UpdateSecuenciaRadicadoById(m *SecuenciaRadicado) (err error) {
	o := orm.NewOrm()
	m.FechaModificacion = time.Now()
	_, err = o.Update(m)
	return
}

func DeleteSecuenciaRadicado(id int) (err error) {
	o := orm.NewOrm()
	v := SecuenciaRadicado{Id: id}
	if err = o.Read(&v); err == nil {
		v.Activo = false
		v.FechaModificacion = time.Now()
		_, err = o.Update(&v, "Activo", "FechaModificacion")
	}
	return
}
