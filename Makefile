.PHONY: test build clean generate

test:
	@echo " --- Generando codigo con SQLC"
	sqlc generate
	@echo " - "
	@echo " --- Limpiando contenedores y volumenes previos"
	docker compose down -v --remove-orphans
	@echo " - "
	@echo " --- Levantando contenedor de Base de Datos"
	docker compose up -d db
	@echo " - "
	@echo "Esperando a que la base de datos este lista para recibir conexion..."
	@until docker exec postgres_db pg_isready -U postgres -d mi_base_de_datos > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo " --- Base de datos lista para recibir conexiones"
	@echo " - "
	@echo " --- Aplicando migraciones / esquema a la base de datos"
	atlas migrate apply --dir "file://db/migrations" --url "postgres://postgres:postgres_password@localhost:5432/mi_base_de_datos?sslmode=disable" || true

	@echo " --- Ejecucion de Tests "
	@echo " - "
	@echo " Corriendo test con paquete de testing de GO "
	go test -v ./...
	@echo " - Finalizaron los tests "
	@echo " --- Limpiando contenedores y volumenes creados en las pruebas "
	docker compose down -v
	@echo " - "
	@echo " Proceso finalizado del TP2 "