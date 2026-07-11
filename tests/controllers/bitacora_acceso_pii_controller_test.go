package controllers_test

import (
	"net/http"
	"testing"
)

// Ruta renombrada de guion_bajo a guion (lineamientos-endpoints.md): /bitacora-acceso-pii.
func TestGetAllBitacoraAccesoPii(t *testing.T) {
	response, err := http.Get(baseURL + "/bitacora-acceso-pii")
	if err != nil {
		t.Error("Error GetAll bitacora-acceso-pii:", err.Error())
		t.Fail()
		return
	}
	if response.StatusCode != 200 {
		t.Error("Error en GetAll bitacora-acceso-pii, se esperaba 200 y se obtuvo", response.StatusCode)
		t.Fail()
	} else {
		t.Log("GetAll bitacora-acceso-pii finalizado correctamente (OK)")
	}
}

func TestPostBitacoraAccesoPii(t *testing.T) {
	usuarioId := crearUsuario(t, "ADM")
	body := map[string]interface{}{
		"usuario":      map[string]interface{}{"id": usuarioId},
		"recurso_tipo": "solicitud_beneficio",
		"accion":       "CONSULTA",
	}
	status, result := postJSON(t, "/bitacora-acceso-pii", body)
	if status != 201 {
		t.Error("Error en Post bitacora-acceso-pii, se esperaba 201 y se obtuvo", status, result)
		t.Fail()
		return
	}
	t.Log("Post bitacora-acceso-pii finalizado correctamente (OK), id creado:", idFrom(t, result))
}
