# Diseño: Búsqueda y paginación

**Fecha:** 2026-05-29
**Estado:** Aprobado para escribir el plan de implementación
**Contexto:** 6º sub-proyecto del portafolio. Añade búsqueda por texto y
paginación numerada a las tres listas (personajes, lugares, escenas).

## Objetivo

Que las listas se puedan buscar y paginar, como una app real con muchos datos.
Demuestra manejo de datos a escala (LIMIT/OFFSET + COUNT) en el backend y
controles de UI (búsqueda con debounce + paginación) en el frontend.

## Decisiones tomadas (brainstorming)

| Tema        | Decisión                                                      |
| ----------- | ------------------------------------------------------------- |
| Alcance     | Las 3 listas (characters, locations, scenes), uniforme        |
| Búsqueda    | Param `q`; `WHERE title LIKE ? OR text LIKE ?` (case-insensit.)|
| Paginación  | `page` (def. 1) + `pageSize` (def. 20, máx. 100); numerada     |
| Respuesta   | `{ items, total, page, pageSize }` (cambia el contrato)        |
| Orden       | Por `id` ascendente (sin orden configurable, por ahora)       |
| UI búsqueda | `SearchBar` con debounce ~300 ms; al buscar, vuelve a página 1 |

## API

Los `GET` de listas aceptan query params y devuelven un objeto paginado:

```
GET /api/characters?q=asha&page=1&pageSize=20
→ { "items": [ ... ], "total": 42, "page": 1, "pageSize": 20 }
```

- `q` vacío/ausente: sin filtro de texto.
- `page` < 1 se normaliza a 1; `pageSize` se acota a [1, 100] (default 20).
- Los GET por id (`/characters/{id}`) no cambian.

## Backend (capas)

- **model:** tipo genérico `Page[T]` (`Items []T`, `Total int`, `Page int`,
  `PageSize int`) con tags JSON `items/total/page/pageSize`. Y un struct
  `ListParams` (`Query string`, `Page int`, `PageSize int`) con un método que
  normaliza (page≥1, pageSize en [1,100]) y calcula `Limit()`/`Offset()`.
- **repository:** el `List` de cada recurso pasa a `List(ctx, q string, limit,
  offset int) ([]T, error)` y se añade `Count(ctx, q string) (int, error)`. La
  búsqueda usa `WHERE title LIKE ? OR (text IS NOT NULL AND text LIKE ?)` con el
  patrón `%q%`; sin `q`, no añade WHERE. Mantiene `ORDER BY id ASC LIMIT ?
  OFFSET ?`.
- **service:** el `List(ctx, params model.ListParams)` normaliza params, llama a
  `repo.Count` y `repo.List`, y devuelve `model.Page[T]`.
- **handler:** un helper `parseListParams(r)` lee `q`, `page`, `pageSize` de la
  query string (con defaults y parseo tolerante), llama al service y responde el
  `Page[T]`.
- **scenes:** el `List` paginado mantiene la carga de personajes/lugares de cada
  escena de la página (no de todas).

## Frontend

- **types.ts:** `Paged<T> = { items: T[]; total: number; page: number; pageSize:
  number }`.
- **api/resources.ts:** las funciones `list*` pasan a aceptar
  `(params: { q: string; page: number; pageSize: number })` y devolver
  `Paged<T>`. Se añaden helpers `listAllCharacters()` / `listAllLocations()` que
  llaman con `pageSize` grande (p. ej. 1000) y devuelven `.items`, para los
  desplegables del formulario de escenas.
- **hooks/usePagedList.ts:** hook que gestiona `page`, `query`, `data` (items),
  `total`, `pageSize`, `loading`, `error`; expone `setQuery` (reinicia a página
  1), `setPage`, `reload`. Mantiene la protección contra respuestas obsoletas
  (igual que `useList`).
- **ui/SearchBar.tsx:** input controlado con debounce ~300 ms que llama a
  `onQueryChange`.
- **ui/Pagination.tsx:** "Página X de Y" + botones Anterior/Siguiente
  (deshabilitados en los extremos); `totalPages = max(1, ceil(total/pageSize))`.
- **ResourcePage / ScenesPage:** usan `usePagedList`; renderizan `SearchBar`
  arriba y `Pagination` bajo la tabla. ScenesPage usa `listAll*` para las
  opciones del formulario (con `useList`, sin cambios) y `usePagedList` para la
  tabla de escenas.

## Manejo de errores

- Sin cambios en el mapeo de errores. Un `q`/`page`/`pageSize` no numérico se
  trata como ausente (se usa el default), no como error 400.
- Errores de carga siguen mostrando el `EmptyState` con "Reintentar".

## Pruebas

**Backend:**
- Repositorio: `List` con `q` filtra (usando una BD de prueba no es posible sin
  MySQL, así que se testea la lógica de `ListParams.Normalize/Limit/Offset` y la
  construcción de la cláusula con una función pura `buildSearch(q)` que devuelve
  el fragmento WHERE y los args).
- Service: normaliza page 0 → 1 y pageSize 999 → 100.
- Handler: `parseListParams` con query vacía → defaults; con valores → los lee.

**Frontend:**
- `usePagedList`: al `setQuery` vuelve a `page` 1; `setPage` cambia la página.
- `Pagination`: deshabilita "Anterior" en página 1 y "Siguiente" en la última;
  llama a `onPage` con la página correcta.
- `SearchBar`: tras el debounce, llama a `onQueryChange` con el texto.

Toda la suite previa (backend + frontend) debe seguir en verde; los tests de
listas existentes que esperaban un array se actualizan al nuevo contrato.

## Fuera de alcance (YAGNI)

Orden configurable por columnas, filtros por rango de timeline en escenas,
resaltado de coincidencias, y búsqueda full-text. Quedan para después.

## Criterios de éxito

1. `GET /api/characters?q=&page=1&pageSize=20` devuelve `{ items, total, page,
   pageSize }`; `q` filtra por título/texto; `page`/`pageSize` paginan.
2. Valores fuera de rango se normalizan (no fallan).
3. En el frontend: barra de búsqueda con debounce, controles de página que
   funcionan y se deshabilitan en los extremos, y el formulario de escenas
   sigue mostrando todas las opciones.
4. Todas las pruebas pasan; `tsc` y `go build` limpios.
