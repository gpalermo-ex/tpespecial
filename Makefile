
.PHONY: all generate up run test clean
SQLC_VERSION := v1.31.1
TEMPL_VERSION := v0.3.1118

all: run

generate:
	@echo "> Ejecutando sqlc y templ..."
	@go run github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION) generate
	@go run github.com/a-h/templ/cmd/templ@$(TEMPL_VERSION) generate
	@echo "> Archivos generados correctamente."

up: clean
	@echo "> Levantando contenedores"
	@docker compose up -d --wait

run: up generate
	@echo "> Ejecutando main.go..."
	@trap '$(MAKE) clean' EXIT; go run main.go

test: up generate
	@echo "> Ejecutando tests..."
	@trap '$(MAKE) clean' EXIT; go test -v ./tests/...

clean:
	@echo "> Borrando contenedores viejos..."
	@docker compose down -v


