package controllers_test

import (
	"net/http"
	"testing"
)

func TestGetAllBeneficio(t *testing.T) {
	response, err := http.Get(baseURL + "/beneficio")
	if err != nil {
		t.Error("Error GetAll beneficio:", err.Error())
		t.Fail()
		return
	}
	if response.StatusCode != 200 {
		t.Error("Error en GetAll beneficio, se esperaba 200 y se obtuvo", response.StatusCode)
		t.Fail()
	} else {
		t.Log("GetAll beneficio finalizado correctamente (OK)")
	}
}

func TestPostBeneficio(t *testing.T) {
	empresaId := crearEmpresa(t)
	usuarioCreadorId := crearUsuario(t, "EMP")
	id := crearBeneficio(t, empresaId, usuarioCreadorId)
	t.Log("Post beneficio finalizado correctamente (OK), id creado:", id)
}
