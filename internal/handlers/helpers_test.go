package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type sampleData struct {
	Mensaje string `json:"mensaje"`
}

func TestSendJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	SendJSON(rr, http.StatusOK, sampleData{Mensaje: "éxito"})

	if rr.Code != http.StatusOK {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Se esperaba Content-Type 'application/json', se obtuvo '%s'", contentType)
	}

	var res sampleData
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("Fallo al deserializar respuesta: %v", err)
	}

	if res.Mensaje != "éxito" {
		t.Errorf("Se esperaba 'éxito', se obtuvo '%s'", res.Mensaje)
	}
}

func TestSendError(t *testing.T) {
	rr := httptest.NewRecorder()
	SendError(rr, http.StatusBadRequest, "campo obligatorio faltante")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusBadRequest, rr.Code)
	}

	var res ResponseError
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("Fallo al deserializar error: %v", err)
	}

	if res.Error != "campo obligatorio faltante" {
		t.Errorf("Mensaje de error inesperado: '%s'", res.Error)
	}
}

func TestSendHTMXError(t *testing.T) {
	rr := httptest.NewRecorder()
	SendHTMXError(rr, http.StatusBadRequest, "Error <script>alert(1)</script>")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusBadRequest, rr.Code)
	}
	if got := rr.Header().Get("HX-Retarget"); got != "#htmx-error" {
		t.Errorf("Se esperaba HX-Retarget '#htmx-error', se obtuvo %q", got)
	}
	if got := rr.Header().Get("HX-Reswap"); got != "innerHTML" {
		t.Errorf("Se esperaba HX-Reswap 'innerHTML', se obtuvo %q", got)
	}
	body := rr.Body.String()
	if strings.Contains(body, "<script>") || !strings.Contains(body, "&lt;script&gt;") {
		t.Errorf("Se esperaba el mensaje escapado contra XSS, se obtuvo %q", body)
	}
}

func TestIsHTMXPartial(t *testing.T) {
	reqNormal := httptest.NewRequest(http.MethodGet, "/publicaciones", nil)
	if isHTMXPartial(reqNormal) {
		t.Error("una petición normal no debe tratarse como parcial HTMX")
	}

	reqPartial := httptest.NewRequest(http.MethodGet, "/publicaciones?q=cafe", nil)
	reqPartial.Header.Set("HX-Request", "true")
	if !isHTMXPartial(reqPartial) {
		t.Error("una petición con HX-Request: true sin HX-Boosted debe tratarse como parcial")
	}

	reqBoosted := httptest.NewRequest(http.MethodGet, "/publicaciones", nil)
	reqBoosted.Header.Set("HX-Request", "true")
	reqBoosted.Header.Set("HX-Boosted", "true")
	if isHTMXPartial(reqBoosted) {
		t.Error("una navegación con hx-boost debe devolver la página completa")
	}
}

func TestParseJSON(t *testing.T) {
	t.Run("JSON válido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"mensaje":"hola"}`))
		var data sampleData
		if err := ParseJSON(req, &data); err != nil {
			t.Fatalf("ParseJSON falló con JSON válido: %v", err)
		}
		if data.Mensaje != "hola" {
			t.Errorf("Se esperaba 'hola', se obtuvo '%s'", data.Mensaje)
		}
	})

	t.Run("Body vacío", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		var data sampleData
		if err := ParseJSON(req, &data); err == nil {
			t.Error("Se esperaba error con body vacío, pero fue nil")
		}
	})

	t.Run("JSON con campo no permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"mensaje":"hola","extra":"campo"}`))
		var data sampleData
		if err := ParseJSON(req, &data); err == nil {
			t.Error("Se esperaba error por campo no permitido, pero fue nil")
		}
	})
}
