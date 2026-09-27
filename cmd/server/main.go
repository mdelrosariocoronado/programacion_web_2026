package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // driver de PostgreSQL

	"emprendimientos.com/servidor-go/db/sqlc"
	"emprendimientos.com/servidor-go/internal/handlers"
	"emprendimientos.com/servidor-go/internal/repositories"
	"emprendimientos.com/servidor-go/internal/services"
)

func main() {

	// configurac de base de datos desde var de entorno
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "user")
	dbPass := getEnv("DB_PASSWORD", "xyz")
	dbName := getEnv("DB_NAME", "db")

	// construir dns para interpretar los drivers de postgre
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	//conectar a postgresql
	db, error := sql.Open("postgres", dns)
	if error != nil {
		log.Fatalf("Error configurando la base de datos: %v", error)
	}
	defer db.Close()

	//verificacion de conexion
	if error := db.Ping(); error != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", error)
	}
	fmt.Println("--- Conexion exitosa a PostgreSql")

	//--Manejo de capas y su logica

	// USUARIOS
	queries := sqlc.New(db)
	userRepo := repositories.NewUserRepository(queries)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// EMPRENDIMIENTOS
	empRepo := repositories.NewEmprendimientoRepository(queries)
	empService := services.NewEmprendimientoService(empRepo)
	empHandler := handlers.NewEmprendimientoHandler(empService)

	//-- manejo del enrutamiento (mux es para el mapeo de rutas especificas para la capa de presentacion)

	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	//captura cualquier ruta que no coincida con una definida -> ponemos el inicio estatico
	router.HandleFunc("/", usuarioHandler.HandleUsers)

	// endpoint del formulario. verificar si lo dejamos al formulario
	router.HandleFunc("/users", usuarioHandler.HandleUsers)

	//Healthcheck endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	//servidor para archivos estaticos
	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Rutas de Emprendimientos
	router.HandleFunc("POST /api/emprendimientos", empHandler.Create)
	router.HandleFunc("GET /api/emprendimientos", empHandler.GetAll)
	router.HandleFunc("GET /api/emprendimientos/{id}", empHandler.GetByID)
	router.HandleFunc("PUT /api/emprendimientos/{id}", empHandler.Update)
	router.HandleFunc("DELETE /api/emprendimientos/{id}", empHandler.Delete)

	router.HandleFunc("/", userHandler.HandleUsers)
	router.HandleFunc("/users", userHandler.HandleUsers)

	// endpoints para Usuarios
	router.HandleFunc("POST /api/usuarios", usuarioHandler.Create)
	router.HandleFunc("GET /api/usuarios", usuarioHandler.List)
	router.HandleFunc("GET /api/usuarios/{id}", usuarioHandler.GetByID)
	router.HandleFunc("PUT /api/usuarios/{id}", usuarioHandler.Update)
	router.HandleFunc("DELETE /api/usuarios/{id}", usuarioHandler.Delete)

	//ARRANCAR EL SV HTTP
	port := getEnv("PORT", "8080")
	fmt.Printf("--- Servidor corriendo en http://localhost:%s\n", port)

	if error := http.ListenAndServe(":"+port, router); error != nil {
		log.Fatalf("Error al iniciar el servidor: %v", error)
	}

}

// funcion auxiliar para leer variables de entorno con valor por defecto
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
