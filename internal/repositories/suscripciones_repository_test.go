package repositories

import (
	"context"
	"testing"

	db "emprendimientos.com/servidor-go/db/sqlc"
)

func TestQueries_Suscripciones_CRUD(t *testing.T) {
	_, queries := setUpTestDB(t)
	contexto := context.Background()

	// usuario suscriptor para FK
	suscriptor, err := queries.CreateUsuario(contexto, db.CreateUsuarioParams{
		NombreCompleto: "suscriptor prueba",
		Email:          "prueba@gmail.com",
		Rol: "emprendedor",
	})
	if err != nil {
		t.Fatalf("setup suscriptor falló: %v", err)
	}


	// emprendimiento para cumplir FK
	emprendimiento, err := queries.CreateEmprendimiento(contexto, db.CreateEmprendimientoParams{
		Nombre: "emprendimiento prueba",
		Rubro: "peluqueria",

	})
	if err != nil {
		t.Fatalf("setup autor falló: %v", err)
	}

	var createSuscripcionID int64



	// CREATE Y READ
	t.Run("Create and Read Suscripcion", func(t *testing.T) {
		sub, err := queries.CreateSuscripcion(contexto, db.CreateSuscripcionParams{
			IDUsuario: suscriptor.IDUsuario,
			IDEmprendimiento: emprendimiento.IDEmprendimiento,

		})
		if err != nil {
			t.Fatalf("CreateSuscripcion falló: %v", err)
		}

		if sub.IDSuscripcion == 0 {
			t.Errorf("Se esperaba ID de suscripción mayor a 0, se obtuvo: %d", sub.IDSuscripcion)
		}

		createSuscripcionID = int64(sub.IDSuscripcion)

		fetched, err := queries.ListSuscripcionesByUsuario(contexto, int32(createSuscripcionID))
		if err != nil {
			t.Fatalf("GetSuscripcion falló: %v", err)
		}

		if len(fetched) == 0 {
			t.Errorf("El ID de usuario asociado no coincide")
		}
	})

	// 2. DELETE Y VERIFICACIÓN
	t.Run("Delete Suscripcion", func(t *testing.T) {
		err := queries.DeleteSuscripcion(contexto, db.DeleteSuscripcionParams{
			IDUsuario: suscriptor.IDUsuario,
			IDEmprendimiento: emprendimiento.IDEmprendimiento,
		})

		if err != nil {
			t.Fatalf("DeleteSuscripcion falló: %v", err)
		}
		
		suscripciones, err := queries.ListSuscripcionesByUsuario(contexto, suscriptor.IDUsuario)
		if err != nil {
			t.Fatalf("ListSuscripcionesByUsuario falló tras eliminar: %v", err)
		}

		if len(suscripciones) != 0 {
			t.Errorf("Se esperaba lista vacía tras eliminar la suscripción, se encontraron: %d elementos", len(suscripciones))
		}
	})
}
