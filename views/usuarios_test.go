package views

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
	"testing"

	"emprendimientos.com/servidor-go/db/sqlc"
)

func TestUsuarioForm(t *testing.T) {
	var buf bytes.Buffer
	err := UsuarioForm().Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering UsuarioForm: %v", err)
	}

	html := buf.String()
	expectedSubstrings := []string{
		`<form method="POST" action="/usuarios">`,
		`name="nombre_completo"`,
		`name="email"`,
		`name="clave"`,
		`<select name="rol"`,
		`name="id_emprendimiento"`,
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}

func TestUsuarioList_Empty(t *testing.T) {
	var buf bytes.Buffer
	err := UsuarioList([]sqlc.Usuario{}).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering empty UsuarioList: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "No hay usuarios registrados") {
		t.Errorf("expected 'No hay usuarios registrados' in html, got: %s", html)
	}
	if !strings.Contains(html, "<table") {
		t.Errorf("expected table in html, got: %s", html)
	}
}

func TestUsuarioList_WithItems(t *testing.T) {
	usuarios := []sqlc.Usuario{
		{
			IDUsuario:      1,
			NombreCompleto: "Juan Perez",
			Email:          "juan@example.com",
			Rol:            "cliente",
			IDEmprendimiento: sql.NullInt32{
				Int32: 42,
				Valid: true,
			},
		},
		{
			IDUsuario:      2,
			NombreCompleto: "Maria Gomez",
			Email:          "maria@example.com",
			Rol:            "emprendedor",
			IDEmprendimiento: sql.NullInt32{
				Valid: false,
			},
		},
	}

	var buf bytes.Buffer
	err := UsuarioList(usuarios).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering UsuarioList with items: %v", err)
	}

	html := buf.String()
	if strings.Contains(html, "No hay usuarios registrados") {
		t.Errorf("did not expect 'No hay usuarios registrados' in html")
	}

	expectedSubstrings := []string{
		"Juan Perez",
		"juan@example.com",
		"cliente",
		"42",
		"Maria Gomez",
		"maria@example.com",
		"emprendedor",
		"-",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}

func TestUsuariosPage(t *testing.T) {
	var buf bytes.Buffer
	err := UsuariosPage([]sqlc.Usuario{}).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering UsuariosPage: %v", err)
	}

	html := buf.String()
	expectedSubstrings := []string{
		"Gestión de Usuarios",
		`<form method="POST" action="/usuarios">`,
		"No hay usuarios registrados",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}
