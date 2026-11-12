# colores para mayor legibilidad del MAKE TEST
CYAN    := \033[1;36m
GREEN   := \033[1;32m
YELLOW  := \033[1;33m
RED     := \033[1;31m
RESET   := \033[0m

# cargo variables del .env 
-include .env
export

DB_USER ?= user
DB_PASSWORD ?= xyz
DB_NAME ?= db
DB_PORT ?= 5432
DB_HOST ?= localhost

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: test build clean generate e2e docker-up docker-down test-api run run-app

generate:
	@printf "$(CYAN)___________________ Generando codigo con SQLC y Templ ___________________$(RESET)\n"
	@if command -v sqlc >/dev/null 2>&1; then sqlc generate; else echo "sqlc CLI no encontrado en PATH, usando código ya generado en db/sqlc/"; fi
	@if command -v templ >/dev/null 2>&1; then templ generate; else go run github.com/a-h/templ/cmd/templ generate; fi

test:
	@printf "$(CYAN)___________________ Generando codigo con SQLC y Templ ___________________$(RESET)\n"
	@if command -v sqlc >/dev/null 2>&1; then sqlc generate; else echo "sqlc CLI no encontrado en PATH, usando código ya generado en db/sqlc/"; fi
	@if command -v templ >/dev/null 2>&1; then templ generate; else go run github.com/a-h/templ/cmd/templ generate; fi
	@printf " "
	
	@printf "$(YELLOW)___________________ Limpiando contenedores y volumenes previos ___________________$(RESET)\n"
	docker compose down -v --remove-orphans
	@printf " "

	@printf "$(CYAN)___________________ Levantando contenedor de Base de Datos ___________________$(RESET)\n"
	docker compose up -d db
	@printf "Esperando a que la base de datos este lista para recibir conexiones"
	@until docker compose exec -T db pg_isready -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do \
		sleep 1; \
	done

	@printf "$(GREEN)- Base de datos lista para recibir conexiones $(RESET)\n"
	@printf " "

	@printf "$(YELLOW)___________________ Aplicando migraciones / esquema a la base de datos ___________________$(RESET)\n"
	atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)" 
	@printf " "

	@printf "$(CYAN)___________________ EJECUCION DE TESTS ___________________$(RESET)\n"
	@printf " Corriendo test con paquete de testing de GO "

	go test -v ./... | sed \
		-e "s/PASS/$$(printf '$(GREEN)PASS$(RESET)')/g" \
		-e "s/FAIL/$$(printf '$(RED)FAIL$(RESET)')/g"

	@printf "$(CYAN) - Finalizaron los tests $(RESET)\n"
	@printf " "

	@printf "$(YELLOW)________________________ Limpiando contenedores y volumenes creados en las pruebas ____________$(RESET)\n"
	docker compose down -v --remove-orphans
	@printf " "

	@printf "$(GREEN)________________________ Proceso finalizado ____________$(RESET)\n"

run:
	@printf "$(CYAN)Generando vistas con Templ...$(RESET)\n"
	@if command -v templ >/dev/null 2>&1; then templ generate; else go run github.com/a-h/templ/cmd/templ generate; fi
	@printf "$(GREEN)Levantando base de datos...$(RESET)\n"
	docker compose up -d db
	@until docker compose exec -T db pg_isready -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do \
		sleep 1; \
	done
	atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"
	@echo "Iniciando servidor..."
	go run cmd/server/main.go

run-app:
	@printf "$(CYAN)___________________ 1. Tareas previas: SQLC y Templ ___________________$(RESET)\n"
	@if command -v sqlc >/dev/null 2>&1; then sqlc generate; else echo "sqlc CLI no encontrado en PATH, usando código ya generado en db/sqlc/"; fi
	@if command -v templ >/dev/null 2>&1; then templ generate; else go run github.com/a-h/templ/cmd/templ generate; fi
	@printf "$(CYAN)___________________ Compilando aplicacion Go ___________________$(RESET)\n"
	go build -o bin/server ./cmd/server
	@printf "$(YELLOW)___________________ Limpiando contenedores y volumenes previos ___________________$(RESET)\n"
	docker compose down -v --remove-orphans
	@printf "$(CYAN)___________________ Levantando contenedor de Base de Datos ___________________$(RESET)\n"
	docker compose up -d db
	@printf "Esperando a que la base de datos este lista para recibir conexiones...\n"
	@until docker compose exec -T db pg_isready -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do \
		sleep 1; \
	done
	@printf "$(GREEN)- Base de datos lista para recibir conexiones$(RESET)\n"
	@printf "$(YELLOW)___________________ Aplicando migraciones a la base de datos ___________________$(RESET)\n"
	atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"
	@printf "$(GREEN)___________________ 2. Ejecutando la aplicacion (Ctrl+C para detener y limpiar) ___________________$(RESET)\n"
	@trap 'printf "\n$(YELLOW)___________________ 3. Tareas posteriores: Limpiando contenedores y volumenes ___________________$(RESET)\n"; docker compose down -v --remove-orphans; rm -rf bin; printf "$(GREEN)Limpieza completada.$(RESET)\n"' EXIT INT TERM; \
	./bin/server

e2e:
	@chmod +x requests.sh
	@./requests.sh

docker-up:
	@printf "$(GREEN)Levantando BD y Servidor Go en Docker...$(RESET)\n"
	docker compose up -d --build
	@until docker compose exec -T db pg_isready -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do \
		sleep 1; \
	done
	atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"
	@printf "$(GREEN)Todo listo. Servidor corriendo en http://localhost:8080$(RESET)\n"

docker-down:
	@printf "$(YELLOW)Deteniendo contenedores Docker...$(RESET)\n"
	docker compose down


test-api:
	@chmod +x test-api.sh
	@./test-api.sh