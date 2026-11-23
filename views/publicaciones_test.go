package views

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
	"testing"

	"emprendimientos.com/servidor-go/db/sqlc"
)

func TestPublicacionForm(t *testing.T) {
	var buf bytes.Buffer
	err := PublicacionForm().Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering PublicacionForm: %v", err)
	}

	html := buf.String()
	expectedSubstrings := []string{
		`hx-post="/publicaciones"`,
		`hx-target="#lista-publicaciones"`,
		`hx-swap="outerHTML"`,
		`name="id_emprendimiento"`,
		`name="titulo"`,
		`<select name="tipo"`,
		`name="precio"`,
		`name="contenido"`,
		`name="imagen_url"`,
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}

func TestPublicacionFormOOB(t *testing.T) {
	var buf bytes.Buffer
	err := PublicacionFormOOB().Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering PublicacionFormOOB: %v", err)
	}

	html := buf.String()
	for _, sub := range []string{`id="form-publicaciones"`, `hx-swap-oob="true"`, `hx-post="/publicaciones"`} {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}

func TestPublicacionList_Empty(t *testing.T) {
	var buf bytes.Buffer
	err := PublicacionList([]sqlc.Publicacione{}).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering empty PublicacionList: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "No hay publicaciones registradas") {
		t.Errorf("expected 'No hay publicaciones registradas' in html, got: %s", html)
	}
	if !strings.Contains(html, "<table") {
		t.Errorf("expected table in html, got: %s", html)
	}
}

func TestPublicacionList_WithItems(t *testing.T) {
	pubs := []sqlc.Publicacione{
		{
			IDPublicacion:    1,
			IDEmprendimiento: 10,
			Titulo:           "Promo 2x1",
			Tipo:             "promocion",
			Contenido: sql.NullString{
				String: "Descuento especial de viernes",
				Valid:  true,
			},
			Precio: sql.NullString{
				String: "1500.50",
				Valid:  true,
			},
		},
		{
			IDPublicacion:    2,
			IDEmprendimiento: 20,
			Titulo:           "Servicio de Asesoría",
			Tipo:             "servicio",
			Contenido: sql.NullString{
				Valid: false,
			},
			Precio: sql.NullString{
				Valid: false,
			},
		},
	}

	var buf bytes.Buffer
	err := PublicacionList(pubs).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering PublicacionList with items: %v", err)
	}

	html := buf.String()
	if strings.Contains(html, "No hay publicaciones registradas") {
		t.Errorf("did not expect 'No hay publicaciones registradas' in html")
	}

	expectedSubstrings := []string{
		`id="lista-publicaciones"`,
		"Promo 2x1",
		"promocion",
		"10",
		"Descuento especial de viernes",
		"$1500.50",
		"Servicio de Asesoría",
		"servicio",
		"20",
		"-", // fallback for null Precio and Contenido
		`hx-delete="/publicaciones/1"`,
		`hx-delete="/publicaciones/2"`,
		`hx-target="closest tr"`,
		`hx-swap="outerHTML"`,
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}

func TestPublicacionesPage(t *testing.T) {
	var buf bytes.Buffer
	err := PublicacionesPage([]sqlc.Publicacione{}, "promo").Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering PublicacionesPage: %v", err)
	}

	html := buf.String()
	expectedSubstrings := []string{
		"Gestión de Publicaciones",
		`hx-get="/publicaciones"`,
		`hx-trigger="keyup changed delay:500ms, search"`,
		`value="promo"`,
		`hx-post="/publicaciones"`,
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}
