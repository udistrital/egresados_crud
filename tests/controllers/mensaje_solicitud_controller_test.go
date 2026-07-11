package controllers_test

import (
	"net/http"
	"testing"
)

// Ruta renombrada de guion_bajo a guion (lineamientos-endpoints.md): /mensaje-solicitud.
func TestGetAllMensajeSolicitud(t *testing.T) {
	response, err := http.Get(baseURL + "/mensaje-solicitud")
	if err != nil {
		t.Error("Error GetAll mensaje-solicitud:", err.Error())
		t.Fail()
		return
	}
	if response.StatusCode != 200 {
		t.Error("Error en GetAll mensaje-solicitud, se esperaba 200 y se obtuvo", response.StatusCode)
		t.Fail()
	} else {
		t.Log("GetAll mensaje-solicitud finalizado correctamente (OK)")
	}
}

func TestPostMensajeSolicitud(t *testing.T) {
	solicitudId, usuarioId := crearSolicitudCompleta(t)
	body := map[string]interface{}{
		"solicitud_beneficio": map[string]interface{}{"id": solicitudId},
		"usuario":             map[string]interface{}{"id": usuarioId},
		"mensaje":             "Mensaje de prueba de integracion",
	}
	status, result := postJSON(t, "/mensaje-solicitud", body)
	if status != 201 {
		t.Error("Error en Post mensaje-solicitud, se esperaba 201 y se obtuvo", status, result)
		t.Fail()
		return
	}
	t.Log("Post mensaje-solicitud finalizado correctamente (OK), id creado:", idFrom(t, result))
}
