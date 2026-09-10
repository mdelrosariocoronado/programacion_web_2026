package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // Driver de PostgreSQL

	"emprendimientos.com/servidor-go/db/sqlc"
	"emprendimientos.com/servidor-go/internal/handlers"
	"emprendimientos.com/servidor-go/internal/repositories"
	"emprendimientos.com/servidor-go/internal/services"
)

func main (){

  // configurac de base de datos desde var de entorno
  dbHost := getEnv("DB_HOST", "localhost")
  dbPort := getEnv("DB_PORT", "5432")
  dbUser := getEnv("DB_USER", "user")         
  dbPass := getEnv("DB_PASSWORD", "xyz")       
  dbName := getEnv("DB_NAME", "db")            

  // construir dns
  dns:= fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

  //conectar a postgresql
  db, error:= sql.Open("postgres", dns)
  if error!= nil{
    log.Fatalf("Error configurando la base de datos: %v", error)
  }
  defer db.Close()

  //verificacion de conexion
  if error:= db.Ping(); error!= nil{
    log.Fatalf("No se pudo conectar a la base de datos: %v", error)
  }
  fmt.Println("--- Conexion exitosa a PostgreSql")






  //--Manejo de capas y su logica

  // DATOS
  
  queries := sqlc.New(db)
  userRepo:= repositories.NewUserRepository(queries)

  //LOGICA
  userService := services.NewUserService(userRepo)

  // PRESENTACION
  userHandler:= handlers.NewUserHandler(userService)


  //-- manejo del enrutamiento (mux es para el mapeo de rutas especificas para la capa de presentacion)
  
  router:= http.NewServeMux()

  fileServer:= http.FileServer(http.Dir("./static"))
  router.Handle("/static/", http.StripPrefix("/static/", fileServer))

  router.HandleFunc("/", userHandler.HandleUsers)

  // endpoint del formulario. verificar si lo dejamos al formulario
  router.HandleFunc("/users", userHandler.HandleUsers)

  //Healthcheck endpoint
  router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(http.StatusOK)
    w.Write([]byte ("OK"))
  })

  
  //ARRANCAR EL SV HTTP
  port:= getEnv("PORT", "8080")
  fmt.Printf("--- Servidor corriendo en http://localhost:%s\n", port)

  if error:= http.ListenAndServe(":"+port, router); error!=nil{
    log.Fatalf("Error al iniciar el servidor: %v", error)
  }

}


//funcion auxiliar para leer variables de entorno con valor por defecto
func getEnv(key, fallback string) string{
  if value, exists:= os.LookupEnv(key); exists {
    return value
  }
  return fallback
}
