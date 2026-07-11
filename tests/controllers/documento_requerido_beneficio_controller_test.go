package controllers_test

import (
	"net/http"
	"testing"
)

// Ruta renombrada de guion_bajo a guion (lineamientos-endpoints.md): /documento-requerido-beneficio.
func TestGetAllDocumentoRequeridoBeneficio(t *testing.T) {
	response, err := http.Get(baseURL + "/documento-requerido-beneficio")
	if err != nil {
		t.Error("Error GetAll documento-requerido-beneficio:", err.Error())
		t.Fail()
		return
	}
	if response.StatusCode != 200 {
		t.Error("Error en GetAll documento-requerido-beneficio, se esperaba 200 y se obtuvo", response.StatusCode)
		t.Fail()
	} else {
		t.Log("GetAll documento-requerido-beneficio finalizado correctamente (OK)")
	}
}

func TestPostDocumentoRequeridoBeneficio(t *testing.T) {
	empresaId := crearEmpresa(t)
	usuarioCreadorId := crearUsuario(t, "EMP")
	beneficioId := crearBeneficio(t, empresaId, usuarioCreadorId)
	id := crearDocumentoRequeridoBeneficio(t, beneficioId)
	t.Log("Post documento-requerido-beneficio finalizado correctamente (OK), id creado:", id)
}
