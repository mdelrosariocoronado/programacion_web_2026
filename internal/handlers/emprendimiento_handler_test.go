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
		// Simular valor extraído por ServeMux
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

	t.Run("CreateFromForm falla con datos inválidos", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/emprendimientos", strings.NewReader("nombre=&rubro="))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler.CreateFromForm(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Se esperaba código %d con nombre vacío, se obtuvo %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("DeleteFromHTMX falla con ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/emprendimientos/abc", nil)
		req.SetPathValue("id", "abc")
		rr := httptest.NewRecorder()

		handler.DeleteFromHTMX(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Se esperaba código %d con ID inválido, se obtuvo %d", http.StatusBadRequest, rr.Code)
		}
	})
}

func TestFiltrarEmprendimientos(t *testing.T) {
	items := []services.EmprendimientoResponse{
		{IDEmprendimiento: 1, Nombre: "Cafetería Sierra", Rubro: "Gastronomía", Descripcion: "Café de especialidad"},
		{IDEmprendimiento: 2, Nombre: "Taller Norte", Rubro: "Artesanías", Contacto: "taller@tandil.com"},
	}

	if got := filtrarEmprendimientos(items, ""); len(got) != 2 {
		t.Fatalf("sin filtro se esperaban 2 elementos, se obtuvieron %d", len(got))
	}
	if got := filtrarEmprendimientos(items, "ESPECIALIDAD"); len(got) != 1 || got[0].IDEmprendimiento != 1 {
		t.Errorf("filtro por descripción inesperado: %+v", got)
	}
	if got := filtrarEmprendimientos(items, "artesan"); len(got) != 1 || got[0].IDEmprendimiento != 2 {
		t.Errorf("filtro por rubro inesperado: %+v", got)
	}
}
