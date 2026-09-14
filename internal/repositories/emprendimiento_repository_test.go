package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	db "emprendimientos.com/servidor-go/db/sqlc"
)

func TestQueries_Emprendimientos_CRUD(t *testing.T) {
	_, queries := setUpTestDB(t)
	contexto := context.Background()

	var empID int32

	// Create y Read 
	t.Run("Create and Read Emprendimiento", func(t *testing.T) {
		emp, err := queries.CreateEmprendimiento(contexto, db.CreateEmprendimientoParams{
			Nombre:      "Café Central",
			Url:         sql.NullString{String: "cafe-central", Valid: true},
			Rubro:       "gastronomia",
			Descripcion: sql.NullString{String: "Especialidad en café", Valid: true},
			Activo:      sql.NullBool{Bool: true, Valid: true},
		})
		if err != nil {
			t.Fatalf("CreateEmprendimiento falló: %v", err)
		}

		if emp.IDEmprendimiento == 0 {
			t.Errorf("Se esperaba ID mayor a 0, se obtuvo %d", emp.IDEmprendimiento)
		}

		empID = emp.IDEmprendimiento

		fetched, err := queries.GetEmprendimiento(contexto, empID)
		if err != nil {
			t.Fatalf("GetEmprendimiento falló: %v", err)
		}

		if fetched.Nombre != "Café Central" {
			t.Errorf("Esperado 'Café Central', obtenido '%s'", fetched.Nombre)
		}
	})

	// List (muchos)
	t.Run("List Emprendimientos", func(t *testing.T) {
		lista, err := queries.ListEmprendimientos(contexto)
		if err != nil {
			t.Fatalf("ListEmprendimientos falló: %v", err)
		}

		if len(lista) == 0 {
			t.Errorf("Se esperaba al menos 1 emprendimiento en la lista")
		}
	})

	// UPdate
	t.Run("Update Emprendimiento", func(t *testing.T) {
		err := queries.UpdateEmprendimiento(contexto, db.UpdateEmprendimientoParams{
			IDEmprendimiento: empID,
			Nombre:           "Café Central Tandil",
			Url:              sql.NullString{String: "cafe-central-tandil", Valid: true},
			Rubro:            "gastronomia",
			Descripcion:      sql.NullString{String: "Cafetería y tostador", Valid: true},
			Activo:           sql.NullBool{Bool: true, Valid: true},
		})
		if err != nil {
			t.Fatalf("UpdateEmprendimiento falló: %v", err)
		}

		updated, err := queries.GetEmprendimiento(contexto, empID)
		if err != nil {
			t.Fatalf("GetEmprendimiento tras update falló: %v", err)
		}

		if updated.Nombre != "Café Central Tandil" {
			t.Errorf("El nombre no se actualizó correctamente")
		}
	})

	// Delete
	t.Run("Delete Emprendimiento", func(t *testing.T) {
		err := queries.DeleteEmprendimiento(contexto, empID)
		if err != nil {
			t.Fatalf("DeleteEmprendimiento falló: %v", err)
		}

		_, err = queries.GetEmprendimiento(contexto, empID)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Se esperaba sql.ErrNoRows al consultar un emprendimiento eliminado, se obtuvo: %v", err)
		}
	})
}