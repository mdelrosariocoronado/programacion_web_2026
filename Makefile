

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
	@echo "____________ Generando codigo con SQLC ____________"
	sqlc generate
	@echo " "
	
	@echo "____________ Limpiando contenedores y volumenes previos ____________"
	docker compose down -v --remove-orphans
	@echo " "

	@echo "____________ Levantando contenedor de Base de Datos ____________"
	docker compose up -d db
	@echo "Esperando a que la base de datos este lista para recibir conexiones"
	@until docker compose exec -T db pg_isready -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do \
		sleep 1; \
	done

	@echo "- Base de datos lista para recibir conexiones"
	@echo " "

	@echo "____________ Aplicando migraciones / esquema a la base de datos ____________"
	atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)" 
	@echo " "

	@echo "____________ Ejecucion de Tests ____________"
	@echo " Corriendo test con paquete de testing de GO "

	go test -v ./...

	@echo " - Finalizaron los tests "
	@echo " "

	@echo "____________ Limpiando contenedores y volumenes creados en las pruebas ____________"
	docker compose down -v --remove-orphans
	@echo " "

	@echo "____________ Proceso finalizado del TP2 ____________"