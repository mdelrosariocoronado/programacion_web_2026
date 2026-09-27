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