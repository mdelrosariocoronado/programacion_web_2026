package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"emprendimientos.com/servidor-go/internal/services"
)

func TestEmprendimientoHandler_Validaciones(t *testing.T) {
	service := services.NewEmprendimientoService(nil)
	handler := NewEmprendimientoHandler(service)

	t.Run("Create falla con body inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/emprendimientos", strings.NewReader(`{}`))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Se esperaba código %d ante validación fallida, se obtuvo %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("GetByID falla con ID no numérico", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/emprendimientos/abc", nil)
		// Simular valor extraído por ServeMux de Go 1.22
		req.SetPathValue("id", "abc")
		rr := httptest.NewRecorder()

		handler.GetByID(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Se esperaba código %d con ID no numérico, se obtuvo %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Delete falla con ID menor o igual a cero", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/emprendimientos/0", nil)
		req.SetPathValue("id", "0")
		rr := httptest.NewRecorder()

		handler.Delete(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Se esperaba código %d con ID 0, se obtuvo %d", http.StatusBadRequest, rr.Code)
		}
	})
}
