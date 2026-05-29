# Story API — Go + MySQL

> 🌐 Idiomas / Languages: **English** · [Español](README.md)

A full-stack web application to manage story **characters**, **locations** and
**scenes** (and the relationships between them). A JSON REST API written in Go
with the standard library (`net/http`, `database/sql`) and MySQL, plus a React +
TypeScript admin dashboard — served together as a single service.

## 🚀 Live demo

**Admin dashboard:** <!-- Paste your public Railway URL here after deploying, e.g. https://story-go-mysql-production.up.railway.app -->

### Screenshots

![Characters dashboard](docs/screenshots/characters.png)
![Modal form](docs/screenshots/form.png)

_(Replace these with your own screenshots in `docs/screenshots/`.)_

## Tech stack

- **Backend:** Go (net/http, database/sql) + MySQL.
- **Frontend:** React 19 + TypeScript + Vite + CSS Modules.
- **Tests:** Vitest + React Testing Library (frontend), `go test` (backend).
- **Deployment:** Docker (multi-stage build) on Railway, single service.

## Highlights

- **Layered architecture** with dependency injection (handler → service →
  repository); domain errors isolate database and HTTP details.
- **Reusable React design system** in CSS Modules: data tables, modal forms,
  toast notifications, loading skeletons, empty states and a responsive layout.
- **Automated tests** on both ends and a branch-based Git workflow.
- **Single-service deployment:** the Go binary serves the API under `/api/*`
  and the compiled frontend on every other route (SPA fallback). Idempotent
  schema migrations run on startup, so a fresh database initializes itself.

## Architecture

```
cmd/server/main.go        Composition root: config → db → migrate → repos → services → handlers → server
internal/
  config/                 Environment-variable configuration (with defaults)
  model/                  Domain types and request DTOs
  apperror/               Domain errors (NotFound, DuplicateTitle, Validation)
  storage/                MySQL connection (pool + ping) and embedded migrations
  web/                    JSON helpers and domain-error → HTTP mapping
  repository/             Data access (SQL); maps driver errors to domain errors
  service/                Business logic and validation
  handler/                HTTP layer, router, and static/SPA serving
web/                      React + TypeScript frontend (Vite)
```

Request flow:

```
handler  →  service  →  repository  →  MySQL
(HTTP)      (rules)      (SQL)
```

## Requirements

- Go 1.22+ (tested with 1.26)
- Node.js 20+ (for the frontend)
- MySQL

## Run locally

1. **Database** — start MySQL and create the database:
   ```sql
   CREATE DATABASE story_go_db;
   ```
   No manual table creation needed: the server creates tables on startup if they
   don't exist (idempotent migrations embedded in `internal/storage/migrations/`).

2. **Backend** — from the repo root:
   ```bash
   go mod tidy
   go run ./cmd/server      # serves http://localhost:8080
   ```

3. **Frontend** — in another terminal:
   ```bash
   cd web
   npm install
   npm run dev              # opens http://localhost:5173 (proxies /api → :8080)
   ```

## Configuration

The app is configured via environment variables; every one has a sensible local
default, so it runs out of the box.

| Variable      | Default       | Description                |
| ------------- | ------------- | -------------------------- |
| `PORT`        | _(unset)_     | If set, overrides the listen port (used by hosts like Railway) |
| `SERVER_ADDR` | `:8080`       | HTTP listen address        |
| `WEB_DIR`     | `web/dist`    | Compiled frontend directory to serve |
| `DB_USER`     | `root`        | MySQL user                 |
| `DB_PASSWORD` | _(empty)_     | MySQL password             |
| `DB_HOST`     | `127.0.0.1`   | MySQL host                 |
| `DB_PORT`     | `3306`        | MySQL port                 |
| `DB_NAME`     | `story_go_db` | Database name              |

## API

Each resource (`characters`, `locations`, `scenes`) supports the same set of
operations (in production they live under the `/api` prefix):

| Method   | Path                | Action          |
| -------- | ------------------- | --------------- |
| `POST`   | `/{resource}`       | Create          |
| `GET`    | `/{resource}`       | List            |
| `GET`    | `/{resource}/{id}`  | Get by ID       |
| `PUT`    | `/{resource}/{id}`  | Update          |
| `DELETE` | `/{resource}/{id}`  | Delete          |

Errors are returned as JSON (`{ "error": "title is required" }`) with status
codes `400`, `404`, `405`, `409` (duplicate title) and `500`.

## Deployment (Railway)

The app deploys as a **single service**: the Go binary serves the API under
`/api/*` and the compiled frontend (`web/dist`) on all other routes. Tables are
created automatically on startup (idempotent migrations).

1. Push the repository to GitHub.
2. On [Railway](https://railway.app): **New Project → Deploy from GitHub repo**.
3. Add a **MySQL** plugin to the project.
4. In the service variables, reference the MySQL ones:
   - `DB_HOST=${{MySQL.MYSQLHOST}}`
   - `DB_PORT=${{MySQL.MYSQLPORT}}`
   - `DB_USER=${{MySQL.MYSQLUSER}}`
   - `DB_PASSWORD=${{MySQL.MYSQLPASSWORD}}`
   - `DB_NAME=${{MySQL.MYSQLDATABASE}}`
5. Railway detects the `Dockerfile`, builds it, and publishes a URL. Every push
   to `main` redeploys automatically.
