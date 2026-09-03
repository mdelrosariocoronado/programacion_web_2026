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