package controllers_test

import (
	"net/http"
	"testing"
)

func TestGetAllUsuario(t *testing.T) {
	response, err := http.Get(baseURL + "/usuario")
	if err != nil {
		t.Error("Error GetAll usuario:", err.Error())
		t.Fail()
		return
	}
	if response.StatusCode != 200 {
		t.Error("Error en GetAll usuario, se esperaba 200 y se obtuvo", response.StatusCode)
		t.Fail()
	} else {
		t.Log("GetAll usuario finalizado correctamente (OK)")
	}
}

func TestPostUsuario(t *testing.T) {
	id := crearUsuario(t, "EGR")
	t.Log("Post usuario finalizado correctamente (OK), id creado:", id)
}
