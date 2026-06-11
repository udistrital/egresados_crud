package controllers

import (
	"errors"
	"strings"

	"github.com/beego/beego/v2/server/web"
)

// parseGetAllParams extrae los parámetros estándar de listado de los *_crud
// institucionales: query (k:v,k:v con dot-notation), fields, sortby, order,
// limit (default 10; 0 = todos) y offset.
func parseGetAllParams(c *web.Controller) (query map[string]string, fields []string,
	sortby []string, order []string, offset int64, limit int64, err error) {
	query = make(map[string]string)
	limit = 10

	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	if v, e := c.GetInt64("limit"); e == nil {
		limit = v
	}
	if v, e := c.GetInt64("offset"); e == nil {
		offset = v
	}
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				err = errors.New("query inválido: se espera k:v,k:v")
				return
			}
			query[kv[0]] = kv[1]
		}
	}
	return
}
