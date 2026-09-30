package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	db "emprendimientos.com/servidor-go/db/sqlc"
	"emprendimientos.com/servidor-go/internal/services"
)

type mockPublicacionService struct {
	createFunc func(ctx context.Context, input services.CreatePublicacionInput) (db.Publicacione, error)
	getByIDFunc func(ctx context.Context, id int32) (db.Publicacione, error)
	listByEmprendimientoFunc func(ctx context.Context, idEmprendimiento int32) ([]db.Publicacione, error)
	updateFunc func(ctx context.Context, id int32, input services.UpdatePublicacionInput) error
	deleteFunc func(ctx context.Context, id int32) error
}

func (m *mockPublicacionService) Create(ctx context.Context, input services.CreatePublicacionInput) (db.Publicacione, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return db.Publicacione{IDPublicacion: 1, Titulo: input.Titulo}, nil
}

func (m *mockPublicacionService) GetByID(ctx context.Context, id int32) (db.Publicacione, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return db.Publicacione{IDPublicacion: id, Titulo: "Publicacion Mock"}, nil
}

func (m *mockPublicacionService) ListByEmprendimiento(ctx context.Context, idEmprendimiento int32) ([]db.Publicacione, error) {
	if m.listByEmprendimientoFunc != nil {
		return m.listByEmprendimientoFunc(ctx, idEmprendimiento)
	}
	return []db.Publicacione{{IDPublicacion: 1, Titulo: "Pub 1"}}, nil
}

func (m *mockPublicacionService) Update(ctx context.Context, id int32, input services.UpdatePublicacionInput) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil
}

func (m *mockPublicacionService) Delete(ctx context.Context, id int32) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func TestPublicacionHandler_Create(t *testing.T) {
	t.Run("exito 201 Created", func(t *testing.T) {
		handler := NewPublicacionHandler(&mockPublicacionService{})

		body := `{"id_emprendimiento":1,"titulo":"Oferta Nueva","tipo":"promocion"}`
		req := httptest.NewRequest(http.MethodPost, "/api/publicaciones", strings.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("se esperaba status %d, se obtuvo %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("falla 400 Bad Request con error de validacion", func(t *testing.T) {
		mockSvc := &mockPublicacionService{
			createFunc: func(ctx context.Context, input services.CreatePublicacionInput) (db.Publicacione, error) {
				return db.Publicacione{}, services.ErrTituloRequerido
			},
		}
		handler := NewPublicacionHandler(mockSvc)

		body := `{"id_emprendimiento":1,"titulo":"","tipo":"promocion"}`
		req := httptest.NewRequest(http.MethodPost, "/api/publicaciones", strings.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("se esperaba status %d, se obtuvo %d", http.StatusBadRequest, rr.Code)
		}
	})
}

func TestPublicacionHandler_GetByID(t *testing.T) {
	t.Run("exito 200 OK", func(t *testing.T) {
		handler := NewPublicacionHandler(&mockPublicacionService{})

		req := httptest.NewRequest(http.MethodGet, "/api/publicaciones/5", nil)
		req.SetPathValue("id", "5")
		rr := httptest.NewRecorder()

		handler.GetByID(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("se esperaba status %d, se obtuvo %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("falla 404 Not Found si no existe", func(t *testing.T) {
		mockSvc := &mockPublicacionService{
			getByIDFunc: func(ctx context.Context, id int32) (db.Publicacione, error) {
				return db.Publicacione{}, services.ErrPublicacionNoEncontrada
			},
		}
		handler := NewPublicacionHandler(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/api/publicaciones/99", nil)
		req.SetPathValue("id", "99")
		rr := httptest.NewRecorder()

		handler.GetByID(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("se esperaba status %d, se obtuvo %d", http.StatusNotFound, rr.Code)
		}
	})
}

func TestPublicacionHandler_ListByEmprendimiento(t *testing.T) {
	t.Run("exito 200 OK", func(t *testing.T) {
		handler := NewPublicacionHandler(&mockPublicacionService{})

		req := httptest.NewRequest(http.MethodGet, "/api/emprendimientos/2/publicaciones", nil)
		req.SetPathValue("id", "2")
		rr := httptest.NewRecorder()

		handler.ListByEmprendimiento(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("se esperaba status %d, se obtuvo %d", http.StatusOK, rr.Code)
		}
	})
}

func TestPublicacionHandler_Update(t *testing.T) {
	t.Run("exito 200 OK", func(t *testing.T) {
		handler := NewPublicacionHandler(&mockPublicacionService{})

		body := `{"titulo":"Titulo Actualizado","tipo":"servicio"}`
		req := httptest.NewRequest(http.MethodPut, "/api/publicaciones/1", strings.NewReader(body))
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()

		handler.Update(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("se esperaba status %d, se obtuvo %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("falla 404 Not Found si no existe", func(t *testing.T) {
		mockSvc := &mockPublicacionService{
			updateFunc: func(ctx context.Context, id int32, input services.UpdatePublicacionInput) error {
				return services.ErrPublicacionNoEncontrada
			},
		}
		handler := NewPublicacionHandler(mockSvc)

		body := `{"titulo":"Titulo","tipo":"servicio"}`
		req := httptest.NewRequest(http.MethodPut, "/api/publicaciones/404", strings.NewReader(body))
		req.SetPathValue("id", "404")
		rr := httptest.NewRecorder()

		handler.Update(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("se esperaba status %d, se obtuvo %d", http.StatusNotFound, rr.Code)
		}
	})
}

func TestPublicacionHandler_Delete(t *testing.T) {
	t.Run("exito 204 No Content", func(t *testing.T) {
		handler := NewPublicacionHandler(&mockPublicacionService{})

		req := httptest.NewRequest(http.MethodDelete, "/api/publicaciones/1", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()

		handler.Delete(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("se esperaba status %d, se obtuvo %d", http.StatusNoContent, rr.Code)
		}
	})

	t.Run("falla 404 Not Found si no existe", func(t *testing.T) {
		mockSvc := &mockPublicacionService{
			deleteFunc: func(ctx context.Context, id int32) error {
				return services.ErrPublicacionNoEncontrada
			},
		}
		handler := NewPublicacionHandler(mockSvc)

		req := httptest.NewRequest(http.MethodDelete, "/api/publicaciones/88", nil)
		req.SetPathValue("id", "88")
		rr := httptest.NewRecorder()

		handler.Delete(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("se esperaba status %d, se obtuvo %d", http.StatusNotFound, rr.Code)
		}
	})
}
