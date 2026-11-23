# Portal de Emprendimientos Tandil — Trabajo Práctico N°6 (Entrega Final: Interfaces Dinámicas con HTMX)

## Integrantes del grupo
*   María Del Rosario Coronado
*   Brisa Fernandez Lema
*   Benjamín Knudsen

---

## Documentación del Proyecto

Aplicación web desarrollada en Go que modela una red centralizada de emprendimientos locales en Tandil, conectando emprendedores y clientes mediante publicaciones y suscripciones.

El sistema implementa el **modelo de tres capas** (**Capa de Presentación SSR + HTMX con `templ` $\leftrightarrow$ Capa de Lógica de Negocio (`handlers` $\rightarrow$ `services`) $\rightarrow$ Capa de Acceso a Datos (`repositories` + `sqlc` $\rightarrow$ PostgreSQL)**), ofreciendo una experiencia interactiva y fluida sin recargas completas de página y sin depender de JavaScript en el cliente.

### Estructura del Proyecto
```text
.
├── cmd/server/main.go          # Punto de entrada HTTP, inyección de dependencias y enrutamiento
├── db/
│   ├── migrations/             # Migraciones DDL versionadas con Atlas
│   ├── queries/                # Consultas SQL puras para sqlc
│   ├── schema/                 # Definición del esquema relacional PostgreSQL
│   └── sqlc/                   # Código Go fuertemente tipado generado por sqlc
├── internal/
│   ├── handlers/               # Handlers HTTP (vistas SSR/HTMX y endpoints REST JSON) + tests
│   ├── repositories/           # Capa de acceso a datos desacoplada + tests de integración
│   └── services/               # Lógica de negocio y validaciones de dominio + tests unitarios
├── static/                     # Archivos estáticos (styles.css, tipografías e imágenes)
├── views/                      # Componentes HTML tipados en Templ (*.templ y *_templ.go) + tests
├── docker-compose.yml          # Orquestación de PostgreSQL 16 y aplicación
├── Makefile                    # Automatización (run-app, test, generate, test-api)
└── run-app.sh                  # Script ejecutable equivalente a make run-app
```

### Modelo de Datos
* **`emprendimientos`:** Perfil del negocio, rubro, descripción, URL, logo, contacto y estado.
* **`usuarios`:** Clientes, emprendedores y administradores con validación de email único, clave y rol (`administrador`, `emprendedor`, `cliente`), y relación opcional a un emprendimiento (`id_emprendimiento`).
* **`publicaciones`:** Contenido publicado por cada emprendimiento (`producto`, `promocion`, `servicio` u `otro`), con precio opcional y eliminación en cascada (`ON DELETE CASCADE`) al borrar su emprendimiento.
* **`suscripciones`:** Tabla intermedia que materializa la relación N:M entre usuarios y los emprendimientos que siguen.

### Arquitectura, Herramientas y Patrones Aplicados
* **Interfaces Reactivas con HTMX (TP 6):** En lugar de recargar la página completa (PRG del TP 5) o mantener estado/JSON en JavaScript cliente (TP 4), los formularios y botones interactúan con el servidor mediante atributos declarativos HTML:
  * **Alta sin recarga:** Al enviar un formulario, el servidor persiste la entidad mediante la capa `services`, vuelve a consultar el listado actualizado y devuelve únicamente el fragmento HTML del componente de lista (`EmprendimientoList`, `UsuarioList`, `PublicacionList`).
  * **Limpieza automática de formularios (*Out-of-Band Swap*):** En la misma respuesta `POST`, el servidor envía un segundo fragmento con `hx-swap-oob="true"` (`EmprendimientoFormOOB`, `UsuarioFormOOB`, `PublicacionFormOOB`) para reemplazar el formulario por uno limpio sin recargar la vista ni usar JS.
  * **Borrado sin recarga:** Cada fila de las tablas incluye un botón de borrado con diálogo de confirmación. El servidor elimina el registro y responde `200 OK` con cuerpo vacío, indicando a HTMX que remueva la fila del DOM instantáneamente.
* **Server-Side Rendering (SSR) con `templ` (TP 5):** Componentes HTML compilados a código Go nativo en `views/` (`layout.templ`, `home.templ`, `emprendimientos.templ`, `usuarios.templ`, `publicaciones.templ`) con chequeo de tipos en compilación y escapado automático contra ataques XSS.
* **Arquitectura en Capas y API REST (TP 3):** Separación estricta de responsabilidades (`handlers` $\rightarrow$ `services` $\rightarrow$ `repositories`), manteniendo además los endpoints REST `/api/...` y sus pruebas E2E.
* **SQL Directo tipado con `sqlc` y Migraciones con Atlas (TP 2):** Consultas SQL compiladas a structs y métodos Go seguros, más versionado reproducible del esquema en PostgreSQL.
* **Contenedores con Docker Compose:** Entorno reproducible y automatizado para la base de datos y las pruebas.

## Requisitos Previos

Antes de ejecutar este proyecto, asegurate de contar con las siguientes herramientas instaladas:

* **[Git](https://git-scm.com/downloads)**
* **[Docker](https://docs.docker.com/get-docker/) & Docker Compose** (en ejecución)
* **[Go](https://go.dev/dl/)** (v1.26 o superior, compatible con `templ`)
* **[Atlas CLI](https://atlasgo.io/getting-started/)** (para la aplicación de migraciones DDL)
* **[sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)** *(opcional: si no está en el PATH, `make run-app` utiliza el código pregenerado en `db/sqlc/`)*
* **[templ](https://templ.guide/)** *(se ejecuta automáticamente vía `go run github.com/a-h/templ/cmd/templ generate`)*

---

## Inicio Rápido: Clonar y Ejecutar la Aplicación (TP 6)

Para evaluar la aplicación de punta a punta en un entorno limpio:

**1. Clonar el repositorio y acceder a la rama `tp6`:**
```bash
git clone https://github.com/mdelrosariocoronado/prog_web.git
cd prog_web
git checkout tp6
```

**2. Ejecutar la aplicación con el ciclo completo automatizado:**
```bash
make run-app
```
*(O de forma equivalente ejecutando directamente el script bash: `./run-app.sh`)*

Este comando ejecuta automáticamente las tres etapas requeridas por la cátedra:
1. **Tareas previas:**
   * Ejecución de `sqlc generate` (si está disponible) y compilación de componentes de vista con `templ generate`.
   * Compilación del binario del servidor Go (`go build -o bin/server ./cmd/server`).
   * Limpieza de contenedores y volúmenes previos (`docker compose down -v --remove-orphans`).
   * Levantamiento del contenedor de base de datos PostgreSQL (`docker compose up -d db`) y espera activa hasta que acepte conexiones (`pg_isready`).
   * Aplicación de migraciones DDL con Atlas.
2. **Ejecución de la aplicación:**
   * Inicia el servidor HTTP en **`http://localhost:8080`**.
   * Desde el navegador se puede recorrer el portal en `/`, `/emprendimientos`, `/usuarios` y `/publicaciones`, creando y eliminando registros sin recargas completas de página.
3. **Tareas posteriores (al detener con `Ctrl+C`):**
   * Captura la señal de cierre (`EXIT`, `INT`, `TERM`) y elimina automáticamente todos los contenedores, volúmenes y binarios temporales creados.

---

## Inicio Rápido: Pruebas Automatizadas de la API (TP 3)

Para ejecutar la suite automatizada de pruebas end-to-end sobre los endpoints REST:

```bash
make test-api
```
*(O de forma equivalente ejecutando directamente el script: `./test-api.sh`)*

---

## Guía de Comandos del `Makefile`

* **`make run-app`** *(o `./run-app.sh`)*: Ciclo completo de la entrega (tareas previas + servidor en `http://localhost:8080` + limpieza automática de contenedores y volúmenes al detener).
* **`make test`:** Genera código (`sqlc`/`templ`), levanta PostgreSQL en Docker, aplica migraciones, ejecuta toda la suite de tests en Go (`views`, `handlers`, `services`, `repositories`) y limpia contenedores/volúmenes al finalizar.
* **`make generate`:** Regenera el código Go de `sqlc` (`db/sqlc/`) y de `templ` (`views/*_templ.go`).
* **`make run`:** Genera las vistas `templ`, levanta PostgreSQL en Docker, aplica migraciones e inicia el servidor localmente para desarrollo.

---

## Tabla de Rutas y Endpoints

### Vistas Dinámicas SSR + HTMX con `templ` (TP 5 / TP 6)
| Método | Ruta | Respuesta | Descripción |
| :--- | :--- | :--- | :--- |
| `GET` | `/` | `200 OK` (HTML) | Página de inicio del Portal de Emprendimientos con accesos a los 3 módulos |
| `GET` | `/emprendimientos` | `200 OK` (HTML) | Vista SSR con formulario HTMX de alta y tabla de emprendimientos |
| `POST` | `/emprendimientos` | `200 OK` (HTML) | Crea el emprendimiento vía HTMX y devuelve la lista actualizada + formulario limpio (OOB swap) |
| `DELETE` | `/emprendimientos/{id}` | `200 OK` (vacío) | Elimina el emprendimiento vía HTMX y remueve la fila `<tr>` del DOM sin recargar |
| `GET` | `/usuarios` | `200 OK` (HTML) | Vista SSR con formulario HTMX de alta y tabla de usuarios |
| `POST` | `/usuarios` | `200 OK` (HTML) | Crea el usuario vía HTMX y devuelve la lista actualizada + formulario limpio (OOB swap) |
| `DELETE` | `/usuarios/{id}` | `200 OK` (vacío) | Elimina el usuario vía HTMX y remueve la fila `<tr>` del DOM sin recargar |
| `GET` | `/publicaciones` | `200 OK` (HTML) | Vista SSR con formulario HTMX y listado de publicaciones (admite filtro `?id_emprendimiento={id}`) |
| `POST` | `/publicaciones` | `200 OK` (HTML) | Crea la publicación vía HTMX y devuelve la lista actualizada + formulario limpio (OOB swap) |
| `DELETE` | `/publicaciones/{id}` | `200 OK` (vacío) | Elimina la publicación vía HTMX y remueve la fila `<tr>` del DOM sin recargar |
| `GET` | `/static/` | `200 OK` | Servidor de archivos estáticos (`styles.css`, tipografías e imágenes) |
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