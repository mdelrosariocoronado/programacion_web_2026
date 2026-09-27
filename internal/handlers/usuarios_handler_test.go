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