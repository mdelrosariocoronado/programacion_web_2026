package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	db "emprendimientos.com/servidor-go/db/sqlc"
)

//TEST CRUD en PUBLICACIONES
func TestQueries_Publicaciones_CRUD(t *testing.T) {
	_, queries := setUpTestDB(t) // helper de setup_test.go
	contexto := context.Background()

	// emprendimiento para cumplir FK
	emprendimiento, err := queries.CreateEmprendimiento(contexto, db.CreateEmprendimientoParams{
		Nombre: "emprendimiento prueba",
		Rubro: "peluqueria",

	})
	if err != nil {
		t.Fatalf("setup autor falló: %v", err)
	}

	var createPublicacionID int64


	// CREATE Y READ
	t.Run("Create and Read Publicacion", func(t *testing.T) {
		pub, err := queries.CreatePublicacion(contexto, db.CreatePublicacionParams{
			IDEmprendimiento: emprendimiento.IDEmprendimiento, // Pasamos el ID del usuario recién creado
			Titulo:    "Oferta Especial de Prueba",
			Tipo: "prueba",
		})
		if err != nil {
			t.Fatalf("CreatePublicacion falló: %v", err)
		}

		if pub.IDPublicacion == 0 {
			t.Errorf("Se esperaba ID autogenerado mayor a 0")
		}

		createPublicacionID = int64(pub.IDPublicacion)

		fetched, err := queries.GetPublicacion(contexto, int32(createPublicacionID))
		if err != nil {
			t.Fatalf("GetPublicacion falló: %v", err)
		}

		if fetched.Titulo != "Oferta Especial de Primavera" {
			t.Errorf("Título incorrecto, esperado 'Oferta Especial de Primavera', obtenido '%s'", fetched.Titulo)
		}
	})

	// UPDATE
	t.Run("Update Publicacion", func(t *testing.T) {
		err := queries.UpdatePublicacion(contexto, db.UpdatePublicacionParams{
			IDPublicacion: int32(createPublicacionID),
			Titulo: "Oferta Actualizada",
			Tipo: "evento",
		})
		
		if err != nil {
			t.Fatalf("UpdatePublicacion falló: %v", err)
		}

		updated, _ := queries.GetPublicacion(contexto, int32(createPublicacionID))
		if updated.Titulo != "Oferta Actualizada" {
			t.Errorf("El título no se actualizó correctamente")
		}
	})

	// DELETE
	t.Run("Delete Publicacion", func(t *testing.T) {
		err := queries.DeletePublicacion(contexto, int32(createPublicacionID))
		if err != nil {
			t.Fatalf("DeletePublicacion falló: %v", err)
		}

		_, err = queries.GetPublicacion(contexto, int32(createPublicacionID))
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Se esperaba sql.ErrNoRows al consultar una publicación eliminada, se obtuvo: %v", err)
		}
	})
}