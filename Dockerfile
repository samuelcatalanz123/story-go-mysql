# Etapa 1 — compilar el frontend
FROM node:22-alpine AS frontend
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Etapa 2 — compilar el binario de Go
FROM golang:1.26-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Etapa 3 — imagen final mínima
FROM alpine:3.20
WORKDIR /app
COPY --from=backend /app/server /app/server
COPY --from=frontend /app/web/dist /app/web/dist
ENV WEB_DIR=/app/web/dist
EXPOSE 8080
CMD ["/app/server"]
