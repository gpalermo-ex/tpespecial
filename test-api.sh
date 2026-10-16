#!/bin/bash

# Detener el script inmediatamente si algún comando falla
set -e

echo -e "- Instalando sqlc (si no estaba instalado previamente)...\n"
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

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

echo -e "- Levantando servidor WEB...\n"
# Compilamos la aplicación para despues poder obtener el PID y poder terminar la ejecución
go build -o api_server .
# Ejecutamos en segundo plano
./api_server &
# Guardamos el PID para luego poder terminar la ejecución
APP_PID=$!

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


echo -e "  Ejecutando pruebas con Hurl "
hurl --test requests.hurl

echo -e "Todas las pruebas finalizaron con éxito!"