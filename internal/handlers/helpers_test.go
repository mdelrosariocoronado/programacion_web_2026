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
