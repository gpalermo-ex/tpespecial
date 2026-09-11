-- name: CreateAlbum :one
INSERT INTO albums (id_usuario, titulo, artista, descripcion, genero, duracion_s, ranking) 
VALUES ($1, $2, $3, $4, $5, $6, $7) 
RETURNING id_album, id_usuario, titulo, artista, descripcion, genero, duracion_s, ranking, created_at;

-- name: GetAlbum :one
SELECT id_album, id_usuario, titulo, artista, descripcion, genero, duracion_s, ranking, created_at 
FROM albums 
WHERE id_album = $1;

-- name: ListAlbums :many
SELECT id_album, id_usuario, titulo, artista, descripcion, genero, duracion_s, ranking, created_at 
FROM albums 
ORDER BY id_album;

-- name: UpdateAlbum :exec
UPDATE albums 
SET id_usuario = $2, titulo = $3, artista = $4, descripcion = $5, genero = $6, duracion_s = $7, ranking = $8
WHERE id_album = $1;

-- name: DeleteAlbum :exec
DELETE FROM albums 
WHERE id_album = $1;