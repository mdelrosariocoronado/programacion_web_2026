package views

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"emprendimientos.com/servidor-go/internal/services"
)

func TestEmprendimientoForm(t *testing.T) {
	var buf bytes.Buffer
	err := EmprendimientoForm().Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering EmprendimientoForm: %v", err)
	}

	html := buf.String()
	expectedSubstrings := []string{
		`<form method="POST" action="/emprendimientos">`,
		`name="nombre"`,
		`name="rubro"`,
		`name="descripcion"`,
		`name="contacto"`,
		`name="url"`,
		`name="logo"`,
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}

func TestEmprendimientoList_Empty(t *testing.T) {
	var buf bytes.Buffer
	err := EmprendimientoList([]services.EmprendimientoResponse{}).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering empty EmprendimientoList: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "No hay emprendimientos registrados todavía.") {
		t.Errorf("expected 'No hay emprendimientos registrados todavía.' in html, got: %s", html)
	}
}

func TestEmprendimientoList_WithItems(t *testing.T) {
	items := []services.EmprendimientoResponse{
		{
			IDEmprendimiento: 1,
			Nombre:           "Cafetería Las Orquídeas",
			Rubro:            "Gastronomía",
			Descripcion:      "Café de especialidad en Tandil",
			Contacto:         "2494112233",
			Url:              "https://lasorquideas.com",
			Activo:           true,
		},
		{
			IDEmprendimiento: 2,
			Nombre:           "Taller Creativo",
			Rubro:            "Artesanías",
			Contacto:         "contacto@taller.com",
			Activo:           false,
		},
	}

	var buf bytes.Buffer
	err := EmprendimientoList(items).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering EmprendimientoList with items: %v", err)
	}

	html := buf.String()
	if strings.Contains(html, "No hay emprendimientos registrados todavía.") {
		t.Errorf("did not expect 'No hay emprendimientos registrados todavía.' in html")
	}

	expectedSubstrings := []string{
		"Cafetería Las Orquídeas",
		"Gastronomía",
		"Café de especialidad en Tandil",
		"https://lasorquideas.com",
		"Taller Creativo",
		"Artesanías",
		"badge-activo",
		"badge-inactivo",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}

func TestEmprendimientosPage(t *testing.T) {
	var buf bytes.Buffer
	err := EmprendimientosPage([]services.EmprendimientoResponse{}).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("error rendering EmprendimientosPage: %v", err)
	}

	html := buf.String()
	expectedSubstrings := []string{
		"Gestión de Emprendimientos",
		`<form method="POST" action="/emprendimientos">`,
		"No hay emprendimientos registrados todavía.",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("expected html to contain '%s', got: %s", sub, html)
		}
	}
}
