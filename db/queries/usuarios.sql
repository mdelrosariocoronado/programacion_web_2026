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