# Story API — Go + MySQL

API HTTP en JSON, escrita en Go con la librería estándar (`net/http`,
`database/sql`) y MySQL, para gestionar personajes, lugares, escenas y las
relaciones entre escenas con personajes y lugares.

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
  handler/                Capa HTTP + router
sql/                      Migraciones de esquema
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

```bash
mysql -u root -p story_go_db < sql/001_create_characters.sql
mysql -u root -p story_go_db < sql/002_create_locations.sql
mysql -u root -p story_go_db < sql/003_create_scenes.sql
```

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
