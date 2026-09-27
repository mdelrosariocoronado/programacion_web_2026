# Trabajo practico N°2
## Integrantes del grupo

*   María Del Rosario Coronado
*   Brisa Fernandez Lema
*   Benjamín Knudsen


## Documentación (Persistencia)
La aplicación modela una red de emprendimientos locales en tandil.

### Modelo de Datos

* **`emprendimientos`:** Entidad principal que almacena el perfil del negocio, rubro y datos de contacto.
* **`usuarios`:** Clientes y emprendedores con control de acceso por roles (`administrador`, `emprendedor`, `cliente`).
* **`publicaciones`:** Contenido publicado por cada emprendimiento (tipo `producto`, `promocion`, `servicio` u `otro`), vinculado mediante clave foránea con eliminación en cascada (`ON DELETE CASCADE`).
* **`suscripciones`:** Tabla intermedia que materializa la relación N:M entre usuarios y los emprendimientos que siguen.

### Herramientas y Patrones
* **SQL Directo con `sqlc`:** Se adoptó `sqlc` para compilar consultas SQL puras a código Go seguro, evitando la discordancia de impedancia y la sobrecarga de un ORM pesado.
* **Migraciones con Atlas:** El esquema se versiona mediante scripts DDL fechados en `db/migrations/`, garantizando reproducibilidad y control de versiones en el motor.
* **Entorno con Docker Compose:** La base de datos se ejecuta de forma aislada y estandarizada mediante contenedores.
  
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
`http://localhost:8080/static/` para conocer mas sobre las entidades y sus atributos.



---
*Proyecto incremental desarrollando una web completa, para la materia de Programación Web de la carrera de Ingenieria en Sistemas, UNICEN.*
