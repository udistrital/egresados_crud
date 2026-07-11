package controllers_test

import (
	"net/http"
	"testing"
)

func TestGetAllEgresado(t *testing.T) {
	response, err := http.Get(baseURL + "/egresado")
	if err != nil {
		t.Error("Error GetAll egresado:", err.Error())
		t.Fail()
		return
	}
	if response.StatusCode != 200 {
		t.Error("Error en GetAll egresado, se esperaba 200 y se obtuvo", response.StatusCode)
		t.Fail()
	} else {
		t.Log("GetAll egresado finalizado correctamente (OK)")
	}
}

func TestPostEgresado(t *testing.T) {
	usuarioId := crearUsuario(t, "EGR")
	id := crearEgresado(t, usuarioId)
	t.Log("Post egresado finalizado correctamente (OK), id creado:", id)
}
