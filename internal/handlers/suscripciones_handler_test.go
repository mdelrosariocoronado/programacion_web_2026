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

type mockSuscripcionService struct {
	suscribeFunc      func(ctx context.Context, idUsuario, idEmprendimiento int32) (db.Suscripcione, error)
	listByUsuarioFunc func(ctx context.Context, idUsuario int32) ([]db.Suscripcione, error)
	cancelFunc        func(ctx context.Context, idUsuario, idEmprendimiento int32) error
}

func (m *mockSuscripcionService) Suscribe(ctx context.Context, idUsuario, idEmprendimiento int32) (db.Suscripcione, error) {
	if m.suscribeFunc != nil {
		return m.suscribeFunc(ctx, idUsuario, idEmprendimiento)
	}
	return db.Suscripcione{IDSuscripcion: 1, IDUsuario: idUsuario, IDEmprendimiento: idEmprendimiento}, nil
}

func (m *mockSuscripcionService) ListByUsuario(ctx context.Context, idUsuario int32) ([]db.Suscripcione, error) {
	if m.listByUsuarioFunc != nil {
		return m.listByUsuarioFunc(ctx, idUsuario)
	}
	return []db.Suscripcione{}, nil
}

func (m *mockSuscripcionService) Cancel(ctx context.Context, idUsuario, idEmprendimiento int32) error {
	if m.cancelFunc != nil {
		return m.cancelFunc(ctx, idUsuario, idEmprendimiento)
	}
	return nil
}

type mockEmprendimientoLister struct {
	listFunc func(ctx context.Context) ([]services.EmprendimientoResponse, error)
}

func (m *mockEmprendimientoLister) List(ctx context.Context) ([]services.EmprendimientoResponse, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return []services.EmprendimientoResponse{
		{IDEmprendimiento: 10, Nombre: "Café Tandil", Rubro: "Gastronomía"},
	}, nil
}

func TestSuscripcionHandler_RenderPage(t *testing.T) {
	users := []db.Usuario{
		{IDUsuario: 1, NombreCompleto: "Carlos", Email: "carlos@tandil.com"},
		{IDUsuario: 2, NombreCompleto: "Lucía", Email: "lucia@tandil.com"},
	}

	t.Run("carga_completa_default_primer_usuario", func(t *testing.T) {
		var queriedUserID int32
		handler := NewSuscripcionHandler(
			&mockSuscripcionService{
				listByUsuarioFunc: func(ctx context.Context, idUsuario int32) ([]db.Suscripcione, error) {
					queriedUserID = idUsuario
					return []db.Suscripcione{{IDUsuario: idUsuario, IDEmprendimiento: 10}}, nil
				},
			},
			&MockUserService{
				listarFunc: func(ctx context.Context) ([]db.Usuario, error) {
					return users, nil
				},
			},
			&mockEmprendimientoLister{},
		)

		req := httptest.NewRequest(http.MethodGet, "/suscripciones", nil)
		rr := httptest.NewRecorder()
		handler.RenderPage(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d", rr.Code)
		}
		if queriedUserID != 1 {
			t.Errorf("se esperaba usuario seleccionado por defecto 1, se obtuvo %d", queriedUserID)
		}
		body := rr.Body.String()
		if !strings.Contains(body, "Gestión de Suscripciones") || !strings.Contains(body, "Café Tandil") {
			t.Errorf("contenido inesperado en página completa: %s", body)
		}
	})

	t.Run("cambio_usuario_parcial_htmx", func(t *testing.T) {
		handler := NewSuscripcionHandler(
			&mockSuscripcionService{},
			&MockUserService{
				listarFunc: func(ctx context.Context) ([]db.Usuario, error) {
					return users, nil
				},
			},
			&mockEmprendimientoLister{},
		)

		req := httptest.NewRequest(http.MethodGet, "/suscripciones?id_usuario=2", nil)
		req.Header.Set("HX-Request", "true")
		rr := httptest.NewRecorder()
		handler.RenderPage(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d", rr.Code)
		}
		body := rr.Body.String()
		if strings.Contains(body, "Gestión de Suscripciones") {
			t.Errorf("el cambio parcial HTMX debe devolver solo #panel-suscripciones")
		}
		if !strings.Contains(body, `id="panel-suscripciones"`) || !strings.Contains(body, `hx-post="/suscripciones/2/10"`) {
			t.Errorf("se esperaba panel actualizado para usuario 2, se obtuvo: %s", body)
		}
	})

	t.Run("navegacion_boosted_devuelve_pagina_completa", func(t *testing.T) {
		handler := NewSuscripcionHandler(
			&mockSuscripcionService{},
			&MockUserService{
				listarFunc: func(ctx context.Context) ([]db.Usuario, error) {
					return users, nil
				},
			},
			&mockEmprendimientoLister{},
		)

		req := httptest.NewRequest(http.MethodGet, "/suscripciones", nil)
		req.Header.Set("HX-Request", "true")
		req.Header.Set("HX-Boosted", "true")
		rr := httptest.NewRecorder()
		handler.RenderPage(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "Gestión de Suscripciones") {
			t.Errorf("con HX-Boosted debe renderizarse la página completa")
		}
	})
}

func TestSuscripcionHandler_SubscribeFromHTMX(t *testing.T) {
	t.Run("exito_200_devuelve_panel_y_banner_oob", func(t *testing.T) {
		handler := NewSuscripcionHandler(
			&mockSuscripcionService{
				listByUsuarioFunc: func(ctx context.Context, idUsuario int32) ([]db.Suscripcione, error) {
					return []db.Suscripcione{{IDUsuario: idUsuario, IDEmprendimiento: 10}}, nil
				},
			},
			&MockUserService{},
			&mockEmprendimientoLister{},
		)

		req := httptest.NewRequest(http.MethodPost, "/suscripciones/1/10", nil)
		req.SetPathValue("id_usuario", "1")
		req.SetPathValue("id_emprendimiento", "10")
		rr := httptest.NewRecorder()
		handler.SubscribeFromHTMX(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d", rr.Code)
		}
		body := rr.Body.String()
		if !strings.Contains(body, `id="panel-suscripciones"`) || !strings.Contains(body, `hx-swap-oob="true"`) {
			t.Errorf("se esperaba panel y banner OOB, se obtuvo: %s", body)
		}
	})

	t.Run("conflicto_409_retarget_a_htmx_error", func(t *testing.T) {
		handler := NewSuscripcionHandler(
			&mockSuscripcionService{
				suscribeFunc: func(ctx context.Context, idUsuario, idEmprendimiento int32) (db.Suscripcione, error) {
					return db.Suscripcione{}, services.ErrYaSuscrito
				},
			},
			&MockUserService{},
			&mockEmprendimientoLister{},
		)

		req := httptest.NewRequest(http.MethodPost, "/suscripciones/1/10", nil)
		req.SetPathValue("id_usuario", "1")
		req.SetPathValue("id_emprendimiento", "10")
		rr := httptest.NewRecorder()
		handler.SubscribeFromHTMX(rr, req)

		if rr.Code != http.StatusConflict {
			t.Errorf("se esperaba 409, se obtuvo %d", rr.Code)
		}
		if got := rr.Header().Get("HX-Retarget"); got != "#htmx-error" {
			t.Errorf("se esperaba HX-Retarget '#htmx-error', se obtuvo %q", got)
		}
	})

	t.Run("falla_400_con_id_invalido", func(t *testing.T) {
		handler := NewSuscripcionHandler(&mockSuscripcionService{}, &MockUserService{}, &mockEmprendimientoLister{})

		req := httptest.NewRequest(http.MethodPost, "/suscripciones/abc/10", nil)
		req.SetPathValue("id_usuario", "abc")
		req.SetPathValue("id_emprendimiento", "10")
		rr := httptest.NewRecorder()
		handler.SubscribeFromHTMX(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("se esperaba 400, se obtuvo %d", rr.Code)
		}
	})
}

func TestSuscripcionHandler_CancelFromHTMX(t *testing.T) {
	t.Run("exito_200_devuelve_panel_y_banner_oob", func(t *testing.T) {
		handler := NewSuscripcionHandler(
			&mockSuscripcionService{},
			&MockUserService{},
			&mockEmprendimientoLister{},
		)

		req := httptest.NewRequest(http.MethodDelete, "/suscripciones/1/10", nil)
		req.SetPathValue("id_usuario", "1")
		req.SetPathValue("id_emprendimiento", "10")
		rr := httptest.NewRecorder()
		handler.CancelFromHTMX(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d", rr.Code)
		}
		body := rr.Body.String()
		if !strings.Contains(body, `id="panel-suscripciones"`) || !strings.Contains(body, `hx-swap-oob="true"`) {
			t.Errorf("se esperaba panel y banner OOB, se obtuvo: %s", body)
		}
	})

	t.Run("no_encontrada_404_retarget_a_htmx_error", func(t *testing.T) {
		handler := NewSuscripcionHandler(
			&mockSuscripcionService{
				cancelFunc: func(ctx context.Context, idUsuario, idEmprendimiento int32) error {
					return services.ErrSuscripcionNoEncontrada
				},
			},
			&MockUserService{},
			&mockEmprendimientoLister{},
		)

		req := httptest.NewRequest(http.MethodDelete, "/suscripciones/1/10", nil)
		req.SetPathValue("id_usuario", "1")
		req.SetPathValue("id_emprendimiento", "10")
		rr := httptest.NewRecorder()
		handler.CancelFromHTMX(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("se esperaba 404, se obtuvo %d", rr.Code)
		}
		if got := rr.Header().Get("HX-Retarget"); got != "#htmx-error" {
			t.Errorf("se esperaba HX-Retarget '#htmx-error', se obtuvo %q", got)
		}
	})
}
