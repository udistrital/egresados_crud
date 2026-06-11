package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

// getAllQuery aplica a un QuerySeter el contrato estándar de listado de los
// *_crud institucionales (query/fields/sortby/order/offset/limit) y retorna
// la lista resultante. container debe ser un puntero a slice del modelo
// (p. ej. *[]Beneficio); se usa para que el ORM materialice las filas.
//
// Semántica de query (la de terceros_crud, la variante más completa del SGA):
//   - dot-notation → relaciones del ORM: Empresa.Id → Empresa__Id
//   - sufijo isnull: valor true/1 o false/0
//   - sufijo __in: valores separados por |
//   - sufijo __icontainsall: todas las palabras separadas por | deben aparecer
//     (icontains por cada palabra, sin importar el orden)
//   - resto: Filter directo — acepta los operadores nativos del ORM
//     (__gte, __lte, __icontains, __startswith, ...)
//
// limit 0 = sin límite (DefaultRowsLimit de beego v2 es -1).
func getAllQuery(qs orm.QuerySeter, query map[string]string, fields []string,
	sortby []string, order []string, offset int64, limit int64,
	container interface{}) (ml []interface{}, err error) {

	for k, v := range query {
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, v == "true" || v == "1")
		} else if strings.HasSuffix(k, "__in") {
			qs = qs.Filter(k, strings.Split(v, "|"))
		} else if strings.HasSuffix(k, "__icontainsall") {
			k = strings.TrimSuffix(k, "all")
			for _, word := range strings.Split(v, "|") {
				qs = qs.Filter(k, word)
			}
		} else {
			qs = qs.Filter(k, v)
		}
	}

	var sortFields []string
	if len(sortby) != 0 {
		if len(order) != len(sortby) && len(order) != 1 {
			return nil, errors.New("'sortby' y 'order' no coinciden en tamaño (o 'order' debe ser único)")
		}
		for i, v := range sortby {
			o := order[0]
			if len(order) == len(sortby) {
				o = order[i]
			}
			switch o {
			case "desc":
				sortFields = append(sortFields, "-"+v)
			case "asc":
				sortFields = append(sortFields, v)
			default:
				return nil, errors.New("'order' inválido: debe ser asc o desc")
			}
		}
		qs = qs.OrderBy(sortFields...)
	} else if len(order) != 0 {
		return nil, errors.New("'order' sin 'sortby'")
	}

	if _, err = qs.Limit(limit, offset).All(container, fields...); err != nil {
		return nil, err
	}

	lv := reflect.Indirect(reflect.ValueOf(container))
	for i := 0; i < lv.Len(); i++ {
		item := lv.Index(i)
		if len(fields) == 0 {
			ml = append(ml, item.Interface())
		} else {
			// recortar a los campos pedidos (nombres de campo Go, como en terceros_crud)
			m := make(map[string]interface{})
			for _, fname := range fields {
				f := item.FieldByName(fname)
				if !f.IsValid() {
					return nil, fmt.Errorf("campo desconocido en 'fields': %s", fname)
				}
				m[fname] = f.Interface()
			}
			ml = append(ml, m)
		}
	}
	return ml, nil
}
