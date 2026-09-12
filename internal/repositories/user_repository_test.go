package repositories

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"emprendimientos.com/servidor-go/db/sqlc"
	_ "github.com/lib/pq"
)

// Devuelve el valor de la variable o hace fallar el test si no está definida
func getRequiredEnv(t *testing.T, key string) string {
	val := os.Getenv(key)
	if val == "" {
		t.Fatalf("Error de configuración: la variable de entorno %s no está definida en el .env", key)
	}
	return val
}

func setupTestDB(t *testing.T) *sqlc.Queries {
	// Leemos las variables exportadas por el Makefile desde el .env
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost" // Valores no sensibles (como host o puerto) sí pueden tener fallback
	}
	
	port := os.Getenv("DB_PORT_EXTERNAL")
	if port == "" {
		port = "5432"
	}

	user := getRequiredEnv(t, "DB_USER")
	password := getRequiredEnv(t, "DB_PASSWORD")
	dbname := getRequiredEnv(t, "DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error al abrir la conexión de prueba: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("No se pudo conectar a la base de datos de prueba: %v", err)
	}

	return sqlc.New(db)
}