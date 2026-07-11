package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

// Familia B: pruebas de integración contra el servidor real (bee run / docker run),
// sin mocks. baseURL asume el CRUD corriendo en el puerto por defecto de dev (8080).
const baseURL = "http://localhost:8080/v1"

var seqCounter int64

// uniqueSuffix es corto a propósito: nit/codigo_institucional son VARCHAR(20) y deben
// caber junto con su prefijo ("NIT"/"COD"). El contador atómico evita colisiones entre
// llamadas consecutivas dentro de la misma resolución de reloj de Windows (~15ms).
func uniqueSuffix() string {
	n := atomic.AddInt64(&seqCounter, 1)
	return fmt.Sprintf("%d%03d", time.Now().Unix()%100000, n%1000)
}

func postJSON(t *testing.T, path string, body interface{}) (int, map[string]interface{}) {
	t.Helper()
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("error al convertir el cuerpo a JSON: %v", err)
	}
	resp, err := http.Post(baseURL+path, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("error en POST %s: %v", path, err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error al leer respuesta de POST %s: %v", path, err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		// La respuesta no es un objeto (p. ej. un mensaje de error como string plano).
		result = map[string]interface{}{"raw": string(respBytes)}
	}
	return resp.StatusCode, result
}

func idFrom(t *testing.T, result map[string]interface{}) int {
	t.Helper()
	idFloat, ok := result["id"].(float64)
	if !ok {
		t.Fatalf("respuesta sin id: %v", result)
	}
	return int(idFloat)
}

// crearUsuario crea un usuario prerequisito (EGR/EMP/ADM) vía POST real y devuelve su id.
func crearUsuario(t *testing.T, tipoUsuario string) int {
	t.Helper()
	suf := uniqueSuffix()
	body := map[string]interface{}{
		"nombre":         "Test " + tipoUsuario + " " + suf,
		"correo":         "test." + suf + "@example.com",
		"tipo_usuario":   tipoUsuario,
		"sistema_origen": "LOCAL",
		// uq_usuario_id_externo es UNIQUE(sistema_origen, id_externo); con id_externo
		// vacío (Go zero-value, no NULL) solo un usuario LOCAL podria existir. Se fija
		// un valor único por usuario de prueba para no colisionar entre pruebas.
		"id_externo": "TEST-" + suf,
	}
	status, result := postJSON(t, "/usuario", body)
	if status != 201 {
		t.Fatalf("no se pudo crear usuario prerequisito: status %d, body %v", status, result)
	}
	return idFrom(t, result)
}

// crearEgresado crea un egresado (subtipo EGR) ligado a usuarioId.
func crearEgresado(t *testing.T, usuarioId int) int {
	t.Helper()
	suf := uniqueSuffix()
	body := map[string]interface{}{
		"usuario":              map[string]interface{}{"id": usuarioId},
		"codigo_institucional": "COD" + suf,
	}
	status, result := postJSON(t, "/egresado", body)
	if status != 201 {
		t.Fatalf("no se pudo crear egresado prerequisito: status %d, body %v", status, result)
	}
	return idFrom(t, result)
}

// crearEmpresa crea una empresa prerequisito. estado_empresa_id es una referencia lógica
// (sin FK, C-6), cualquier entero sirve para efectos de esta prueba de integración.
func crearEmpresa(t *testing.T) int {
	t.Helper()
	suf := uniqueSuffix()
	body := map[string]interface{}{
		"nit":               "NIT" + suf,
		"razon_social":      "Empresa Test " + suf,
		"estado_empresa_id": 1,
	}
	status, result := postJSON(t, "/empresa", body)
	if status != 201 {
		t.Fatalf("no se pudo crear empresa prerequisito: status %d, body %v", status, result)
	}
	return idFrom(t, result)
}

// crearBeneficio crea un beneficio publicado por empresaId y creado por usuarioCreadorId.
func crearBeneficio(t *testing.T, empresaId, usuarioCreadorId int) int {
	t.Helper()
	suf := uniqueSuffix()
	body := map[string]interface{}{
		"empresa":                map[string]interface{}{"id": empresaId},
		"categoria_beneficio_id": 1,
		"estado_beneficio_id":    1,
		"titulo":                 "Beneficio Test " + suf,
		"descripcion":            "Descripcion de prueba",
		"condiciones":            "Condiciones de prueba",
		"fecha_inicio":           "2026-01-01T00:00:00Z",
		"fecha_fin":              "2026-12-31T00:00:00Z",
		"cupos_total":            10,
		"cupos_disponibles":      10,
		"usuario_creador":        map[string]interface{}{"id": usuarioCreadorId},
	}
	status, result := postJSON(t, "/beneficio", body)
	if status != 201 {
		t.Fatalf("no se pudo crear beneficio prerequisito: status %d, body %v", status, result)
	}
	return idFrom(t, result)
}

// crearSolicitudBeneficio crea una solicitud del egresadoId sobre beneficioId.
func crearSolicitudBeneficio(t *testing.T, egresadoId, beneficioId int) int {
	t.Helper()
	body := map[string]interface{}{
		"egresado":  map[string]interface{}{"id": egresadoId},
		"beneficio": map[string]interface{}{"id": beneficioId},
	}
	status, result := postJSON(t, "/solicitud-beneficio", body)
	if status != 201 {
		t.Fatalf("no se pudo crear solicitud_beneficio prerequisito: status %d, body %v", status, result)
	}
	return idFrom(t, result)
}

// crearDocumentoRequeridoBeneficio crea un documento requerido para beneficioId.
func crearDocumentoRequeridoBeneficio(t *testing.T, beneficioId int) int {
	t.Helper()
	body := map[string]interface{}{
		"beneficio":   map[string]interface{}{"id": beneficioId},
		"nombre":      "Hoja de vida",
		"descripcion": "Documento requerido de prueba",
	}
	status, result := postJSON(t, "/documento-requerido-beneficio", body)
	if status != 201 {
		t.Fatalf("no se pudo crear documento_requerido_beneficio prerequisito: status %d, body %v", status, result)
	}
	return idFrom(t, result)
}

// crearSolicitudCompleta arma la cadena completa egresado→empresa→beneficio→solicitud y
// devuelve (solicitudId, usuarioCreadorId) para las pruebas que solo necesitan una solicitud
// válida sin repetir el armado en cada archivo.
func crearSolicitudCompleta(t *testing.T) (solicitudId int, usuarioCreadorId int) {
	t.Helper()
	usuarioEgr := crearUsuario(t, "EGR")
	egresadoId := crearEgresado(t, usuarioEgr)
	empresaId := crearEmpresa(t)
	usuarioCreadorId = crearUsuario(t, "EMP")
	beneficioId := crearBeneficio(t, empresaId, usuarioCreadorId)
	solicitudId = crearSolicitudBeneficio(t, egresadoId, beneficioId)
	return solicitudId, usuarioCreadorId
}
