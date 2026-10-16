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
docker compose up -d --wait

echo "4. Tareas Previas: Esperando que la base de datos esté lista"
echo "Esperando..."
until docker exec $(docker compose ps -q db) pg_isready -U admin > /dev/null 2>&1; do
  sleep 1
done
echo "¡Base de datos corriendo!"

sleep 2

echo "6. Ejecutando los tests"
go test -v

echo "7. Tareas Posteriores: Borrando contenedores y volúmenes"
docker compose down -v

echo "¡Todo el proceso terminó con éxito!"