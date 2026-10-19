#!/bin/bash
#
# Script de pruebas automatizadas E2E de la API (TP 3)
# Cumple con todo el ciclo de vida requerido:
# 1) Tareas previas: ejecución sqlc, compilación, borrar contenedores y volúmenes previos,
#    levantar contenedores, esperar que la base esté corriendo, migraciones, levantar aplicación.
# 2) Prueba de cada uno de los endpoints (requests.sh).
# 3) Tareas posteriores: borrar contenedores y volúmenes.
#

set -euo pipefail

# Colores para la salida en consola
CYAN="\033[1;36m"
GREEN="\033[1;32m"
YELLOW="\033[1;33m"
RED="\033[1;31m"
BOLD="\033[1m"
RESET="\033[0m"

# Cargar variables de entorno si existe .env
if [ -f .env ]; then
  set -a
  source .env
  set +a
fi

DB_USER="${DB_USER:-user}"
DB_PASSWORD="${DB_PASSWORD:-xyz}"
DB_NAME="${DB_NAME:-db}"
DB_PORT="${DB_PORT_EXTERNAL:-5432}"
DB_HOST="${DB_HOST:-localhost}"
SERVER_PORT="${SERVER_PORT:-8080}"
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# TAREAS POSTERIORES (Garantizadas ante éxito, error o cancelación)
cleanup() {
  echo ""
  printf "${YELLOW}___________________ Limpiando contenedores y volumenes ___________________${RESET}\n"
  docker compose down -v --remove-orphans >/dev/null 2>&1 || true
  printf "${GREEN}Entorno de pruebas limpiado correctamente.${RESET}\n"
}
trap cleanup EXIT

echo ""
printf "${BOLD}${CYAN}======================================================================${RESET}\n"
printf "${BOLD}${CYAN}          INICIANDO SUITE DE PRUEBAS AUTOMATIZADA                     ${RESET}\n"
printf "${BOLD}${CYAN}======================================================================${RESET}\n"
echo ""

# 1) TAREAS PREVIAS __________________________________________________

# Limpieza inicial de contenedores y volúmenes previos
printf "${YELLOW}--- Limpiando contenedores y volumenes previos ---${RESET}\n"
docker compose down -v --remove-orphans
echo ""

# Ejecución de SQLC
printf "${CYAN}--- Ejecutando generador SQLC ---${RESET}\n"
if command -v sqlc >/dev/null 2>&1; then
  sqlc generate
  printf "${GREEN}Código generado exitosamente con sqlc local.${RESET}\n"
elif command -v docker >/dev/null 2>&1 && docker info >/dev/null 2>&1; then
  printf "${YELLOW}sqlc local no encontrado, usando imagen docker oficial sqlc/sqlc...${RESET}\n"
  docker run --rm -v "$(pwd):/src" -w /src sqlc/sqlc generate
  printf "${GREEN}Código generado con sqlc en Docker.${RESET}\n"
else
  printf "${YELLOW}Aviso: sqlc no está en el PATH del host. Se utilizará el código generado en db/sqlc/.${RESET}\n"
fi
echo ""

# Validación de compilación
printf "${CYAN}--- Validando compilacion de la aplicacion Go ---${RESET}\n"
if command -v go >/dev/null 2>&1; then
  go build -o /dev/null ./cmd/server
  printf "${GREEN}Compilación de Go validada con éxito.${RESET}\n"
fi
echo ""

# Levantar contenedor de base de datos
printf "${CYAN}--- Levantando contenedor de Base de Datos PostgreSQL ---${RESET}\n"
docker compose up -d db
echo ""

# Esperar que la base de datos esté corriendo y lista
printf "${CYAN}--- Esperando disponibilidad de la Base de Datos ---${RESET}\n"
until docker compose exec -T db pg_isready -U "${DB_USER}" -d "${DB_NAME}" > /dev/null 2>&1; do
  printf "."
  sleep 1
done
printf "\n${GREEN}Base de datos lista para recibir conexiones.${RESET}\n\n"

# Aplicar migraciones con Atlas
printf "${YELLOW}--- Aplicando migraciones de esquema con Atlas ---${RESET}\n"
if command -v atlas >/dev/null 2>&1; then
  atlas migrate apply --dir "file://db/migrations" --url "${DB_URL}"
else
  docker run --rm -v "$(pwd)/db/migrations:/migrations" --network host arigaio/atlas:latest migrate apply --dir "file:///migrations" --url "${DB_URL}"
fi
echo ""

# Levantar la aplicación web en Docker
printf "${CYAN}--- Levantando aplicacion web en contenedor Docker ---${RESET}\n"
docker compose up -d --build web
echo ""

# Esperar a que la aplicación esté lista y responda en /health
printf "${CYAN}--- Esperando que la aplicacion web este lista en http://localhost:${SERVER_PORT}/health ---${RESET}\n"
MAX_RETRIES=30
COUNT=0
until curl -s "http://localhost:${SERVER_PORT}/health" > /dev/null 2>&1; do
  sleep 1
  COUNT=$((COUNT + 1))
  printf "."
  if [ "$COUNT" -ge "$MAX_RETRIES" ]; then
    printf "\n${RED}ERROR: La aplicación no respondió después de ${MAX_RETRIES} segundos.${RESET}\n"
    docker compose logs web
    exit 1
  fi
done
printf "\n${GREEN}Aplicacion web lista y respondiendo en el puerto ${SERVER_PORT}.${RESET}\n\n"

# 2) PRUEBA DE CADA UNO DE LOS ENDPOINTS _____________________________
printf "${BOLD}${CYAN}======================================================================${RESET}\n"
printf "${BOLD}${CYAN}         2) EJECUTANDO PRUEBAS DE ENDPOINTS (requests.sh)            ${RESET}\n"
printf "${BOLD}${CYAN}======================================================================${RESET}\n"
echo ""

chmod +x requests.sh
./requests.sh

# 3) TAREAS POSTERIORES: se ejecutan automáticamente a través de la función cleanup en el 'trap EXIT'.
