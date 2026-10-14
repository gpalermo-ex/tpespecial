#!/bin/bash


# ========================================================= PREGUNTAR ESTO =========================================================
# README.md que ponemos?
# Mucho más detallado, como una documentación de todo el proyecto y también como correr los tests
# De última para más orden se puede hacer otro archivo para la documentación e indicarlo en el readme junto a las instrucciones de ejecución

# Detener el script inmediatamente si algún comando falla
set -e

echo -e "- Instalando sqlc (si no estaba instalado previamente)...\n"
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# ========================================================= PREGUNTAR ESTO =========================================================
# Tengo que agregar esto porque sino me dice que el comando sqlc no lo reconoce, está bien? Lo raro es que en el anterior script no
# lo pusimos y el profe no dijo nada, pero creo que queda mejor así.
# Si, está mejor así, go instala en un directorio suyo, así le indicamos a bash el path de sqlc
export PATH=$PATH:$(go env GOPATH)/bin

echo -e "- Ejecutando sqlc...\n"
sqlc generate

echo -e "- Borrando contenedores y volúmenes viejos...\n"
docker compose down -v

echo -e "- Levantando contenedor de postgres...\n"
# INVESTIGAR DOCKER COMPOSE WAIT
docker compose up -d --wait

echo -e "- Comprobando que la base de datos esté lista..."
until docker exec postgres-container pg_isready -U admin > /dev/null 2>&1; do
  sleep 1
done
echo -e "- Base de datos corriendo\n"

echo -e "- Creando la tabla albums dentro de la base de datos...\n"
# ========================================================= PREGUNTAR ESTO =========================================================
# Está bien agregar la tabla álbums así dentro de la base de datos?
# No, hacerlo dentro del docker compose, hay un parámetro, algo de volume
docker exec -i postgres-container psql -U admin -d base_catalogo < db/schema/schema.sql

echo -e "- Levantando servidor WEB...\n"
# Compilamos la aplicación para despues poder obtener el PID y poder terminar la ejecución
go build -o api_server .
# Ejecutamos en segundo plano
./api_server &
# Guardamos el PID para luego poder terminar la ejecución
APP_PID=$!
sleep 2

# Esto se ejecuta al terminar el script, ya sea de forma erronea o no 
cleanup() {
    echo -e "- Deteniendo el servidor Go (PID: $APP_PID)...\n"
    kill $APP_PID || true
    echo -e "- Borrando contenedores y volúmenes residuales...\n"
    docker compose down -v
    echo -e "- Borrando ejecutable...\n"
    rm -f api_server
    echo -e "- Limpieza finalizada.\n"
}
# La instrucción trap vincula la función cleanup al evento EXIT
trap cleanup EXIT

# Esperamos unos segundos por las dudas si no llegó a prenderse el server WEB
sleep 2


# ========================================================= CORREGIR ESTO =========================================================
# Se deben controlar los resultados de los comandos, se puede hacer comparando los json devueltos o usando hurl es más fácil
# con varios tests y sus asserts
echo -e "- PRUEBA DE API ENDPOINTS:\n"

echo "- Probando POST /albums (crear álbum)..."
curl -X POST http://localhost:8080/albums -H "Content-Type: application/json" -d '{"id_usuario" : 1,"titulo" : "Album 1","artista" : "Artista 1","descripcion" : "Descripción ejemplo","genero" : "Rock","duracion_s" : 3400, "ranking" : 4}'
echo "" 

echo "- Probando GET /albums (listar todos)..."
curl -s -X GET http://localhost:8080/albums
echo ""

echo "- Probando GET /albums/1 (obtener por IDAlbum 1)..."
curl -s -X GET http://localhost:8080/albums/1
echo ""

echo "- Probando PUT /albums/1 (actualizar álbum)..."
curl -X PUT http://localhost:8080/albums/1 -H "Content-Type: application/json" -d '{"id_usuario" : 1,"titulo" : "Album A","artista" : "Artista 2","descripcion" : "Descripción ejemplo 2","genero" : "Pop","duracion_s" : 3400,"ranking" : 4}'
echo ""

echo "- Probando GET /albums/1 (verificar actualización)..."
curl -s -X GET http://localhost:8080/albums/1
echo ""

echo "- Probando DELETE /albums/1 (borrar álbum)..."
# Debe devolver error 204, No Content
curl -s -i -X DELETE http://localhost:8080/albums/1
echo ""

echo "- Probando GET /albums (verificar que está vacío)..."
curl -s -X GET http://localhost:8080/albums
echo ""

echo -e "- Todas las pruebas de curl finalizaron con éxito.\n"
