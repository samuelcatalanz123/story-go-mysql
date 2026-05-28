# Panel de administración React + TypeScript — Plan de implementación

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Construir un panel de administración web (React + TypeScript) que permita crear, listar, editar y borrar personajes, lugares y escenas consumiendo la Story API de Go.

**Architecture:** App de una sola página con React Router para navegar entre tres secciones. Los datos se piden con `fetch` a través de una pequeña capa `api/` y un hook `useList`. Un proxy de Vite reenvía `/api/*` al backend Go en `:8080` para evitar CORS. Personajes y Lugares comparten un componente genérico `ResourcePage`; Escenas tiene su propio formulario por sus campos extra y relaciones.

**Tech Stack:** Vite, React 18, TypeScript, React Router, CSS Modules, Vitest + React Testing Library.

---

## Estructura de archivos

```
web/
  index.html
  package.json            (generado por Vite, luego editado)
  tsconfig.json           (generado por Vite)
  vite.config.ts          Proxy /api → :8080 + config de Vitest
  src/
    main.tsx              Punto de entrada; monta React + Router
    App.tsx               Layout (menú) + rutas
    App.module.css
    setupTests.ts         Configuración de las pruebas
    types.ts              Tipos TS que reflejan los modelos de Go
    api/
      client.ts           Envoltura de fetch (URL, JSON, errores)
      client.test.ts
      resources.ts        Funciones CRUD por recurso
    hooks/
      useList.ts          Hook de listado (cargando/error/datos)
    components/
      Layout.tsx
      Layout.module.css
      DataTable.tsx       Tabla reutilizable
      DataTable.module.css
      DataTable.test.tsx
      ResourceForm.tsx    Formulario título + texto
      ResourceForm.test.tsx
      ResourcePage.tsx    Página genérica (lista + alta/edición/borrado)
      SceneForm.tsx       Formulario de escena (campos extra + relaciones)
    pages/
      CharactersPage.tsx
      LocationsPage.tsx
      ScenesPage.tsx
```

Nota sobre commits: el repo no tiene git todavía. La Tarea 1 lo inicializa. Si
prefieres no usar git, puedes omitir los pasos "Commit" de cada tarea.

---

### Task 1: Andamiaje del proyecto y proxy

**Files:**
- Create: `web/` (proyecto Vite)
- Modify: `web/vite.config.ts`
- Create: `web/src/setupTests.ts`
- Modify: `web/package.json`

- [ ] **Step 1: Inicializar git (si no existe)**

Run desde la raíz del repo:
```bash
git init
git add -A
git commit -m "chore: snapshot inicial antes del frontend"
```
Expected: un commit creado (o "nothing to commit" si ya estaba inicializado).

- [ ] **Step 2: Crear el proyecto Vite**

Run desde la raíz del repo:
```bash
npm create vite@latest web -- --template react-ts
cd web && npm install
```
Expected: se crea la carpeta `web/` con la plantilla React + TypeScript y se instalan dependencias.

- [ ] **Step 3: Instalar las dependencias del proyecto**

Run dentro de `web/`:
```bash
npm install react-router-dom
npm install -D vitest jsdom @testing-library/react @testing-library/jest-dom @testing-library/user-event
```
Expected: instalación correcta sin errores.

- [ ] **Step 4: Configurar el proxy y Vitest en `web/vite.config.ts`**

Reemplaza TODO el contenido de `web/vite.config.ts` por:
```ts
/// <reference types="vitest/config" />
import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: "./src/setupTests.ts",
  },
});
```

- [ ] **Step 5: Crear `web/src/setupTests.ts`**

```ts
import "@testing-library/jest-dom";
```

- [ ] **Step 6: Añadir el script de test en `web/package.json`**

Dentro del objeto `"scripts"` añade la línea del test (deja las demás como están):
```json
"scripts": {
  "dev": "vite",
  "build": "tsc -b && vite build",
  "preview": "vite preview",
  "test": "vitest run"
}
```

- [ ] **Step 7: Verificar que el servidor de desarrollo arranca**

Run dentro de `web/`:
```bash
npm run dev
```
Expected: Vite imprime una URL local (p. ej. `http://localhost:5173`). Ábrela; debe verse la página por defecto de Vite. Detén el servidor con `Ctrl+C`.

- [ ] **Step 8: Commit**

```bash
cd .. && git add -A && git commit -m "chore: scaffold web (Vite + React + TS) con proxy y Vitest"
```

---

### Task 2: Tipos de TypeScript

**Files:**
- Create: `web/src/types.ts`

- [ ] **Step 1: Crear `web/src/types.ts`**

```ts
// Tipos que reflejan los modelos JSON de la API de Go.
// En Go los campos *string (anulables) se representan como string | null.

export type Character = {
  id: number;
  title: string;
  text: string | null;
  createdAt: string;
  updatedAt: string;
};

export type Location = {
  id: number;
  title: string;
  text: string | null;
  createdAt: string;
  updatedAt: string;
};

export type Scene = {
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

export type CharacterRequest = { title: string; text: string | null };
export type LocationRequest = { title: string; text: string | null };
export type SceneRequest = {
  title: string;
  text: string | null;
  startTimeline: number;
  endTimeline: number;
  characterIds: number[];
  locationIds: number[];
};
```

- [ ] **Step 2: Verificar que compila**

Run dentro de `web/`:
```bash
npx tsc --noEmit
```
Expected: sin errores.

- [ ] **Step 3: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): tipos del dominio"
```

---

### Task 3: Cliente fetch (`client.ts`) — TDD

**Files:**
- Test: `web/src/api/client.test.ts`
- Create: `web/src/api/client.ts`

- [ ] **Step 1: Escribir la prueba que falla**

Crea `web/src/api/client.test.ts`:
```ts
import { describe, it, expect, vi, afterEach } from "vitest";
import { apiFetch } from "./client";

afterEach(() => vi.unstubAllGlobals());

describe("apiFetch", () => {
  it("devuelve el JSON cuando la respuesta es exitosa", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ id: 1, title: "Asha" }),
      }),
    );
    const data = await apiFetch<{ id: number; title: string }>("/characters/1");
    expect(data).toEqual({ id: 1, title: "Asha" });
  });

  it("lanza el mensaje de error de la API cuando falla", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        status: 409,
        json: async () => ({ error: "title already exists" }),
      }),
    );
    await expect(apiFetch("/characters")).rejects.toThrow("title already exists");
  });
});
```

- [ ] **Step 2: Ejecutar la prueba para ver que falla**

Run dentro de `web/`:
```bash
npx vitest run src/api/client.test.ts
```
Expected: FALLA con un error de importación ("does not provide an export named 'apiFetch'") o similar.

- [ ] **Step 3: Implementación mínima en `web/src/api/client.ts`**

```ts
const BASE = "/api";

// Hace una petición a la API y devuelve el JSON ya parseado.
// Si la respuesta no es exitosa, lanza un Error con el mensaje de la API.
export async function apiFetch<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  if (!res.ok) {
    let message = `Error ${res.status}`;
    try {
      const body = await res.json();
      if (body && typeof body.error === "string") message = body.error;
    } catch {
      // La respuesta de error no traía JSON; usamos el mensaje por defecto.
    }
    throw new Error(message);
  }

  // 204 No Content (p. ej. en DELETE): no hay cuerpo que parsear.
  if (res.status === 204) return undefined as T;
  return (await res.json()) as T;
}
```

- [ ] **Step 4: Ejecutar la prueba para ver que pasa**

Run dentro de `web/`:
```bash
npx vitest run src/api/client.test.ts
```
Expected: PASA (2 tests).

- [ ] **Step 5: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): cliente fetch con manejo de errores"
```

---

### Task 4: Funciones CRUD por recurso (`resources.ts`)

**Files:**
- Create: `web/src/api/resources.ts`

- [ ] **Step 1: Crear `web/src/api/resources.ts`**

```ts
import { apiFetch } from "./client";
import type {
  Character,
  Location,
  Scene,
  CharacterRequest,
  LocationRequest,
  SceneRequest,
} from "../types";

// --- Personajes ---
export const listCharacters = () => apiFetch<Character[]>("/characters");
export const createCharacter = (body: CharacterRequest) =>
  apiFetch<Character>("/characters", { method: "POST", body: JSON.stringify(body) });
export const updateCharacter = (id: number, body: CharacterRequest) =>
  apiFetch<Character>(`/characters/${id}`, { method: "PUT", body: JSON.stringify(body) });
export const deleteCharacter = (id: number) =>
  apiFetch<void>(`/characters/${id}`, { method: "DELETE" });

// --- Lugares ---
export const listLocations = () => apiFetch<Location[]>("/locations");
export const createLocation = (body: LocationRequest) =>
  apiFetch<Location>("/locations", { method: "POST", body: JSON.stringify(body) });
export const updateLocation = (id: number, body: LocationRequest) =>
  apiFetch<Location>(`/locations/${id}`, { method: "PUT", body: JSON.stringify(body) });
export const deleteLocation = (id: number) =>
  apiFetch<void>(`/locations/${id}`, { method: "DELETE" });

// --- Escenas ---
export const listScenes = () => apiFetch<Scene[]>("/scenes");
export const createScene = (body: SceneRequest) =>
  apiFetch<Scene>("/scenes", { method: "POST", body: JSON.stringify(body) });
export const updateScene = (id: number, body: SceneRequest) =>
  apiFetch<Scene>(`/scenes/${id}`, { method: "PUT", body: JSON.stringify(body) });
export const deleteScene = (id: number) =>
  apiFetch<void>(`/scenes/${id}`, { method: "DELETE" });
```

- [ ] **Step 2: Verificar que compila**

Run dentro de `web/`:
```bash
npx tsc --noEmit
```
Expected: sin errores.

- [ ] **Step 3: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): funciones CRUD por recurso"
```

---

### Task 5: Hook `useList`

**Files:**
- Create: `web/src/hooks/useList.ts`

- [ ] **Step 1: Crear `web/src/hooks/useList.ts`**

```ts
import { useCallback, useEffect, useState } from "react";

export type ListState<T> = {
  data: T[];
  loading: boolean;
  error: string | null;
  reload: () => void;
};

// Carga una lista al montar el componente y expone estado de
// carga/error/datos, además de una función para recargar.
// `loader` debe ser una referencia estable (p. ej. una función de api/resources.ts).
export function useList<T>(loader: () => Promise<T[]>): ListState<T> {
  const [data, setData] = useState<T[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    setError(null);
    loader()
      .then((items) => setData(items))
      .catch((e: unknown) =>
        setError(e instanceof Error ? e.message : "Error desconocido"),
      )
      .finally(() => setLoading(false));
  }, [loader]);

  useEffect(() => {
    load();
  }, [load]);

  return { data, loading, error, reload: load };
}
```

- [ ] **Step 2: Verificar que compila**

Run dentro de `web/`:
```bash
npx tsc --noEmit
```
Expected: sin errores.

- [ ] **Step 3: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): hook useList"
```

---

### Task 6: DataTable — TDD

**Files:**
- Test: `web/src/components/DataTable.test.tsx`
- Create: `web/src/components/DataTable.tsx`
- Create: `web/src/components/DataTable.module.css`

- [ ] **Step 1: Escribir la prueba que falla**

Crea `web/src/components/DataTable.test.tsx`:
```tsx
import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { DataTable } from "./DataTable";

describe("DataTable", () => {
  it("renderiza una fila por cada elemento", () => {
    render(
      <DataTable
        columns={[{ header: "Título", render: (r) => r.title }]}
        rows={[
          { id: 1, title: "Asha" },
          { id: 2, title: "Bo" },
        ]}
      />,
    );
    expect(screen.getByText("Asha")).toBeInTheDocument();
    expect(screen.getByText("Bo")).toBeInTheDocument();
  });

  it("muestra un mensaje cuando no hay elementos", () => {
    render(<DataTable columns={[]} rows={[]} />);
    expect(screen.getByText("No hay elementos todavía.")).toBeInTheDocument();
  });
});
```

- [ ] **Step 2: Ejecutar la prueba para ver que falla**

Run dentro de `web/`:
```bash
npx vitest run src/components/DataTable.test.tsx
```
Expected: FALLA (no existe `DataTable`).

- [ ] **Step 3: Crear `web/src/components/DataTable.module.css`**

```css
.table {
  width: 100%;
  border-collapse: collapse;
}
.table th,
.table td {
  border: 1px solid #ddd;
  padding: 0.5rem;
  text-align: left;
}
.table th {
  background: #f5f5f5;
}
.actions button {
  margin-right: 0.5rem;
}
```

- [ ] **Step 4: Crear `web/src/components/DataTable.tsx`**

```tsx
import type { ReactNode } from "react";
import styles from "./DataTable.module.css";

export type Column<T> = {
  header: string;
  render: (row: T) => ReactNode;
};

type Props<T> = {
  columns: Column<T>[];
  rows: T[];
  onEdit?: (row: T) => void;
  onDelete?: (row: T) => void;
};

export function DataTable<T extends { id: number }>({
  columns,
  rows,
  onEdit,
  onDelete,
}: Props<T>) {
  if (rows.length === 0) {
    return <p>No hay elementos todavía.</p>;
  }

  const hasActions = Boolean(onEdit || onDelete);

  return (
    <table className={styles.table}>
      <thead>
        <tr>
          {columns.map((c) => (
            <th key={c.header}>{c.header}</th>
          ))}
          {hasActions && <th>Acciones</th>}
        </tr>
      </thead>
      <tbody>
        {rows.map((row) => (
          <tr key={row.id}>
            {columns.map((c) => (
              <td key={c.header}>{c.render(row)}</td>
            ))}
            {hasActions && (
              <td className={styles.actions}>
                {onEdit && <button onClick={() => onEdit(row)}>Editar</button>}
                {onDelete && <button onClick={() => onDelete(row)}>Borrar</button>}
              </td>
            )}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
```

- [ ] **Step 5: Ejecutar la prueba para ver que pasa**

Run dentro de `web/`:
```bash
npx vitest run src/components/DataTable.test.tsx
```
Expected: PASA (2 tests).

- [ ] **Step 6: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): componente DataTable"
```

---

### Task 7: ResourceForm — TDD

**Files:**
- Test: `web/src/components/ResourceForm.test.tsx`
- Create: `web/src/components/ResourceForm.tsx`

- [ ] **Step 1: Escribir la prueba que falla**

Crea `web/src/components/ResourceForm.test.tsx`:
```tsx
import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { ResourceForm } from "./ResourceForm";

describe("ResourceForm", () => {
  it("muestra error y no envía si el título está vacío", () => {
    const onSubmit = vi.fn();
    render(<ResourceForm onSubmit={onSubmit} onCancel={() => {}} />);
    fireEvent.click(screen.getByText("Guardar"));
    expect(screen.getByRole("alert")).toHaveTextContent("El título es obligatorio");
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it("envía título y texto cuando son válidos", () => {
    const onSubmit = vi.fn();
    render(<ResourceForm onSubmit={onSubmit} onCancel={() => {}} />);
    fireEvent.change(screen.getByLabelText("Título"), {
      target: { value: "Asha" },
    });
    fireEvent.change(screen.getByLabelText("Texto"), {
      target: { value: "Una piloto" },
    });
    fireEvent.click(screen.getByText("Guardar"));
    expect(onSubmit).toHaveBeenCalledWith({ title: "Asha", text: "Una piloto" });
  });
});
```

- [ ] **Step 2: Ejecutar la prueba para ver que falla**

Run dentro de `web/`:
```bash
npx vitest run src/components/ResourceForm.test.tsx
```
Expected: FALLA (no existe `ResourceForm`).

- [ ] **Step 3: Crear `web/src/components/ResourceForm.tsx`**

```tsx
import { useState } from "react";
import type { FormEvent } from "react";

export type ResourceFormValues = { title: string; text: string };

type Props = {
  initial?: ResourceFormValues;
  onSubmit: (values: ResourceFormValues) => void;
  onCancel: () => void;
  submitting?: boolean;
  error?: string | null;
};

export function ResourceForm({ initial, onSubmit, onCancel, submitting, error }: Props) {
  const [title, setTitle] = useState(initial?.title ?? "");
  const [text, setText] = useState(initial?.text ?? "");
  const [localError, setLocalError] = useState<string | null>(null);

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (title.trim() === "") {
      setLocalError("El título es obligatorio");
      return;
    }
    setLocalError(null);
    onSubmit({ title: title.trim(), text });
  }

  const shownError = localError ?? error;

  return (
    <form onSubmit={handleSubmit}>
      {shownError && <p role="alert">{shownError}</p>}
      <div>
        <label htmlFor="rf-title">Título</label>
        <input
          id="rf-title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
      </div>
      <div>
        <label htmlFor="rf-text">Texto</label>
        <textarea
          id="rf-text"
          value={text}
          onChange={(e) => setText(e.target.value)}
        />
      </div>
      <button type="submit" disabled={submitting}>
        Guardar
      </button>
      <button type="button" onClick={onCancel}>
        Cancelar
      </button>
    </form>
  );
}
```

- [ ] **Step 4: Ejecutar la prueba para ver que pasa**

Run dentro de `web/`:
```bash
npx vitest run src/components/ResourceForm.test.tsx
```
Expected: PASA (2 tests).

- [ ] **Step 5: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): formulario ResourceForm con validación"
```

---

### Task 8: ResourcePage genérica (lista + alta/edición/borrado)

**Files:**
- Create: `web/src/components/ResourcePage.tsx`

- [ ] **Step 1: Crear `web/src/components/ResourcePage.tsx`**

```tsx
import { useState } from "react";
import { useList } from "../hooks/useList";
import { DataTable } from "./DataTable";
import type { Column } from "./DataTable";
import { ResourceForm } from "./ResourceForm";
import type { ResourceFormValues } from "./ResourceForm";

// Forma mínima que deben cumplir los elementos de esta página.
type ResourceItem = {
  id: number;
  title: string;
  text: string | null;
  updatedAt: string;
};

type RequestBody = { title: string; text: string | null };

type Props<T extends ResourceItem> = {
  heading: string;
  list: () => Promise<T[]>;
  create: (body: RequestBody) => Promise<T>;
  update: (id: number, body: RequestBody) => Promise<T>;
  remove: (id: number) => Promise<void>;
};

// `null` = sólo lista; "new" = creando; un objeto = editando ese elemento.
type Editing<T> = null | "new" | T;

export function ResourcePage<T extends ResourceItem>({
  heading,
  list,
  create,
  update,
  remove,
}: Props<T>) {
  const { data, loading, error, reload } = useList(list);
  const [editing, setEditing] = useState<Editing<T>>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const columns: Column<T>[] = [
    { header: "ID", render: (r) => r.id },
    { header: "Título", render: (r) => r.title },
    { header: "Texto", render: (r) => r.text ?? "—" },
    {
      header: "Actualizado",
      render: (r) => new Date(r.updatedAt).toLocaleString(),
    },
  ];

  function toBody(values: ResourceFormValues): RequestBody {
    return { title: values.title, text: values.text.trim() === "" ? null : values.text };
  }

  async function handleSubmit(values: ResourceFormValues) {
    setSubmitting(true);
    setFormError(null);
    try {
      if (editing === "new") {
        await create(toBody(values));
      } else if (editing) {
        await update(editing.id, toBody(values));
      }
      setEditing(null);
      reload();
    } catch (e: unknown) {
      setFormError(e instanceof Error ? e.message : "Error desconocido");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(row: T) {
    if (!window.confirm(`¿Borrar "${row.title}"?`)) return;
    try {
      await remove(row.id);
      reload();
    } catch (e: unknown) {
      alert(e instanceof Error ? e.message : "Error desconocido");
    }
  }

  if (editing !== null) {
    const initial =
      editing === "new"
        ? undefined
        : { title: editing.title, text: editing.text ?? "" };
    return (
      <section>
        <h1>{editing === "new" ? `Nuevo: ${heading}` : `Editar: ${heading}`}</h1>
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
      </section>
    );
  }

  return (
    <section>
      <h1>{heading}</h1>
      <button onClick={() => setEditing("new")}>Nuevo</button>
      {loading && <p>Cargando…</p>}
      {error && <p role="alert">Error: {error}</p>}
      {!loading && !error && (
        <DataTable
          columns={columns}
          rows={data}
          onEdit={(row) => setEditing(row)}
          onDelete={handleDelete}
        />
      )}
    </section>
  );
}
```

- [ ] **Step 2: Verificar que compila**

Run dentro de `web/`:
```bash
npx tsc --noEmit
```
Expected: sin errores.

- [ ] **Step 3: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): ResourcePage genérica"
```

---

### Task 9: Páginas de Personajes y Lugares

**Files:**
- Create: `web/src/pages/CharactersPage.tsx`
- Create: `web/src/pages/LocationsPage.tsx`

- [ ] **Step 1: Crear `web/src/pages/CharactersPage.tsx`**

```tsx
import { ResourcePage } from "../components/ResourcePage";
import {
  listCharacters,
  createCharacter,
  updateCharacter,
  deleteCharacter,
} from "../api/resources";

export function CharactersPage() {
  return (
    <ResourcePage
      heading="Personajes"
      list={listCharacters}
      create={createCharacter}
      update={updateCharacter}
      remove={deleteCharacter}
    />
  );
}
```

- [ ] **Step 2: Crear `web/src/pages/LocationsPage.tsx`**

```tsx
import { ResourcePage } from "../components/ResourcePage";
import {
  listLocations,
  createLocation,
  updateLocation,
  deleteLocation,
} from "../api/resources";

export function LocationsPage() {
  return (
    <ResourcePage
      heading="Lugares"
      list={listLocations}
      create={createLocation}
      update={updateLocation}
      remove={deleteLocation}
    />
  );
}
```

- [ ] **Step 3: Verificar que compila**

Run dentro de `web/`:
```bash
npx tsc --noEmit
```
Expected: sin errores.

- [ ] **Step 4: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): páginas de personajes y lugares"
```

---

### Task 10: Formulario de escenas (`SceneForm`)

**Files:**
- Create: `web/src/components/SceneForm.tsx`

- [ ] **Step 1: Crear `web/src/components/SceneForm.tsx`**

```tsx
import { useState } from "react";
import type { FormEvent } from "react";
import type { Character, Location } from "../types";

export type SceneFormValues = {
  title: string;
  text: string;
  startTimeline: number;
  endTimeline: number;
  characterIds: number[];
  locationIds: number[];
};

type Props = {
  initial?: SceneFormValues;
  characterOptions: Character[];
  locationOptions: Location[];
  onSubmit: (values: SceneFormValues) => void;
  onCancel: () => void;
  submitting?: boolean;
  error?: string | null;
};

// Lee los valores seleccionados de un <select multiple> como números.
function selectedNumbers(select: HTMLSelectElement): number[] {
  return Array.from(select.selectedOptions).map((o) => Number(o.value));
}

export function SceneForm({
  initial,
  characterOptions,
  locationOptions,
  onSubmit,
  onCancel,
  submitting,
  error,
}: Props) {
  const [title, setTitle] = useState(initial?.title ?? "");
  const [text, setText] = useState(initial?.text ?? "");
  const [start, setStart] = useState(String(initial?.startTimeline ?? 0));
  const [end, setEnd] = useState(String(initial?.endTimeline ?? 0));
  const [characterIds, setCharacterIds] = useState<number[]>(
    initial?.characterIds ?? [],
  );
  const [locationIds, setLocationIds] = useState<number[]>(
    initial?.locationIds ?? [],
  );
  const [localError, setLocalError] = useState<string | null>(null);

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (title.trim() === "") {
      setLocalError("El título es obligatorio");
      return;
    }
    setLocalError(null);
    onSubmit({
      title: title.trim(),
      text,
      startTimeline: Number(start),
      endTimeline: Number(end),
      characterIds,
      locationIds,
    });
  }

  const shownError = localError ?? error;

  return (
    <form onSubmit={handleSubmit}>
      {shownError && <p role="alert">{shownError}</p>}
      <div>
        <label htmlFor="sf-title">Título</label>
        <input id="sf-title" value={title} onChange={(e) => setTitle(e.target.value)} />
      </div>
      <div>
        <label htmlFor="sf-text">Texto</label>
        <textarea id="sf-text" value={text} onChange={(e) => setText(e.target.value)} />
      </div>
      <div>
        <label htmlFor="sf-start">Inicio (timeline)</label>
        <input
          id="sf-start"
          type="number"
          value={start}
          onChange={(e) => setStart(e.target.value)}
        />
      </div>
      <div>
        <label htmlFor="sf-end">Fin (timeline)</label>
        <input
          id="sf-end"
          type="number"
          value={end}
          onChange={(e) => setEnd(e.target.value)}
        />
      </div>
      <div>
        <label htmlFor="sf-characters">Personajes</label>
        <select
          id="sf-characters"
          multiple
          value={characterIds.map(String)}
          onChange={(e) => setCharacterIds(selectedNumbers(e.target))}
        >
          {characterOptions.map((c) => (
            <option key={c.id} value={c.id}>
              {c.title}
            </option>
          ))}
        </select>
      </div>
      <div>
        <label htmlFor="sf-locations">Lugares</label>
        <select
          id="sf-locations"
          multiple
          value={locationIds.map(String)}
          onChange={(e) => setLocationIds(selectedNumbers(e.target))}
        >
          {locationOptions.map((l) => (
            <option key={l.id} value={l.id}>
              {l.title}
            </option>
          ))}
        </select>
      </div>
      <button type="submit" disabled={submitting}>
        Guardar
      </button>
      <button type="button" onClick={onCancel}>
        Cancelar
      </button>
    </form>
  );
}
```

- [ ] **Step 2: Verificar que compila**

Run dentro de `web/`:
```bash
npx tsc --noEmit
```
Expected: sin errores.

- [ ] **Step 3: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): formulario de escenas con relaciones"
```

---

### Task 11: Página de Escenas

**Files:**
- Create: `web/src/pages/ScenesPage.tsx`

- [ ] **Step 1: Crear `web/src/pages/ScenesPage.tsx`**

```tsx
import { useState } from "react";
import { useList } from "../hooks/useList";
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
  listCharacters,
  listLocations,
} from "../api/resources";

type Editing = null | "new" | Scene;

export function ScenesPage() {
  const scenes = useList(listScenes);
  const characters = useList(listCharacters);
  const locations = useList(listLocations);

  const [editing, setEditing] = useState<Editing>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const columns: Column<Scene>[] = [
    { header: "ID", render: (s) => s.id },
    { header: "Título", render: (s) => s.title },
    { header: "Inicio", render: (s) => s.startTimeline },
    { header: "Fin", render: (s) => s.endTimeline },
    {
      header: "Personajes",
      render: (s) => s.characters.map((c) => c.title).join(", ") || "—",
    },
    {
      header: "Lugares",
      render: (s) => s.locations.map((l) => l.title).join(", ") || "—",
    },
  ];

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
      } else if (editing) {
        await updateScene(editing.id, body);
      }
      setEditing(null);
      scenes.reload();
    } catch (e: unknown) {
      setFormError(e instanceof Error ? e.message : "Error desconocido");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(scene: Scene) {
    if (!window.confirm(`¿Borrar "${scene.title}"?`)) return;
    try {
      await deleteScene(scene.id);
      scenes.reload();
    } catch (e: unknown) {
      alert(e instanceof Error ? e.message : "Error desconocido");
    }
  }

  if (editing !== null) {
    const initial: SceneFormValues | undefined =
      editing === "new"
        ? undefined
        : {
            title: editing.title,
            text: editing.text ?? "",
            startTimeline: editing.startTimeline,
            endTimeline: editing.endTimeline,
            characterIds: editing.characters.map((c) => c.id),
            locationIds: editing.locations.map((l) => l.id),
          };
    return (
      <section>
        <h1>{editing === "new" ? "Nueva escena" : "Editar escena"}</h1>
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
      </section>
    );
  }

  return (
    <section>
      <h1>Escenas</h1>
      <button onClick={() => setEditing("new")}>Nueva</button>
      {scenes.loading && <p>Cargando…</p>}
      {scenes.error && <p role="alert">Error: {scenes.error}</p>}
      {!scenes.loading && !scenes.error && (
        <DataTable
          columns={columns}
          rows={scenes.data}
          onEdit={(row) => setEditing(row)}
          onDelete={handleDelete}
        />
      )}
    </section>
  );
}
```

- [ ] **Step 2: Verificar que compila**

Run dentro de `web/`:
```bash
npx tsc --noEmit
```
Expected: sin errores.

- [ ] **Step 3: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): página de escenas"
```

---

### Task 12: Layout, rutas y arranque

**Files:**
- Create: `web/src/components/Layout.tsx`
- Create: `web/src/components/Layout.module.css`
- Modify: `web/src/App.tsx`
- Create: `web/src/App.module.css`
- Modify: `web/src/main.tsx`

- [ ] **Step 1: Crear `web/src/components/Layout.module.css`**

```css
.shell {
  display: flex;
  min-height: 100vh;
  font-family: system-ui, sans-serif;
}
.sidebar {
  width: 200px;
  background: #1f2933;
  padding: 1rem;
}
.sidebar h2 {
  color: #fff;
  font-size: 1rem;
}
.sidebar a {
  display: block;
  color: #cbd2d9;
  text-decoration: none;
  padding: 0.5rem 0;
}
.sidebar a.active {
  color: #fff;
  font-weight: bold;
}
.content {
  flex: 1;
  padding: 1.5rem;
}
```

- [ ] **Step 2: Crear `web/src/components/Layout.tsx`**

```tsx
import { NavLink, Outlet } from "react-router-dom";
import styles from "./Layout.module.css";

export function Layout() {
  const linkClass = ({ isActive }: { isActive: boolean }) =>
    isActive ? styles.active : "";
  return (
    <div className={styles.shell}>
      <nav className={styles.sidebar}>
        <h2>Story Admin</h2>
        <NavLink to="/characters" className={linkClass}>
          Personajes
        </NavLink>
        <NavLink to="/locations" className={linkClass}>
          Lugares
        </NavLink>
        <NavLink to="/scenes" className={linkClass}>
          Escenas
        </NavLink>
      </nav>
      <main className={styles.content}>
        <Outlet />
      </main>
    </div>
  );
}
```

- [ ] **Step 3: Crear `web/src/App.module.css`**

```css
/* Estilos globales mínimos para los formularios y botones. */
form div {
  margin-bottom: 0.75rem;
}
form label {
  display: block;
  font-weight: 600;
  margin-bottom: 0.25rem;
}
form input,
form textarea,
form select {
  width: 100%;
  max-width: 400px;
  padding: 0.4rem;
}
button {
  margin-right: 0.5rem;
  cursor: pointer;
}
[role="alert"] {
  color: #b91c1c;
}
```

- [ ] **Step 4: Reemplazar `web/src/App.tsx`**

Reemplaza TODO el contenido por:
```tsx
import { Routes, Route, Navigate } from "react-router-dom";
import { Layout } from "./components/Layout";
import { CharactersPage } from "./pages/CharactersPage";
import { LocationsPage } from "./pages/LocationsPage";
import { ScenesPage } from "./pages/ScenesPage";
import "./App.module.css";

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Navigate to="/characters" replace />} />
        <Route path="characters" element={<CharactersPage />} />
        <Route path="locations" element={<LocationsPage />} />
        <Route path="scenes" element={<ScenesPage />} />
      </Route>
    </Routes>
  );
}
```

- [ ] **Step 5: Reemplazar `web/src/main.tsx`**

Reemplaza TODO el contenido por:
```tsx
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </StrictMode>,
);
```

- [ ] **Step 6: Verificar que compila**

Run dentro de `web/`:
```bash
npx tsc --noEmit
```
Expected: sin errores. (Si `App.css` o `index.css` por defecto causan problemas visuales, puedes vaciarlos; no es obligatorio.)

- [ ] **Step 7: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): layout, rutas y arranque de la app"
```

---

### Task 13: Verificación final (pruebas + extremo a extremo)

**Files:** ninguno (solo verificación)

- [ ] **Step 1: Ejecutar todas las pruebas**

Run dentro de `web/`:
```bash
npm test
```
Expected: PASAN todas (client, DataTable, ResourceForm).

- [ ] **Step 2: Comprobación de tipos completa**

Run dentro de `web/`:
```bash
npx tsc --noEmit
```
Expected: sin errores.

- [ ] **Step 3: Prueba manual extremo a extremo**

1. En una terminal, arranca la API: `go run ./cmd/server` (desde la raíz del repo; requiere MySQL en marcha con la base `story_go_db` creada según el README).
2. En otra terminal: `cd web && npm run dev`.
3. Abre la URL local. Verifica:
   - La sección **Personajes** lista, crea, edita y borra.
   - Lo mismo en **Lugares**.
   - En **Escenas**, al crear una escena puedes seleccionar personajes y lugares ya creados, y la tabla los muestra.
   - Crear un título duplicado muestra el mensaje de error (`409`) en el formulario.

Expected: todo el ciclo CRUD funciona contra la API real.

- [ ] **Step 4: Commit final**

```bash
cd .. && git add -A && git commit -m "test(web): verificación final del panel de administración"
```

---

## Notas de verificación (self-review del plan)

- **Cobertura del spec:** proxy (Task 1), tipos (Task 2), cliente/errores (Task 3),
  CRUD por recurso (Task 4), hook (Task 5), tabla (Task 6), formulario+validación
  (Task 7), páginas personajes/lugares (Tasks 8-9), escenas con relaciones
  (Tasks 10-11), navegación/rutas (Task 12), pruebas de ejemplo (Tasks 3,6,7,13).
- **Sin placeholders:** todo el código está completo en cada paso.
- **Consistencia de tipos:** `ResourceFormValues`, `SceneFormValues`, `Column<T>`,
  `apiFetch`, y las funciones de `resources.ts` se usan con los mismos nombres y
  firmas en todas las tareas.
