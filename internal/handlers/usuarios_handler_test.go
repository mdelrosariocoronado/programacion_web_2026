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
type MockUserService struct{}

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
}

func TestUserHandler_CreateFromForm(t *testing.T) {
	mockService := &MockUserService{}
	userHandler := NewUserHandler(mockService)

	t.Run("exito_303_SeeOther_y_redirect", func(t *testing.T) {
		formData := "nombre_completo=Carlos&email=carlos@example.com&clave=123456&rol=cliente&id_emprendimiento=5"
		req := httptest.NewRequest(http.MethodPost, "/usuarios", strings.NewReader(formData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		userHandler.CreateFromForm(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusSeeOther, rr.Code)
		}
		if loc := rr.Header().Get("Location"); loc != "/usuarios" {
			t.Errorf("Se esperaba redirección a '/usuarios', se obtuvo '%s'", loc)
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

