-- schema.sql
CREATE TABLE albums (
	id_album SERIAL PRIMARY KEY,
	id_usuario INT NOT NULL,
	titulo VARCHAR(255) NOT NULL,
	artista VARCHAR(255) NOT NULL,
	descripcion VARCHAR(255) NOT NULL,
	genero VARCHAR(255) NOT NULL,
	duracion_s INT NOT NULL,
	ranking INT CHECK (ranking >= 1 AND ranking <= 5) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);