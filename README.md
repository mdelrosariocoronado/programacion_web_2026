# Trabajo practico N°2

## Integrantes del grupo
*   María Del Rosario Coronado
*   Brisa Fernandez Lema
*   Benjamín Knudsen


## Documentación del Proyecto

Aplicación web desarrollada en Go que modela una red de emprendimientos locales en Tandil. Implementa una arquitectura por capas limpia (**Handlers $\rightarrow$ Services $\rightarrow$ Repositories $\rightarrow$ Database**) con persistencia en PostgreSQL mediante consultas tipadas y compiladas con `sqlc`.

### Modelo de Datos
* **`emprendimientos`:** Perfil del negocio, rubro y datos de contacto.
* **`usuarios`:** Clientes y emprendedores con control de acceso por roles (`administrador`, `emprendedor`, `cliente`).
* **`publicaciones`:** Contenido publicado por cada emprendimiento (`producto`, `promocion`, `servicio` u `otro`), vinculado mediante clave foránea con eliminación en cascada (`ON DELETE CASCADE`).
* **`suscripciones`:** Tabla intermedia que materializa la relación N:M entre usuarios y los emprendimientos que siguen.

### Herramientas y Patrones
* **SQL Directo con `sqlc`:** Se adoptó `sqlc` para compilar consultas SQL puras a código Go seguro, evitando la discordancia de impedancia y la sobrecarga de un ORM pesado.
* **Migraciones con Atlas:** El esquema se versiona mediante scripts DDL fechados en `db/migrations/`, garantizando reproducibilidad y control de versiones en el motor.
* **Entorno con Docker Compose:** Ejecución de base de datos y aplicación web en contenedores reproducibles.

## Requisitos Previos
Antes de ejecutar este proyecto, asegúrate de tener instalado el siguiente software en tu computadora:

* **[Go](https://go.dev/dl/)** (v1.22 o superior)
* **[Docker](https://docs.docker.com/get-docker/) & Docker Compose** (en ejecución)
* **[Atlas CLI](https://atlasgo.io/getting-started/)** (para aplicar migraciones DDL)
* **[Git](https://git-scm.com/downloads)**
* **cURL** (para ejecutar las pruebas end-to-end)
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



## Guía de Ejecución

### Opción A: Ejecución con Docker (Recomendada)

Permite levantar tanto la base de datos PostgreSQL como la API de Go dentro de contenedores de forma aislada.

**Forma rápida con Makefile:**
```bash
# Construye, levanta los contenedores y aplica migraciones automáticamente:
make docker-up

# Para detener todos los contenedores:
make docker-down
```

**O de forma manual paso a paso:**
1. **Construir y levantar todos los servicios:**
   ```bash
   docker compose up -d --build
   ```

2. **Aplicar migraciones:**
   ```bash
   atlas migrate apply --dir "file://db/migrations" --url "postgres://user:xyz@localhost:5432/db?sslmode=disable"
   ```

3. **Verificar estado:**
   El servidor estará disponible en `http://localhost:8080/health`.

4. **Detener contenedores:**
   ```bash
   docker compose down
   ```


---

### Opción B: Ejecución sin Docker (Local Nativo)

En este modo, el servidor Go se ejecuta directamente sobre el sistema operativo:

1. **Asegurar que PostgreSQL esté corriendo:**
   Puedes usar una instancia local de PostgreSQL o levantar únicamente el contenedor de la BD:
   ```bash
   docker compose up -d db
   ```

2. **Aplicar migraciones:**
   ```bash
   atlas migrate apply --dir "file://db/migrations" --url "postgres://user:xyz@localhost:5432/db?sslmode=disable"
   ```

3. **Iniciar el servidor web:**
   ```bash
   make run
   # o alternativamente:
   go run cmd/server/main.go
   ```

4. **Acceso web:**
   - API y Health Check: `http://localhost:8080/health`
   - Interfaz web estática: `http://localhost:8080/static/`

## Tabla de Endpoints Disponibles

### 1. Sistema y Vistas Estáticas
| Método | Ruta | Código Éxito | Errores Posibles | Descripción |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/health` | `200 OK` | `500` | Chequeo de estado del servidor |
| `GET` | `/static/` | `200 OK` | `404` | Archivos estáticos de la interfaz web |
| `GET` | `/` | `200 OK` | `500` | Página de inicio |
| `GET` | `/users` | `200 OK` | `500` | Formulario de gestión de usuarios |

### 2. Emprendimientos (`/api/emprendimientos`)
| Método | Ruta | Código Éxito | Errores Posibles | Descripción |
| :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/emprendimientos` | `201 Created` | `400`, `500` | Crear nuevo emprendimiento |
| `GET` | `/api/emprendimientos` | `200 OK` | `500` | Listar todos los emprendimientos |
| `GET` | `/api/emprendimientos/{id}` | `200 OK` | `400`, `404`, `500` | Obtener emprendimiento por ID |
| `PUT` | `/api/emprendimientos/{id}` | `200 OK` | `400`, `404`, `500` | Actualizar datos de un emprendimiento |
| `DELETE` | `/api/emprendimientos/{id}` | `204 No Content` | `400`, `404`, `500` | Eliminar emprendimiento (borra en cascada publicaciones) |

### 3. Usuarios (`/api/usuarios`)
| Método | Ruta | Código Éxito | Errores Posibles | Descripción |
| :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/usuarios` | `201 Created` | `400`, `500` | Crear nuevo usuario |
| `GET` | `/api/usuarios` | `200 OK` | `500` | Listar todos los usuarios |
| `GET` | `/api/usuarios/{id}` | `200 OK` | `400`, `404`, `500` | Obtener usuario por ID |
| `PUT` | `/api/usuarios/{id}` | `200 OK` | `400`, `404`, `500` | Actualizar datos de un usuario |
| `DELETE` | `/api/usuarios/{id}` | `204 No Content` | `400`, `404`, `500` | Eliminar usuario |

### 4. Publicaciones (`/api/publicaciones`)
| Método | Ruta | Código Éxito | Errores Posibles | Descripción |
| :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/publicaciones` | `201 Created` | `400`, `500` | Crear publicación (`producto`, `promocion`, `servicio`, `otro`) |
| `GET` | `/api/publicaciones/{id}` | `200 OK` | `400`, `404`, `500` | Obtener publicación por ID |
| `GET` | `/api/emprendimientos/{id}/publicaciones` | `200 OK` | `400`, `500` | Listar publicaciones de un emprendimiento |
| `PUT` | `/api/publicaciones/{id}` | `200 OK` | `400`, `404`, `500` | Actualizar datos de una publicación |
| `DELETE` | `/api/publicaciones/{id}` | `204 No Content` | `400`, `404`, `500` | Eliminar publicación por ID |

### 5. Suscripciones (`/api/suscripciones` - Relación N:M)
| Método | Ruta | Código Éxito | Errores Posibles | Descripción |
| :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/suscripciones` | `201 Created` | `400`, `409`, `500` | Suscribir usuario a un emprendimiento |
| `GET` | `/api/usuarios/{id}/suscripciones` | `200 OK` | `400`, `500` | Listar suscripciones activas de un usuario |
| `DELETE` | `/api/suscripciones` | `204 No Content` | `400`, `404`, `500` | Cancelar/eliminar una suscripción existente |

---

## Pruebas de la Aplicación

### 1. Pruebas Unitarias y de Integración (Go Testing)
El proyecto incluye un pipeline automatizado en el `Makefile` que genera código con `sqlc`, levanta la base de datos de test temporal, aplica migraciones y corre la suite completa de Go:

```bash
make test
```

### 2. Pruebas End-to-End (`requests.sh`)
Para validar el funcionamiento completo de la API de punta a punta con peticiones HTTP reales:

1. Asegúrate de tener el servidor corriendo en `http://localhost:8080` (vía `make run` o Docker).
2. Ejecuta la suite de pruebas mediante cualquiera de las siguientes opciones:
   ```bash
   # Opción 1: A través del Makefile
   make e2e

   # Opción 2: Ejecutando directamente el script
   chmod +x requests.sh
   ./requests.sh
   ```

El script valida de manera secuencial:
1. **Creación de registros (201 Created):** Emprendimientos, usuarios, publicaciones y suscripciones con captura dinámica de IDs.
2. **Consultas y listados (200 OK):** Verificación de lectura de todos los recursos y colecciones.
3. **Modificaciones (200 OK):** Actualización de datos mediante peticiones `PUT`.
4. **Casos de error (400 Bad Request / 404 Not Found):** Validaciones de datos y consultas sobre IDs inexistentes.
5. **Eliminaciones (204 No Content y verificación 404 posterior):** Bajas de recursos y verificación estricta de que ya no se encuentran disponibles.

## Ejecución del Servidor Web, CORRESPONDIENTE AL TP1

**Iniciar el servidor**
Ejecuta el archivo principal de Go para levantar el servidor web:
`make run` 

**Visualizar la página del TP1**
Abre tu navegador web de preferencia y dirígete a la siguiente dirección:
`http://localhost:8080/static/` para conocer mas sobre las entidades y sus atributos.



---
*Proyecto incremental desarrollando una web completa, para la materia de Programación Web de la carrera de Ingenieria en Sistemas, UNICEN.*