package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	db "emprendimientos.com/servidor-go/db/sqlc"
)

// TestPublicacionRepository_CRUD valida el ciclo completo CRUD y reglas de integridad para PublicacionRepo
func TestPublicacionRepository_CRUD(t *testing.T) {
	_, queries := setUpTestDB(t) // helper de setup_test.go
	contexto := context.Background()

	repo := NewPublicacionRepository(queries)

	// emprendimiento para cumplir FK
	emprendimiento, err := queries.CreateEmprendimiento(contexto, db.CreateEmprendimientoParams{
		Nombre: "emprendimiento prueba",
		Rubro:  "peluqueria",
	})
	if err != nil {
		t.Fatalf("setup emprendimiento fallo: %v", err)
	}

	var createPublicacionID int32

	// CREATE Y READ
	t.Run("Create and Read Publicacion", func(t *testing.T) {
		pub, err := repo.Create(contexto, db.CreatePublicacionParams{
			IDEmprendimiento: emprendimiento.IDEmprendimiento,
			Titulo:           "Oferta Especial de Prueba",
			Tipo:             "promocion",
		})
		if err != nil {
			t.Fatalf("repo.Create fallo: %v", err)
		}

		if pub.IDPublicacion == 0 {
			t.Errorf("Se esperaba ID autogenerado mayor a 0")
		}

		createPublicacionID = pub.IDPublicacion

		fetched, err := repo.GetByID(contexto, createPublicacionID)
		if err != nil {
			t.Fatalf("repo.GetByID fallo: %v", err)
		}

		if fetched.Titulo != "Oferta Especial de Prueba" {
			t.Errorf("Titulo incorrecto, esperado 'Oferta Especial de Prueba', obtenido '%s'", fetched.Titulo)
		}
	})

	// LIST
	t.Run("List Publicaciones", func(t *testing.T) {
		lista, err := repo.ListByEmprendimiento(contexto, emprendimiento.IDEmprendimiento)
		if err != nil {
			t.Fatalf("repo.ListByEmprendimiento fallo: %v", err)
		}

		if len(lista) == 0 {
			t.Errorf("Se esperaba al menos 1 publicacion en la lista")
		}
	})

	// UPDATE
	t.Run("Update Publicacion", func(t *testing.T) {
		nuevoTitulo := "Oferta Actualizada de Primavera"

		err := repo.Update(contexto, db.UpdatePublicacionParams{
			IDPublicacion: createPublicacionID,
			Titulo:        nuevoTitulo,
			Contenido:     sql.NullString{String: "Nuevo contenido descriptivo", Valid: true},
			ImagenUrl:     sql.NullString{String: "http://ejemplo.com/imagen.jpg", Valid: true},
			Tipo:          "promocion",
			Precio:        sql.NullString{String: "1500.00", Valid: true},
		})
		if err != nil {
			t.Fatalf("repo.Update fallo: %v", err)
		}

		// Releer desde la base de datos para confirmar que persistio
		updated, err := repo.GetByID(contexto, createPublicacionID)
		if err != nil {
			t.Fatalf("repo.GetByID tras update fallo: %v", err)
		}

		if updated.Titulo != nuevoTitulo {
			t.Errorf("El titulo no se actualizo correctamente. Esperado: '%s', obtenido: '%s'", nuevoTitulo, updated.Titulo)
		}
	})

	// DELETE
	t.Run("Delete Publicacion", func(t *testing.T) {
		err := repo.Delete(contexto, createPublicacionID)
		if err != nil {
			t.Fatalf("repo.Delete fallo: %v", err)
		}

		_, err = repo.GetByID(contexto, createPublicacionID)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Se esperaba sql.ErrNoRows al consultar una publicacion eliminada, se obtuvo: %v", err)
		}
	})

	// CASCADE DELETE
	t.Run("Cascade Delete Publicaciones al eliminar Emprendimiento", func(t *testing.T) {
		// emprendimiento temporal para el test
		emp, err := queries.CreateEmprendimiento(contexto, db.CreateEmprendimientoParams{
			Nombre: "Negocio Temporal",
			Rubro:  "servicios",
		})
		if err != nil {
			t.Fatalf("setup negocio temporal fallo: %v", err)
		}

		// publicacion relacionada usando el repositorio
		pub, err := repo.Create(contexto, db.CreatePublicacionParams{
			IDEmprendimiento: emp.IDEmprendimiento,
			Titulo:           "Post que debe desaparecer",
			Tipo:             "otro",
		})
		if err != nil {
			t.Fatalf("repo.Create fallo: %v", err)
		}

		// eliminar emprendimiento (origen de la relacion)
		err = queries.DeleteEmprendimiento(contexto, emp.IDEmprendimiento)
		if err != nil {
			t.Fatalf("DeleteEmprendimiento fallo: %v", err)
		}

		// verificar que se borro en cascada la publicacion
		_, err = repo.GetByID(contexto, pub.IDPublicacion)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Se esperaba sql.ErrNoRows en la publicacion tras borrar el negocio (ON DELETE CASCADE), se obtuvo: %v", err)
		}
	})
}
