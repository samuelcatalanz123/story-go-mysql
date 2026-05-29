# Diseño: Autenticación con JWT (registro + login)

**Fecha:** 2026-05-28
**Estado:** Aprobado para escribir el plan de implementación
**Contexto:** 5º sub-proyecto del portafolio. Añade autenticación a la Story
API + panel React.

## Objetivo

Añadir registro e inicio de sesión para proteger las operaciones de escritura
de la API, demostrando conocimientos de autenticación/seguridad (una función
muy valorada para empleos junior). Se construye íntegramente en local.

## Decisiones tomadas (brainstorming)

| Tema           | Decisión                                                     |
| -------------- | ------------------------------------------------------------ |
| Mecanismo      | **JWT** (token en cabecera `Authorization: Bearer <token>`)  |
| Caducidad      | 24 h, firmado con `JWT_SECRET` (HS256)                       |
| Usuarios       | Registro + login abiertos; **datos compartidos**             |
| Qué se protege | Solo escrituras (POST/PUT/DELETE); los GET siguen públicos    |
| Contraseñas    | Hasheadas con **bcrypt**                                      |
| UI sin sesión  | Los botones de gestión se **ocultan**; aviso para iniciar sesión |

## Backend

### Base de datos
Migración `internal/storage/migrations/004_create_users.sql` (idempotente):
```sql
CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Capas
- **model:** `User` (id, email, createdAt — nunca el hash en JSON),
  `RegisterRequest` (email, password), `LoginRequest` (email, password),
  `AuthResponse` (token, user).
- **apperror:** nuevo error de dominio `Unauthorized` (→ HTTP 401). Email
  duplicado reutiliza el mapeo a 409 (Conflict).
- **repository:** `UserRepository` con `CreateUser(email, hash) (User, error)`
  y `GetUserByEmail(email) (User, hash, error)`; traduce el error de email
  duplicado del driver a dominio.
- **service:** `AuthService`:
  - `Register(req)` valida (email no vacío, password mínimo 8 caracteres),
    hashea con bcrypt, guarda, devuelve `AuthResponse` con token.
  - `Login(req)` busca por email, compara bcrypt; si falla → `Unauthorized`;
    si ok → emite JWT y devuelve `AuthResponse`.
- **auth (JWT):** helper para emitir y validar tokens (HS256, claims `sub`=id y
  `exp`=24h), leyendo el secreto de `config.JWTSecret`.
- **handler:** `AuthHandler` con `POST /auth/register` y `POST /auth/login`
  (públicos); responden `AuthResponse` en JSON.

### Middleware
`handler.RequireAuth(next)` lee `Authorization: Bearer <token>`, valida el JWT
y, si falla (ausente/inválido/expirado), responde 401 vía el mapeo de errores
existente. Si es válido, deja pasar (opcionalmente inyecta el user id en el
contexto; no es necesario porque los datos son compartidos).

### Router
- `GET` de cada recurso: público (sin cambios).
- `POST/PUT/DELETE` de cada recurso: envueltos con `RequireAuth`.
- `/auth/register`, `/auth/login`: públicos.
- Se mantiene el prefijo `/api` (StripPrefix) y el servido del frontend.

### Config
Nuevo campo `JWTSecret` (`env("JWT_SECRET", "<default-dev>")`). En producción
(Railway) deberá definirse un secreto fuerte.

### Dependencias nuevas
`github.com/golang-jwt/jwt/v5` y `golang.org/x/crypto/bcrypt`.

## Frontend

- **types.ts:** `User`, `AuthResponse`, `RegisterRequest`, `LoginRequest`.
- **api/auth.ts:** `register(body)`, `login(body)`.
- **api/client.ts:** adjunta `Authorization: Bearer <token>` si hay token en
  `localStorage` (clave `token`).
- **AuthContext + useAuth():** estado `{ user, token }` persistido en
  `localStorage`; métodos `login`, `register`, `logout`, e `isAuthenticated`.
  `ToastProvider` y `AuthProvider` envuelven la app en `main.tsx`.
- **Páginas** `/login` y `/register` (públicas) con formularios (reutilizan
  `Field`, `Button`); al éxito guardan sesión y redirigen a `/characters`.
- **Layout:** si hay sesión, muestra el email + botón "Cerrar sesión"; si no,
  enlace "Iniciar sesión".
- **Listas (ResourcePage / ScenesPage):** el botón "Nuevo" y las acciones
  Editar/Borrar solo se renderizan si `isAuthenticated`; si no, se muestra un
  aviso "Inicia sesión para gestionar" con enlace a `/login`. Los GET siguen
  mostrando los datos a cualquiera.
- **Manejo de 401 en acciones:** si una escritura devuelve 401 (token expirado),
  se hace `logout` y se redirige a `/login` con un toast.

## Manejo de errores

- 401 Unauthorized: credenciales inválidas o token ausente/expirado.
- 409 Conflict: email ya registrado.
- 400: validación (email vacío, password < 8).
- El frontend muestra el mensaje del campo `error` y, en acciones protegidas,
  redirige a login cuando corresponde.

## Pruebas

**Backend:**
- `AuthService.Register` hashea (el hash no es igual al password y `bcrypt`
  lo verifica) y rechaza password corta.
- `AuthService.Login` con password incorrecta → error Unauthorized; correcta →
  token no vacío.
- Helper JWT: un token emitido se valida y devuelve el id; un token manipulado
  o expirado falla.
- `RequireAuth`: sin cabecera → 401; token inválido → 401; token válido → 200
  (llama al next).

**Frontend:**
- El formulario de login muestra error si los campos están vacíos.
- `client.ts` adjunta la cabecera `Authorization` cuando hay token guardado.

Toda la suite previa (backend + 14 frontend) debe seguir en verde.

## Seguridad

- bcrypt (coste por defecto) para contraseñas; nunca se devuelven hashes.
- JWT firmado con secreto de entorno; caducidad 24 h.
- Email único a nivel de BD.
- `JWT_SECRET` fuerte obligatorio en producción (documentado en README/deploy).

## Fuera de alcance (YAGNI)

Recuperación de contraseña, verificación por email, roles/permisos, refresh
tokens, y datos privados por usuario (multi-inquilino).

## Criterios de éxito

1. `POST /api/auth/register` y `/api/auth/login` devuelven un token JWT válido.
2. `GET` de recursos funciona sin token; `POST/PUT/DELETE` sin token → 401.
3. Con token válido, las escrituras funcionan.
4. En el frontend: registro/login/logout funcionan; los controles de gestión
   solo aparecen logueado; los datos se ven sin sesión.
5. Todas las pruebas (backend nuevas + previas, y 14+ frontend) pasan; `tsc` y
   `go build` limpios.
