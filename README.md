# Catálogo Musical

Aplicación web en Go para registrar y consultar valoraciones de álbumes. La página muestra los registros guardados en PostgreSQL y permite crear y eliminar registros desde una interfaz HTML actualizada dinámicamente con HTMX.

## Funcionalidades

- Listar los registros ordenados por su identificador interno.
- Crear un registro con álbum, usuario, título, artista, descripción, género, duración y rating.
- Validar que los campos tengan valores correctos.
- Mostrar mensajes de errores por pantalla cuando haga falta.
- Eliminar un registro desde su fila.

Hay además consultas como `GetRegister` (para obtener de a un registro por su ID) y `UpdateRegister` (para actualizar registros existentes) que también están definidas y generan funciones Go mediante SQLC, pero actualmente no tienen handlers ni controles en la interfaz web. Decidimos dejarlas para una posible futura actualización.)

## Requisitos

- Go compatible con la versión indicada en `go.mod`.
- Docker con el plugin Docker Compose.
- GNU Make.
- Acceso a Internet la primera vez que se ejecuta `make generate`, porque los comandos descargan las versiones fijadas de SQLC y Templ.

## Ejecución

Desde la carpeta del proyecto:

```sh
make run
```

El objetivo `run` ejecuta, en orden:

1. `make up`, que depende de `make clean`: detiene y elimina los contenedores y volúmenes de Compose.
2. `docker compose up -d --wait`, que inicia PostgreSQL en segundo plano y espera a que esté saludable.
3. `make generate`, que genera el código de SQLC y las vistas Templ.
4. `go run main.go`, que inicia el servidor web en `http://localhost:8080`.

Al terminar o interrumpir el servidor con `Ctrl+C`, la receta ejecuta `make clean` mediante un `trap`. `clean` corre `docker compose down -v`, que elimina también el volumen de la base de datos.

> **Importante:** tanto `make run` como `make up` limpian los volúmenes al comenzar. Al finalizar `make run`, también se eliminan. Por lo tanto, los datos no se conservan entre ejecuciones. No uses estos objetivos si necesitas preservar los datos de una base local.

## Comandos de Make

| Comando | Función |
| --- | --- |
| `make` o `make all` | Ejecuta el objetivo predeterminado `run`. |
| `make generate` | Ejecuta SQLC `v1.31.1` y Templ `v0.3.1118` para regenerar código. |
| `make clean` | Ejecuta `docker compose down -v`; detiene los servicios y elimina los volúmenes. |
| `make up` | Limpia el entorno y luego levanta PostgreSQL esperando el healthcheck. |
| `make run` | Prepara la base, genera el código, inicia la aplicación y limpia Docker al salir. |
| `make test` | Prepara la base, genera el código, ejecuta las pruebas de `tests/` y limpia Docker al terminar. |

Las versiones de los generadores están fijadas en `SQLC_VERSION` y `TEMPL_VERSION` dentro del `Makefile`, para evitar que una ejecución use automáticamente una versión distinta.

## Pruebas

```sh
make test
```

Las pruebas de `tests/CRUD_test.go` se conectan a PostgreSQL y ejercitan las operaciones de creación, lectura, listado, actualización y borrado. El objetivo `test` prepara la base y ejecuta `go test -v ./tests/...`; al terminar, su `trap` llama a `make clean`, incluso si las pruebas fallan.

Los archivos `tests/requests.hurl`, `tests/test-api.sh` y `tests/test.sh` contienen material adicional para probar la API. Para las pruebas de integración, PostgreSQL debe poder iniciarse mediante Docker Compose.

## Tecnologías

- **Go** implementa el servidor HTTP, los handlers y la validación del formulario.
- **PostgreSQL** almacena los registros.
- **Docker Compose** ejecuta PostgreSQL en un contenedor y espera a que pase su healthcheck.
- **SQLC** transforma consultas SQL escritas manualmente en funciones y tipos Go generados y tipados.
- **Templ** genera componentes Go que renderizan el HTML del layout, el formulario y la lista.
- **HTMX** envía el formulario y las solicitudes de borrado sin recargar la página completa, e intercambia fragmentos HTML devueltos por el servidor.
- **pgx** es el driver usado por `database/sql` para conectarse a PostgreSQL.
- **Make** automatiza la generación, el arranque y las pruebas.

## Flujo y arquitectura interna

La aplicación Go corre en el host, escuchando en el puerto `8080`. PostgreSQL corre en Docker y Compose publica el puerto `5432` del contenedor como `5432` en el host. Por eso el servidor Go se conecta a `localhost:5432`; no hace falta incluir la aplicación web en Docker para esta configuración.

Al iniciar, `main.go` abre la conexión a la base, comprueba que PostgreSQL responda con `Ping` y crea el conjunto de consultas SQLC. Luego registra estos endpoints:

| Método | Ruta | Comportamiento |
| --- | --- | --- |
| `GET` | `/` | Consulta los registros y renderiza la página completa. |
| `POST` | `/albums` | Valida los campos, crea el registro y devuelve la lista actualizada. Si la conversión numérica o el rango del rating fallan, devuelve un mensaje para el destino HTMX `#form-error`. |
| `DELETE` | `/albums/{id}` | Elimina el registro identificado por su `id_registro`. |

El formulario Templ envía `POST /albums` con `hx-target="#register-list"` y `hx-swap="outerHTML"`. En caso exitoso, el servidor consulta la lista de nuevo y devuelve el fragmento `RegisterList`; HTMX reemplaza el contenedor existente. Cada botón de borrado apunta al `id_registro` de su fila y, si la solicitud tiene éxito, HTMX elimina esa fila del documento.

El mensaje de validación se ubica dentro de `#register-list`. El handler indica a HTMX que lo use como destino mediante `HX-Retarget` y que inserte el texto dentro del elemento con `HX-Reswap`. Los errores de base de datos durante la creación actualmente se registran en el log y redirigen a `/`.

## Base de datos

El esquema está en `db/schema/schema.sql` y crea la tabla `albums` con estas columnas:

| Columna | Uso |
| --- | --- |
| `id_registro` | Clave primaria autoincremental que identifica cada registro. |
| `id_album` | Identificador del álbum asociado. |
| `id_usuario` | Usuario que registra la valoración. |
| `titulo`, `artista`, `descripcion`, `genero` | Datos descriptivos del álbum. |
| `duracion_min` | Duración expresada en minutos. |
| `rating` | Valor entero obligatorio entre 1 y 5, validado también por una restricción `CHECK`. |
| `created_at` | Fecha y hora asignadas por PostgreSQL al crear el registro. |

Las consultas están en `db/queries/queries.sql`. `sqlc.yaml` le indica a SQLC que use el motor PostgreSQL, lea el esquema y las consultas de esas carpetas y escriba el código generado en `db/sqlc/`. Los archivos generados no se editan manualmente: se modifican las consultas o el esquema y luego se vuelve a ejecutar `make generate`.

Las vistas fuente de Templ están en `views/*.templ`; los archivos `*_templ.go` son generados y tampoco se editan manualmente.

## Estructura del proyecto

```text
.
|-- main.go                 # servidor HTTP, conexión y handlers
|-- Makefile                # generación, entorno, ejecución y pruebas
|-- docker-compose.yml      # servicio PostgreSQL
|-- sqlc.yaml               # configuración de SQLC
|-- db/
|   |-- schema/schema.sql   # esquema de PostgreSQL
|   |-- queries/queries.sql # consultas SQL fuente de SQLC
|   `-- sqlc/               # código Go generado por SQLC
|-- views/
|   |-- *.templ             # componentes fuente
|   `-- *_templ.go          # código generado por Templ
|-- tests/                  # pruebas y solicitudes de ejemplo
`-- static/                 # archivos estáticos del proyecto
```