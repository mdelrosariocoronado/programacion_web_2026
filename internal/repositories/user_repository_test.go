package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	db "programacion_web_2026/db/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//primer test que aisla la conexión y configuracion inicial con bd

func setUpTestDB( t *testing.T) (*sql.DB, *db.Queries){
	t.Helper() // para que los logs indiquen la línea del test caller

	dsn := "postgres://user:xyz@localhost:5432/db?sslmode=disable"
	conection, error := sql.Open("pgx", dsn)

	if error != nil {
		t.Fatalf(" -- SETUP TEST: error al intentar abrir la conexion: %v", error)
	}

	contexto, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	if error := conection.PingContext(contexto); error!= nil{
		conection.Close()
		t.Fatalf(" -- SETUP TEST: la base de datos no responde ante el Ping: %v", error)
	}

	//CIERRE AUTOMATICO 

	t.Cleanup(func(){
		conection.Close()
	})

	return conection, db.New(conection)

}


// suite de pruebas principal CRUD (CREATE, READ, UPDATE, DELETE)

func TestQueries_CRUD(t *testing.T){
	_, queries := setUpTestDB(t) //el helper declarado arriba que abre y prepara la bd
	contexto:= context.Background()

	var createID int64

	//TEST CREATE Y READ
	t.Run("Create and Read User", func(t *testing.T){
		newUser, error := queries.CreateUser( contexto, db.CreateUserParams{
				NombreCompleto: "Rosario Coronado",
				Email: "ro@gmail.com",
		})
		if error != nil {
			t.Fatalf("CreateUser fallo: %v", error)
		}

		if newUser.ID == 0 {
			t.Errorf("Se esperaba ID autogenerado que sea mayor a 0 , se obtuvo: %d", newUser.ID)
		}

		createID = newUser.ID // guarda para las otras pruebas

		fetchedUser, error := queries.GetUserByID(contexto, createID)

		if error != nil {

			t.Fatalf("GetUserByID fallo: %v", error)
		}

		if fetchedUser.Email != "ro@gmail.com"{
			t.Error(" Email incorrecto, se esperaba 'ro@gmail.com' pero se obtuvo '$s'", fetchedUser.Email)
		}



	})

	//TEST UPDATE

	t.Run("Update User", func(t *testing.T) {
		error:= queries.UpdateUser(contexto,  db.UpdateUserParams{}
			ID: createID,
			NombreCompleto: "M. del Rosario Coronado",
			Email: "ro@gmail.com",
		})

		if error!= nil {
			t.Fatalf("UpdateUser fallo en: %v", error)

		}

		updatedUser, _ := queries.GetUserByID(contexto, createID)

		if updatedUser.NombreCompleto != "M. del Rosario Coronado" {
			t.Errorf("El nombre completo no se actualizo correctamente")
		}


	})

	//TEST DELETE
	t.Run("Delete", func (t *testing.T){
		error:= queries.DeleteUser(contexto, createID)

		if error!= nil {
			t.Fatalf("DeleteUser fallo: %v", error) 
		}

		_, error = queries.GetUserByID(contexto, createdID)
		if !errors.Is(error, sql.ErrNoRows) {
			t.Errorf("se esperaba sql.ErrNoRows al consultar por el usuario eliminado, se obtuvo: %v", err)
		}

	})


}


