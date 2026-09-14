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

.PHONY: test build clean generate

test:
	@printf "$(CYAN)___________________ Generando codigo con SQLC ___________________$(RESET)\n"
	sqlc generate
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

	@printf "$(GREEN)________________________ Proceso finalizado del TP2 ____________$(RESET)\n"