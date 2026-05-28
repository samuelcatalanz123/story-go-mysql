# Diseño: Panel de administración (React + TypeScript) para la Story API

**Fecha:** 2026-05-28
**Estado:** Aprobado para escribir el plan de implementación

## Objetivo

Crear una interfaz web que sirva de "cara visual" para la Story API (Go +
MySQL). La interfaz será un **panel de administración CRUD** para los tres
recursos del backend: **personajes**, **lugares** y **escenas**.

CRUD = Create, Read, Update, Delete, que se corresponden con los métodos HTTP
de la API: `POST`, `GET`, `PUT`, `DELETE`.

## Decisiones tomadas (brainstorming)

| Tema        | Decisión                          | Motivo                                            |
| ----------- | --------------------------------- | ------------------------------------------------- |
| Tipo de UI  | Panel de administración (CRUD)    | Directo y educativo para empezar con React        |
| Estilos     | CSS plano / CSS Modules           | Máximo aprendizaje, sin "magia" oculta            |
| Navegación  | React Router (URLs reales)        | Estándar profesional; URLs compartibles           |
| Datos/API   | `fetch` nativo + hooks propios    | Se ve todo el ciclo carga/éxito/error sin librerías |
| Conexión    | Proxy de Vite (`/api` → `:8080`)  | Evita CORS sin modificar el backend de Go         |
| Testing     | Vitest + React Testing Library    | Base de pruebas con ejemplos, sin cobertura total |

## Stack técnico

- **Vite** como herramienta de construcción y servidor de desarrollo.
- **React 18 + TypeScript**.
- **React Router** para el enrutado.
- **CSS Modules** para estilos encapsulados por componente.
- **Vitest + React Testing Library** para pruebas.
- Sin librería de estado ni de data-fetching: estado local con hooks propios.

## Conexión con la API (CORS y proxy)

El backend de Go corre en `http://localhost:8080` y **no tiene CORS
configurado** (verificado en el código). El servidor de Vite corre en otro
puerto (5173), y el navegador bloquearía las llamadas entre puertos distintos.

**Solución:** configurar un proxy en `vite.config.ts`. Todo lo que el código
pida bajo `/api` se reenvía a `http://localhost:8080`, eliminando el prefijo
`/api`. Para el navegador todo parece venir del mismo origen, así que no hay
problema de CORS y **no se modifica el backend**.

Ejemplo conceptual:

```
fetch("/api/characters")  →  (proxy Vite)  →  http://localhost:8080/characters
```

## Estructura de carpetas

Carpeta nueva `web/` dentro del repo. El código Go no se toca.

```
web/
  index.html
  package.json
  tsconfig.json
  vite.config.ts          Proxy /api → :8080
  src/
    main.tsx              Punto de entrada; monta React + Router
    App.tsx               Layout (menú) + definición de rutas
    types.ts              Tipos TS que reflejan los modelos de Go
    api/
      client.ts           Envoltura de fetch (URL, JSON, errores)
      resources.ts        Funciones CRUD por recurso
    hooks/
      useList.ts          Hook de listado (cargando/error/datos)
    components/
      Layout.tsx          Menú de navegación + área de contenido
      DataTable.tsx       Tabla reutilizable de listado
      ResourceForm.tsx    Formulario crear/editar (título + texto)
      SceneForm.tsx       Formulario de escena (campos extra + relaciones)
    pages/
      CharactersPage.tsx
      LocationsPage.tsx
      ScenesPage.tsx
```

## Tipos de TypeScript (contrato con la API)

Reflejan exactamente los modelos JSON del backend. Nota: en Go los campos
`*string` (que pueden ser nulos) se representan como `string | null`, y las
fechas llegan como cadenas de texto.

```typescript
type Character = {
  id: number;
  title: string;
  text: string | null;
  createdAt: string;
  updatedAt: string;
};

// Location es idéntico a Character.

type Scene = {
  id: number;
  title: string;
  text: string | null;
  startTimeline: number;
  endTimeline: number;
  characters: Character[];
  locations: Location[];
  createdAt: string;
  updatedAt: string;
};

// Payloads de creación/actualización:
type CharacterRequest = { title: string; text: string | null };
type LocationRequest  = { title: string; text: string | null };
type SceneRequest = {
  title: string;
  text: string | null;
  startTimeline: number;
  endTimeline: number;
  characterIds: number[];
  locationIds: number[];
};
```

## Flujo de datos

```
Página (pages/*)  →  hook (useList) / función api  →  client.ts (fetch)
                  →  proxy Vite  →  API Go  →  MySQL
```

1. La página pide datos a través del hook o de las funciones de `api/`.
2. `client.ts` arma la URL (`/api/...`), envía la petición y parsea el JSON.
3. Si la respuesta es un error (4xx/5xx), `client.ts` lanza una excepción con
   el mensaje del campo `error` que devuelve la API.
4. El componente muestra estado de *cargando*, los datos, o el *error*.

## Funcionalidad por sección

### Personajes y Lugares (idénticos)
- Tabla con la lista (id, título, texto, fecha de actualización).
- Botón "Nuevo" que abre el formulario.
- Acciones por fila: "Editar" y "Borrar" (con confirmación).
- Formulario: `título` (obligatorio) y `texto` (opcional).

### Escenas
- Igual que arriba, más:
  - Campos numéricos `startTimeline` y `endTimeline`.
  - Dos selectores múltiples para elegir personajes y lugares relacionados.
    Las opciones se obtienen llamando a `/api/characters` y `/api/locations`.
- La vista de detalle/lista muestra los personajes y lugares asociados.

## Manejo de errores

La API responde con `{ "error": "mensaje" }` y códigos HTTP:

- `400` validación o JSON inválido
- `404` recurso inexistente
- `405` método no permitido
- `409` título duplicado
- `500` error inesperado

`client.ts` detecta los códigos no exitosos y propaga el mensaje. Los
formularios muestran el error junto al campo o como aviso (ej. título
duplicado → `409`), y las listas muestran un mensaje si falla la carga.

## Pruebas

Configurar **Vitest + React Testing Library** con pruebas de ejemplo:

- El `ResourceForm` muestra un error si se envía con el título vacío.
- `DataTable` renderiza las filas que recibe.
- (Opcional) `client.ts` lanza el error correcto ante una respuesta 4xx.

El objetivo es dejar la base y enseñar el concepto, no lograr cobertura total.

## Fuera de alcance (YAGNI)

- Autenticación / login (la API no la tiene).
- Paginación, búsqueda y filtros avanzados.
- Vista de línea de tiempo gráfica o grafo de relaciones (sería otro proyecto).
- Despliegue a producción / build optimizado (nos centramos en desarrollo local).

## Criterios de éxito

1. `npm run dev` levanta la interfaz y se conecta a la API local vía proxy.
2. Se pueden crear, listar, editar y borrar personajes, lugares y escenas.
3. Al crear una escena se pueden asociar personajes y lugares existentes.
4. Los errores de la API se muestran de forma comprensible.
5. Las pruebas de ejemplo pasan con `npm test`.
