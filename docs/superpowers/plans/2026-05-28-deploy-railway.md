# Despliegue en Railway (servicio único) — Plan de implementación

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Preparar la app para desplegarse como un único servicio en Railway, donde el binario de Go sirve la API (`/api/*`) y el frontend React compilado (`/*`), con migraciones automáticas y construcción vía Docker.

**Architecture:** Cambios mínimos en el backend Go (servir estáticos con SPA fallback, enrutar `/api/*`, leer `PORT`, migrar al arrancar), un cambio mínimo en el proxy de Vite, un Dockerfile multi-etapa, y un README de portafolio. Los pasos de cuentas (GitHub, Railway) los hace el usuario, guiado.

**Tech Stack:** Go 1.26 (net/http, database/sql, embed), React/Vite, Docker, Railway, MySQL.

---

## Estructura de archivos

```
internal/config/config.go              (mod) leer PORT y WEB_DIR
internal/storage/migrate.go            (nuevo) runner de migraciones (embed)
internal/storage/migrate_test.go       (nuevo) test de splitStatements
internal/storage/migrations/           (nuevo) esquema idempotente embebido
  001_create_characters.sql
  002_create_locations.sql
  003_create_scenes.sql
internal/handler/spa.go                (nuevo) SPAHandler + WithFrontend
internal/handler/spa_test.go           (nuevo) tests del SPA y enrutado /api
cmd/server/main.go                     (mod) migrar al arrancar + WithFrontend
web/vite.config.ts                     (mod) el proxy /api ya no reescribe
Dockerfile                             (nuevo) build multi-etapa
.dockerignore                          (nuevo)
README.md                              (mod) demo, capturas, despliegue
docs/screenshots/.gitkeep              (nuevo) carpeta para capturas
sql/001_create_characters.sql          (eliminar)
sql/002_create_locations.sql           (eliminar)
sql/003_create_scenes.sql              (eliminar)
```

Trabajamos en la rama `feat/deploy` (ya creada). Commits ahí. Comandos Go desde
la raíz del repo; comandos npm desde `web/`.

---

### Task 1: Config — leer PORT y WEB_DIR

**Files:**
- Modify: `internal/config/config.go`

- [ ] **Step 1: Reemplazar TODO `internal/config/config.go`**

```go
// Package config loads runtime configuration from environment variables,
// falling back to sensible development defaults.
package config

import (
	"fmt"
	"os"
)

// Config holds every setting the application needs to start.
type Config struct {
	ServerAddr string
	WebDir     string
	DB         DBConfig
}

// DBConfig holds the MySQL connection settings.
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// DSN builds the MySQL data source name from the configured fields.
func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}

// Load reads configuration from the environment, applying defaults that
// match a local MySQL instance so the project runs out of the box.
func Load() Config {
	return Config{
		ServerAddr: serverAddr(),
		WebDir:     env("WEB_DIR", "web/dist"),
		DB: DBConfig{
			User:     env("DB_USER", "root"),
			Password: env("DB_PASSWORD", ""),
			Host:     env("DB_HOST", "127.0.0.1"),
			Port:     env("DB_PORT", "3306"),
			Name:     env("DB_NAME", "story_go_db"),
		},
	}
}

// serverAddr prefers the PORT variable that platforms like Railway inject,
// falling back to SERVER_ADDR (default :8080) for local development.
func serverAddr() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return ":" + port
	}
	return env("SERVER_ADDR", ":8080")
}

// env returns the value of the environment variable or a fallback default.
func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
```

- [ ] **Step 2: Verificar compilación**

Run desde la raíz: `go build ./...`
Expected: sin errores.

- [ ] **Step 3: Commit**

```bash
git add internal/config/config.go && git commit -m "feat: leer PORT y WEB_DIR de la configuración"
```

---

### Task 2: Migraciones idempotentes embebidas

**Files:**
- Create: `internal/storage/migrations/001_create_characters.sql`
- Create: `internal/storage/migrations/002_create_locations.sql`
- Create: `internal/storage/migrations/003_create_scenes.sql`
- Delete: `sql/001_create_characters.sql`, `sql/002_create_locations.sql`, `sql/003_create_scenes.sql`

- [ ] **Step 1: Crear `internal/storage/migrations/001_create_characters.sql`**

```sql
CREATE TABLE IF NOT EXISTS characters (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL UNIQUE,
  text TEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

- [ ] **Step 2: Crear `internal/storage/migrations/002_create_locations.sql`**

```sql
CREATE TABLE IF NOT EXISTS locations (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL UNIQUE,
  text TEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

- [ ] **Step 3: Crear `internal/storage/migrations/003_create_scenes.sql`**

```sql
CREATE TABLE IF NOT EXISTS scenes (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL UNIQUE,
  text TEXT NULL,
  start_timeline INT NOT NULL,
  end_timeline INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS scene_characters (
  scene_id BIGINT UNSIGNED NOT NULL,
  character_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (scene_id, character_id),
  CONSTRAINT fk_scene_characters_scene
    FOREIGN KEY (scene_id) REFERENCES scenes(id) ON DELETE CASCADE,
  CONSTRAINT fk_scene_characters_character
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS scene_locations (
  scene_id BIGINT UNSIGNED NOT NULL,
  location_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (scene_id, location_id),
  CONSTRAINT fk_scene_locations_scene
    FOREIGN KEY (scene_id) REFERENCES scenes(id) ON DELETE CASCADE,
  CONSTRAINT fk_scene_locations_location
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE
);
```

- [ ] **Step 4: Eliminar los `.sql` antiguos de la raíz**

Run desde la raíz:
```bash
rm sql/001_create_characters.sql sql/002_create_locations.sql sql/003_create_scenes.sql
rmdir sql 2>/dev/null || true
```

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "refactor: esquema idempotente movido a internal/storage/migrations"
```

---

### Task 3: Runner de migraciones (TDD) + integración en main

**Files:**
- Test: `internal/storage/migrate_test.go`
- Create: `internal/storage/migrate.go`
- Modify: `cmd/server/main.go`

- [ ] **Step 1: Escribir el test que falla — `internal/storage/migrate_test.go`**

```go
package storage

import "testing"

func TestSplitStatements(t *testing.T) {
	in := "CREATE TABLE a (id INT);\n\nCREATE TABLE b (id INT);\n"
	got := splitStatements(in)
	if len(got) != 2 {
		t.Fatalf("esperaba 2 sentencias, obtuve %d: %v", len(got), got)
	}
	if got[0] != "CREATE TABLE a (id INT)" {
		t.Fatalf("sentencia 0 inesperada: %q", got[0])
	}
	if got[1] != "CREATE TABLE b (id INT)" {
		t.Fatalf("sentencia 1 inesperada: %q", got[1])
	}
}

func TestSplitStatementsIgnoresEmpty(t *testing.T) {
	if got := splitStatements("   ;\n;  "); len(got) != 0 {
		t.Fatalf("esperaba 0 sentencias, obtuve %v", got)
	}
}
```

- [ ] **Step 2: Ejecutar y ver que falla**

Run desde la raíz: `go test ./internal/storage/`
Expected: FALLA (no existe `splitStatements`, no compila).

- [ ] **Step 3: Crear `internal/storage/migrate.go`**

```go
package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate runs every embedded .sql migration in filename order. The schema
// uses CREATE TABLE IF NOT EXISTS, so running it repeatedly is safe and a
// fresh database gets initialized automatically on startup.
func Migrate(ctx context.Context, db *sql.DB) error {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		content, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		for _, stmt := range splitStatements(string(content)) {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return fmt.Errorf("exec %s: %w", name, err)
			}
		}
	}
	return nil
}

// splitStatements divides a SQL file into individual statements by ';',
// discarding empty fragments and surrounding whitespace. (The MySQL driver
// executes one statement per Exec call.)
func splitStatements(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	stmts := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			stmts = append(stmts, s)
		}
	}
	return stmts
}
```

- [ ] **Step 4: Ejecutar y ver que pasa**

Run desde la raíz: `go test ./internal/storage/`
Expected: PASA (2 tests).

- [ ] **Step 5: Integrar en `cmd/server/main.go`** — añadir la migración tras abrir la BD

Localiza en `run()` este bloque:
```go
	db, err := storage.NewMySQL(cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()
```
y déjalo así (añade el bloque de migración justo debajo del `defer db.Close()`):
```go
	db, err := storage.NewMySQL(cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Crea el esquema si no existe (idempotente). Así una base de datos
	// nueva (p. ej. en Railway) queda lista al arrancar.
	migCtx, migCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer migCancel()
	if err := storage.Migrate(migCtx, db); err != nil {
		return err
	}
```
(`context` y `time` ya están importados en este archivo.)

- [ ] **Step 6: Verificar compilación**

Run desde la raíz: `go build ./...`
Expected: sin errores.

- [ ] **Step 7: Commit**

```bash
git add -A && git commit -m "feat: ejecutar migraciones idempotentes al arrancar"
```

---

### Task 4: Servir el frontend (SPA) y enrutar /api (TDD)

**Files:**
- Test: `internal/handler/spa_test.go`
- Create: `internal/handler/spa.go`
- Modify: `cmd/server/main.go`

- [ ] **Step 1: Escribir el test que falla — `internal/handler/spa_test.go`**

```go
package handler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"story-go-mysql/internal/handler"
)

func TestSPAHandlerServesExistingFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("INDEX"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "app.js"), []byte("JS"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := handler.SPAHandler(dir)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/app.js", nil))
	if rec.Body.String() != "JS" {
		t.Fatalf("esperaba el contenido del archivo, obtuve %q", rec.Body.String())
	}
}

func TestSPAHandlerFallsBackToIndex(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("INDEX"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := handler.SPAHandler(dir)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/scenes", nil))
	if rec.Body.String() != "INDEX" {
		t.Fatalf("esperaba el index.html como fallback, obtuve %q", rec.Body.String())
	}
}

func TestWithFrontendRoutesAPI(t *testing.T) {
	api := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/characters" {
			_, _ = w.Write([]byte("API"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("INDEX"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := handler.WithFrontend(api, dir)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/characters", nil))
	if rec.Body.String() != "API" {
		t.Fatalf("esperaba que /api/characters enrutara a la API, obtuve %q", rec.Body.String())
	}
}
```

- [ ] **Step 2: Ejecutar y ver que falla**

Run desde la raíz: `go test ./internal/handler/`
Expected: FALLA (no existen `SPAHandler` ni `WithFrontend`).

- [ ] **Step 3: Crear `internal/handler/spa.go`**

```go
package handler

import (
	"net/http"
	"os"
	"path/filepath"
)

// SPAHandler serves static files from dir. Requests that don't map to an
// existing file fall back to index.html, so client-side routing (React
// Router) keeps working when the user reloads a deep link.
func SPAHandler(dir string) http.Handler {
	fileServer := http.FileServer(http.Dir(dir))
	index := filepath.Join(dir, "index.html")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(dir, filepath.Clean(r.URL.Path))
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, index)
	})
}

// WithFrontend composes a top-level handler: requests under /api go to the
// API (with the /api prefix stripped); everything else is served by the SPA
// handler when webDir exists. If webDir is empty or missing, only the API is
// served — useful in local development, where Vite serves the frontend.
func WithFrontend(api http.Handler, webDir string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api))
	if webDir != "" {
		if info, err := os.Stat(webDir); err == nil && info.IsDir() {
			mux.Handle("/", SPAHandler(webDir))
		}
	}
	return mux
}
```

- [ ] **Step 4: Ejecutar y ver que pasa**

Run desde la raíz: `go test ./internal/handler/`
Expected: PASA (3 tests nuevos; los tests previos del paquete, si los hay, siguen pasando).

- [ ] **Step 5: Integrar en `cmd/server/main.go`** — envolver el router con WithFrontend

Localiza este bloque en `run()`:
```go
	// Handlers (HTTP) and router.
	router := handler.Router(
		handler.NewCharacterHandler(characterSvc),
		handler.NewLocationHandler(locationSvc),
		handler.NewSceneHandler(sceneSvc),
	)

	server := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
```
y déjalo así (envuelve `router` con `handler.WithFrontend` usando `cfg.WebDir`):
```go
	// Handlers (HTTP) and router.
	router := handler.Router(
		handler.NewCharacterHandler(characterSvc),
		handler.NewLocationHandler(locationSvc),
		handler.NewSceneHandler(sceneSvc),
	)

	// El binario sirve la API en /api/* y el frontend compilado en el resto.
	app := handler.WithFrontend(router, cfg.WebDir)

	server := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           app,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
```

- [ ] **Step 6: Verificar compilación y tests**

Run desde la raíz: `go build ./... && go test ./...`
Expected: compila; tests del backend pasan.

- [ ] **Step 7: Commit**

```bash
git add -A && git commit -m "feat: servir frontend con SPA fallback y enrutar /api"
```

---

### Task 5: Ajustar el proxy de Vite (no reescribir /api)

**Files:**
- Modify: `web/vite.config.ts`

- [ ] **Step 1: Editar `web/vite.config.ts`** — quitar la línea `rewrite`

Localiza el bloque del proxy:
```ts
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
```
y déjalo así (sin `rewrite`, para que `/api/...` se reenvíe tal cual y coincida con producción):
```ts
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
```

- [ ] **Step 2: Verificar build del frontend**

Run dentro de `web/`: `npm run build`
Expected: build correcto.

- [ ] **Step 3: Commit**

```bash
git add web/vite.config.ts && git commit -m "chore(web): el proxy /api ya no reescribe (igual que producción)"
```

---

### Task 6: Dockerfile multi-etapa y .dockerignore

**Files:**
- Create: `Dockerfile`
- Create: `.dockerignore`

- [ ] **Step 1: Crear `Dockerfile` en la raíz**

```dockerfile
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
```

- [ ] **Step 2: Crear `.dockerignore` en la raíz**

```
**/node_modules
web/dist
.git
.gitignore
docs
*.md
```

- [ ] **Step 3: Verificar el build de Docker (si Docker está disponible)**

Run desde la raíz:
```bash
docker build -t story-app . && echo "DOCKER BUILD OK"
```
Expected: termina con `DOCKER BUILD OK`. Si Docker no está instalado en este equipo, anótalo y omite este paso (Railway construirá la imagen igualmente).

- [ ] **Step 4: Commit**

```bash
git add Dockerfile .dockerignore && git commit -m "build: Dockerfile multi-etapa (frontend + Go) y .dockerignore"
```

---

### Task 7: README de portafolio + carpeta de capturas

**Files:**
- Modify: `README.md`
- Create: `docs/screenshots/.gitkeep`

- [ ] **Step 1: Crear la carpeta de capturas**

Run desde la raíz:
```bash
mkdir -p docs/screenshots && touch docs/screenshots/.gitkeep
```

- [ ] **Step 2: Añadir secciones al `README.md`**

Inserta este bloque justo después del título principal `# Story API — Go + MySQL` y su párrafo introductorio (antes de `## Arquitectura`):

```markdown
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
```

Y añade al final del `README.md` esta sección de despliegue:

```markdown
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
5. Railway detecta el `Dockerfile`, construye y publica una URL. Cada `push`
   a `main` vuelve a desplegar.

### Desarrollo local

Sigue necesitando un MySQL local (ver más abajo). Arranca la API con
`go run ./cmd/server` y el frontend con `cd web && npm run dev`.
```

- [ ] **Step 3: Actualizar la sección "Crear tablas" del README** — ahora es automática

Localiza la sección que indica ejecutar los `sql/00x.sql` a mano (bloques con
`mysql -u root -p story_go_db < sql/...`) y reemplaza esa sección por:

```markdown
## Crear tablas

No hace falta ningún paso manual: al arrancar, el servidor crea las tablas si
no existen (migraciones embebidas en `internal/storage/migrations/`).
```

- [ ] **Step 4: Commit**

```bash
git add -A && git commit -m "docs: README de portafolio (demo, capturas, despliegue) + carpeta de capturas"
```

---

### Task 8: Verificación final + guía de despliegue manual

**Files:** ninguno (verificación e instrucciones)

- [ ] **Step 1: Verificación automática completa**

Run desde la raíz:
```bash
go build ./... && go test ./... && echo "GO OK"
cd web && npm test && npx tsc --noEmit && npm run build && echo "WEB OK"
```
Expected: `GO OK` y `WEB OK`. (Tests backend: storage + handler; frontend: 14.)

- [ ] **Step 2: Smoke test local del servicio único (requiere MySQL local)**

```bash
# build del frontend para tener web/dist
cd web && npm run build && cd ..
# arranca la API (sirve también el frontend porque existe web/dist)
go run ./cmd/server &
sleep 2
curl -s -o /dev/null -w "GET / -> %{http_code}\n" http://localhost:8080/
curl -s -o /dev/null -w "GET /api/characters -> %{http_code}\n" http://localhost:8080/api/characters
# detén el servidor (Ctrl+C o kill del proceso)
```
Expected: `GET / -> 200` (sirve el index del frontend) y `GET /api/characters -> 200`.

- [ ] **Step 3: Verificar Docker (si está disponible)**

```bash
docker build -t story-app . && echo "DOCKER OK"
```
Expected: `DOCKER OK`, o nota de que Docker no está disponible localmente.

- [ ] **Step 4: Guía de despliegue (pasos interactivos del usuario)**

Estos pasos los ejecuta el usuario (requieren sus cuentas). El asistente los
acompaña pero no puede autenticarse por él:

1. **GitHub:** `gh auth login`, luego `gh repo create story-go-mysql --public --source=. --remote=origin --push` (sube `main`).
2. **Railway:** crear cuenta, **New Project → Deploy from GitHub repo**, elegir el repo.
3. Añadir **plugin MySQL** y definir las 5 variables `DB_*` referenciando las del MySQL (ver README).
4. Esperar el build; abrir la URL pública y comprobar el panel.
5. Hacer capturas y guardarlas en `docs/screenshots/`; pegar la URL en el README; commit + push.

- [ ] **Step 5: Commit final (si hubo ajustes)**

```bash
git add -A && git commit -m "test: verificación final previa al despliegue" --allow-empty
```

---

## Notas de verificación (self-review del plan)

- **Cobertura del spec:** servir frontend + SPA (Task 4), enrutar `/api` (Task 4),
  PORT (Task 1), WEB_DIR (Tasks 1 y 4), migraciones idempotentes embebidas
  (Tasks 2-3), cambio de proxy Vite (Task 5), Dockerfile + .dockerignore
  (Task 6), variables de entorno MySQL y GitHub/Railway (README Task 7 + guía
  Task 8), README de portafolio con demo y capturas (Task 7), pruebas (Tasks 3,
  4, 8). 
- **Sin placeholders:** todo el código está completo. El único hueco es la URL
  de la demo en el README, que es contenido que el usuario rellena tras
  desplegar (paso interactivo), no un placeholder de código.
- **Consistencia de tipos/firmas:** `Config.WebDir`, `serverAddr()`,
  `storage.Migrate(ctx, db)`, `splitStatements(string) []string`,
  `handler.SPAHandler(dir string) http.Handler`,
  `handler.WithFrontend(api http.Handler, webDir string) http.Handler` se usan
  con los mismos nombres y firmas en todas las tareas y en `main.go`.
- **Compatibilidad local:** el proxy de Vite ahora reenvía `/api/*` sin
  reescribir y el servidor Go lo enruta con StripPrefix; en local sin `web/dist`
  el servidor solo expone la API (Vite sirve el frontend). Las 14 pruebas del
  frontend no dependen del proxy real, así que siguen pasando.
```
