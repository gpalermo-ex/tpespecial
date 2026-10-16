# Catálogo Musical API

Este proyecto es una API REST desarrollada en Go para gestionar un catálogo de álbumes musicales. Permite listar, crear, consultar, actualizar y eliminar registros de álbumes almacenados en una base de datos PostgreSQL.

La aplicación está pensada como un backend simple y educativo para practicar CRUD, conexión a base de datos, generación de queries con SQLC y creación de endpoints HTTP con Go.

## ¿Qué hace la aplicación?

La API permite administrar información sobre álbumes, incluyendo:

- identificador del álbum
- identificador del usuario
- título
- artista
- descripción
- género
- duración en segundos
- ranking (valor entre 1 y 5)
- fecha de creación

Con esto, el sistema puede servir como base para un catálogo musical, una app de recomendaciones o un backend para una interfaz web.

## Tecnologías utilizadas

- Go
- PostgreSQL
- Docker + Docker Compose
- SQLC para generar código de acceso a la base de datos
- Net/HTTP de Go para los endpoints REST
- pgx como driver de PostgreSQL

## Estructura del proyecto

- `main.go`: servidor HTTP principal y handlers de las rutas.
- `db/schema/schema.sql`: esquema de la tabla `albums`.
- `db/queries/queries.sql`: consultas SQL del CRUD.
- `db/sqlc/`: código generado automáticamente por SQLC.
- `static/index.html`: ejemplo de frontend estático para la vista del catálogo.
- `docker-compose.yml`: configuración del contenedor de PostgreSQL.
- `CRUD_test.go`: pruebas para validar el CRUD.
- `test-api.sh` y `test.sh`: scripts útiles para pruebas y ejecución.

## Base de datos

La base de datos se estructura con la siguiente tabla:

```sql
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
```

La tabla guarda cada álbum con su metadata principal y un ranking que limita la puntuación a valores válidos entre 1 y 5.

## Endpoints de la API

El servidor corre en `:8080` y expone los siguientes endpoints:

### Obtener todos los álbumes

- `GET /albums`

Respuesta: JSON con la lista de álbumes.

### Crear un álbum

- `POST /albums`

### Obtener un álbum por id

- `GET /albums/:id`

### Actualizar un álbum

- `PUT /albums/:id`

### Eliminar un álbum

- `DELETE /albums/:id`

## Cómo ejecutar el proyecto

### 1. Levantar la base de datos

```bash
docker compose up -d
```

Esto levanta un contenedor PostgreSQL con la configuración definida en `docker-compose.yml`.

### 2. Ejecutar la API

```bash
go run .
```

La aplicación se inicia en el puerto `8080`.

### 3. Probar los endpoints

Ejemplo con curl:

```bash
curl http://localhost:8080/albums
```

```bash
curl -X POST http://localhost:8080/albums \
  -H "Content-Type: application/json" \
    -d '{
        "id_usuario": 1,
        "titulo": "Detour",
        "artista": "Kim Petras",
        "descripcion": "Álbum de música pop",
        "genero": "Pop",
        "duracion_s": 2203,
        "ranking": 5
    }'
```

## Tests

El proyecto incluye pruebas del CRUD en `CRUD_test.go`, donde se valida la creación, lectura, listado, actualización y eliminación de álbumes usando las funciones generadas por sqlc o sino `test-api.sh` donde se valida todo esto pero desde la API REST utilizando la herramienta HURL.

Para ejecutarlos:

```bash
# Para testear el CRUD_test.go
chmod +x test.sh
./test.sh

# Para testear API REST
chmod +x test-api.sh
./test-api.sh
```



