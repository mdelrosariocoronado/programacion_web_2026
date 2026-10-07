#!/bin/bash
# =============================================================================
#  requests.sh - Pruebas End-to-End de la API de Emprendimientos
#
#  Ejecutar:  chmod +x requests.sh && ./requests.sh
#
#  Requisitos:
#    - El servidor debe estar corriendo en http://localhost:8080
#    - La base de datos debe aceptar registros nuevos
#    - curl debe estar instalado
#
#  NOTA: Este script CREA y ELIMINA registros reales en la base de datos.
#        Ejecutar solo en un entorno de desarrollo.
# =============================================================================

set -euo pipefail

# ---------------------------------------------------------------------------
# Configuracion
# ---------------------------------------------------------------------------
BASE_URL="http://localhost:8080"
CONTENT_TYPE="Content-Type: application/json"

# Colores para la salida
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[1;33m"
CYAN="\033[0;36m"
BOLD="\033[1m"
RESET="\033[0m"

# Contadores de resultados
PASS=0
FAIL=0
TOTAL=0

# Archivo temporal para capturar respuestas
RESPONSE_FILE=$(mktemp /tmp/api_response_XXXXXX.json)
trap "rm -f $RESPONSE_FILE" EXIT

# ---------------------------------------------------------------------------
# Funcion auxiliar: assert_status
#   $1 = numero de test (ej: "1.1")
#   $2 = descripcion del test
#   $3 = codigo HTTP esperado
#   $4... = comando curl completo
# ---------------------------------------------------------------------------
assert_status() {
    local test_num="$1"
    local description="$2"
    local expected="$3"
    shift 3

    # Ejecutar curl y capturar el codigo HTTP + el cuerpo de respuesta
    local http_code
    http_code=$("$@" -s -o "$RESPONSE_FILE" -w "%{http_code}")

    TOTAL=$((TOTAL + 1))

    if [ "$http_code" = "$expected" ]; then
        PASS=$((PASS + 1))
        printf "  ${GREEN}[PASS]${RESET} %-4s %-50s -> %s\n" "$test_num" "$description" "$http_code"
    else
        FAIL=$((FAIL + 1))
        printf "  ${RED}[FAIL]${RESET} %-4s %-50s -> %s (esperado: %s)\n" "$test_num" "$description" "$http_code" "$expected"
        # Mostrar cuerpo de respuesta para depuracion
        printf "         ${YELLOW}Respuesta: %s${RESET}\n" "$(cat $RESPONSE_FILE)"
    fi
}

# ---------------------------------------------------------------------------
# Funcion auxiliar: extraer ID numerico de una respuesta JSON
#   $1 = nombre del campo (ej: "id_emprendimiento")
#   Usa el contenido de $RESPONSE_FILE
# ---------------------------------------------------------------------------
extract_id() {
    local field="$1"
    grep -o "\"${field}\":[0-9]*" "$RESPONSE_FILE" | head -1 | grep -o '[0-9]*$'
}

# ---------------------------------------------------------------------------
# Verificacion previa: el servidor esta corriendo?
# ---------------------------------------------------------------------------
echo ""
printf "${BOLD}============================================${RESET}\n"
printf "${BOLD}  PRUEBAS E2E - API Emprendimientos${RESET}\n"
printf "${BOLD}  Servidor: ${CYAN}%s${RESET}\n" "$BASE_URL"
printf "${BOLD}============================================${RESET}\n"
echo ""

# Chequear que el servidor responde
HEALTH_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/health" 2>/dev/null || echo "000")
if [ "$HEALTH_CODE" != "200" ]; then
    printf "${RED}ERROR: El servidor no responde en %s (codigo: %s)${RESET}\n" "$BASE_URL" "$HEALTH_CODE"
    printf "${YELLOW}Asegurate de que el servidor este corriendo con: make run${RESET}\n"
    exit 1
fi
printf "${GREEN}Servidor activo (health check OK)${RESET}\n\n"

# =============================================================================
# SECCION 1: CREACION DE REGISTROS (201 Created)
# =============================================================================
printf "${BOLD}${CYAN}--- SECCION 1: CREACION DE REGISTROS (201 Created) ---${RESET}\n\n"

# 1.1 Crear Emprendimiento "TechCorp"
assert_status "1.1" "Crear emprendimiento TechCorp" "201" \
    curl -X POST "${BASE_URL}/api/emprendimientos" \
    -H "$CONTENT_TYPE" \
    -d '{
        "nombre": "TechCorp Soluciones",
        "rubro": "Tecnologia",
        "descripcion": "Empresa de desarrollo de software",
        "url": "https://techcorp.com",
        "contacto": "info@techcorp.com",
        "activo": true
    }'
EMP_ID_1=$(extract_id "id_emprendimiento")

# 1.2 Crear Emprendimiento "FoodPlace"
assert_status "1.2" "Crear emprendimiento FoodPlace" "201" \
    curl -X POST "${BASE_URL}/api/emprendimientos" \
    -H "$CONTENT_TYPE" \
    -d '{
        "nombre": "FoodPlace Gastronomia",
        "rubro": "Gastronomia",
        "descripcion": "Restaurante de comida saludable",
        "url": "https://foodplace.com",
        "contacto": "hola@foodplace.com",
        "activo": true
    }'
EMP_ID_2=$(extract_id "id_emprendimiento")

# 1.3 Crear Usuario emprendedor (vinculado al emprendimiento 1)
assert_status "1.3" "Crear usuario emprendedor" "201" \
    curl -X POST "${BASE_URL}/api/usuarios" \
    -H "$CONTENT_TYPE" \
    -d "{
        \"nombre_completo\": \"Juan Perez\",
        \"email\": \"juan.perez@email.com\",
        \"clave\": \"clave123segura\",
        \"rol\": \"emprendedor\",
        \"id_emprendimiento\": ${EMP_ID_1}
    }"
USR_ID_1=$(extract_id "id_usuario")

# 1.4 Crear Usuario cliente
assert_status "1.4" "Crear usuario cliente" "201" \
    curl -X POST "${BASE_URL}/api/usuarios" \
    -H "$CONTENT_TYPE" \
    -d '{
        "nombre_completo": "Maria Garcia",
        "email": "maria.garcia@email.com",
        "clave": "password456",
        "rol": "cliente"
    }'
USR_ID_2=$(extract_id "id_usuario")

# 1.5 Crear Publicacion tipo producto
assert_status "1.5" "Crear publicacion tipo producto" "201" \
    curl -X POST "${BASE_URL}/api/publicaciones" \
    -H "$CONTENT_TYPE" \
    -d "{
        \"id_emprendimiento\": ${EMP_ID_1},
        \"titulo\": \"Nuevo producto estrella\",
        \"contenido\": \"Descripcion detallada del producto de prueba\",
        \"imagen_url\": \"https://ejemplo.com/imagen.png\",
        \"tipo\": \"producto\",
        \"precio\": 2500.00
    }"
PUB_ID_1=$(extract_id "id_publicacion")

# 1.6 Crear Publicacion tipo servicio
assert_status "1.6" "Crear publicacion tipo servicio" "201" \
    curl -X POST "${BASE_URL}/api/publicaciones" \
    -H "$CONTENT_TYPE" \
    -d "{
        \"id_emprendimiento\": ${EMP_ID_1},
        \"titulo\": \"Servicio de consultoria\",
        \"contenido\": \"Consultoria profesional en desarrollo web\",
        \"imagen_url\": \"https://ejemplo.com/servicio.png\",
        \"tipo\": \"servicio\",
        \"precio\": 5000.50
    }"
PUB_ID_2=$(extract_id "id_publicacion")

# 1.7 Crear Suscripcion (usuario cliente -> emprendimiento 1)
assert_status "1.7" "Crear suscripcion" "201" \
    curl -X POST "${BASE_URL}/api/suscripciones" \
    -H "$CONTENT_TYPE" \
    -d "{
        \"id_usuario\": ${USR_ID_2},
        \"id_emprendimiento\": ${EMP_ID_1}
    }"

echo ""
printf "  ${YELLOW}IDs capturados: EMP1=%s, EMP2=%s, USR1=%s, USR2=%s, PUB1=%s, PUB2=%s${RESET}\n" \
    "$EMP_ID_1" "$EMP_ID_2" "$USR_ID_1" "$USR_ID_2" "$PUB_ID_1" "$PUB_ID_2"
echo ""

# =============================================================================
# SECCION 2: CONSULTAS Y LISTADOS (200 OK)
# =============================================================================
printf "${BOLD}${CYAN}--- SECCION 2: CONSULTAS Y LISTADOS (200 OK) ---${RESET}\n\n"

# 2.1 Listar todos los emprendimientos
assert_status "2.1" "Listar todos los emprendimientos" "200" \
    curl -X GET "${BASE_URL}/api/emprendimientos"

# 2.2 Obtener emprendimiento por ID
assert_status "2.2" "Obtener emprendimiento por ID (${EMP_ID_1})" "200" \
    curl -X GET "${BASE_URL}/api/emprendimientos/${EMP_ID_1}"

# 2.3 Listar todos los usuarios
assert_status "2.3" "Listar todos los usuarios" "200" \
    curl -X GET "${BASE_URL}/api/usuarios"

# 2.4 Obtener usuario por ID
assert_status "2.4" "Obtener usuario por ID (${USR_ID_1})" "200" \
    curl -X GET "${BASE_URL}/api/usuarios/${USR_ID_1}"

# 2.5 Obtener publicacion por ID
assert_status "2.5" "Obtener publicacion por ID (${PUB_ID_1})" "200" \
    curl -X GET "${BASE_URL}/api/publicaciones/${PUB_ID_1}"

# 2.6 Listar publicaciones de un emprendimiento
assert_status "2.6" "Listar publicaciones del emprendimiento ${EMP_ID_1}" "200" \
    curl -X GET "${BASE_URL}/api/emprendimientos/${EMP_ID_1}/publicaciones"

# 2.7 Listar suscripciones de un usuario
assert_status "2.7" "Listar suscripciones del usuario ${USR_ID_2}" "200" \
    curl -X GET "${BASE_URL}/api/usuarios/${USR_ID_2}/suscripciones"

echo ""

# =============================================================================
# SECCION 3: MODIFICACIONES (200 OK)
# =============================================================================
printf "${BOLD}${CYAN}--- SECCION 3: MODIFICACIONES (200 OK) ---${RESET}\n\n"

# 3.1 Actualizar emprendimiento
assert_status "3.1" "Actualizar emprendimiento (${EMP_ID_1})" "200" \
    curl -X PUT "${BASE_URL}/api/emprendimientos/${EMP_ID_1}" \
    -H "$CONTENT_TYPE" \
    -d '{
        "nombre": "TechCorp Soluciones Avanzadas",
        "rubro": "Tecnologia e Innovacion",
        "descripcion": "Empresa lider en desarrollo de software",
        "url": "https://techcorp-avanzadas.com",
        "contacto": "contacto@techcorp.com",
        "activo": true
    }'

# 3.2 Actualizar usuario
assert_status "3.2" "Actualizar usuario (${USR_ID_1})" "200" \
    curl -X PUT "${BASE_URL}/api/usuarios/${USR_ID_1}" \
    -H "$CONTENT_TYPE" \
    -d "{
        \"nombre_completo\": \"Juan Carlos Perez\",
        \"email\": \"juan.carlos@email.com\",
        \"clave\": \"nuevaClave789\",
        \"rol\": \"emprendedor\",
        \"id_emprendimiento\": ${EMP_ID_1}
    }"

# 3.3 Actualizar publicacion
assert_status "3.3" "Actualizar publicacion (${PUB_ID_1})" "200" \
    curl -X PUT "${BASE_URL}/api/publicaciones/${PUB_ID_1}" \
    -H "$CONTENT_TYPE" \
    -d '{
        "titulo": "Producto estrella actualizado",
        "contenido": "Descripcion mejorada del producto",
        "imagen_url": "https://ejemplo.com/imagen_v2.png",
        "tipo": "producto",
        "precio": 3200.00
    }'

echo ""

# =============================================================================
# SECCION 4: CASOS DE ERROR (400 Bad Request / 404 Not Found)
# =============================================================================
printf "${BOLD}${CYAN}--- SECCION 4: CASOS DE ERROR (400 / 404) ---${RESET}\n\n"

# --- Errores de validacion (400) ---

# 4.1 Emprendimiento sin nombre
assert_status "4.1" "POST emprendimiento sin nombre -> 400" "400" \
    curl -X POST "${BASE_URL}/api/emprendimientos" \
    -H "$CONTENT_TYPE" \
    -d '{
        "nombre": "",
        "rubro": "Tecnologia"
    }'

# 4.2 Emprendimiento sin rubro
assert_status "4.2" "POST emprendimiento sin rubro -> 400" "400" \
    curl -X POST "${BASE_URL}/api/emprendimientos" \
    -H "$CONTENT_TYPE" \
    -d '{
        "nombre": "Emprendimiento Test",
        "rubro": ""
    }'

# 4.3 Usuario con email invalido
assert_status "4.3" "POST usuario con email invalido -> 400" "400" \
    curl -X POST "${BASE_URL}/api/usuarios" \
    -H "$CONTENT_TYPE" \
    -d '{
        "nombre_completo": "Test User",
        "email": "email-sin-arroba",
        "clave": "password123",
        "rol": "cliente"
    }'

# 4.4 Usuario con clave corta (menos de 6 caracteres)
assert_status "4.4" "POST usuario con clave corta -> 400" "400" \
    curl -X POST "${BASE_URL}/api/usuarios" \
    -H "$CONTENT_TYPE" \
    -d '{
        "nombre_completo": "Test User",
        "email": "test.user@email.com",
        "clave": "abc",
        "rol": "cliente"
    }'

# 4.5 Usuario con rol invalido
assert_status "4.5" "POST usuario con rol invalido -> 400" "400" \
    curl -X POST "${BASE_URL}/api/usuarios" \
    -H "$CONTENT_TYPE" \
    -d '{
        "nombre_completo": "Test User",
        "email": "test.user@email.com",
        "clave": "password123",
        "rol": "superadmin"
    }'

# 4.6 Publicacion sin titulo
assert_status "4.6" "POST publicacion sin titulo -> 400" "400" \
    curl -X POST "${BASE_URL}/api/publicaciones" \
    -H "$CONTENT_TYPE" \
    -d "{
        \"id_emprendimiento\": ${EMP_ID_1},
        \"titulo\": \"\",
        \"tipo\": \"producto\",
        \"precio\": 100.00
    }"

# 4.7 Publicacion con tipo invalido
assert_status "4.7" "POST publicacion con tipo invalido -> 400" "400" \
    curl -X POST "${BASE_URL}/api/publicaciones" \
    -H "$CONTENT_TYPE" \
    -d "{
        \"id_emprendimiento\": ${EMP_ID_1},
        \"titulo\": \"Publicacion test\",
        \"tipo\": \"tipo_inexistente\",
        \"precio\": 100.00
    }"

# 4.8 Publicacion con precio negativo
assert_status "4.8" "POST publicacion con precio negativo -> 400" "400" \
    curl -X POST "${BASE_URL}/api/publicaciones" \
    -H "$CONTENT_TYPE" \
    -d "{
        \"id_emprendimiento\": ${EMP_ID_1},
        \"titulo\": \"Publicacion test\",
        \"tipo\": \"producto\",
        \"precio\": -50.00
    }"

# 4.9 Publicacion con id_emprendimiento invalido
assert_status "4.9" "POST publicacion con id_emprendimiento <= 0 -> 400" "400" \
    curl -X POST "${BASE_URL}/api/publicaciones" \
    -H "$CONTENT_TYPE" \
    -d '{
        "id_emprendimiento": 0,
        "titulo": "Publicacion test",
        "tipo": "producto",
        "precio": 100.00
    }'

# 4.10 Suscripcion con IDs invalidos
assert_status "4.10" "POST suscripcion con ids invalidos -> 400" "400" \
    curl -X POST "${BASE_URL}/api/suscripciones" \
    -H "$CONTENT_TYPE" \
    -d '{
        "id_usuario": 0,
        "id_emprendimiento": -1
    }'

# --- Errores de recurso no encontrado (404) ---

# 4.11 GET emprendimiento inexistente
assert_status "4.11" "GET emprendimiento inexistente (99999) -> 404" "404" \
    curl -X GET "${BASE_URL}/api/emprendimientos/99999"

# 4.12 GET publicacion inexistente
assert_status "4.12" "GET publicacion inexistente (99999) -> 404" "404" \
    curl -X GET "${BASE_URL}/api/publicaciones/99999"

# 4.13 GET usuario inexistente
assert_status "4.13" "GET usuario inexistente (99999) -> 404" "404" \
    curl -X GET "${BASE_URL}/api/usuarios/99999"

# 4.14 PUT emprendimiento inexistente
assert_status "4.14" "PUT emprendimiento inexistente (99999) -> 404" "404" \
    curl -X PUT "${BASE_URL}/api/emprendimientos/99999" \
    -H "$CONTENT_TYPE" \
    -d '{
        "nombre": "No existe",
        "rubro": "Test"
    }'

# 4.15 PUT publicacion inexistente
assert_status "4.15" "PUT publicacion inexistente (99999) -> 404" "404" \
    curl -X PUT "${BASE_URL}/api/publicaciones/99999" \
    -H "$CONTENT_TYPE" \
    -d '{
        "titulo": "No existe",
        "tipo": "producto",
        "precio": 100.00
    }'

# 4.16 PUT usuario inexistente
assert_status "4.16" "PUT usuario inexistente (99999) -> 404" "404" \
    curl -X PUT "${BASE_URL}/api/usuarios/99999" \
    -H "$CONTENT_TYPE" \
    -d '{
        "nombre_completo": "No existe",
        "email": "noexiste@test.com",
        "clave": "clave123",
        "rol": "cliente"
    }'

# 4.17 DELETE emprendimiento inexistente
assert_status "4.17" "DELETE emprendimiento inexistente (99999) -> 404" "404" \
    curl -X DELETE "${BASE_URL}/api/emprendimientos/99999"

# 4.18 DELETE publicacion inexistente
assert_status "4.18" "DELETE publicacion inexistente (99999) -> 404" "404" \
    curl -X DELETE "${BASE_URL}/api/publicaciones/99999"

# 4.19 DELETE usuario inexistente
assert_status "4.19" "DELETE usuario inexistente (99999) -> 404" "404" \
    curl -X DELETE "${BASE_URL}/api/usuarios/99999"

# 4.20 DELETE suscripcion inexistente
assert_status "4.20" "DELETE suscripcion inexistente -> 404" "404" \
    curl -X DELETE "${BASE_URL}/api/suscripciones" \
    -H "$CONTENT_TYPE" \
    -d '{
        "id_usuario": 99999,
        "id_emprendimiento": 99999
    }'

echo ""

# =============================================================================
# SECCION 5: ELIMINACIONES (204 No Content + verificacion 404)
# =============================================================================
printf "${BOLD}${CYAN}--- SECCION 5: ELIMINACIONES (204 + verificacion 404) ---${RESET}\n\n"

# 5.1 Eliminar suscripcion -> 204
assert_status "5.1a" "Eliminar suscripcion -> 204" "204" \
    curl -X DELETE "${BASE_URL}/api/suscripciones" \
    -H "$CONTENT_TYPE" \
    -d "{
        \"id_usuario\": ${USR_ID_2},
        \"id_emprendimiento\": ${EMP_ID_1}
    }"

# 5.1 Verificar que la suscripcion ya no existe -> 404
assert_status "5.1b" "Verificar suscripcion eliminada -> 404" "404" \
    curl -X DELETE "${BASE_URL}/api/suscripciones" \
    -H "$CONTENT_TYPE" \
    -d "{
        \"id_usuario\": ${USR_ID_2},
        \"id_emprendimiento\": ${EMP_ID_1}
    }"

# 5.2 Eliminar publicacion 1 -> 204
assert_status "5.2a" "Eliminar publicacion ${PUB_ID_1} -> 204" "204" \
    curl -X DELETE "${BASE_URL}/api/publicaciones/${PUB_ID_1}"

# Verificar que la publicacion 1 ya no existe -> 404
assert_status "5.2b" "Verificar publicacion ${PUB_ID_1} eliminada -> 404" "404" \
    curl -X GET "${BASE_URL}/api/publicaciones/${PUB_ID_1}"

# 5.3 Eliminar publicacion 2 -> 204
assert_status "5.3a" "Eliminar publicacion ${PUB_ID_2} -> 204" "204" \
    curl -X DELETE "${BASE_URL}/api/publicaciones/${PUB_ID_2}"

# Verificar que la publicacion 2 ya no existe -> 404
assert_status "5.3b" "Verificar publicacion ${PUB_ID_2} eliminada -> 404" "404" \
    curl -X GET "${BASE_URL}/api/publicaciones/${PUB_ID_2}"

# 5.4 Eliminar usuario 1 -> 204
assert_status "5.4a" "Eliminar usuario ${USR_ID_1} -> 204" "204" \
    curl -X DELETE "${BASE_URL}/api/usuarios/${USR_ID_1}"

# Verificar que el usuario 1 ya no existe -> 404
assert_status "5.4b" "Verificar usuario ${USR_ID_1} eliminado -> 404" "404" \
    curl -X GET "${BASE_URL}/api/usuarios/${USR_ID_1}"

# 5.5 Eliminar usuario 2 -> 204
assert_status "5.5a" "Eliminar usuario ${USR_ID_2} -> 204" "204" \
    curl -X DELETE "${BASE_URL}/api/usuarios/${USR_ID_2}"

# Verificar que el usuario 2 ya no existe -> 404
assert_status "5.5b" "Verificar usuario ${USR_ID_2} eliminado -> 404" "404" \
    curl -X GET "${BASE_URL}/api/usuarios/${USR_ID_2}"

# 5.6 Eliminar emprendimiento 1 -> 204
assert_status "5.6a" "Eliminar emprendimiento ${EMP_ID_1} -> 204" "204" \
    curl -X DELETE "${BASE_URL}/api/emprendimientos/${EMP_ID_1}"

# Verificar que el emprendimiento 1 ya no existe -> 404
assert_status "5.6b" "Verificar emprendimiento ${EMP_ID_1} eliminado -> 404" "404" \
    curl -X GET "${BASE_URL}/api/emprendimientos/${EMP_ID_1}"

# 5.7 Eliminar emprendimiento 2 -> 204
assert_status "5.7a" "Eliminar emprendimiento ${EMP_ID_2} -> 204" "204" \
    curl -X DELETE "${BASE_URL}/api/emprendimientos/${EMP_ID_2}"

# Verificar que el emprendimiento 2 ya no existe -> 404
assert_status "5.7b" "Verificar emprendimiento ${EMP_ID_2} eliminado -> 404" "404" \
    curl -X GET "${BASE_URL}/api/emprendimientos/${EMP_ID_2}"

echo ""

# =============================================================================
# RESUMEN FINAL
# =============================================================================
printf "${BOLD}============================================${RESET}\n"
if [ "$FAIL" -eq 0 ]; then
    printf "${BOLD}  ${GREEN}RESULTADOS: %d/%d pasaron, %d fallaron${RESET}\n" "$PASS" "$TOTAL" "$FAIL"
else
    printf "${BOLD}  ${RED}RESULTADOS: %d/%d pasaron, %d fallaron${RESET}\n" "$PASS" "$TOTAL" "$FAIL"
fi
printf "${BOLD}============================================${RESET}\n"
echo ""

# Codigo de salida: 0 si todos pasaron, 1 si hubo fallos
exit "$FAIL"
