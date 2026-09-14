-- name: CreateSuscripcion :one
INSERT INTO suscripciones (id_usuario, id_emprendimiento)
VALUES ($1, $2)
RETURNING id_suscripcion, id_usuario, id_emprendimiento, fecha_sub;

-- name: GetSuscripcionByUserAndEmprendimiento :one
SELECT id_suscripcion, id_usuario, id_emprendimiento, fecha_sub
FROM suscripciones
WHERE id_usuario = $1 AND id_emprendimiento = $2;

-- name: ListSuscripcionesByUsuario :many
SELECT id_suscripcion, id_usuario, id_emprendimiento, fecha_sub
FROM suscripciones
WHERE id_usuario = $1;

-- name: UpdateSuscripcion :exec
UPDATE suscripciones
SET fecha_sub = $3
WHERE id_usuario = $1 AND id_emprendimiento = $2;
-- name: DeleteSuscripcion :exec
DELETE FROM suscripciones
WHERE id_usuario = $1 AND id_emprendimiento = $2;