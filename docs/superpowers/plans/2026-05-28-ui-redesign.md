# Rediseño profesional de la UI — Plan de implementación

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Convertir el panel de administración en una pieza de portafolio con aspecto "SaaS minimalista" e interacciones pulidas (modales, toasts, skeletons, estados vacíos), usando CSS Modules + tokens de diseño.

**Architecture:** Se añade un sistema de tokens CSS global y una carpeta `web/src/ui/` con componentes primitivos reutilizables (Button, Modal, ConfirmDialog, Toast, etc.). Las pantallas existentes (`ResourcePage`, `ScenesPage`, `DataTable`, formularios, `Layout`) se refactorizan para usar esos primitivos. La capa de datos (`api/`, `useList`) no cambia.

**Tech Stack:** React 19, TypeScript, react-router-dom 7, CSS Modules, Vitest + React Testing Library.

---

## Estructura de archivos

```
web/src/
  styles/tokens.css                 (nuevo) tokens + reset + estilos base de controles
  ui/                               (nuevo) primitivos
    Button.tsx / .module.css / Button.test.tsx
    Field.tsx / .module.css
    Badge.tsx / .module.css
    Skeleton.tsx / .module.css
    PageHeader.tsx / .module.css
    EmptyState.tsx / .module.css / EmptyState.test.tsx
    Modal.tsx / .module.css / Modal.test.tsx
    ConfirmDialog.tsx / .module.css / ConfirmDialog.test.tsx
    Toast.tsx / .module.css / Toast.test.tsx
  components/
    DataTable.tsx (mod) / DataTable.module.css (mod)
    ResourceForm.tsx (mod) / ResourceForm.module.css (nuevo)
    SceneForm.tsx (mod) / SceneForm.module.css (nuevo)
    ResourcePage.tsx (mod)
    Layout.tsx (mod) / Layout.module.css (mod)
  pages/ScenesPage.tsx (mod)
  App.tsx (mod: quitar import de App.module.css)
  main.tsx (mod: importar tokens.css y envolver en ToastProvider)
  App.module.css (eliminar)
```

Nota: trabajamos en la rama `feat/ui-redesign` (ya creada). Commits ahí.
Comandos desde `web/` salvo los `git commit` (con `cd ..`).

---

### Task 1: Tokens de diseño y estilos base

**Files:**
- Create: `web/src/styles/tokens.css`
- Modify: `web/src/main.tsx`
- Modify: `web/src/App.tsx`
- Delete: `web/src/App.module.css`

- [ ] **Step 1: Crear `web/src/styles/tokens.css`**

```css
:root {
  /* Color */
  --color-bg: #f9fafb;
  --color-surface: #ffffff;
  --color-border: #e5e7eb;
  --color-text: #111827;
  --color-text-muted: #6b7280;
  --color-primary: #4f46e5;
  --color-primary-hover: #4338ca;
  --color-primary-soft: #eef2ff;
  --color-danger: #dc2626;
  --color-danger-hover: #b91c1c;
  --color-success: #16a34a;

  /* Espaciado (base 4px) */
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-5: 24px;
  --space-6: 32px;
  --space-7: 48px;
  --space-8: 64px;

  /* Radios */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;

  /* Sombras */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.08);

  /* Tipografía */
  --font-sans: "Inter", system-ui, -apple-system, "Segoe UI", Roboto, sans-serif;
  --text-sm: 0.8125rem;
  --text-base: 0.9375rem;
  --text-lg: 1.125rem;
  --text-xl: 1.375rem;
  --text-2xl: 1.75rem;
}

/* Reset mínimo */
*,
*::before,
*::after {
  box-sizing: border-box;
}
html,
body,
#root {
  height: 100%;
}
body {
  margin: 0;
  font-family: var(--font-sans);
  font-size: var(--text-base);
  color: var(--color-text);
  background: var(--color-bg);
  -webkit-font-smoothing: antialiased;
}
h1,
h2,
h3,
p {
  margin: 0;
}

/* Controles de formulario (estilo coherente en toda la app) */
input,
textarea,
select {
  font-family: inherit;
  font-size: var(--text-base);
  color: var(--color-text);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--space-2) var(--space-3);
  width: 100%;
}
input:focus,
textarea:focus,
select:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-soft);
}
textarea {
  min-height: 80px;
  resize: vertical;
}
```

- [ ] **Step 2: Reemplazar TODO el contenido de `web/src/main.tsx`**

(En esta tarea solo añadimos el import de tokens; el `ToastProvider` se añade en la Task 7.)
```tsx
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App";
import "./styles/tokens.css";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </StrictMode>,
);
```

- [ ] **Step 3: Reemplazar TODO el contenido de `web/src/App.tsx`** (quita el import de `App.module.css`)

```tsx
import { Routes, Route, Navigate } from "react-router-dom";
import { Layout } from "./components/Layout";
import { CharactersPage } from "./pages/CharactersPage";
import { LocationsPage } from "./pages/LocationsPage";
import { ScenesPage } from "./pages/ScenesPage";

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

- [ ] **Step 4: Eliminar `web/src/App.module.css`**

Run: `rm web/src/App.module.css`

- [ ] **Step 5: Verificar**

Run dentro de `web/`: `npx tsc --noEmit`
Expected: sin errores.

- [ ] **Step 6: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): sistema de tokens de diseño y estilos base"
```

---

### Task 2: Button (TDD)

**Files:**
- Test: `web/src/ui/Button.test.tsx`
- Create: `web/src/ui/Button.tsx`, `web/src/ui/Button.module.css`

- [ ] **Step 1: Escribir la prueba que falla — `web/src/ui/Button.test.tsx`**

```tsx
import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { Button } from "./Button";

describe("Button", () => {
  it("aplica la clase de la variante indicada", () => {
    render(<Button variant="danger">Borrar</Button>);
    expect(screen.getByRole("button", { name: "Borrar" }).className).toMatch(/danger/);
  });

  it("respeta el atributo disabled", () => {
    render(<Button disabled>Guardar</Button>);
    expect(screen.getByRole("button", { name: "Guardar" })).toBeDisabled();
  });
});
```

- [ ] **Step 2: Ejecutar y ver que falla**

Run: `npx vitest run src/ui/Button.test.tsx`
Expected: FALLA (no existe `Button`).

- [ ] **Step 3: Crear `web/src/ui/Button.module.css`**

```css
.button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s, color 0.15s;
}
.button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.sm {
  padding: var(--space-1) var(--space-3);
  font-size: var(--text-sm);
}
.md {
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-base);
}
.primary {
  background: var(--color-primary);
  color: #fff;
}
.primary:hover:not(:disabled) {
  background: var(--color-primary-hover);
}
.secondary {
  background: var(--color-surface);
  color: var(--color-text);
  border-color: var(--color-border);
}
.secondary:hover:not(:disabled) {
  background: var(--color-bg);
}
.danger {
  background: var(--color-danger);
  color: #fff;
}
.danger:hover:not(:disabled) {
  background: var(--color-danger-hover);
}
.ghost {
  background: transparent;
  color: var(--color-text-muted);
}
.ghost:hover:not(:disabled) {
  background: var(--color-bg);
  color: var(--color-text);
}
```

- [ ] **Step 4: Crear `web/src/ui/Button.tsx`**

```tsx
import type { ButtonHTMLAttributes } from "react";
import styles from "./Button.module.css";

type Variant = "primary" | "secondary" | "danger" | "ghost";
type Size = "sm" | "md";

type Props = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: Variant;
  size?: Size;
};

export function Button({ variant = "primary", size = "md", className, ...rest }: Props) {
  const classes = [styles.button, styles[variant], styles[size], className]
    .filter(Boolean)
    .join(" ");
  return <button className={classes} {...rest} />;
}
```

- [ ] **Step 5: Ejecutar y ver que pasa**

Run: `npx vitest run src/ui/Button.test.tsx`
Expected: PASA (2 tests).

- [ ] **Step 6: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): componente Button"
```

---

### Task 3: Primitivos presentacionales (Field, Badge, Skeleton, PageHeader)

**Files:**
- Create: `web/src/ui/Field.tsx`, `Field.module.css`
- Create: `web/src/ui/Badge.tsx`, `Badge.module.css`
- Create: `web/src/ui/Skeleton.tsx`, `Skeleton.module.css`
- Create: `web/src/ui/PageHeader.tsx`, `PageHeader.module.css`

- [ ] **Step 1: Crear `web/src/ui/Field.tsx`**

```tsx
import type { ReactNode } from "react";
import styles from "./Field.module.css";

type Props = {
  label: string;
  htmlFor: string;
  error?: string | null;
  children: ReactNode;
};

export function Field({ label, htmlFor, error, children }: Props) {
  return (
    <div className={styles.field}>
      <label htmlFor={htmlFor} className={styles.label}>
        {label}
      </label>
      {children}
      {error && <span className={styles.error}>{error}</span>}
    </div>
  );
}
```

- [ ] **Step 2: Crear `web/src/ui/Field.module.css`**

```css
.field {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  margin-bottom: var(--space-4);
}
.label {
  font-weight: 600;
  font-size: var(--text-sm);
  color: var(--color-text);
}
.error {
  color: var(--color-danger);
  font-size: var(--text-sm);
}
```

- [ ] **Step 3: Crear `web/src/ui/Badge.tsx`**

```tsx
import type { ReactNode } from "react";
import styles from "./Badge.module.css";

type Props = { children: ReactNode; tone?: "neutral" | "primary" };

export function Badge({ children, tone = "neutral" }: Props) {
  return <span className={`${styles.badge} ${styles[tone]}`}>{children}</span>;
}
```

- [ ] **Step 4: Crear `web/src/ui/Badge.module.css`**

```css
.badge {
  display: inline-block;
  padding: 2px var(--space-2);
  border-radius: var(--radius-sm);
  font-size: var(--text-sm);
  font-weight: 600;
  margin: 2px;
}
.neutral {
  background: var(--color-bg);
  color: var(--color-text-muted);
  border: 1px solid var(--color-border);
}
.primary {
  background: var(--color-primary-soft);
  color: var(--color-primary);
}
```

- [ ] **Step 5: Crear `web/src/ui/Skeleton.tsx`**

```tsx
import styles from "./Skeleton.module.css";

type Props = { width?: string; height?: string };

export function Skeleton({ width = "100%", height = "1rem" }: Props) {
  return <span className={styles.skeleton} style={{ width, height }} />;
}

export function SkeletonRows({ rows = 3, cols = 4 }: { rows?: number; cols?: number }) {
  return (
    <div className={styles.rows}>
      {Array.from({ length: rows }).map((_, r) => (
        <div key={r} className={styles.row}>
          {Array.from({ length: cols }).map((_, c) => (
            <Skeleton key={c} height="1.25rem" />
          ))}
        </div>
      ))}
    </div>
  );
}
```

- [ ] **Step 6: Crear `web/src/ui/Skeleton.module.css`**

```css
.skeleton {
  display: inline-block;
  background: linear-gradient(90deg, #eee 25%, #f5f5f5 37%, #eee 63%);
  background-size: 400% 100%;
  border-radius: var(--radius-sm);
  animation: shimmer 1.4s ease infinite;
}
@keyframes shimmer {
  0% {
    background-position: 100% 0;
  }
  100% {
    background-position: 0 0;
  }
}
.rows {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: var(--space-3) 0;
}
.row {
  display: flex;
  gap: var(--space-4);
}
.row > * {
  flex: 1;
}
```

- [ ] **Step 7: Crear `web/src/ui/PageHeader.tsx`**

```tsx
import type { ReactNode } from "react";
import styles from "./PageHeader.module.css";

type Props = { title: string; action?: ReactNode };

export function PageHeader({ title, action }: Props) {
  return (
    <header className={styles.header}>
      <h1 className={styles.title}>{title}</h1>
      {action && <div>{action}</div>}
    </header>
  );
}
```

- [ ] **Step 8: Crear `web/src/ui/PageHeader.module.css`**

```css
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-5);
}
.title {
  font-size: var(--text-2xl);
}
```

- [ ] **Step 9: Verificar**

Run dentro de `web/`: `npx tsc --noEmit`
Expected: sin errores.

- [ ] **Step 10: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): primitivos Field, Badge, Skeleton, PageHeader"
```

---

### Task 4: EmptyState (TDD)

**Files:**
- Test: `web/src/ui/EmptyState.test.tsx`
- Create: `web/src/ui/EmptyState.tsx`, `EmptyState.module.css`

- [ ] **Step 1: Escribir la prueba que falla — `web/src/ui/EmptyState.test.tsx`**

```tsx
import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { EmptyState } from "./EmptyState";

describe("EmptyState", () => {
  it("renderiza título, mensaje y acción", () => {
    render(
      <EmptyState
        title="Sin datos"
        message="No hay nada todavía"
        action={<button>Crear</button>}
      />,
    );
    expect(screen.getByText("Sin datos")).toBeInTheDocument();
    expect(screen.getByText("No hay nada todavía")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Crear" })).toBeInTheDocument();
  });
});
```

- [ ] **Step 2: Ejecutar y ver que falla**

Run: `npx vitest run src/ui/EmptyState.test.tsx`
Expected: FALLA (no existe `EmptyState`).

- [ ] **Step 3: Crear `web/src/ui/EmptyState.tsx`**

```tsx
import type { ReactNode } from "react";
import styles from "./EmptyState.module.css";

type Props = { title: string; message?: string; action?: ReactNode };

export function EmptyState({ title, message, action }: Props) {
  return (
    <div className={styles.empty}>
      <h3 className={styles.title}>{title}</h3>
      {message && <p className={styles.message}>{message}</p>}
      {action && <div className={styles.action}>{action}</div>}
    </div>
  );
}
```

- [ ] **Step 4: Crear `web/src/ui/EmptyState.module.css`**

```css
.empty {
  text-align: center;
  padding: var(--space-7) var(--space-4);
  color: var(--color-text-muted);
  background: var(--color-surface);
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-lg);
}
.title {
  font-size: var(--text-lg);
  color: var(--color-text);
  margin-bottom: var(--space-2);
}
.message {
  margin-bottom: var(--space-4);
}
.action {
  display: flex;
  justify-content: center;
}
```

- [ ] **Step 5: Ejecutar y ver que pasa**

Run: `npx vitest run src/ui/EmptyState.test.tsx`
Expected: PASA (1 test).

- [ ] **Step 6: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): componente EmptyState"
```

---

### Task 5: Modal (TDD)

**Files:**
- Test: `web/src/ui/Modal.test.tsx`
- Create: `web/src/ui/Modal.tsx`, `Modal.module.css`

- [ ] **Step 1: Escribir la prueba que falla — `web/src/ui/Modal.test.tsx`**

```tsx
import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { Modal } from "./Modal";

describe("Modal", () => {
  it("no renderiza nada si open=false", () => {
    render(
      <Modal open={false} onClose={() => {}} title="T">
        contenido
      </Modal>,
    );
    expect(screen.queryByText("contenido")).not.toBeInTheDocument();
  });

  it("llama a onClose al pulsar Escape", () => {
    const onClose = vi.fn();
    render(
      <Modal open onClose={onClose} title="T">
        contenido
      </Modal>,
    );
    fireEvent.keyDown(document, { key: "Escape" });
    expect(onClose).toHaveBeenCalled();
  });

  it("llama a onClose al hacer clic en el overlay", () => {
    const onClose = vi.fn();
    render(
      <Modal open onClose={onClose} title="T">
        contenido
      </Modal>,
    );
    fireEvent.click(screen.getByTestId("modal-overlay"));
    expect(onClose).toHaveBeenCalled();
  });
});
```

- [ ] **Step 2: Ejecutar y ver que falla**

Run: `npx vitest run src/ui/Modal.test.tsx`
Expected: FALLA (no existe `Modal`).

- [ ] **Step 3: Crear `web/src/ui/Modal.tsx`**

```tsx
import { useEffect, useRef } from "react";
import type { ReactNode } from "react";
import styles from "./Modal.module.css";

type Props = {
  open: boolean;
  onClose: () => void;
  title: string;
  children: ReactNode;
};

export function Modal({ open, onClose, title, children }: Props) {
  const dialogRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") onClose();
    }
    document.addEventListener("keydown", onKey);
    dialogRef.current?.focus();
    return () => document.removeEventListener("keydown", onKey);
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div
      className={styles.overlay}
      data-testid="modal-overlay"
      onClick={onClose}
    >
      <div
        ref={dialogRef}
        className={styles.dialog}
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-title"
        tabIndex={-1}
        onClick={(e) => e.stopPropagation()}
      >
        <h2 id="modal-title" className={styles.title}>
          {title}
        </h2>
        <div>{children}</div>
      </div>
    </div>
  );
}
```

- [ ] **Step 4: Crear `web/src/ui/Modal.module.css`**

```css
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(17, 24, 39, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-4);
  z-index: 50;
}
.dialog {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  width: 100%;
  max-width: 480px;
  max-height: 90vh;
  overflow-y: auto;
  padding: var(--space-5);
}
.dialog:focus {
  outline: none;
}
.title {
  font-size: var(--text-xl);
  margin-bottom: var(--space-4);
}
```

- [ ] **Step 5: Ejecutar y ver que pasa**

Run: `npx vitest run src/ui/Modal.test.tsx`
Expected: PASA (3 tests).

- [ ] **Step 6: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): componente Modal accesible"
```

---

### Task 6: ConfirmDialog (TDD)

**Files:**
- Test: `web/src/ui/ConfirmDialog.test.tsx`
- Create: `web/src/ui/ConfirmDialog.tsx`, `ConfirmDialog.module.css`

- [ ] **Step 1: Escribir la prueba que falla — `web/src/ui/ConfirmDialog.test.tsx`**

```tsx
import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { ConfirmDialog } from "./ConfirmDialog";

describe("ConfirmDialog", () => {
  it("confirmar llama a onConfirm y cancelar a onCancel", () => {
    const onConfirm = vi.fn();
    const onCancel = vi.fn();
    render(
      <ConfirmDialog
        open
        title="Borrar"
        message="¿Seguro?"
        confirmLabel="Borrar"
        onConfirm={onConfirm}
        onCancel={onCancel}
      />,
    );
    fireEvent.click(screen.getByRole("button", { name: "Borrar" }));
    expect(onConfirm).toHaveBeenCalled();
    fireEvent.click(screen.getByRole("button", { name: "Cancelar" }));
    expect(onCancel).toHaveBeenCalled();
  });
});
```

- [ ] **Step 2: Ejecutar y ver que falla**

Run: `npx vitest run src/ui/ConfirmDialog.test.tsx`
Expected: FALLA (no existe `ConfirmDialog`).

- [ ] **Step 3: Crear `web/src/ui/ConfirmDialog.tsx`**

```tsx
import { Modal } from "./Modal";
import { Button } from "./Button";
import styles from "./ConfirmDialog.module.css";

type Props = {
  open: boolean;
  title: string;
  message: string;
  confirmLabel?: string;
  onConfirm: () => void;
  onCancel: () => void;
};

export function ConfirmDialog({
  open,
  title,
  message,
  confirmLabel = "Confirmar",
  onConfirm,
  onCancel,
}: Props) {
  return (
    <Modal open={open} onClose={onCancel} title={title}>
      <p className={styles.message}>{message}</p>
      <div className={styles.actions}>
        <Button variant="secondary" onClick={onCancel}>
          Cancelar
        </Button>
        <Button variant="danger" onClick={onConfirm}>
          {confirmLabel}
        </Button>
      </div>
    </Modal>
  );
}
```

- [ ] **Step 4: Crear `web/src/ui/ConfirmDialog.module.css`**

```css
.message {
  color: var(--color-text-muted);
  margin-bottom: var(--space-5);
}
.actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-2);
}
```

- [ ] **Step 5: Ejecutar y ver que pasa**

Run: `npx vitest run src/ui/ConfirmDialog.test.tsx`
Expected: PASA (1 test).

- [ ] **Step 6: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): componente ConfirmDialog"
```

---

### Task 7: Toast (TDD) + integración en main.tsx

**Files:**
- Test: `web/src/ui/Toast.test.tsx`
- Create: `web/src/ui/Toast.tsx`, `Toast.module.css`
- Modify: `web/src/main.tsx`

- [ ] **Step 1: Escribir la prueba que falla — `web/src/ui/Toast.test.tsx`**

```tsx
import { render, screen, act } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { ToastProvider, useToast } from "./Toast";

function Trigger() {
  const toast = useToast();
  return <button onClick={() => toast.success("Guardado")}>lanzar</button>;
}

describe("Toast", () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => vi.useRealTimers());

  it("muestra el toast y se auto-cierra", () => {
    render(
      <ToastProvider>
        <Trigger />
      </ToastProvider>,
    );
    act(() => {
      screen.getByRole("button", { name: "lanzar" }).click();
    });
    expect(screen.getByText("Guardado")).toBeInTheDocument();
    act(() => {
      vi.advanceTimersByTime(4000);
    });
    expect(screen.queryByText("Guardado")).not.toBeInTheDocument();
  });
});
```

- [ ] **Step 2: Ejecutar y ver que falla**

Run: `npx vitest run src/ui/Toast.test.tsx`
Expected: FALLA (no existe `Toast`).

- [ ] **Step 3: Crear `web/src/ui/Toast.tsx`**

```tsx
import { createContext, useCallback, useContext, useMemo, useRef, useState } from "react";
import type { ReactNode } from "react";
import styles from "./Toast.module.css";

type ToastTone = "success" | "error";
type ToastItem = { id: number; tone: ToastTone; message: string };

type ToastApi = {
  success: (message: string) => void;
  error: (message: string) => void;
};

const ToastContext = createContext<ToastApi | null>(null);

const TOAST_DURATION_MS = 4000;

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([]);
  const nextId = useRef(0);

  const remove = useCallback((id: number) => {
    setToasts((list) => list.filter((t) => t.id !== id));
  }, []);

  const push = useCallback(
    (tone: ToastTone, message: string) => {
      const id = ++nextId.current;
      setToasts((list) => [...list, { id, tone, message }]);
      setTimeout(() => remove(id), TOAST_DURATION_MS);
    },
    [remove],
  );

  const api = useMemo<ToastApi>(
    () => ({
      success: (m: string) => push("success", m),
      error: (m: string) => push("error", m),
    }),
    [push],
  );

  return (
    <ToastContext.Provider value={api}>
      {children}
      <div className={styles.container} aria-live="polite">
        {toasts.map((t) => (
          <div key={t.id} className={`${styles.toast} ${styles[t.tone]}`} role="status">
            {t.message}
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export function useToast(): ToastApi {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToast debe usarse dentro de <ToastProvider>");
  return ctx;
}
```

- [ ] **Step 4: Crear `web/src/ui/Toast.module.css`**

```css
.container {
  position: fixed;
  bottom: var(--space-5);
  right: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  z-index: 100;
}
.toast {
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
  color: #fff;
  font-weight: 600;
  font-size: var(--text-sm);
  min-width: 220px;
}
.success {
  background: var(--color-success);
}
.error {
  background: var(--color-danger);
}
```

- [ ] **Step 5: Ejecutar y ver que pasa**

Run: `npx vitest run src/ui/Toast.test.tsx`
Expected: PASA (1 test).

- [ ] **Step 6: Envolver la app — reemplazar TODO `web/src/main.tsx`**

```tsx
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { ToastProvider } from "./ui/Toast";
import App from "./App";
import "./styles/tokens.css";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <ToastProvider>
        <App />
      </ToastProvider>
    </BrowserRouter>
  </StrictMode>,
);
```

- [ ] **Step 7: Verificar tipos**

Run dentro de `web/`: `npx tsc --noEmit`
Expected: sin errores.

- [ ] **Step 8: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): sistema de notificaciones (Toast) e integración en la app"
```

---

### Task 8: Reestilizar DataTable

**Files:**
- Modify: `web/src/components/DataTable.tsx`
- Modify: `web/src/components/DataTable.module.css`

- [ ] **Step 1: Reemplazar TODO `web/src/components/DataTable.tsx`**

```tsx
import type { ReactNode } from "react";
import { Button } from "../ui/Button";
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
              <td>
                <div className={styles.actions}>
                  {onEdit && (
                    <Button variant="ghost" size="sm" onClick={() => onEdit(row)}>
                      Editar
                    </Button>
                  )}
                  {onDelete && (
                    <Button variant="danger" size="sm" onClick={() => onDelete(row)}>
                      Borrar
                    </Button>
                  )}
                </div>
              </td>
            )}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
```

- [ ] **Step 2: Reemplazar TODO `web/src/components/DataTable.module.css`**

```css
.table {
  width: 100%;
  border-collapse: collapse;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}
.table th,
.table td {
  padding: var(--space-3) var(--space-4);
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}
.table th {
  background: var(--color-bg);
  font-size: var(--text-sm);
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.03em;
}
.table tbody tr:last-child td {
  border-bottom: none;
}
.table tbody tr:hover {
  background: var(--color-bg);
}
.actions {
  display: flex;
  gap: var(--space-2);
}
```

- [ ] **Step 3: Verificar (tipos + prueba existente de DataTable)**

Run dentro de `web/`:
```bash
npx tsc --noEmit && npx vitest run src/components/DataTable.test.tsx
```
Expected: tsc sin errores; la prueba existente de DataTable sigue PASANDO (2 tests).

- [ ] **Step 4: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): reestilizar DataTable con tokens y Button"
```

---

### Task 9: Actualizar formularios (ResourceForm y SceneForm)

**Files:**
- Modify: `web/src/components/ResourceForm.tsx`
- Create: `web/src/components/ResourceForm.module.css`
- Modify: `web/src/components/SceneForm.tsx`
- Create: `web/src/components/SceneForm.module.css`

- [ ] **Step 1: Reemplazar TODO `web/src/components/ResourceForm.tsx`**

```tsx
import { useState } from "react";
import type { FormEvent } from "react";
import { Field } from "../ui/Field";
import { Button } from "../ui/Button";
import styles from "./ResourceForm.module.css";

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
      {shownError && (
        <p role="alert" className={styles.alert}>
          {shownError}
        </p>
      )}
      <Field label="Título" htmlFor="rf-title">
        <input id="rf-title" value={title} onChange={(e) => setTitle(e.target.value)} />
      </Field>
      <Field label="Texto" htmlFor="rf-text">
        <textarea id="rf-text" value={text} onChange={(e) => setText(e.target.value)} />
      </Field>
      <div className={styles.actions}>
        <Button type="button" variant="secondary" onClick={onCancel}>
          Cancelar
        </Button>
        <Button type="submit" disabled={submitting}>
          Guardar
        </Button>
      </div>
    </form>
  );
}
```

- [ ] **Step 2: Crear `web/src/components/ResourceForm.module.css`**

```css
.alert {
  color: var(--color-danger);
  margin-bottom: var(--space-3);
}
.actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-2);
  margin-top: var(--space-4);
}
```

- [ ] **Step 3: Reemplazar TODO `web/src/components/SceneForm.tsx`**

```tsx
import { useState } from "react";
import type { FormEvent } from "react";
import type { Character, Location } from "../types";
import { Field } from "../ui/Field";
import { Button } from "../ui/Button";
import styles from "./SceneForm.module.css";

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
  const [characterIds, setCharacterIds] = useState<number[]>(initial?.characterIds ?? []);
  const [locationIds, setLocationIds] = useState<number[]>(initial?.locationIds ?? []);
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
      {shownError && (
        <p role="alert" className={styles.alert}>
          {shownError}
        </p>
      )}
      <Field label="Título" htmlFor="sf-title">
        <input id="sf-title" value={title} onChange={(e) => setTitle(e.target.value)} />
      </Field>
      <Field label="Texto" htmlFor="sf-text">
        <textarea id="sf-text" value={text} onChange={(e) => setText(e.target.value)} />
      </Field>
      <div className={styles.row}>
        <Field label="Inicio (timeline)" htmlFor="sf-start">
          <input
            id="sf-start"
            type="number"
            value={start}
            onChange={(e) => setStart(e.target.value)}
          />
        </Field>
        <Field label="Fin (timeline)" htmlFor="sf-end">
          <input
            id="sf-end"
            type="number"
            value={end}
            onChange={(e) => setEnd(e.target.value)}
          />
        </Field>
      </div>
      <Field label="Personajes" htmlFor="sf-characters">
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
      </Field>
      <Field label="Lugares" htmlFor="sf-locations">
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
      </Field>
      <div className={styles.actions}>
        <Button type="button" variant="secondary" onClick={onCancel}>
          Cancelar
        </Button>
        <Button type="submit" disabled={submitting}>
          Guardar
        </Button>
      </div>
    </form>
  );
}
```

- [ ] **Step 4: Crear `web/src/components/SceneForm.module.css`**

```css
.alert {
  color: var(--color-danger);
  margin-bottom: var(--space-3);
}
.row {
  display: flex;
  gap: var(--space-4);
}
.row > * {
  flex: 1;
}
.actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-2);
  margin-top: var(--space-4);
}
```

- [ ] **Step 5: Verificar (tipos + prueba existente de ResourceForm)**

Run dentro de `web/`:
```bash
npx tsc --noEmit && npx vitest run src/components/ResourceForm.test.tsx
```
Expected: tsc sin errores; la prueba existente de ResourceForm sigue PASANDO (2 tests).

- [ ] **Step 6: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): formularios con Field y Button"
```

---

### Task 10: Refactor de ResourcePage (modal, confirm, toast, skeleton, empty)

**Files:**
- Modify: `web/src/components/ResourcePage.tsx`

- [ ] **Step 1: Reemplazar TODO `web/src/components/ResourcePage.tsx`**

```tsx
import { useState } from "react";
import { useList } from "../hooks/useList";
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
import { useToast } from "../ui/Toast";

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
  const { data, loading, error, reload } = useList(list);
  const toast = useToast();
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
      toast.error(e instanceof Error ? e.message : "Error desconocido");
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
      <PageHeader title={heading} action={<Button onClick={openNew}>Nuevo</Button>} />

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
          title="Aún no hay nada aquí"
          message="Crea el primer elemento para empezar."
          action={<Button onClick={openNew}>Nuevo</Button>}
        />
      )}
      {!loading && !error && data.length > 0 && (
        <DataTable
          columns={columns}
          rows={data}
          onEdit={(row) => {
            setFormError(null);
            setEditing(row);
          }}
          onDelete={(row) => setDeleting(row)}
        />
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
Expected: sin errores.

- [ ] **Step 3: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): ResourcePage con modal, confirmación, toasts, skeleton y empty state"
```

---

### Task 11: Refactor de ScenesPage (+ badges)

**Files:**
- Modify: `web/src/pages/ScenesPage.tsx`

- [ ] **Step 1: Reemplazar TODO `web/src/pages/ScenesPage.tsx`**

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
import { Modal } from "../ui/Modal";
import { ConfirmDialog } from "../ui/ConfirmDialog";
import { Button } from "../ui/Button";
import { PageHeader } from "../ui/PageHeader";
import { SkeletonRows } from "../ui/Skeleton";
import { EmptyState } from "../ui/EmptyState";
import { Badge } from "../ui/Badge";
import { useToast } from "../ui/Toast";

type Editing = null | "new" | Scene;

export function ScenesPage() {
  const scenes = useList(listScenes);
  const characters = useList(listCharacters);
  const locations = useList(listLocations);
  const toast = useToast();

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
        s.locations.length
          ? s.locations.map((l) => <Badge key={l.id}>{l.title}</Badge>)
          : "—",
    },
  ];

  function openNew() {
    setFormError(null);
    setEditing("new");
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
      toast.error(e instanceof Error ? e.message : "Error desconocido");
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
      <PageHeader title="Escenas" action={<Button onClick={openNew}>Nueva</Button>} />

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
          title="Aún no hay escenas"
          message="Crea la primera escena para empezar."
          action={<Button onClick={openNew}>Nueva</Button>}
        />
      )}
      {!scenes.loading && !scenes.error && scenes.data.length > 0 && (
        <DataTable
          columns={columns}
          rows={scenes.data}
          onEdit={(row) => {
            setFormError(null);
            setEditing(row);
          }}
          onDelete={(row) => setDeleting(row)}
        />
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

- [ ] **Step 2: Verificar tipos**

Run dentro de `web/`: `npx tsc --noEmit`
Expected: sin errores.

- [ ] **Step 3: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): ScenesPage con modal, confirmación, toasts, skeleton, empty y badges"
```

---

### Task 12: Rediseño del Layout (responsive + iconos)

**Files:**
- Modify: `web/src/components/Layout.tsx`
- Modify: `web/src/components/Layout.module.css`

- [ ] **Step 1: Reemplazar TODO `web/src/components/Layout.tsx`**

```tsx
import { NavLink, Outlet } from "react-router-dom";
import styles from "./Layout.module.css";

function UserIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <circle cx="12" cy="8" r="4" />
      <path d="M4 20c0-4 4-6 8-6s8 2 8 6" />
    </svg>
  );
}
function PinIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M12 21s-7-6-7-11a7 7 0 0 1 14 0c0 5-7 11-7 11z" />
      <circle cx="12" cy="10" r="2.5" />
    </svg>
  );
}
function FilmIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <rect x="3" y="4" width="18" height="16" rx="2" />
      <path d="M3 9h18M3 15h18M9 4v16M15 4v16" />
    </svg>
  );
}

const links = [
  { to: "/characters", label: "Personajes", icon: <UserIcon /> },
  { to: "/locations", label: "Lugares", icon: <PinIcon /> },
  { to: "/scenes", label: "Escenas", icon: <FilmIcon /> },
];

export function Layout() {
  return (
    <div className={styles.shell}>
      <nav className={styles.sidebar}>
        <div className={styles.brand}>Story Admin</div>
        <ul className={styles.nav}>
          {links.map((l) => (
            <li key={l.to}>
              <NavLink
                to={l.to}
                className={({ isActive }) =>
                  `${styles.link} ${isActive ? styles.active : ""}`
                }
              >
                <span className={styles.icon}>{l.icon}</span>
                {l.label}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
      <main className={styles.content}>
        <div className={styles.container}>
          <Outlet />
        </div>
      </main>
    </div>
  );
}
```

- [ ] **Step 2: Reemplazar TODO `web/src/components/Layout.module.css`**

```css
.shell {
  display: flex;
  min-height: 100vh;
}
.sidebar {
  width: 220px;
  flex-shrink: 0;
  background: var(--color-surface);
  border-right: 1px solid var(--color-border);
  padding: var(--space-5) var(--space-3);
}
.brand {
  font-weight: 700;
  font-size: var(--text-lg);
  padding: 0 var(--space-3) var(--space-5);
}
.nav {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}
.link {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
  text-decoration: none;
  font-weight: 600;
}
.link:hover {
  background: var(--color-bg);
  color: var(--color-text);
}
.active {
  background: var(--color-primary-soft);
  color: var(--color-primary);
}
.icon {
  display: inline-flex;
}
.content {
  flex: 1;
  padding: var(--space-6);
  overflow-x: auto;
}
.container {
  max-width: 960px;
  margin: 0 auto;
}

@media (max-width: 720px) {
  .shell {
    flex-direction: column;
  }
  .sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid var(--color-border);
    padding: var(--space-3);
  }
  .nav {
    flex-direction: row;
    flex-wrap: wrap;
  }
  .brand {
    padding: var(--space-2);
  }
  .content {
    padding: var(--space-4);
  }
}
```

- [ ] **Step 3: Verificar tipos**

Run dentro de `web/`: `npx tsc --noEmit`
Expected: sin errores.

- [ ] **Step 4: Commit**

```bash
cd .. && git add -A && git commit -m "feat(web): rediseño del Layout (sidebar con iconos + responsive)"
```

---

### Task 13: Verificación final

**Files:** ninguno (solo verificación)

- [ ] **Step 1: Ejecutar todas las pruebas**

Run dentro de `web/`: `npm test`
Expected: PASAN todas. Conteo esperado: las 6 originales (client x2, DataTable x2, ResourceForm x2) + Button x2 + EmptyState x1 + Modal x3 + ConfirmDialog x1 + Toast x1 = **17 pruebas**.

- [ ] **Step 2: Comprobación de tipos y build**

Run dentro de `web/`:
```bash
npx tsc --noEmit && npm run build
```
Expected: sin errores; build correcto.

- [ ] **Step 3: Prueba manual extremo a extremo**

1. Arranca la API: `go run ./cmd/server` (raíz del repo; requiere MySQL).
2. `cd web && npm run dev`; abre la URL local.
3. Verifica visualmente:
   - El layout tiene barra lateral con iconos y estado activo; aspecto limpio.
   - "Nuevo" abre un **modal** con el formulario (no reemplaza la página).
   - Crear/editar muestra un **toast** de éxito; un título duplicado muestra toast de error y el error en el formulario.
   - "Borrar" abre un **ConfirmDialog**; al confirmar, toast de éxito y la fila desaparece.
   - Al cargar, se ven **skeletons**; si una lista está vacía, se ve el **EmptyState**.
   - En Escenas, personajes y lugares se muestran como **badges**.
   - Reduce el ancho de la ventana: el layout es usable (sidebar arriba).

Expected: todo funciona y se ve profesional.

- [ ] **Step 4: Commit final (si hubo ajustes)**

```bash
cd .. && git add -A && git commit -m "test(web): verificación final del rediseño" --allow-empty
```

---

## Notas de verificación (self-review del plan)

- **Cobertura del spec:** tokens (Task 1), componentes ui (Tasks 2-7), DataTable
  reestilizado (Task 8), formularios con Field/Button (Task 9), ResourcePage
  (Task 10), ScenesPage + badges (Task 11), Layout responsive (Task 12),
  verificación + manual (Task 13). Toasts integrados en main (Task 7).
- **Sin placeholders:** todo el código está completo.
- **Consistencia de tipos:** `Button` (variant/size), `Modal` (open/onClose/title),
  `ConfirmDialog` (open/title/message/confirmLabel/onConfirm/onCancel),
  `useToast` (success/error), `Field` (label/htmlFor/error), `SkeletonRows`
  (rows/cols), `Column<T>`, `ResourceFormValues`, `SceneFormValues` se usan con
  nombres y firmas idénticos en todas las tareas.
- **Pruebas existentes:** DataTable y ResourceForm conservan sus selectores
  (`getByText`, `getByLabelText`, `role="alert"`), así que sus pruebas siguen
  pasando tras el refactor.
```
