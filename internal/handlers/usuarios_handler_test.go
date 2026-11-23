package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	db "emprendimientos.com/servidor-go/db/sqlc"
)

//mock para testear funcionalidad, parcheando las funcionalidades
type MockUserService struct {
	listarFunc func(ctx context.Context) ([]db.Usuario, error)
}

func (m *MockUserService) CrearUsuario(ctx context.Context, nombre, email, clave, rol string, idEmprendimiento *int32) (db.Usuario, error) {
	return db.Usuario{
		IDUsuario:      1,
		NombreCompleto: nombre,
		Email:          email,
		Rol:            rol,
	}, nil
}

func (m *MockUserService) ObtenerUsuario(ctx context.Context, id int32) (db.Usuario, error) {
	return db.Usuario{}, nil
}

func (m *MockUserService) ObtenerUsuarioPorEmail(ctx context.Context, email string) (db.Usuario, error) {
	return db.Usuario{}, nil
}

func (m *MockUserService) ListarUsuarios(ctx context.Context) ([]db.Usuario, error) {
	if m.listarFunc != nil {
		return m.listarFunc(ctx)
	}
	return []db.Usuario{}, nil
}

func (m *MockUserService) ActualizarUsuario(ctx context.Context, id int32, nombre, email, clave, rol string, idEmprendimiento *int32) error {
	return nil
}

func (m *MockUserService) EliminarUsuario(ctx context.Context, id int32) error {
	return nil
}


func TestUserHandler_Create(t *testing.T) {
	//  Crear una instancia de UserHandler
	mockService := &MockUserService{}
	userHandler := NewUserHandler(mockService)

	// Simular la petición HTTP POST
	userJSON := `{"nombre_completo":"Carlos","email":"carlos@example.com","rol":"cliente"}`
	req := httptest.NewRequest(http.MethodPost, "/api/usuarios", strings.NewReader(userJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	//Ejecutar el Handler
	handler := http.HandlerFunc(userHandler.Create)
	handler.ServeHTTP(rr, req)

	

	//  Validar que la respuesta sea 201 Created
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler devolvió código %v, se esperaba %v. Respuesta: %s", 
			status, http.StatusCreated, rr.Body.String())
	}
}

func TestUserHandler_RenderPage(t *testing.T) {
	t.Run("carga_completa", func(t *testing.T) {
		mockService := &MockUserService{}
		userHandler := NewUserHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/usuarios", nil)
		rr := httptest.NewRecorder()

		userHandler.RenderPage(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("RenderPage devolvió código %v, se esperaba %v. Respuesta: %s",
				status, http.StatusOK, rr.Body.String())
		}

		body := rr.Body.String()
		if !strings.Contains(body, "Gestión de Usuarios") {
			t.Errorf("Se esperaba que la respuesta contenga 'Gestión de Usuarios', se obtuvo: %s", body)
		}
	})

	t.Run("busqueda_parcial_htmx", func(t *testing.T) {
		mockService := &MockUserService{
			listarFunc: func(ctx context.Context) ([]db.Usuario, error) {
				return []db.Usuario{
					{IDUsuario: 1, NombreCompleto: "Carlos Ruiz", Email: "carlos@tandil.com", Rol: "cliente"},
					{IDUsuario: 2, NombreCompleto: "Ana López", Email: "ana@tandil.com", Rol: "emprendedor"},
				}, nil
			},
		}
		userHandler := NewUserHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/usuarios?q=carlos", nil)
		req.Header.Set("HX-Request", "true")
		rr := httptest.NewRecorder()

		userHandler.RenderPage(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d", rr.Code)
		}
		body := rr.Body.String()
		if strings.Contains(body, "<!doctype html>") || strings.Contains(body, "Gestión de Usuarios") {
			t.Errorf("la búsqueda HTMX debe devolver solo el fragmento de lista, no la página completa")
		}
		if !strings.Contains(body, "Carlos Ruiz") || strings.Contains(body, "Ana López") {
			t.Errorf("se esperaba solo a Carlos Ruiz filtrado, se obtuvo: %s", body)
		}
	})

	t.Run("navegacion_boosted_devuelve_pagina_completa", func(t *testing.T) {
		mockService := &MockUserService{}
		userHandler := NewUserHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/usuarios", nil)
		req.Header.Set("HX-Request", "true")
		req.Header.Set("HX-Boosted", "true")
		rr := httptest.NewRecorder()

		userHandler.RenderPage(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "Gestión de Usuarios") {
			t.Errorf("con HX-Boosted debe renderizarse la página completa")
		}
	})
}

func TestUserHandler_CreateFromForm(t *testing.T) {
	mockService := &MockUserService{}
	userHandler := NewUserHandler(mockService)

	t.Run("exito_200_OK_fragmento_lista_y_form_OOB", func(t *testing.T) {
		formData := "nombre_completo=Carlos&email=carlos@example.com&clave=123456&rol=cliente&id_emprendimiento=5"
		req := httptest.NewRequest(http.MethodPost, "/usuarios", strings.NewReader(formData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		userHandler.CreateFromForm(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusOK, rr.Code)
		}
		body := rr.Body.String()
		if !strings.Contains(body, `id="lista-usuarios"`) {
			t.Errorf("Se esperaba el fragmento de lista de usuarios, se obtuvo: %s", body)
		}
		if !strings.Contains(body, `hx-swap-oob="true"`) {
			t.Errorf("Se esperaba el formulario OOB para limpiar los inputs, se obtuvo: %s", body)
		}
	})

	t.Run("falla_400_con_id_emprendimiento_invalido", func(t *testing.T) {
		formData := "nombre_completo=Carlos&email=carlos@example.com&clave=123456&rol=cliente&id_emprendimiento=abc"
		req := httptest.NewRequest(http.MethodPost, "/usuarios", strings.NewReader(formData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		userHandler.CreateFromForm(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Se esperaba código %d ante id_emprendimiento inválido, se obtuvo %d", http.StatusBadRequest, rr.Code)
		}
	})
}

func TestUserHandler_DeleteFromHTMX(t *testing.T) {
	mockService := &MockUserService{}
	userHandler := NewUserHandler(mockService)

	t.Run("exito_200_OK_vacio", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/usuarios/1", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()

		userHandler.DeleteFromHTMX(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusOK, rr.Code)
		}
		if rr.Body.Len() != 0 {
			t.Errorf("Se esperaba respuesta vacía para HTMX outerHTML swap, se obtuvo: %q", rr.Body.String())
		}
	})

	t.Run("falla_400_con_id_invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/usuarios/abc", nil)
		req.SetPathValue("id", "abc")
		rr := httptest.NewRecorder()

		userHandler.DeleteFromHTMX(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Se esperaba código %d ante ID inválido, se obtuvo %d", http.StatusBadRequest, rr.Code)
		}
	})
}

