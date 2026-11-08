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

echo "> 3. Tareas Previas: Levantando base de datos "
# docker compose up -d inicia el contenedor en segundo plano
docker compose up -d --wait

trap 'echo "> Borrando contenedores y volúmenes..." && docker compose down -v' EXIT

echo "> 4. Ejecutando los tests"
go test -v ./tests/...

echo "Todo el proceso terminó con éxito."