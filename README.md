# Trabajo practico N°3

## Integrantes del grupo
*   María Del Rosario Coronado
*   Brisa Fernandez Lema
*   Benjamín Knudsen

---

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

Antes de ejecutar este proyecto, asegúrate de contar con las siguientes herramientas instaladas:

* **[Git](https://git-scm.com/downloads)**
* **[Docker](https://docs.docker.com/get-docker/) & Docker Compose** (en ejecución)
* **[Go](https://go.dev/dl/)** (v1.22 o superior)
* **cURL** (para la ejecución de las pruebas HTTP)
* **[Atlas CLI](https://atlasgo.io/getting-started/)** (para la aplicación de migraciones DDL)
* **[sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)** (para la generación tipada de repositorios)

---

## Inicio Rápido: Clonar y Ejecutar Pruebas (TP 3)

Para evaluar el proyecto de punta a punta en un entorno limpio, sigue estos pasos:

**1. Clonar el repositorio y acceder a la rama `tp3`:**
```bash
git clone https://github.com/mdelrosariocoronado/prog_web.git
cd prog_web
git checkout tp3
```

**2. Configurar las variables de entorno:**
```bash
cp .env.example .env
```

**3. Ejecutar la suite completa de pruebas de la API:**
```bash
make test-api
```
*(O de forma equivalente ejecutando directamente el script: `./test-api.sh`)*

Este único comando realiza automáticamente todo el ciclo requerido:
* **Tareas previas:** Generación de código con `sqlc`, validación de compilación, limpieza de contenedores y volúmenes previos, levantamiento de la base de datos en Docker, espera hasta que PostgreSQL esté listo, aplicación de migraciones DDL con Atlas, y levantamiento de la aplicación web esperando a que esté disponible.
* **Pruebas de la API:** Ejecuta la suite completa de pruebas end-to-end con `curl` (`requests.sh`) cubriendo cada uno de los endpoints.
* **Tareas posteriores:** Limpieza y borrado de todos los contenedores y volúmenes temporales creados.

---

## Guía de Ejecución del Servidor

Si deseas levantar el servidor para interactuar manualmente con la API o las vistas web:

### Opción 1: Ejecución completa en Docker (Recomendada)
Levanta tanto PostgreSQL como el servidor Go en contenedores aislados:
```bash
# Construye, levanta los servicios y aplica migraciones:
make docker-up

# Para detener todos los contenedores:
make docker-down
```
El servidor estará disponible en `http://localhost:8080`.

### Opción 2: Ejecución Local Nativa (Servidor Go local + BD en Docker)
Ejecuta la base de datos en Docker y el binario Go directamente en tu máquina:
```bash
# Levanta la BD en Docker, aplica migraciones e inicia el servidor Go:
make run

# Para detener todos los contenedores:
make docker-down
```

### Vistas Web del Proyecto (TP 1)
Con el servidor en ejecución, podés acceder desde cualquier navegador web a:
* **Catálogo e información estática:** `http://localhost:8080/static/` (presentación de entidades y atributos del TP1).
* **Formulario de gestión:** `http://localhost:8080/users`.
* **Chequeo de salud del servicio:** `http://localhost:8080/health`.

---

## Tabla de Endpoints de la API

### Sistema y Vistas Estáticas
| Método | Ruta | Código Éxito | Errores Posibles | Descripción |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/health` | `200 OK` | `500` | Chequeo de estado del servidor |
| `GET` | `/static/` | `200 OK` | `404` | Archivos estáticos de la interfaz web |
| `GET` | `/` | `200 OK` | `500` | Página de inicio |
| `GET` | `/users` | `200 OK` | `500` | Formulario de gestión de usuarios |

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