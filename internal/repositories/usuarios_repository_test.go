package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "github.com/lib/pq" // Driver de PostgreSQL

	db "emprendimientos.com/servidor-go/db/sqlc"
)



// suite de pruebas principal CRUD (CREATE, READ, UPDATE, DELETE) de USUARIO

func TestQueries_CRUD(t *testing.T){
	_, queries := setUpTestDB(t) //el helper declarado arriba que abre y prepara la bd
	contexto:= context.Background()

	var createID int64

	//TEST CREATE Y READ
	t.Run("Create and Read User", func(t *testing.T){
		newUser, err := queries.CreateUsuario( contexto, db.CreateUsuarioParams{
				NombreCompleto: "Rosario Coronado",
				Email: "ro@gmail.com",
				Clave: "12345",
				Rol: "cliente",
		})
		if err != nil {
			t.Fatalf("CreateUser fallo: %v", err)
		}

		if newUser.IDUsuario == 0 {
			t.Errorf("Se esperaba ID autogenerado que sea mayor a 0 , se obtuvo: %d", newUser.IDUsuario)
		}

		createID = int64(newUser.IDUsuario) // guarda para las otras pruebas

		fetchedUser, err := queries.GetUsuario(contexto, int32(createID))

		if err != nil {

			t.Fatalf("GetUserByID fallo: %v", err)
		}

		if fetchedUser.Email != "ro@gmail.com"{
			t.Errorf(" Email incorrecto, se esperaba 'ro@gmail.com' pero se obtuvo '%s'", fetchedUser.Email)
		}



	})

	//TEST UPDATE

	t.Run("Update User", func(t *testing.T) {
		err:= queries.UpdateUsuario(contexto,  db.UpdateUsuarioParams{
			IDUsuario: int32(createID),
			NombreCompleto: "M. del Rosario Coronado",
			Email: "ro@gmail.com",
			Clave: "12345",
			Rol: "cliente",
		})

		if err!= nil {
			t.Fatalf("UpdateUser fallo en: %v", err)

		}

		updatedUser, _ := queries.GetUsuario(contexto, int32(createID))

		if updatedUser.NombreCompleto != "M. del Rosario Coronado" {
			t.Errorf("El nombre completo no se actualizo correctamente")
		}


	})

	//TEST DELETE
	t.Run("Delete", func (t *testing.T){
		err:= queries.DeleteUsuario(contexto, int32(createID))

		if err!= nil {
			t.Fatalf("DeleteUser fallo: %v", err) 
		}

		_, err = queries.GetUsuario(contexto, int32(createID))
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("se esperaba sql.ErrNoRows al consultar por el usuario eliminado, se obtuvo: %v", err)
		}

	})


}


