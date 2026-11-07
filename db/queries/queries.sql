-- name: CreateRegister :one
INSERT INTO albums (id_album, id_usuario, titulo, artista, descripcion, genero, duracion_min, rating) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
RETURNING id_album, id_usuario, titulo, artista, descripcion, genero, duracion_min, rating, created_at;

-- name: GetRegister :one
SELECT id_registro, id_album, id_usuario, titulo, artista, descripcion, genero, duracion_min, rating, created_at 
FROM albums 
WHERE id_registro = $1;

-- name: ListRegisters :many
SELECT id_registro, id_album, id_usuario, titulo, artista, descripcion, genero, duracion_min, rating, created_at 
FROM albums
ORDER BY id_registro;

-- name: UpdateRegister :exec
UPDATE albums
SET id_album = $2, id_usuario = $3, titulo = $4, artista = $5, descripcion = $6, genero = $7, duracion_min = $8, rating = $9
WHERE id_registro = $1;

-- name: DeleteRegister :exec
DELETE FROM albums 
WHERE id_registro = $1;