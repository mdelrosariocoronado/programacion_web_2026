
-- name: CreateEmprendimiento :one
INSERT INTO emprendimientos (nombre, url, rubro, descripcion, logo, contacto, activo)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id_emprendimiento, nombre, url, rubro, descripcion, logo, contacto, activo, created_at;

-- name: GetEmprendimiento :one
SELECT id_emprendimiento, nombre, url, rubro, descripcion, logo, contacto, activo, created_at
FROM emprendimientos
WHERE id_emprendimiento = $1;

-- name: ListEmprendimientos :many
SELECT id_emprendimiento, nombre, url, rubro, descripcion, logo, contacto, activo, created_at
FROM emprendimientos
ORDER BY nombre;

-- name: UpdateEmprendimiento :exec
UPDATE emprendimientos
SET nombre = $2, url = $3, rubro = $4, descripcion = $5, logo = $6, contacto = $7, activo = $8
WHERE id_emprendimiento = $1;

-- name: DeleteEmprendimiento :exec
DELETE FROM emprendimientos
WHERE id_emprendimiento = $1;

-- name: CreateUsuario :one
INSERT INTO usuarios (nombre_completo, email, clave, rol, id_emprendimiento)
VALUES ($1, $2, $3, $4, $5)
RETURNING id_usuario, nombre_completo, email, clave, rol, id_emprendimiento, created_at;

-- name: GetUsuario :one
SELECT id_usuario, nombre_completo, email, clave, rol, id_emprendimiento, created_at
FROM usuarios
WHERE id_usuario = $1;

-- name: ListUsuarios :many
SELECT id_usuario, nombre_completo, email, clave, rol, id_emprendimiento, created_at
FROM usuarios;

-- name: UpdateUsuario :exec
UPDATE usuarios
SET nombre_completo = $2, email = $3, clave = $4, rol = $5, id_emprendimiento = $6
WHERE id_usuario = $1;

-- name: DeleteUsuario :exec
DELETE FROM usuarios
WHERE id_usuario = $1;

-- name: CreatePublicacion :one
INSERT INTO publicaciones (id_emprendimiento, titulo, contenido, imagen_url, tipo, precio)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id_publicacion, id_emprendimiento, titulo, contenido, imagen_url, tipo, precio, created_at;

-- name: GetPublicacion :one
SELECT id_publicacion, id_emprendimiento, titulo, contenido, imagen_url, tipo, precio, created_at
FROM publicaciones
WHERE id_publicacion = $1;

-- name: ListPublicacionesByEmprendimiento :many
SELECT id_publicacion, id_emprendimiento, titulo, contenido, imagen_url, tipo, precio, created_at
FROM publicaciones
WHERE id_emprendimiento = $1
ORDER BY created_at DESC;

-- name: DeletePublicacion :exec
DELETE FROM publicaciones
WHERE id_publicacion = $1;

-- name: CreateSuscripcion :one
INSERT INTO suscripciones (id_usuario, id_emprendimiento)
VALUES ($1, $2)
RETURNING id_suscripcion, id_usuario, id_emprendimiento, fecha_sub;

-- name: ListSuscripcionesByUsuario :many
SELECT id_suscripcion, id_usuario, id_emprendimiento, fecha_sub
FROM suscripciones
WHERE id_usuario = $1;

-- name: DeleteSuscripcion :exec
DELETE FROM suscripciones
WHERE id_usuario = $1 AND id_emprendimiento = $2;