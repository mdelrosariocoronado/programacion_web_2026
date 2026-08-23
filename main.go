package main

import (
  "fmt"
  "net/http"
)

func main() {
  //Directorio que contiene los archivos estáticos.
  staticDir := "./static"

  //Manejador de servidor de archivos. (Configura el header automaticamente)
  fileServer := http.FileServer(http.Dir(staticDir))
  http.Handle("/", fileServer)

  //Deficion del puerto y muestra un mensaje.
  port := ":8080"
  fmt.Printf("Servidor ESTÁTICO escuchando en http://localhost%s\n", port)
  fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

  //Inicio el servidor.
  err := http.ListenAndServe(port, nil)
  if err != nil {
    fmt.Printf("Error al iniciar el servidor: %s\n", err)
  }
}