package controllers_test

import (
	"net/http"
	"testing"
)

// Ruta renombrada de guion_bajo a guion (lineamientos-endpoints.md): /documento-solicitud.
func TestGetAllDocumentoSolicitud(t *testing.T) {
	response, err := http.Get(baseURL + "/documento-solicitud")
	if err != nil {
		t.Error("Error GetAll documento-solicitud:", err.Error())
		t.Fail()
		return
	}
	if response.StatusCode != 200 {
		t.Error("Error en GetAll documento-solicitud, se esperaba 200 y se obtuvo", response.StatusCode)
		t.Fail()
	} else {
		t.Log("GetAll documento-solicitud finalizado correctamente (OK)")
	}
}

func TestPostDocumentoSolicitud(t *testing.T) {
	solicitudId, _ := crearSolicitudCompleta(t)
	usuarioEgr := crearUsuario(t, "EMP")
	empresaId := crearEmpresa(t)
	beneficioId := crearBeneficio(t, empresaId, usuarioEgr)
	documentoRequeridoId := crearDocumentoRequeridoBeneficio(t, beneficioId)

	body := map[string]interface{}{
		"solicitud_beneficio":      map[string]interface{}{"id": solicitudId},
		"documento_requerido":      map[string]interface{}{"id": documentoRequeridoId},
		"nombre_archivo":           "hoja_de_vida.pdf",
		"enlace_gestor_documental": "uid-test-" + uniqueSuffix(),
	}
	status, result := postJSON(t, "/documento-solicitud", body)
	if status != 201 {
		t.Error("Error en Post documento-solicitud, se esperaba 201 y se obtuvo", status, result)
		t.Fail()
		return
	}
	t.Log("Post documento-solicitud finalizado correctamente (OK), id creado:", idFrom(t, result))
}
