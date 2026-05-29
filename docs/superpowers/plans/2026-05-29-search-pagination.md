# Búsqueda y paginación — Plan de implementación

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Añadir búsqueda por texto (`q`) y paginación numerada (`page`/`pageSize`) a las tres listas, cambiando la respuesta a `{ items, total, page, pageSize }`, con búsqueda con debounce y controles de página en el frontend.

**Architecture:** Backend Go: tipo genérico `model.Page[T]` y `model.ListParams` (normaliza/calcula limit/offset); repositorios con `List(q,limit,offset)`+`Count(q)` y un helper puro `buildSearch`; servicios devuelven `Page[T]`; handlers leen la query. Frontend: tipo `Paged<T>`, hook `usePagedList`, componentes `SearchBar` (debounce) y `Pagination`, y las páginas con búsqueda + paginación.

**Tech Stack:** Go (database/sql, generics), MySQL; React 19 + TypeScript, Vitest + Testing Library.

---

## Estructura de archivos

```
Backend:
  internal/model/model.go              (mod: Page[T], ListParams)
  internal/model/model_test.go         (nuevo)
  internal/repository/repository.go    (mod: buildSearch)
  internal/repository/repository_test.go (nuevo)
  internal/repository/character.go     (mod: List+Count)
  internal/repository/location.go      (mod: List+Count)
  internal/repository/scene.go         (mod: ListIDs+Count)
  internal/service/character.go        (mod: List(params)→Page)
  internal/service/location.go         (mod)
  internal/service/scene.go            (mod)
  internal/handler/handler.go          (mod: parseListParams)
  internal/handler/listparams_test.go  (nuevo)
  internal/handler/character.go        (mod: List)
  internal/handler/location.go         (mod: List)
  internal/handler/scene.go            (mod: List)

Frontend (web/src):
  types.ts                             (mod: Paged<T>)
  api/resources.ts                     (mod: list args + listAll helpers)
  hooks/usePagedList.ts                (nuevo) + usePagedList.test.tsx
  ui/SearchBar.tsx + .module.css + test
  ui/Pagination.tsx + .module.css + test
  components/ResourcePage.tsx          (mod)
  pages/ScenesPage.tsx                 (mod)
```

Rama `feat/search-pagination` (ya creada). Comandos Go desde la raíz; npm desde `web/`.

---

### Task 1: model.Page/ListParams y buildSearch (TDD)

**Files:**
- Test: `internal/model/model_test.go`, `internal/repository/repository_test.go`
- Modify: `internal/model/model.go`, `internal/repository/repository.go`

- [ ] **Step 1: Escribir el test de model — `internal/model/model_test.go`**

```go
package model

import "testing"

func TestListParamsNormalize(t *testing.T) {
	got := ListParams{Page: 0, PageSize: 999}.Normalize()
	if got.Page != 1 {
		t.Fatalf("page=%d, esperaba 1", got.Page)
	}
	if got.PageSize != 100 {
		t.Fatalf("pageSize=%d, esperaba 100", got.PageSize)
	}
	if d := (ListParams{}).Normalize(); d.PageSize != 20 {
		t.Fatalf("default pageSize=%d, esperaba 20", d.PageSize)
	}
}

func TestListParamsLimitOffset(t *testing.T) {
	p := ListParams{Page: 3, PageSize: 20}.Normalize()
	if p.Limit() != 20 || p.Offset() != 40 {
		t.Fatalf("limit=%d offset=%d, esperaba 20 y 40", p.Limit(), p.Offset())
	}
}
```

- [ ] **Step 2: Ejecutar y ver que falla**

Run: `go test ./internal/model/`
Expected: FALLA (no existe `ListParams`).

- [ ] **Step 3: Añadir tipos al final de `internal/model/model.go`**

```go
// Page is a paginated slice of results returned by list endpoints.
type Page[T any] struct {
	Items    []T `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// Pagination defaults and bounds for list endpoints.
const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// ListParams holds search and pagination parameters for list endpoints.
type ListParams struct {
	Query    string
	Page     int
	PageSize int
}

// Normalize clamps the params to safe values: Page >= 1 and PageSize within
// [1, MaxPageSize] (defaulting to DefaultPageSize). It returns a copy.
func (p ListParams) Normalize() ListParams {
	out := p
	if out.Page < 1 {
		out.Page = 1
	}
	if out.PageSize <= 0 {
		out.PageSize = DefaultPageSize
	}
	if out.PageSize > MaxPageSize {
		out.PageSize = MaxPageSize
	}
	return out
}

// Limit returns the SQL LIMIT (call after Normalize).
func (p ListParams) Limit() int { return p.PageSize }

// Offset returns the SQL OFFSET (call after Normalize).
func (p ListParams) Offset() int { return (p.Page - 1) * p.PageSize }
```

- [ ] **Step 4: Ejecutar y ver que pasa**

Run: `go test ./internal/model/`
Expected: PASA (2 tests).

- [ ] **Step 5: Escribir el test del repositorio — `internal/repository/repository_test.go`**

```go
package repository

import "testing"

func TestBuildSearchEmpty(t *testing.T) {
	clause, args := buildSearch("")
	if clause != "" {
		t.Fatalf("clause=%q, esperaba vacío", clause)
	}
	if args != nil {
		t.Fatalf("args=%v, esperaba nil", args)
	}
}

func TestBuildSearchNonEmpty(t *testing.T) {
	clause, args := buildSearch("asha")
	if clause == "" {
		t.Fatal("esperaba una cláusula WHERE")
	}
	if len(args) != 2 || args[0] != "%asha%" || args[1] != "%asha%" {
		t.Fatalf("args=%v, esperaba dos patrones %%asha%%", args)
	}
}
```

- [ ] **Step 6: Ejecutar y ver que falla**

Run: `go test ./internal/repository/`
Expected: FALLA (no existe `buildSearch`).

- [ ] **Step 7: Añadir `buildSearch` al final de `internal/repository/repository.go`**

```go
// buildSearch returns the WHERE clause (with a trailing space) and args for a
// title/text search. An empty query yields an empty clause and nil args, so
// callers can concatenate " ORDER BY ... LIMIT ? OFFSET ?" right after it.
func buildSearch(q string) (string, []any) {
	if q == "" {
		return "", nil
	}
	pattern := "%" + q + "%"
	return "WHERE title LIKE ? OR (text IS NOT NULL AND text LIKE ?) ", []any{pattern, pattern}
}
```

- [ ] **Step 8: Ejecutar y ver que pasa**

Run: `go test ./internal/repository/`
Expected: PASA (2 tests).

- [ ] **Step 9: Commit**

```bash
git add -A && git commit -m "feat: model Page/ListParams y helper buildSearch (TDD)"
```

---

### Task 2: Repos y servicios de characters y locations

**Files:**
- Modify: `internal/repository/character.go`, `internal/repository/location.go`, `internal/service/character.go`, `internal/service/location.go`

- [ ] **Step 1: En `internal/repository/character.go`, reemplazar el método `List`** por estos dos (List paginado + Count)

```go
// List returns a page of characters matching q (empty q = no filter),
// ordered by ID, limited to limit rows starting at offset.
func (r *CharacterRepository) List(ctx context.Context, q string, limit, offset int) ([]model.Character, error) {
	where, args := buildSearch(q)
	query := "SELECT id, title, text, created_at, updated_at FROM characters " +
		where + "ORDER BY id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, translate(err)
	}
	defer rows.Close()

	characters := []model.Character{}
	for rows.Next() {
		var c model.Character
		if err := rows.Scan(&c.ID, &c.Title, &c.Text, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		characters = append(characters, c)
	}
	return characters, rows.Err()
}

// Count returns the number of characters matching q.
func (r *CharacterRepository) Count(ctx context.Context, q string) (int, error) {
	where, args := buildSearch(q)
	var n int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM characters "+where, args...).Scan(&n); err != nil {
		return 0, translate(err)
	}
	return n, nil
}
```

- [ ] **Step 2: En `internal/repository/location.go`, reemplazar `List`** por (idéntico patrón, tabla `locations`):

```go
// List returns a page of locations matching q, ordered by ID.
func (r *LocationRepository) List(ctx context.Context, q string, limit, offset int) ([]model.Location, error) {
	where, args := buildSearch(q)
	query := "SELECT id, title, text, created_at, updated_at FROM locations " +
		where + "ORDER BY id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, translate(err)
	}
	defer rows.Close()

	locations := []model.Location{}
	for rows.Next() {
		var l model.Location
		if err := rows.Scan(&l.ID, &l.Title, &l.Text, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, rows.Err()
}

// Count returns the number of locations matching q.
func (r *LocationRepository) Count(ctx context.Context, q string) (int, error) {
	where, args := buildSearch(q)
	var n int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM locations "+where, args...).Scan(&n); err != nil {
		return 0, translate(err)
	}
	return n, nil
}
```
(Si el `location.go` actual no tenía un `List` con esa firma exacta, reemplaza el `List` existente por completo. Mantén el resto de métodos intactos.)

- [ ] **Step 3: En `internal/service/character.go`, reemplazar el método `List`** por:

```go
// List returns a page of characters matching the given params.
func (s *CharacterService) List(ctx context.Context, params model.ListParams) (model.Page[model.Character], error) {
	p := params.Normalize()
	total, err := s.repo.Count(ctx, p.Query)
	if err != nil {
		return model.Page[model.Character]{}, err
	}
	items, err := s.repo.List(ctx, p.Query, p.Limit(), p.Offset())
	if err != nil {
		return model.Page[model.Character]{}, err
	}
	return model.Page[model.Character]{Items: items, Total: total, Page: p.Page, PageSize: p.PageSize}, nil
}
```

- [ ] **Step 4: En `internal/service/location.go`, reemplazar `List`** por (tipo `Location`):

```go
// List returns a page of locations matching the given params.
func (s *LocationService) List(ctx context.Context, params model.ListParams) (model.Page[model.Location], error) {
	p := params.Normalize()
	total, err := s.repo.Count(ctx, p.Query)
	if err != nil {
		return model.Page[model.Location]{}, err
	}
	items, err := s.repo.List(ctx, p.Query, p.Limit(), p.Offset())
	if err != nil {
		return model.Page[model.Location]{}, err
	}
	return model.Page[model.Location]{Items: items, Total: total, Page: p.Page, PageSize: p.PageSize}, nil
}
```
(El campo del repositorio en `LocationService` se llama `repo`, igual que en `CharacterService`. Si difiere, usa el nombre real.)

- [ ] **Step 5: Verificar compilación**

Run: `go build ./internal/repository/ ./internal/model/`
Expected: sin errores. (El paquete service y handler aún no compilan hasta la Task 4; eso es esperado.)

- [ ] **Step 6: Commit**

```bash
git add -A && git commit -m "feat: paginación/búsqueda en repos y servicios de characters/locations"
```

---

### Task 3: Repo y servicio de scenes

**Files:**
- Modify: `internal/repository/scene.go`, `internal/service/scene.go`

- [ ] **Step 1: En `internal/repository/scene.go`, reemplazar el método `ListIDs`** por (paginado/filtrado) y añadir `Count`:

```go
// ListIDs returns a page of scene IDs matching q, ordered ascending. The
// service composes full scenes from these IDs via GetByID.
func (r *SceneRepository) ListIDs(ctx context.Context, q string, limit, offset int) ([]uint64, error) {
	where, args := buildSearch(q)
	query := "SELECT id FROM scenes " + where + "ORDER BY id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, translate(err)
	}
	defer rows.Close()

	ids := []uint64{}
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// Count returns the number of scenes matching q.
func (r *SceneRepository) Count(ctx context.Context, q string) (int, error) {
	where, args := buildSearch(q)
	var n int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM scenes "+where, args...).Scan(&n); err != nil {
		return 0, translate(err)
	}
	return n, nil
}
```

- [ ] **Step 2: En `internal/service/scene.go`, reemplazar el método `List`** por:

```go
// List returns a page of scenes matching the given params, each populated
// with its relations.
func (s *SceneService) List(ctx context.Context, params model.ListParams) (model.Page[model.Scene], error) {
	p := params.Normalize()
	total, err := s.scenes.Count(ctx, p.Query)
	if err != nil {
		return model.Page[model.Scene]{}, err
	}
	ids, err := s.scenes.ListIDs(ctx, p.Query, p.Limit(), p.Offset())
	if err != nil {
		return model.Page[model.Scene]{}, err
	}
	scenes := make([]model.Scene, 0, len(ids))
	for _, id := range ids {
		scene, err := s.scenes.GetByID(ctx, id)
		if err != nil {
			return model.Page[model.Scene]{}, err
		}
		scenes = append(scenes, scene)
	}
	return model.Page[model.Scene]{Items: scenes, Total: total, Page: p.Page, PageSize: p.PageSize}, nil
}
```

- [ ] **Step 3: Verificar compilación de repos/servicios**

Run: `go build ./internal/repository/`
Expected: sin errores. (service/handler completarán en Task 4.)

- [ ] **Step 4: Commit**

```bash
git add -A && git commit -m "feat: paginación/búsqueda en repo y servicio de scenes"
```

---

### Task 4: Handlers (parseListParams TDD) y List paginado

**Files:**
- Modify: `internal/handler/handler.go`, `internal/handler/character.go`, `internal/handler/location.go`, `internal/handler/scene.go`
- Test: `internal/handler/listparams_test.go`

- [ ] **Step 1: Escribir el test — `internal/handler/listparams_test.go`**

```go
package handler

import (
	"net/http/httptest"
	"testing"
)

func TestParseListParamsDefaults(t *testing.T) {
	p := parseListParams(httptest.NewRequest("GET", "/characters", nil))
	if p.Query != "" || p.Page != 0 || p.PageSize != 0 {
		t.Fatalf("inesperado: %+v", p)
	}
}

func TestParseListParamsReadsValues(t *testing.T) {
	p := parseListParams(httptest.NewRequest("GET", "/characters?q=asha&page=2&pageSize=5", nil))
	if p.Query != "asha" || p.Page != 2 || p.PageSize != 5 {
		t.Fatalf("inesperado: %+v", p)
	}
}
```

- [ ] **Step 2: Ejecutar y ver que falla**

Run: `go test ./internal/handler/`
Expected: FALLA (no existe `parseListParams`). (También fallará el build del paquete por las firmas de List; lo arreglan los pasos siguientes.)

- [ ] **Step 3: Añadir `parseListParams` a `internal/handler/handler.go`**

Añade `"story-go-mysql/internal/model"` a los imports y esta función:
```go
// parseListParams reads q/page/pageSize from the query string. Missing or
// non-numeric page/pageSize stay at 0, which the service normalizes to
// sensible defaults.
func parseListParams(r *http.Request) model.ListParams {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	return model.ListParams{Query: q.Get("q"), Page: page, PageSize: pageSize}
}
```
(`strconv` ya está importado en handler.go.)

- [ ] **Step 4: En `internal/handler/character.go`, reemplazar el método `List`** por:

```go
func (h *CharacterHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := h.svc.List(r.Context(), parseListParams(r))
	if err != nil {
		web.RespondError(w, characterResource, err)
		return
	}
	web.JSON(w, http.StatusOK, page)
}
```

- [ ] **Step 5: En `internal/handler/location.go`, reemplazar el método `List`** por (recurso location):

```go
func (h *LocationHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := h.svc.List(r.Context(), parseListParams(r))
	if err != nil {
		web.RespondError(w, locationResource, err)
		return
	}
	web.JSON(w, http.StatusOK, page)
}
```
(Usa la constante de recurso que ya exista en ese archivo, p. ej. `locationResource`.)

- [ ] **Step 6: En `internal/handler/scene.go`, reemplazar el método `List`** por (recurso scene):

```go
func (h *SceneHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := h.svc.List(r.Context(), parseListParams(r))
	if err != nil {
		web.RespondError(w, sceneResource, err)
		return
	}
	web.JSON(w, http.StatusOK, page)
}
```
(Usa la constante de recurso que ya exista, p. ej. `sceneResource`.)

- [ ] **Step 7: Ejecutar y ver que pasa**

Run: `go build ./... && go test ./...`
Expected: compila; `parseListParams` 2 tests pasan; el resto sigue en verde.

- [ ] **Step 8: Commit**

```bash
git add -A && git commit -m "feat: handlers de listas con búsqueda/paginación"
```

---

### Task 5: Frontend — tipo Paged y api/resources

**Files:**
- Modify: `web/src/types.ts`, `web/src/api/resources.ts`

- [ ] **Step 1: Añadir el tipo al final de `web/src/types.ts`**

```ts
export type Paged<T> = {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
};
```

- [ ] **Step 2: Reemplazar TODO `web/src/api/resources.ts`**

```ts
import { apiFetch } from "./client";
import type {
  Character,
  Location,
  Scene,
  CharacterRequest,
  LocationRequest,
  SceneRequest,
  Paged,
} from "../types";

export type ListArgs = { q: string; page: number; pageSize: number };

function listQuery({ q, page, pageSize }: ListArgs): string {
  return new URLSearchParams({
    q,
    page: String(page),
    pageSize: String(pageSize),
  }).toString();
}

// --- Personajes ---
export const listCharacters = (args: ListArgs) =>
  apiFetch<Paged<Character>>(`/characters?${listQuery(args)}`);
export const createCharacter = (body: CharacterRequest) =>
  apiFetch<Character>("/characters", { method: "POST", body: JSON.stringify(body) });
export const updateCharacter = (id: number, body: CharacterRequest) =>
  apiFetch<Character>(`/characters/${id}`, { method: "PUT", body: JSON.stringify(body) });
export const deleteCharacter = (id: number) =>
  apiFetch<void>(`/characters/${id}`, { method: "DELETE" });

// --- Lugares ---
export const listLocations = (args: ListArgs) =>
  apiFetch<Paged<Location>>(`/locations?${listQuery(args)}`);
export const createLocation = (body: LocationRequest) =>
  apiFetch<Location>("/locations", { method: "POST", body: JSON.stringify(body) });
export const updateLocation = (id: number, body: LocationRequest) =>
  apiFetch<Location>(`/locations/${id}`, { method: "PUT", body: JSON.stringify(body) });
export const deleteLocation = (id: number) =>
  apiFetch<void>(`/locations/${id}`, { method: "DELETE" });

// --- Escenas ---
export const listScenes = (args: ListArgs) =>
  apiFetch<Paged<Scene>>(`/scenes?${listQuery(args)}`);
export const createScene = (body: SceneRequest) =>
  apiFetch<Scene>("/scenes", { method: "POST", body: JSON.stringify(body) });
export const updateScene = (id: number, body: SceneRequest) =>
  apiFetch<Scene>(`/scenes/${id}`, { method: "PUT", body: JSON.stringify(body) });
export const deleteScene = (id: number) =>
  apiFetch<void>(`/scenes/${id}`, { method: "DELETE" });

// Helpers que traen TODOS los elementos (para los desplegables del formulario
// de escenas), usando un pageSize grande.
const ALL: ListArgs = { q: "", page: 1, pageSize: 1000 };
export const listAllCharacters = () =>
  apiFetch<Paged<Character>>(`/characters?${listQuery(ALL)}`).then((p) => p.items);
export const listAllLocations = () =>
  apiFetch<Paged<Location>>(`/locations?${listQuery(ALL)}`).then((p) => p.items);
```

- [ ] **Step 3: Verificar tipos** (habrá errores en páginas hasta Tasks 7-8; verifica solo este archivo de forma aislada saltando luego)

Run dentro de `web/`: `npx tsc --noEmit` → mostrará errores en ResourcePage/ScenesPage (esperado hasta las Tasks 7-8). El archivo `api/resources.ts` en sí no debe tener errores propios.

- [ ] **Step 4: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): tipo Paged y api de listas con búsqueda/paginación"
```

---

### Task 6: Hook usePagedList y componentes SearchBar/Pagination (TDD)

**Files:**
- Create: `web/src/hooks/usePagedList.ts`, `web/src/hooks/usePagedList.test.tsx`
- Create: `web/src/ui/SearchBar.tsx`, `SearchBar.module.css`, `SearchBar.test.tsx`
- Create: `web/src/ui/Pagination.tsx`, `Pagination.module.css`, `Pagination.test.tsx`

- [ ] **Step 1: Crear `web/src/hooks/usePagedList.ts`**

```ts
import { useCallback, useEffect, useRef, useState } from "react";
import type { Paged } from "../types";

const PAGE_SIZE = 20;

export type PagedListState<T> = {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
  query: string;
  loading: boolean;
  error: string | null;
  setQuery: (q: string) => void;
  setPage: (page: number) => void;
  reload: () => void;
};

// Carga una lista paginada y con búsqueda. `loader` debe ser una referencia
// estable (una función de api/resources.ts). Descarta respuestas obsoletas.
export function usePagedList<T>(
  loader: (args: { q: string; page: number; pageSize: number }) => Promise<Paged<T>>,
): PagedListState<T> {
  const [data, setData] = useState<T[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [query, setQueryState] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const requestId = useRef(0);

  const load = useCallback(() => {
    const id = ++requestId.current;
    setLoading(true);
    setError(null);
    loader({ q: query, page, pageSize: PAGE_SIZE })
      .then((res) => {
        if (id !== requestId.current) return;
        setData(res.items);
        setTotal(res.total);
      })
      .catch((e: unknown) => {
        if (id === requestId.current)
          setError(e instanceof Error ? e.message : "Error desconocido");
      })
      .finally(() => {
        if (id === requestId.current) setLoading(false);
      });
  }, [loader, query, page]);

  useEffect(() => {
    load();
    return () => {
      requestId.current++;
    };
  }, [load]);

  const setQuery = useCallback((q: string) => {
    setQueryState(q);
    setPage(1); // nueva búsqueda → primera página
  }, []);

  return {
    data,
    total,
    page,
    pageSize: PAGE_SIZE,
    query,
    loading,
    error,
    setQuery,
    setPage,
    reload: load,
  };
}
```

- [ ] **Step 2: Escribir el test — `web/src/hooks/usePagedList.test.tsx`**

```tsx
import { renderHook, act, waitFor } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { usePagedList } from "./usePagedList";

describe("usePagedList", () => {
  it("al buscar vuelve a la página 1", async () => {
    const loader = vi.fn().mockResolvedValue({ items: [], total: 0, page: 1, pageSize: 20 });
    const { result } = renderHook(() => usePagedList(loader));
    await waitFor(() => expect(result.current.loading).toBe(false));

    act(() => result.current.setPage(3));
    expect(result.current.page).toBe(3);

    act(() => result.current.setQuery("asha"));
    expect(result.current.page).toBe(1);
    expect(result.current.query).toBe("asha");
  });
});
```

- [ ] **Step 3: Ejecutar el test del hook**

Run dentro de `web/`: `npx vitest run src/hooks/usePagedList.test.tsx`
Expected: PASA (1 test).

- [ ] **Step 4: Crear `web/src/ui/SearchBar.tsx`**

```tsx
import { useEffect, useState } from "react";
import styles from "./SearchBar.module.css";

type Props = { onQueryChange: (q: string) => void; placeholder?: string };

// Input de búsqueda con debounce: espera ~300 ms tras la última tecla antes de
// avisar, para no consultar la API en cada pulsación.
export function SearchBar({ onQueryChange, placeholder = "Buscar…" }: Props) {
  const [value, setValue] = useState("");

  useEffect(() => {
    const timer = setTimeout(() => onQueryChange(value.trim()), 300);
    return () => clearTimeout(timer);
  }, [value, onQueryChange]);

  return (
    <input
      className={styles.search}
      type="search"
      value={value}
      placeholder={placeholder}
      aria-label="Buscar"
      onChange={(e) => setValue(e.target.value)}
    />
  );
}
```

- [ ] **Step 5: Crear `web/src/ui/SearchBar.module.css`**

```css
.search {
  max-width: 320px;
  margin-bottom: var(--space-4);
}
```

- [ ] **Step 6: Escribir el test — `web/src/ui/SearchBar.test.tsx`**

```tsx
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { SearchBar } from "./SearchBar";

describe("SearchBar", () => {
  it("llama a onQueryChange con el texto tras el debounce", async () => {
    const onQueryChange = vi.fn();
    render(<SearchBar onQueryChange={onQueryChange} />);
    fireEvent.change(screen.getByLabelText("Buscar"), { target: { value: "asha" } });
    await waitFor(() => expect(onQueryChange).toHaveBeenCalledWith("asha"));
  });
});
```

- [ ] **Step 7: Crear `web/src/ui/Pagination.tsx`**

```tsx
import { Button } from "./Button";
import styles from "./Pagination.module.css";

type Props = {
  page: number;
  pageSize: number;
  total: number;
  onPage: (page: number) => void;
};

export function Pagination({ page, pageSize, total, onPage }: Props) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  return (
    <div className={styles.pagination}>
      <Button variant="secondary" size="sm" disabled={page <= 1} onClick={() => onPage(page - 1)}>
        Anterior
      </Button>
      <span className={styles.info}>
        Página {page} de {totalPages}
      </span>
      <Button
        variant="secondary"
        size="sm"
        disabled={page >= totalPages}
        onClick={() => onPage(page + 1)}
      >
        Siguiente
      </Button>
    </div>
  );
}
```

- [ ] **Step 8: Crear `web/src/ui/Pagination.module.css`**

```css
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  margin-top: var(--space-4);
}
.info {
  font-size: var(--text-sm);
  color: var(--color-text-muted);
}
```

- [ ] **Step 9: Escribir el test — `web/src/ui/Pagination.test.tsx`**

```tsx
import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { Pagination } from "./Pagination";

describe("Pagination", () => {
  it("deshabilita Anterior en la primera página", () => {
    render(<Pagination page={1} pageSize={20} total={50} onPage={() => {}} />);
    expect(screen.getByRole("button", { name: "Anterior" })).toBeDisabled();
    expect(screen.getByRole("button", { name: "Siguiente" })).not.toBeDisabled();
  });

  it("deshabilita Siguiente en la última página", () => {
    render(<Pagination page={3} pageSize={20} total={50} onPage={() => {}} />);
    expect(screen.getByRole("button", { name: "Siguiente" })).toBeDisabled();
  });

  it("llama a onPage al avanzar", () => {
    const onPage = vi.fn();
    render(<Pagination page={1} pageSize={20} total={50} onPage={onPage} />);
    fireEvent.click(screen.getByRole("button", { name: "Siguiente" }));
    expect(onPage).toHaveBeenCalledWith(2);
  });
});
```

- [ ] **Step 10: Ejecutar las pruebas nuevas**

Run dentro de `web/`: `npx vitest run src/ui/SearchBar.test.tsx src/ui/Pagination.test.tsx src/hooks/usePagedList.test.tsx`
Expected: PASAN (1 + 3 + 1 = 5 tests).

- [ ] **Step 11: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): usePagedList, SearchBar (debounce) y Pagination"
```

---

### Task 7: ResourcePage con búsqueda y paginación

**Files:**
- Modify: `web/src/components/ResourcePage.tsx`

- [ ] **Step 1: Reemplazar TODO `web/src/components/ResourcePage.tsx`**

```tsx
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { usePagedList } from "../hooks/usePagedList";
import { DataTable } from "./DataTable";
import type { Column } from "./DataTable";
import { ResourceForm } from "./ResourceForm";
import type { ResourceFormValues } from "./ResourceForm";
import { Modal } from "../ui/Modal";
import { ConfirmDialog } from "../ui/ConfirmDialog";
import { Button } from "../ui/Button";
import { PageHeader } from "../ui/PageHeader";
import { SkeletonRows } from "../ui/Skeleton";
import { EmptyState } from "../ui/EmptyState";
import { SearchBar } from "../ui/SearchBar";
import { Pagination } from "../ui/Pagination";
import { useToast } from "../ui/Toast";
import { useAuth } from "../auth/AuthContext";
import { ApiError } from "../api/client";
import type { Paged } from "../types";

type ResourceItem = {
  id: number;
  title: string;
  text: string | null;
  updatedAt: string;
};

type RequestBody = { title: string; text: string | null };
type ListArgs = { q: string; page: number; pageSize: number };

type Props<T extends ResourceItem> = {
  heading: string;
  list: (args: ListArgs) => Promise<Paged<T>>;
  create: (body: RequestBody) => Promise<T>;
  update: (id: number, body: RequestBody) => Promise<T>;
  remove: (id: number) => Promise<void>;
};

type Editing<T> = null | "new" | T;

function formatDate(value: string): string {
  const d = new Date(value);
  return isNaN(d.getTime()) ? value : d.toLocaleString();
}

export function ResourcePage<T extends ResourceItem>({
  heading,
  list,
  create,
  update,
  remove,
}: Props<T>) {
  const { data, total, page, pageSize, loading, error, setQuery, setPage, reload } =
    usePagedList(list);
  const toast = useToast();
  const { isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();
  const [editing, setEditing] = useState<Editing<T>>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [deleting, setDeleting] = useState<T | null>(null);

  const columns: Column<T>[] = [
    { header: "ID", render: (r) => r.id },
    { header: "Título", render: (r) => r.title },
    { header: "Texto", render: (r) => r.text ?? "—" },
    { header: "Actualizado", render: (r) => formatDate(r.updatedAt) },
  ];

  function openNew() {
    setFormError(null);
    setEditing("new");
  }

  function isUnauthorized(e: unknown): boolean {
    if (e instanceof ApiError && e.status === 401) {
      logout();
      toast.error("Tu sesión expiró, inicia sesión de nuevo");
      navigate("/login");
      return true;
    }
    return false;
  }

  function toBody(values: ResourceFormValues): RequestBody {
    return { title: values.title, text: values.text.trim() === "" ? null : values.text };
  }

  async function handleSubmit(values: ResourceFormValues) {
    setSubmitting(true);
    setFormError(null);
    try {
      if (editing === "new") {
        await create(toBody(values));
        toast.success("Creado correctamente");
      } else if (editing) {
        await update(editing.id, toBody(values));
        toast.success("Cambios guardados");
      }
      setEditing(null);
      reload();
    } catch (e: unknown) {
      if (isUnauthorized(e)) return;
      const msg = e instanceof Error ? e.message : "Error desconocido";
      setFormError(msg);
      toast.error(msg);
    } finally {
      setSubmitting(false);
    }
  }

  async function confirmDelete() {
    if (!deleting) return;
    try {
      await remove(deleting.id);
      toast.success("Eliminado");
      reload();
    } catch (e: unknown) {
      if (!isUnauthorized(e)) {
        toast.error(e instanceof Error ? e.message : "Error desconocido");
      }
    } finally {
      setDeleting(null);
    }
  }

  const initial =
    editing && editing !== "new"
      ? { title: editing.title, text: editing.text ?? "" }
      : undefined;

  return (
    <section>
      <PageHeader
        title={heading}
        action={
          isAuthenticated ? (
            <Button onClick={openNew}>Nuevo</Button>
          ) : (
            <Link to="/login">Inicia sesión para gestionar</Link>
          )
        }
      />

      <SearchBar onQueryChange={setQuery} />

      {loading && <SkeletonRows rows={4} cols={4} />}
      {error && (
        <EmptyState
          title="No se pudo cargar"
          message={error}
          action={
            <Button variant="secondary" onClick={reload}>
              Reintentar
            </Button>
          }
        />
      )}
      {!loading && !error && data.length === 0 && (
        <EmptyState
          title="No hay resultados"
          message="Prueba con otra búsqueda o crea el primer elemento."
          action={isAuthenticated ? <Button onClick={openNew}>Nuevo</Button> : undefined}
        />
      )}
      {!loading && !error && data.length > 0 && (
        <>
          <DataTable
            columns={columns}
            rows={data}
            onEdit={
              isAuthenticated
                ? (row) => {
                    setFormError(null);
                    setEditing(row);
                  }
                : undefined
            }
            onDelete={isAuthenticated ? (row) => setDeleting(row) : undefined}
          />
          <Pagination page={page} pageSize={pageSize} total={total} onPage={setPage} />
        </>
      )}

      <Modal
        open={editing !== null}
        onClose={() => {
          setEditing(null);
          setFormError(null);
        }}
        title={editing === "new" ? `Nuevo: ${heading}` : `Editar: ${heading}`}
      >
        <ResourceForm
          initial={initial}
          onSubmit={handleSubmit}
          onCancel={() => {
            setEditing(null);
            setFormError(null);
          }}
          submitting={submitting}
          error={formError}
        />
      </Modal>

      <ConfirmDialog
        open={deleting !== null}
        title="Confirmar borrado"
        message={deleting ? `¿Seguro que quieres borrar "${deleting.title}"?` : ""}
        confirmLabel="Borrar"
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(null)}
      />
    </section>
  );
}
```

- [ ] **Step 2: Verificar tipos**

Run dentro de `web/`: `npx tsc --noEmit`
Expected: sin errores en ResourcePage (ScenesPage puede seguir con error hasta la Task 8).

- [ ] **Step 3: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): búsqueda y paginación en ResourcePage"
```

---

### Task 8: ScenesPage con búsqueda y paginación

**Files:**
- Modify: `web/src/pages/ScenesPage.tsx`

- [ ] **Step 1: Reemplazar TODO `web/src/pages/ScenesPage.tsx`**

```tsx
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useList } from "../hooks/useList";
import { usePagedList } from "../hooks/usePagedList";
import { DataTable } from "../components/DataTable";
import type { Column } from "../components/DataTable";
import { SceneForm } from "../components/SceneForm";
import type { SceneFormValues } from "../components/SceneForm";
import type { Scene } from "../types";
import {
  listScenes,
  createScene,
  updateScene,
  deleteScene,
  listAllCharacters,
  listAllLocations,
} from "../api/resources";
import { Modal } from "../ui/Modal";
import { ConfirmDialog } from "../ui/ConfirmDialog";
import { Button } from "../ui/Button";
import { PageHeader } from "../ui/PageHeader";
import { SkeletonRows } from "../ui/Skeleton";
import { EmptyState } from "../ui/EmptyState";
import { Badge } from "../ui/Badge";
import { SearchBar } from "../ui/SearchBar";
import { Pagination } from "../ui/Pagination";
import { useToast } from "../ui/Toast";
import { useAuth } from "../auth/AuthContext";
import { ApiError } from "../api/client";

type Editing = null | "new" | Scene;

export function ScenesPage() {
  const scenes = usePagedList(listScenes);
  const characters = useList(listAllCharacters);
  const locations = useList(listAllLocations);
  const toast = useToast();
  const { isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();

  const [editing, setEditing] = useState<Editing>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [deleting, setDeleting] = useState<Scene | null>(null);

  const columns: Column<Scene>[] = [
    { header: "ID", render: (s) => s.id },
    { header: "Título", render: (s) => s.title },
    { header: "Inicio", render: (s) => s.startTimeline },
    { header: "Fin", render: (s) => s.endTimeline },
    {
      header: "Personajes",
      render: (s) =>
        s.characters.length
          ? s.characters.map((c) => (
              <Badge key={c.id} tone="primary">
                {c.title}
              </Badge>
            ))
          : "—",
    },
    {
      header: "Lugares",
      render: (s) =>
        s.locations.length ? s.locations.map((l) => <Badge key={l.id}>{l.title}</Badge>) : "—",
    },
  ];

  function openNew() {
    setFormError(null);
    setEditing("new");
  }

  function isUnauthorized(e: unknown): boolean {
    if (e instanceof ApiError && e.status === 401) {
      logout();
      toast.error("Tu sesión expiró, inicia sesión de nuevo");
      navigate("/login");
      return true;
    }
    return false;
  }

  async function handleSubmit(values: SceneFormValues) {
    setSubmitting(true);
    setFormError(null);
    const body = {
      title: values.title,
      text: values.text.trim() === "" ? null : values.text,
      startTimeline: values.startTimeline,
      endTimeline: values.endTimeline,
      characterIds: values.characterIds,
      locationIds: values.locationIds,
    };
    try {
      if (editing === "new") {
        await createScene(body);
        toast.success("Escena creada");
      } else if (editing) {
        await updateScene(editing.id, body);
        toast.success("Cambios guardados");
      }
      setEditing(null);
      scenes.reload();
    } catch (e: unknown) {
      if (isUnauthorized(e)) return;
      const msg = e instanceof Error ? e.message : "Error desconocido";
      setFormError(msg);
      toast.error(msg);
    } finally {
      setSubmitting(false);
    }
  }

  async function confirmDelete() {
    if (!deleting) return;
    try {
      await deleteScene(deleting.id);
      toast.success("Escena eliminada");
      scenes.reload();
    } catch (e: unknown) {
      if (!isUnauthorized(e)) {
        toast.error(e instanceof Error ? e.message : "Error desconocido");
      }
    } finally {
      setDeleting(null);
    }
  }

  const initial: SceneFormValues | undefined =
    editing && editing !== "new"
      ? {
          title: editing.title,
          text: editing.text ?? "",
          startTimeline: editing.startTimeline,
          endTimeline: editing.endTimeline,
          characterIds: editing.characters.map((c) => c.id),
          locationIds: editing.locations.map((l) => l.id),
        }
      : undefined;

  return (
    <section>
      <PageHeader
        title="Escenas"
        action={
          isAuthenticated ? (
            <Button onClick={openNew}>Nueva</Button>
          ) : (
            <Link to="/login">Inicia sesión para gestionar</Link>
          )
        }
      />

      <SearchBar onQueryChange={scenes.setQuery} />

      {scenes.loading && <SkeletonRows rows={4} cols={6} />}
      {scenes.error && (
        <EmptyState
          title="No se pudo cargar"
          message={scenes.error}
          action={
            <Button variant="secondary" onClick={scenes.reload}>
              Reintentar
            </Button>
          }
        />
      )}
      {!scenes.loading && !scenes.error && scenes.data.length === 0 && (
        <EmptyState
          title="No hay resultados"
          message="Prueba con otra búsqueda o crea la primera escena."
          action={isAuthenticated ? <Button onClick={openNew}>Nueva</Button> : undefined}
        />
      )}
      {!scenes.loading && !scenes.error && scenes.data.length > 0 && (
        <>
          <DataTable
            columns={columns}
            rows={scenes.data}
            onEdit={
              isAuthenticated
                ? (row) => {
                    setFormError(null);
                    setEditing(row);
                  }
                : undefined
            }
            onDelete={isAuthenticated ? (row) => setDeleting(row) : undefined}
          />
          <Pagination
            page={scenes.page}
            pageSize={scenes.pageSize}
            total={scenes.total}
            onPage={scenes.setPage}
          />
        </>
      )}

      <Modal
        open={editing !== null}
        onClose={() => {
          setEditing(null);
          setFormError(null);
        }}
        title={editing === "new" ? "Nueva escena" : "Editar escena"}
      >
        <SceneForm
          initial={initial}
          characterOptions={characters.data}
          locationOptions={locations.data}
          onSubmit={handleSubmit}
          onCancel={() => {
            setEditing(null);
            setFormError(null);
          }}
          submitting={submitting}
          error={formError}
        />
      </Modal>

      <ConfirmDialog
        open={deleting !== null}
        title="Confirmar borrado"
        message={deleting ? `¿Seguro que quieres borrar "${deleting.title}"?` : ""}
        confirmLabel="Borrar"
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(null)}
      />
    </section>
  );
}
```

- [ ] **Step 2: Verificar tipos y pruebas**

Run dentro de `web/`: `npx tsc --noEmit && npx vitest run`
Expected: sin errores; todas las pruebas pasan.

- [ ] **Step 3: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): búsqueda y paginación en ScenesPage"
```

---

### Task 9: Verificación final + smoke test

**Files:** ninguno

- [ ] **Step 1: Verificación automática completa**

Run desde la raíz:
```bash
go build ./... && go test ./... && echo "GO OK"
cd web && npm test && npx tsc --noEmit && npm run build && echo "WEB OK"
```
Expected: `GO OK` (model, repository, handler, auth, service, storage) y `WEB OK`.

- [ ] **Step 2: Smoke test extremo a extremo (requiere MySQL local)**

```bash
cd /Users/mqr93ea/Repos/story-go-mysql
go run ./cmd/server &
sleep 2
echo "lista paginada:"
curl -s -w "\n" "http://localhost:8080/api/characters?page=1&pageSize=2" | head -c 400
echo "búsqueda:"
curl -s -w "\n" "http://localhost:8080/api/characters?q=asha" | head -c 400
# (detén el servidor)
```
Expected: respuesta con forma `{ "items": [...], "total": N, "page": 1, "pageSize": 2 }`; la búsqueda por `asha` devuelve solo coincidencias.

- [ ] **Step 3: Commit final**

```bash
git add -A && git commit -m "test: verificación final de búsqueda y paginación" --allow-empty
```

---

## Notas de verificación (self-review del plan)

- **Cobertura del spec:** Page/ListParams + buildSearch (Task 1), repos+servicios
  characters/locations (Task 2), scenes (Task 3), handlers + parseListParams
  (Task 4), tipo Paged + api (Task 5), hook + SearchBar + Pagination (Task 6),
  ResourcePage (Task 7), ScenesPage + opciones con listAll (Task 8),
  verificación + smoke (Task 9). Pruebas en Tasks 1, 4, 6.
- **Sin placeholders:** todo el código está completo.
- **Consistencia de tipos/firmas:** `model.ListParams{Query,Page,PageSize}` +
  `Normalize/Limit/Offset`; `model.Page[T]{Items,Total,Page,PageSize}`;
  repos `List(ctx,q,limit,offset)` + `Count(ctx,q)`; scenes `ListIDs(ctx,q,
  limit,offset)`; servicios `List(ctx, model.ListParams) (model.Page[T], error)`;
  `parseListParams(r) model.ListParams`; frontend `Paged<T>`, `ListArgs`,
  `usePagedList`, `SearchBar{onQueryChange}`, `Pagination{page,pageSize,total,
  onPage}`. Mismos nombres en todas las tareas.
- **Compatibilidad:** no hay pruebas previas que dependan del antiguo contrato de
  listas (los tests de DataTable usan filas literales; los del cliente mockean
  fetch). El formulario de escenas mantiene todas sus opciones vía `listAll*`.
```
