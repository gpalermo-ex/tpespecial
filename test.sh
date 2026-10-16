#!/bin/bash

# Esta línea hace que si algún comando falla, el script se detenga inmediatamente
set -e

echo "Instalando (si no está instalado) sqlc"
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

export PATH=$PATH:$(go env GOPATH)/bin

echo "1. Generando código de sqlc"
sqlc generate

echo "2. Tareas Previas: Limpiando entorno viejo"
# docker compose down -v detiene los contenedores y elimina los volúmenes
docker compose down -v

echo "3. Tareas Previas: Levantando base de datos "
# docker compose up -d inicia el contenedor en segundo plano
docker compose up -d

echo "4. Tareas Previas: Esperando que la base de datos esté lista"
# PostgreSQL tarda unos segundos en arrancar. Este bucle le pregunta constantemente 
# "estás listo?" usando el comando interno 'pg_isready' hasta que responde que sí.
echo "Esperando..."
until docker exec $(docker compose ps -q db) pg_isready -U admin > /dev/null 2>&1; do
  sleep 1
done
echo "¡Base de datos corriendo!"

sleep 2

echo "5. Creamos la tabla albums"
docker exec -i postgres-container psql -U admin -d base_catalogo < db/schema/schema.sql



echo "6. Ejecutando los tests"
# Corre todos los tests de tu proyecto usando el paquete testing
go test -v

echo "7. Tareas Posteriores: Borrando contenedores y volúmenes"
# Volvemos a limpiar todo para no dejar rastros en la computadora[cite: 4]
docker compose down -v

echo "¡Todo el proceso terminó con éxito!"