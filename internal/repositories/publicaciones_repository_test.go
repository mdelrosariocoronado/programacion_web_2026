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
			Tipo: "promocion",
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

		if fetched.Titulo != "Oferta Especial de Prueba" {
			t.Errorf("Título incorrecto, esperado 'Oferta Especial de Prueba', obtenido '%s'", fetched.Titulo)
		}
	})


	t.Run("List Publicaciones", func(t *testing.T) {
		lista, err := queries.ListPublicacionesByEmprendimiento(contexto, emprendimiento.IDEmprendimiento)
		if err != nil {
			t.Fatalf("ListPublicaciones falló: %v", err)
		}

		if len(lista) == 0 {
			t.Errorf("Se esperaba al menos 1 publicacion en la lista")
		}
	})

	// UPDATE// UPDATE
	t.Run("Update Publicacion", func(t *testing.T) {
		nuevoTitulo := "Oferta Actualizada de Primavera"

		err := queries.UpdatePublicacion(contexto, db.UpdatePublicacionParams{
			IDPublicacion: int32(createPublicacionID),
			Titulo:        nuevoTitulo,
			Contenido:     sql.NullString{String: "Nuevo contenido descriptivo", Valid: true},
			ImagenUrl:     sql.NullString{String: "http://ejemplo.com/imagen.jpg", Valid: true},
			Tipo:          "promocion", 
			Precio:        sql.NullString{String: "1500.00", Valid: true}, 
		})
		if err != nil {
			t.Fatalf("UpdatePublicacion falló: %v", err)
		}

		// Releer desde la base de datos para confirmar que persistió
		updated, err := queries.GetPublicacion(contexto, int32(createPublicacionID))
		if err != nil {
			t.Fatalf("GetPublicacion tras update falló: %v", err)
		}

		if updated.Titulo != nuevoTitulo {
			t.Errorf("El título no se actualizó correctamente. Esperado: '%s', obtenido: '%s'", nuevoTitulo, updated.Titulo)
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

	t.Run("Cascade Delete Publicaciones al eliminar Emprendimiento", func(t *testing.T) {
	// emprendimiento para el test
	emp, _ := queries.CreateEmprendimiento(contexto, db.CreateEmprendimientoParams{
		Nombre: "Negocio Temporal",
		Rubro:  "servicios",
	})

	//  publicación relacionada con el emprendimiento
	pub, _ := queries.CreatePublicacion(contexto, db.CreatePublicacionParams{
		IDEmprendimiento: emp.IDEmprendimiento,
		Titulo:           "Post que debe desaparecer",
		Tipo:             "otro",
	})

	// eliminar emprendimiento (origen de realcion)
	err := queries.DeleteEmprendimiento(contexto, emp.IDEmprendimiento)
	if err != nil {
		t.Fatalf("DeleteEmprendimiento falló: %v", err)
	}

	// verificar que se borro en cascada sus publicaciones
	_, err = queries.GetPublicacion(contexto, pub.IDPublicacion)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Se esperaba sql.ErrNoRows en la publicación tras borrar el negocio (ON DELETE CASCADE), se obtuvo: %v", err)
	}
})
}