# Diseño: Despliegue en vivo (Railway, servicio único)

**Fecha:** 2026-05-28
**Estado:** Aprobado para escribir el plan de implementación
**Contexto:** Tercer sub-proyecto del portafolio (ver
`portfolio-goal-roadmap` en memoria). Despliega la app full-stack
(Go API + React) en una URL pública.

## Objetivo

Publicar la aplicación en internet con una URL que cualquiera pueda abrir
(demo en vivo), de la forma más simple y barata posible, dejando además el
código en GitHub. Es la mejora de mayor impacto para el objetivo del usuario
(conseguir su primer trabajo de programador).

## Decisiones tomadas (brainstorming)

| Tema           | Decisión                                                |
| -------------- | ------------------------------------------------------- |
| Alojamiento    | Todo-en-uno en **Railway**                              |
| Arquitectura   | **Servicio único**: el binario Go sirve API + frontend  |
| Base de datos  | MySQL gestionado de Railway                             |
| Migraciones    | Automáticas al arrancar (`CREATE TABLE IF NOT EXISTS`)  |
| Conexión code  | Vía **GitHub** (push → auto-deploy)                      |
| Coste          | Prueba gratis; luego ~5 $/mes (se verifica al desplegar) |

## Arquitectura

Un solo servicio en Railway. El binario de Go atiende todo:

```
Internet → Railway service [ Go binary ]
             ├── /api/*  → API existente (StripPrefix "/api" → router actual)
             └── /*      → React compilado (web/dist) con SPA fallback
           Railway MySQL plugin  ←─ conexión por variables de entorno
```

Al venir todo del mismo origen, **no hay CORS** y solo se gestiona una URL.

## Cambios en el backend (Go)

1. **Servir el frontend con SPA fallback.** Nuevo manejador que sirve los
   archivos estáticos de un directorio (`web/dist`, copiado por Docker). Si la
   ruta pedida no corresponde a un archivo existente, devuelve `index.html`
   (necesario para que las rutas de React Router funcionen al recargar).
   - El directorio se localiza por la variable `WEB_DIR` (default `web/dist`).
     Si no existe (p. ej. en local sin build), se omite el servido estático y
     la API sigue funcionando.
2. **Enrutado de nivel superior.** Un `ServeMux` raíz:
   - `/api/` → `http.StripPrefix("/api", apiRouter)` (el router actual, sin
     cambios en sus rutas internas `/characters`, etc.).
   - `/` → manejador SPA.
3. **Puerto de Railway.** `config.Load` preferirá la variable `PORT` (la que
   Railway inyecta); si está, `ServerAddr = ":" + PORT`; si no, mantiene
   `SERVER_ADDR` (default `:8080`). Comportamiento local intacto.
4. **Migraciones al arrancar.** El esquema vivirá como archivos `.sql`
   idempotentes (`CREATE TABLE IF NOT EXISTS`) embebidos en el binario
   (`go:embed`) dentro del paquete `internal/storage`. Tras conectar a la base
   de datos, `storage.Migrate(db)` ejecuta cada archivo en orden. Una base de
   datos nueva queda lista automáticamente. Los `.sql` existentes en `sql/` se
   trasladan a `internal/storage/migrations/` (única fuente de verdad) y se
   ajustan a `IF NOT EXISTS`.

## Cambio en el frontend

Único cambio: en `web/vite.config.ts`, el proxy de `/api` deja de reescribir
(quitar) el prefijo. Así, en **desarrollo** (Vite) y en **producción** (Go con
StripPrefix) el frontend llama a `/api/...` de forma idéntica. El código React
no cambia.

## Construcción (Dockerfile multi-etapa)

`Dockerfile` en la raíz del repo:

1. **Etapa frontend (node):** `npm ci` + `npm run build` en `web/` → genera
   `web/dist`.
2. **Etapa backend (golang):** copia el código Go, compila un binario estático.
3. **Etapa final (imagen mínima):** copia el binario y `web/dist`; define
   `WEB_DIR=/app/web/dist`; arranca el binario.

Railway detecta el `Dockerfile` y lo construye en cada push. Se añade un
`.dockerignore` para no copiar `node_modules`, `dist` local, ni artefactos.

## Base de datos y variables de entorno

- En Railway se añade el plugin **MySQL**.
- La app ya lee `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME`. En el
  servicio de Railway se definen esas variables **referenciando** las del
  plugin MySQL (p. ej. `DB_HOST=${{MySQL.MYSQLHOST}}`). Sin código nuevo.

## GitHub + Railway

1. Subir el repo a GitHub (rama `main` tras integrar este trabajo).
2. En Railway: crear proyecto desde el repo de GitHub.
3. Añadir el plugin MySQL y las variables.
4. Cada `push` a `main` redespliega automáticamente.

## README profesional

Actualizar `README.md`:
- Sección **Demo en vivo** con el enlace público.
- Capturas de pantalla del panel (carpeta `docs/screenshots/`).
- Sección **Despliegue** explicando Railway + Docker.
- Mención del stack y del flujo (Go + React + MySQL en un servicio).

## Pasos interactivos del usuario (no automatizables)

El asistente prepara el código, Dockerfile, README y commits. El usuario hace:
- `gh auth login` y crear el repositorio en GitHub.
- Crear cuenta/login en Railway, conectar el repo, añadir MySQL, pegar
  variables, y capturar las pantallas para el README.

## Pruebas

- Las 14 pruebas del frontend siguen pasando (cambio de proxy no afecta tests).
- Backend: añadir una prueba del manejador SPA (sirve `index.html` cuando la
  ruta no es un archivo; sirve un archivo existente tal cual) usando un
  directorio temporal. Verificar que `/api/...` enruta a la API.
- `go build ./...`, `go test ./...`, `npm test`, `npm run build` y `docker build`
  (si Docker está disponible localmente) deben pasar.

## Fuera de alcance (YAGNI)

- Dominio personalizado / HTTPS propio (Railway ya da una URL con HTTPS).
- CI/CD más allá del auto-deploy de Railway (los tests en CI son un
  sub-proyecto posterior).
- Escalado, métricas, logs avanzados.
- Migraciones con versionado/rollback (basta con creación idempotente).

## Criterios de éxito

1. `docker build` produce una imagen que, con las variables de entorno
   adecuadas, sirve el frontend en `/` y la API en `/api/*`.
2. En local, `go run ./cmd/server` (sin `web/dist`) sigue sirviendo la API; con
   `web/dist` presente, sirve también el frontend.
3. La app arranca contra una base de datos vacía y crea las tablas sola.
4. Tras conectar Railway a GitHub, un `push` despliega y la URL pública
   muestra el panel funcionando con datos reales.
5. El README incluye el enlace a la demo y capturas.
6. Todas las pruebas (frontend y backend) pasan.
