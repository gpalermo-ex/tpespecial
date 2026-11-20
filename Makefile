
.PHONY: all generate up run test clean

all: run

generate:
	@echo "> Ejecutando sqlc y templ..."
	@go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate
	@go run github.com/a-h/templ/cmd/templ@latest generate
	@echo "> Archivos generados correctamente."

up: clean
	@echo "> Levantando contenedores"
	@docker compose up -d --wait

run: up generate
	@echo "> Ejecutando main.go..."
	@go run main.go

test: up generate
	@echo "> Ejecutando tests..."
	@trap '$(MAKE) clean' EXIT; go test -v ./tests/...

clean:
	@echo "> Borrando contenedores viejos..."
	@docker compose down -v


