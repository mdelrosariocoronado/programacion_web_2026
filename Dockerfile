# BUILDER___________________________________

FROM golang:1.26-alpine AS builder

# directorio de trabajo en el contenedor de compilación
WORKDIR /app  

# copiar dependencias primero para aprovechar la caché de capas
COPY go.mod go.sum ./
RUN go mod download

# copiar el resto
COPY . .

# genera el ejecutable nativo compilado de forma estática
RUN CGO_ENABLED=0 go build -o /app/api ./cmd/server

# IMAGEN FINAL (EJECUCION)_____________________________

FROM alpine:3.20

# se usa un usuario no privilegiado por seguridad
RUN adduser -D -u 1000 app

# directorio de trabajo
WORKDIR /app

# copia  binario compilado desde la etapa 'BUILDER'
COPY --from=builder /app/api /usr/local/bin/api

# copia archivos estáticos (CSS e imágenes para las vistas SSR)
COPY --from=builder /app/static /app/static

USER app
EXPOSE 8080

# comando para ejecutar la aplicación
CMD ["api"]

