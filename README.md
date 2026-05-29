# Story API — Go + MySQL

[![CI](https://github.com/samuelcatalanz123/story-go-mysql/actions/workflows/ci.yml/badge.svg)](https://github.com/samuelcatalanz123/story-go-mysql/actions/workflows/ci.yml)

> 🌐 Idiomas / Languages: **Español** · [English](README.en.md)

API HTTP en JSON, escrita en Go con la librería estándar (`net/http`,
`database/sql`) y MySQL, para gestionar personajes, lugares, escenas y las
relaciones entre escenas con personajes y lugares.

## 🚀 Demo en vivo

**Panel de administración:** <!-- Pega aquí la URL pública de Railway tras desplegar, p. ej. https://story-go-mysql-production.up.railway.app -->

> Aplicación full-stack: API en Go + MySQL y frontend en React + TypeScript,
> servidos como un único servicio en Railway.

### Capturas

![Panel de personajes](docs/screenshots/characters.png)
![Formulario en modal](docs/screenshots/form.png)

_(Sustituye las imágenes por tus propias capturas en `docs/screenshots/`.)_

## Stack

- **Backend:** Go (net/http, database/sql) + MySQL.
- **Frontend:** React 19 + TypeScript + Vite + CSS Modules.
- **Pruebas:** Vitest + Testing Library (frontend), `go test` (backend).
- **Despliegue:** Docker (build multi-etapa) en Railway.

## Arquitectura

El proyecto sigue una arquitectura en capas con inyección de dependencias.
Cada capa solo depende de la de abajo y los errores de dominio aíslan los
detalles de base de datos y de HTTP.

```
cmd/server/main.go        Composition root: config → db → repos → services → handlers → server
internal/
  config/                 Configuración por variables de entorno (con defaults)
  model/                  Tipos de dominio y DTOs de request
  apperror/               Errores de dominio (NotFound, DuplicateTitle, Validation)
  storage/                Conexión MySQL (pool + ping)
  web/                    Helpers JSON y mapeo de errores de dominio → HTTP
  repository/             Acceso a datos (SQL); traduce errores del driver a dominio
  service/                Lógica de negocio y validación
  handler/                Capa HTTP + router + servido del frontend (SPA)
  storage/migrations/     Esquema (migraciones idempotentes, embebidas)
web/                      Frontend React + TypeScript (Vite)
```

Flujo de una petición:

```
handler  →  service  →  repository  →  MySQL
(HTTP)      (reglas)     (SQL)
```

## Requisitos

- Go 1.22+ (probado con 1.26)
- MySQL

## Instalar dependencias

```bash
go mod tidy
```

## Crear base de datos

```bash
mysql -u root -p
```

Dentro de MySQL:

```sql
CREATE DATABASE story_go_db;
exit;
```

## Crear tablas

No hace falta ningún paso manual: al arrancar, el servidor crea las tablas si
no existen (migraciones embebidas en `internal/storage/migrations/`).

## Configuración

La aplicación se configura con variables de entorno. Todas tienen un valor
por defecto pensado para desarrollo local, así que puede correr sin definir
ninguna.

| Variable      | Default       | Descripción                  |
| ------------- | ------------- | ---------------------------- |
| `SERVER_ADDR` | `:8080`       | Dirección de escucha HTTP    |
| `DB_USER`     | `root`        | Usuario de MySQL             |
| `DB_PASSWORD` | _(vacío)_     | Contraseña de MySQL          |
| `DB_HOST`     | `127.0.0.1`   | Host de MySQL                |
| `DB_PORT`     | `3306`        | Puerto de MySQL              |
| `DB_NAME`     | `story_go_db` | Nombre de la base de datos   |

## Iniciar servidor

```bash
go run ./cmd/server
```

El servidor corre en `http://localhost:8080`.

## Endpoints

Cada recurso (`characters`, `locations`, `scenes`) soporta el mismo conjunto
de operaciones:

| Método   | Ruta                | Acción                |
| -------- | ------------------- | --------------------- |
| `POST`   | `/{recurso}`        | Crear                 |
| `GET`    | `/{recurso}`        | Listar                |
| `GET`    | `/{recurso}/{id}`   | Obtener por ID        |
| `PUT`    | `/{recurso}/{id}`   | Actualizar            |
| `DELETE` | `/{recurso}/{id}`   | Eliminar              |

### Ejemplos

```bash
# Crear un personaje
curl -X POST http://localhost:8080/characters \
  -H "Content-Type: application/json" \
  -d '{"title":"Asha Ren","text":"Una piloto de las colonias exteriores."}'

# Crear un lugar
curl -X POST http://localhost:8080/locations \
  -H "Content-Type: application/json" \
  -d '{"title":"Estacion Nuevo Amanecer","text":"Una estacion comercial."}'

# Crear una escena que relaciona personajes y lugares
curl -X POST http://localhost:8080/scenes \
  -H "Content-Type: application/json" \
  -d '{"title":"Primer Contacto","text":"La tripulacion descubre una senal.","startTimeline":100,"endTimeline":120,"characterIds":[1],"locationIds":[1]}'

# Obtener una escena (incluye sus personajes y lugares relacionados)
curl http://localhost:8080/scenes/1
```

## Errores

La API responde con errores en JSON:

```json
{ "error": "title is required" }
```

Códigos usados:

- `400 Bad Request` — validación o JSON inválido
- `404 Not Found` — recurso inexistente
- `405 Method Not Allowed` — método no soportado en la ruta
- `409 Conflict` — título duplicado
- `500 Internal Server Error` — error inesperado

## Despliegue (Railway)

La app se despliega como **un único servicio**: el binario de Go sirve la API
bajo `/api/*` y el frontend compilado (`web/dist`) en el resto de rutas. Las
tablas se crean solas al arrancar (migraciones idempotentes).

1. Sube el repositorio a GitHub.
2. En [Railway](https://railway.app): **New Project → Deploy from GitHub repo**.
3. Añade un plugin **MySQL** al proyecto.
4. En las variables del servicio, referencia las del MySQL:
   - `DB_HOST=${{MySQL.MYSQLHOST}}`
   - `DB_PORT=${{MySQL.MYSQLPORT}}`
   - `DB_USER=${{MySQL.MYSQLUSER}}`
   - `DB_PASSWORD=${{MySQL.MYSQLPASSWORD}}`
   - `DB_NAME=${{MySQL.MYSQLDATABASE}}`
   - `JWT_SECRET=<un secreto largo y aleatorio>`
5. Railway detecta el `Dockerfile`, construye y publica una URL. Cada `push`
   a `main` vuelve a desplegar.

### Desarrollo local

Sigue necesitando un MySQL local (ver más abajo). Arranca la API con
`go run ./cmd/server` y el frontend con `cd web && npm run dev`.
