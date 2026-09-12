package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserHandler_Create(t *testing.T) {
	// 1. Crear una instancia de UserHandler
	userHandler := NewUserHandler(nil)

	// 2. Simular la petición HTTP POST
	userJSON := `{"nombre":"Carlos","email":"carlos@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// 3. Asignar el método Create del struct instanciado
	handler := http.HandlerFunc(userHandler.Create)

	handler.ServeHTTP(rr, req)

	// 4. Validar que la respuesta sea 201 Created
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler devolvió código %v, se esperaba %v", status, http.StatusCreated)
	}
}