package controllers_test

import (
	"net/http"
	"testing"
)

// Ruta renombrada de guion_bajo a guion (lineamientos-endpoints.md): /historial-solicitud.
func TestGetAllHistorialSolicitud(t *testing.T) {
	response, err := http.Get(baseURL + "/historial-solicitud")
	if err != nil {
		t.Error("Error GetAll historial-solicitud:", err.Error())
		t.Fail()
		return
	}
	if response.StatusCode != 200 {
		t.Error("Error en GetAll historial-solicitud, se esperaba 200 y se obtuvo", response.StatusCode)
		t.Fail()
	} else {
		t.Log("GetAll historial-solicitud finalizado correctamente (OK)")
	}
}

func TestPostHistorialSolicitud(t *testing.T) {
	solicitudId, usuarioId := crearSolicitudCompleta(t)
	body := map[string]interface{}{
		"solicitud_beneficio": map[string]interface{}{"id": solicitudId},
		"estado_nuevo_id":     1,
		"usuario":             map[string]interface{}{"id": usuarioId},
	}
	status, result := postJSON(t, "/historial-solicitud", body)
	if status != 201 {
		t.Error("Error en Post historial-solicitud, se esperaba 201 y se obtuvo", status, result)
		t.Fail()
		return
	}
	t.Log("Post historial-solicitud finalizado correctamente (OK), id creado:", idFrom(t, result))
}
