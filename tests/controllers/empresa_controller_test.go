package controllers_test

import (
	"net/http"
	"testing"
)

func TestGetAllEmpresa(t *testing.T) {
	response, err := http.Get(baseURL + "/empresa")
	if err != nil {
		t.Error("Error GetAll empresa:", err.Error())
		t.Fail()
		return
	}
	if response.StatusCode != 200 {
		t.Error("Error en GetAll empresa, se esperaba 200 y se obtuvo", response.StatusCode)
		t.Fail()
	} else {
		t.Log("GetAll empresa finalizado correctamente (OK)")
	}
}

func TestPostEmpresa(t *testing.T) {
	id := crearEmpresa(t)
	t.Log("Post empresa finalizado correctamente (OK), id creado:", id)
}
