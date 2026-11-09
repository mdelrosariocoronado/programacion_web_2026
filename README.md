# Trabajo Práctico N°5 - Vistas Dinámicas con Templ y SQLc (SSR)

## Integrantes del grupo
*   María Del Rosario Coronado
*   Brisa Fernandez Lema
*   Benjamín Knudsen

---

## Documentación del Proyecto

Aplicación web desarrollada en Go que modela una red de emprendimientos locales en Tandil. Implementa una arquitectura en tres capas (**Presentación SSR con `templ` $\leftrightarrow$ Handlers $\rightarrow$ Services $\rightarrow$ Repositories (`sqlc`) $\rightarrow$ PostgreSQL**), renderizando las vistas HTML dinámicamente desde el servidor sin depender de JavaScript en el cliente.

### Modelo de Datos
* **`emprendimientos`:** Perfil del negocio, rubro y datos de contacto.
* **`usuarios`:** Clientes y emprendedores con control de acceso por roles (`administrador`, `emprendedor`, `cliente`).
* **`publicaciones`:** Contenido publicado por cada emprendimiento (`producto`, `promocion`, `servicio` u `otro`), vinculado mediante clave foránea con eliminación en cascada (`ON DELETE CASCADE`).
* **`suscripciones`:** Tabla intermedia que materializa la relación N:M entre usuarios y los emprendimientos que siguen.

### Herramientas y Patrones
* **Server-Side Rendering (SSR) con `templ`:** La capa de presentación se define en componentes tipados dentro de `views/` (`layout.templ`, `emprendimientos.templ`, `usuarios.templ`, `publicaciones.templ`) que se compilan a código Go nativo (`*_templ.go`), brindando seguridad de tipos en compilación y escapado automático contra ataques XSS.
* **Patrón Post/Redirect/Get (PRG):** El alta de entidades se realiza mediante formularios HTML clásicos (`method="POST"`). El servidor procesa los datos con `r.ParseForm()`, valida y persiste mediante la capa de servicios y responde con una redirección HTTP (`303 See Other`) hacia la ruta `GET`, evitando reenvíos accidentales al recargar la página.
* **SQL Directo con `sqlc`:** Compilación de consultas SQL puras a código Go seguro y fuertemente tipado.
* **Migraciones con Atlas:** Versionado reproducible del esquema relacional en `db/migrations/`.
* **Entorno con Docker Compose:** Ejecución del motor PostgreSQL y la aplicación en contenedores reproducibles.

## Requisitos Previos

Antes de ejecutar este proyecto, asegúrate de contar con las siguientes herramientas instaladas:

* **[Git](https://git-scm.com/downloads)**
* **[Docker](https://docs.docker.com/get-docker/) & Docker Compose** (en ejecución)
* **[Go](https://go.dev/dl/)** (v1.22 o superior / compatible con `templ` (usamos v1.26 para este tp))
* **[Atlas CLI](https://atlasgo.io/getting-started/)** (para la aplicación de migraciones DDL)
* **[sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)** (para la generación tipada de repositorios)
* **[templ](https://templ.guide/)** (se ejecuta automáticamente vía `go run github.com/a-h/templ/cmd/templ@latest generate`)

---

## Inicio Rápido: Clonar y Ejecutar la Aplicación (TP 5)

Para evaluar la aplicación renderizada en el servidor de punta a punta en un entorno limpio:

**1. Clonar el repositorio y acceder a la rama `tp5`:**
```bash
git clone https://github.com/mdelrosariocoronado/prog_web.git
cd prog_web
git checkout tp5
```

**2. Ejecutar la aplicación con el ciclo completo automatizado:**
```bash
make run-app
```

Este comando ejecuta automáticamente las tres etapas requeridas por la cátedra:
1. **Tareas previas:**
   * Ejecución de `sqlc generate` y compilación de componentes de vista con `templ generate`.
   * Compilación del servidor Go (`go build`).
   * Limpieza de contenedores y volúmenes previos (`docker compose down -v --remove-orphans`).
   * Levantamiento del contenedor de base de datos PostgreSQL (`docker compose up -d db`) y espera activa hasta que acepte conexiones (`pg_isready`).
   * Aplicación de migraciones DDL con Atlas.
2. **Ejecución de la aplicación:**
   * Inicia el servidor HTTP en **`http://localhost:8080`**.
3. **Tareas posteriores (al detener con `Ctrl+C`):**
   * Captura la señal de cierre y elimina automáticamente los contenedores, volúmenes y binarios temporales creados.

---

## Inicio Rápido: Pruebas Automatizadas de la API (TP 3)

Para ejecutar la suite automatizada de pruebas end-to-end sobre los endpoints REST:

```bash
make test-api
```
*(O de forma equivalente ejecutando directamente el script: `./test-api.sh`)*

---

## Guía de Otros Comandos del `Makefile`

* **`make run-app`:** Ciclo completo del TP5 (tareas previas + servidor en `http://localhost:8080` + limpieza posterior de contenedores y volúmenes al salir).
* **`make generate`:** Regenera el código Go de `sqlc` (`db/sqlc/`) y de `templ` (`views/*_templ.go`).
* **`make run`:** Genera las vistas `templ`, levanta PostgreSQL en Docker, aplica migraciones e inicia el servidor localmente.
* **`make test`:** Ejecuta los tests unitarios y de integración en Go levantando y limpiando la BD automáticamente.

---

## Tabla de Rutas y Endpoints

### Vistas Dinámicas SSR con `templ` (TP 5)
| Método | Ruta | Respuesta | Descripción |
| :--- | :--- | :--- | :--- |
| `GET` | `/` | `200 OK` (HTML) | Página principal renderizada en el servidor (Gestión de Emprendimientos) |
| `GET` | `/emprendimientos` | `200 OK` (HTML) | Vista SSR con formulario de alta y tabla de emprendimientos |
| `POST` | `/emprendimientos` | `303 See Other` | Procesa el formulario de nuevo emprendimiento y redirige a `/emprendimientos` (PRG) |
| `GET` | `/usuarios` | `200 OK` (HTML) | Vista SSR con formulario de alta y tabla de usuarios |
| `POST` | `/usuarios` | `303 See Other` | Procesa el formulario de nuevo usuario y redirige a `/usuarios` (PRG) |
| `GET` | `/publicaciones` | `200 OK` (HTML) | Vista SSR con formulario de alta y listado de publicaciones |
| `POST` | `/publicaciones` | `303 See Other` | Procesa el formulario de nueva publicación y redirige a `/publicaciones` (PRG) |
| `GET` | `/static/` | `200 OK` | Servidor de archivos estáticos (hoja de estilos `style.css` e imágenes) |
| `GET` | `/health` | `200 OK` | Chequeo de estado del servidor |

### Emprendimientos (`/api/emprendimientos`)
| Método | Ruta | Código Éxito | Errores Posibles | Descripción |
| :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/emprendimientos` | `201 Created` | `400`, `500` | Crear nuevo emprendimiento |
| `GET` | `/api/emprendimientos` | `200 OK` | `500` | Listar todos los emprendimientos |
| `GET` | `/api/emprendimientos/{id}` | `200 OK` | `400`, `404`, `500` | Obtener emprendimiento por ID |
| `PUT` | `/api/emprendimientos/{id}` | `200 OK` | `400`, `404`, `500` | Actualizar datos de un emprendimiento |
| `DELETE` | `/api/emprendimientos/{id}` | `204 No Content` | `400`, `404`, `500` | Eliminar emprendimiento (borra en cascada publicaciones) |

### Usuarios (`/api/usuarios`)
| Método | Ruta | Código Éxito | Errores Posibles | Descripción |
| :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/usuarios` | `201 Created` | `400`, `500` | Crear nuevo usuario |
| `GET` | `/api/usuarios` | `200 OK` | `500` | Listar todos los usuarios |
| `GET` | `/api/usuarios/{id}` | `200 OK` | `400`, `404`, `500` | Obtener usuario por ID |
| `PUT` | `/api/usuarios/{id}` | `200 OK` | `400`, `404`, `500` | Actualizar datos de un usuario |
| `DELETE` | `/api/usuarios/{id}` | `204 No Content` | `400`, `404`, `500` | Eliminar usuario |

### Publicaciones (`/api/publicaciones`)
| Método | Ruta | Código Éxito | Errores Posibles | Descripción |
| :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/publicaciones` | `201 Created` | `400`, `500` | Crear publicación (`producto`, `promocion`, `servicio`, `otro`) |
| `GET` | `/api/publicaciones/{id}` | `200 OK` | `400`, `404`, `500` | Obtener publicación por ID |
| `GET` | `/api/emprendimientos/{id}/publicaciones` | `200 OK` | `400`, `500` | Listar publicaciones de un emprendimiento |
| `PUT` | `/api/publicaciones/{id}` | `200 OK` | `400`, `404`, `500` | Actualizar datos de una publicación |
| `DELETE` | `/api/publicaciones/{id}` | `204 No Content` | `400`, `404`, `500` | Eliminar publicación por ID |

### Suscripciones (`/api/suscripciones` - Relación N:M)
| Método | Ruta | Código Éxito | Errores Posibles | Descripción |
| :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/suscripciones` | `201 Created` | `400`, `409`, `500` | Suscribir usuario a un emprendimiento |
| `GET` | `/api/usuarios/{id}/suscripciones` | `200 OK` | `400`, `500` | Listar suscripciones activas de un usuario |
| `DELETE` | `/api/suscripciones` | `204 No Content` | `400`, `404`, `500` | Cancelar/eliminar una suscripción existente |

---

## Pruebas de la Aplicación

### 1. Pruebas End-to-End de la API con cURL (TP 3)

La suite de pruebas automatizada en [`requests.sh`](requests.sh) evalúa de manera exhaustiva el ciclo de vida completo de cada recurso de la API mediante peticiones HTTP reales:

1. **Creación de registros (`201 Created`):** Alta de emprendimientos, usuarios, publicaciones y suscripciones, capturando dinámicamente los IDs generados.
2. **Consultas y listados (`200 OK`):** Lectura individual por ID y listados completos de recursos y colecciones relacionadas.
3. **Modificaciones (`200 OK`):** Actualización de atributos de recursos existentes vía `PUT`.
4. **Casos de error (`400 Bad Request` / `404 Not Found`):** Validación de datos inválidos y consultas sobre identificadores inexistentes.
5. **Eliminaciones (`204 No Content` y verificación posterior `404`):** Bajas lógicas/físicas y confirmación estricta de inexistencia posterior.

**Modos de ejecución:**
* **Ciclo completo automatizado (Recomendado para evaluación):**
  ```bash
  make test-api
  # o directamente:
  ./test-api.sh
  ```
* **Con el servidor ya en ejecución:**
  ```bash
  make e2e
  # o directamente:
  chmod +x requests.sh
  ./requests.sh
  ```

### 2. Pruebas Unitarias y de Integración en Go (TP 2)

Pipeline automatizado que genera código con `sqlc`, levanta una base de datos PostgreSQL temporal de pruebas, aplica el esquema DDL y corre la suite interna de tests en Go:

```bash
make test
```

---
*Proyecto incremental desarrollando una web completa, para la materia de Programación Web de la carrera de Ingeniería en Sistemas, UNICEN.*