# Trabajo practico N°2
## Integrantes del Equipo

*   María Del Rosario Coronado
*   Brisa Fernandez Lema
*   Benjamín Knudsen
  
## Requisitos Previos

Antes de ejecutar este proyecto, asegúrate de tener instalado el siguiente software en tu computadora:

* **[Go](https://go.dev/dl/)** (v1.22 o superior)
* **[Docker](https://docs.docker.com/get-docker/) & Docker Compose** (con el servicio iniciado)
* **[Atlas CLI](https://atlasgo.io/getting-started/)** (para la aplicación de migraciones DDL)
* **[Git](https://git-scm.com/downloads)**
* *(Opcional)* **[sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)** (el código ya se encuentra generado en `db/sqlc/`)


## Instalación y Ejecución Local

Para probar este proyecto en tu entorno local, sigue estos pasos en tu terminal:

**1. Clonar el repositorio**
Descarga el código fuente a tu computadora ejecutando:
`git clone https://github.com/mdelrosariocoronado/prog_web.git`

**2. Navegar al directorio del proyecto**
Ingresa a la carpeta principal que se acaba de descargar:
`cd prog_web`

**3.Configura las variables de entorno**
`cp .env.example .env`

**4. Ejecución de Pruebas Automatizadas**
El proyecto incluye un pipeline automatizado en el `Makefile` que levanta la base de datos en Docker, aplica el esquema DDL con Atlas, compila las consultas con `sqlc` y corre toda la suite de tests unitarios y de integración:
`make test`

## Ejecución del Servidor Web, CORRESPONDIENTE AL TP1

**Iniciar el servidor**
Ejecuta el archivo principal de Go para levantar el servidor web:
`make run` 

**Visualizar la página del TP1**
Abre tu navegador web de preferencia y dirígete a la siguiente dirección:
`http://localhost:8080/static/`

---
*Proyecto incremental desarrollando una web completa, para la materia de Programación Web de la carrera de Ingenieria en Sistemas, UNICEN.*
