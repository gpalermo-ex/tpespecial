
-- schema.sql
CREATE TABLE albums (
	id_registro SERIAL PRIMARY KEY,
	id_album INT NOT NULL,
	id_usuario VARCHAR(255) NOT NULL,
	titulo VARCHAR(255) NOT NULL,
	artista VARCHAR(255) NOT NULL,
	descripcion VARCHAR(255) NOT NULL,
	genero VARCHAR(255) NOT NULL,
	duracion_min INT NOT NULL,
	rating INT CHECK (rating >= 1 AND rating <= 5) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);