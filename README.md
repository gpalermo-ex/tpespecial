# Catálogo Musical

Aplicación web para llevar un catálogo personal de álbumes y registrar una valoración para cada uno. Desde la página principal se pueden consultar los registros existentes y cargar nuevas valoraciones mediante un formulario.

Cada registro incluye el identificador del álbum, el usuario, el título, el artista, una descripción, el género, la duración en minutos y una calificación de 1 a 5. La fecha de creación se asigna automáticamente en la base de datos.

## Tecnologías

- Go para el servidor web y la lógica de la aplicación.
- PostgreSQL para almacenar los registros.
- Docker Compose para iniciar la base de datos.
- SQLC para generar código Go tipado a partir de las consultas SQL.
- Templ para generar las vistas HTML desde componentes escritos en Go.
- pgx como driver de PostgreSQL.
- Make para automatizar generación, ejecución y pruebas.

## Organización

- `main.go`: conecta con PostgreSQL, configura las rutas web y procesa el formulario de alta.
- `views/`: componentes Templ para el layout, el formulario y el listado de registros.
- `db/schema/schema.sql`: definición de la tabla `albums`.
- `db/queries/queries.sql`: consultas para crear, listar, consultar, actualizar y eliminar registros.
- `db/sqlc/`: código Go generado por SQLC.
- `docker-compose.yml`: configuración del servicio PostgreSQL y carga inicial del esquema.
- `tests/`: pruebas del CRUD y solicitudes de prueba de la aplicación.
- `Makefile`: comandos automatizados para generar código, iniciar el entorno y ejecutar pruebas.

## Datos almacenados

La tabla `albums` registra una valoración asociada a un álbum y a un usuario. Incluye título, artista, descripción, género, duración en minutos, calificación entre 1 y 5 y fecha de creación. El identificador de registro es generado por PostgreSQL.

## Requisitos

- Go instalado y compatible con la versión indicada en `go.mod`.
- Docker y Docker Compose.
- Make.

## Ejecutar con Make

Desde la carpeta del proyecto, ejecutá:

```bash
make run
```

Este comando prepara el entorno en este orden:

1. Ejecuta `make clean`, que detiene los servicios y elimina los volúmenes de Docker Compose.
2. Inicia PostgreSQL y espera a que esté disponible.
3. Genera el código de SQLC y las vistas Templ.
4. Inicia el servidor Go en el puerto `8080`.

Cuando aparezca el mensaje de que el servidor está activo, abrí [http://localhost:8080](http://localhost:8080) en el navegador.

> **Atención:** `make run` elimina los volúmenes antes de iniciar la base de datos, por lo que borra los datos guardados en ejecuciones anteriores.

Para detener el servidor, usá `Ctrl+C`. Para detener y eliminar también los contenedores y volúmenes, ejecutá `make clean`.

## Comandos disponibles

| Comando | Función |
| --- | --- |
| `make` o `make all` | Ejecuta el objetivo predeterminado: `make run`. |
| `make run` | Limpia los volúmenes, inicia PostgreSQL, genera SQLC y Templ, y ejecuta la aplicación. |
| `make up` | Limpia los volúmenes e inicia PostgreSQL en segundo plano. |
| `make generate` | Regenera el código de SQLC y las vistas Templ. |
| `make test` | Prepara la base de datos, genera SQLC y Templ, ejecuta `go test -v ./tests/...` y limpia los contenedores y volúmenes al terminar. |
| `make clean` | Detiene los contenedores y elimina los volúmenes de Docker Compose. |

`make up` también depende de `make clean`, por lo tanto recrea la base de datos y elimina sus datos anteriores. `make test` realiza esa preparación antes de ejecutar las pruebas y garantiza la limpieza al finalizar, incluso si las pruebas fallan.




