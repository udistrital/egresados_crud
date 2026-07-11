package controllers_test

import (
	"net/http"
	"testing"
)

// Ruta renombrada de guion_bajo a guion (lineamientos-endpoints.md): /usuario-empresa.
func TestGetAllUsuarioEmpresa(t *testing.T) {
	response, err := http.Get(baseURL + "/usuario-empresa")
	if err != nil {
		t.Error("Error GetAll usuario-empresa:", err.Error())
		t.Fail()
		return
	}
	if response.StatusCode != 200 {
		t.Error("Error en GetAll usuario-empresa, se esperaba 200 y se obtuvo", response.StatusCode)
		t.Fail()
	} else {
		t.Log("GetAll usuario-empresa finalizado correctamente (OK)")
	}
}

func TestPostUsuarioEmpresa(t *testing.T) {
	usuarioId := crearUsuario(t, "EMP")
	empresaId := crearEmpresa(t)
	body := map[string]interface{}{
		"usuario": map[string]interface{}{"id": usuarioId},
		"empresa": map[string]interface{}{"id": empresaId},
	}
	status, result := postJSON(t, "/usuario-empresa", body)
	if status != 201 {
		t.Error("Error en Post usuario-empresa, se esperaba 201 y se obtuvo", status, result)
		t.Fail()
		return
	}
	t.Log("Post usuario-empresa finalizado correctamente (OK), id creado:", idFrom(t, result))
}
