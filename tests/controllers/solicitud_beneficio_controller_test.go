package controllers_test

import (
	"net/http"
	"testing"
)

// Ruta renombrada de guion_bajo a guion (lineamientos-endpoints.md): /solicitud-beneficio.
func TestGetAllSolicitudBeneficio(t *testing.T) {
	response, err := http.Get(baseURL + "/solicitud-beneficio")
	if err != nil {
		t.Error("Error GetAll solicitud-beneficio:", err.Error())
		t.Fail()
		return
	}
	if response.StatusCode != 200 {
		t.Error("Error en GetAll solicitud-beneficio, se esperaba 200 y se obtuvo", response.StatusCode)
		t.Fail()
	} else {
		t.Log("GetAll solicitud-beneficio finalizado correctamente (OK)")
	}
}

func TestPostSolicitudBeneficio(t *testing.T) {
	solicitudId, _ := crearSolicitudCompleta(t)
	t.Log("Post solicitud-beneficio finalizado correctamente (OK), id creado:", solicitudId)
}
