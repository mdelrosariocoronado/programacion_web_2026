package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq" // Driver de PostgreSQL

	db "emprendimientos.com/servidor-go/db/sqlc"
)

//primer test que aisla la conexión y configuracion inicial con bd

func setUpTestDB( t *testing.T) (*sql.DB, *db.Queries){
	t.Helper() // para que los logs indiquen la línea del test caller

	dsn := "postgres://user:xyz@localhost:5432/db?sslmode=disable"
	conection, err := sql.Open("postgres", dsn)

	if err != nil {
		t.Fatalf(" -- SETUP TEST: error al intentar abrir la conexion: %v", err)
	}

	contexto, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	if err := conection.PingContext(contexto); err!= nil{
		conection.Close()
		t.Fatalf(" -- SETUP TEST: la base de datos no responde ante el Ping: %v", err)
	}

	//CIERRE AUTOMATICO 

	t.Cleanup(func(){
		conection.Close()
	})

	return conection, db.New(conection)

}