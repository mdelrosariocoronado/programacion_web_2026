package views

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"emprendimientos.com/servidor-go/db/sqlc"
	"emprendimientos.com/servidor-go/internal/services"
)

func TestSuscripcionesPage_SinUsuarios(t *testing.T) {
	var buf bytes.Buffer
	err := SuscripcionesPage(
		[]sqlc.Usuario{},
		0,
		[]services.EmprendimientoResponse{},
		map[int32]sqlc.Suscripcione{},
	).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering SuscripcionesPage sin usuarios: %v", err)
	}

	html := buf.String()
	expectedSubstrings := []string{
		"Gestión de Suscripciones",
		"No hay usuarios registrados.",
		`href="/usuarios"`,
		"Seleccioná o registrá un usuario para gestionar sus suscripciones.",
	}
	for _, sub := range expectedSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}

func TestSuscripcionesPage_SelectorHTMX(t *testing.T) {
	usuarios := []sqlc.Usuario{
		{IDUsuario: 1, NombreCompleto: "Juan Perez", Email: "juan@tandil.com"},
		{IDUsuario: 2, NombreCompleto: "Ana Gomez", Email: "ana@tandil.com"},
	}

	var buf bytes.Buffer
	err := SuscripcionesPage(
		usuarios,
		2,
		[]services.EmprendimientoResponse{},
		map[int32]sqlc.Suscripcione{},
	).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering SuscripcionesPage con selector: %v", err)
	}

	html := buf.String()
	expectedSubstrings := []string{
		`<select name="id_usuario"`,
		`hx-get="/suscripciones"`,
		`hx-target="#panel-suscripciones"`,
		`value="2" selected`,
		"Ana Gomez (ana@tandil.com)",
	}
	for _, sub := range expectedSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}

func TestSuscripcionesPanel_ConSuscripciones(t *testing.T) {
	emps := []services.EmprendimientoResponse{
		{IDEmprendimiento: 10, Nombre: "Café Centro", Rubro: "Gastronomía"},
		{IDEmprendimiento: 20, Nombre: "Taller Serrano", Rubro: "Artesanías"},
	}
	subs := map[int32]sqlc.Suscripcione{
		10: {
			IDSuscripcion:    5,
			IDUsuario:        1,
			IDEmprendimiento: 10,
			FechaSub: sql.NullTime{
				Time:  time.Date(2026, 11, 20, 12, 0, 0, 0, time.UTC),
				Valid: true,
			},
		},
	}

	var buf bytes.Buffer
	err := SuscripcionesPanel(1, emps, subs).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering SuscripcionesPanel: %v", err)
	}

	html := buf.String()
	expectedSubstrings := []string{
		`id="panel-suscripciones"`,
		"Suscripciones activas: 1 de 2 emprendimientos",
		"Café Centro",
		"Suscrito (20/11/2026)",
		`hx-delete="/suscripciones/1/10"`,
		`hx-confirm="¿Estás seguro de que querés cancelar esta suscripción?"`,
		"Taller Serrano",
		"No suscrito",
		`hx-post="/suscripciones/1/20"`,
		`href="/publicaciones?id_emprendimiento=10"`,
	}
	for _, sub := range expectedSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}
