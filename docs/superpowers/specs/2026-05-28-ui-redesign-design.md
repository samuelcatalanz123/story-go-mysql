# Diseño: Rediseño profesional de la UI (estilo SaaS minimalista)

**Fecha:** 2026-05-28
**Estado:** Aprobado para escribir el plan de implementación
**Proyecto previo:** Panel de administración React + TS (ver
`2026-05-28-react-ui-admin-design.md`). Este rediseño parte de esa base.

## Objetivo

Transformar el panel de administración funcional en una pieza de portafolio
que impresione a un empleador: aspecto profesional "SaaS minimalista" e
interacciones pulidas, manteniendo CSS Modules + un sistema de tokens propio.

Objetivo del usuario: conseguir su primer trabajo de programador. El criterio
de éxito es que la app se vea y se sienta como un producto real a primera vista.

## Decisiones tomadas (brainstorming)

| Tema       | Decisión                                  |
| ---------- | ----------------------------------------- |
| Meta       | Pieza de portafolio para ser contratado   |
| Estilo     | SaaS minimalista y limpio (Linear/Stripe) |
| Método     | CSS Modules + tokens de diseño propios    |
| Alcance    | Visual + interacciones pulidas            |
| Acento     | Índigo                                     |

## Sistema de diseño (tokens)

Archivo `web/src/styles/tokens.css` con variables CSS en `:root`, importado una
vez globalmente. Grupos de tokens:

- **Color:** `--color-bg` (casi blanco), `--color-surface` (blanco),
  `--color-border` (gris suave), `--color-text` (gris muy oscuro),
  `--color-text-muted` (gris medio), `--color-primary` + `--color-primary-hover`
  (índigo), `--color-danger` + `--color-danger-hover` (rojo),
  `--color-success` (verde).
- **Espaciado:** escala `--space-1`..`--space-8` (4px base: 4, 8, 12, 16, 24,
  32, 48, 64).
- **Radios:** `--radius-sm`, `--radius-md`, `--radius-lg`.
- **Sombras:** `--shadow-sm`, `--shadow-md` (sutiles).
- **Tipografía:** `--font-sans` (stack del sistema con Inter como primera
  opción si está disponible), tamaños `--text-sm/base/lg/xl/2xl`, pesos.

Un `reset` mínimo (box-sizing, márgenes) vivirá en un global pequeño junto a
los tokens.

## Componentes de UI (carpeta nueva `web/src/ui/`)

Cada uno con su `.tsx` + `.module.css` + (cuando aplique) test. Una sola
responsabilidad por componente.

| Componente      | Responsabilidad                                                | Interfaz pública (props clave)                              |
| --------------- | -------------------------------------------------------------- | ----------------------------------------------------------- |
| `Button`        | Botón con variantes/tamaños                                    | `variant: "primary"\|"secondary"\|"danger"\|"ghost"`, `size`, `disabled`, estándar de `<button>` |
| `Modal`         | Ventana flotante accesible (overlay, foco, cierre con Esc)     | `open: boolean`, `onClose: () => void`, `title`, `children` |
| `ConfirmDialog` | Confirmación de acción (sobre `Modal`)                         | `open`, `title`, `message`, `confirmLabel`, `onConfirm`, `onCancel` |
| `ToastProvider` | Contexto que gestiona la cola de toasts                        | envuelve la app; expone vía `useToast()`                    |
| `useToast`      | Hook para lanzar notificaciones                                | `toast.success(msg)`, `toast.error(msg)`                    |
| `Skeleton`      | Bloque gris animado de carga                                   | `width`, `height` (o filas para tabla)                      |
| `EmptyState`    | Mensaje + acción cuando no hay datos                           | `title`, `message`, `action?`                               |
| `Field`         | Etiqueta + control de formulario coherente                     | `label`, `htmlFor`, `error?`, `children`                    |
| `Badge`         | Etiqueta tipo "chip"                                           | `children`, `tone?`                                         |

Notas de accesibilidad del `Modal`: `role="dialog"`, `aria-modal="true"`,
`aria-labelledby` apuntando al título, foco inicial dentro del modal y cierre
con `Esc`. (Trampa de foco completa queda como mejora futura; foco inicial +
Esc + overlay clic son suficientes para este alcance.)

## Layout

`Layout` rediseñado:
- Barra lateral (`--color-surface`, borde derecho): marca "Story Admin",
  navegación con icono + texto y estado activo claro (fondo y color de acento).
  Los iconos serán SVG inline simples (sin librería de iconos, para no añadir
  dependencias).
- Cabecera de página: cada página renderiza su título y el botón de acción
  principal alineado a la derecha mediante un pequeño componente `PageHeader`
  (en `ui/`).
- Contenido: ancho máximo (~960px), centrado, con espaciado por tokens.
- Responsive: por debajo de ~720px la barra lateral pasa a horizontal arriba
  (CSS con media query; sin lógica JS de menú hamburguesa en este alcance).

## Cambios en las pantallas existentes

- **ResourcePage** (personajes/lugares): el formulario deja de reemplazar la
  página y se abre en un `Modal`. El borrado usa `ConfirmDialog`. Éxito/error
  de crear/editar/borrar disparan `toast`. La carga muestra `Skeleton` en la
  tabla; lista vacía muestra `EmptyState`. La tabla se reestiliza vía tokens.
- **ScenesPage**: mismo tratamiento (modal, confirm, toasts, skeleton, empty).
  Personajes y lugares de cada escena se muestran como `Badge`.
- **DataTable**: reestilizado con tokens; los botones de fila usan `Button`
  (`ghost`/`danger`). Mantiene su API genérica actual (`columns`, `rows`,
  `onEdit`, `onDelete`).
- **ResourceForm** y **SceneForm**: los controles pasan a usar `Field` y los
  botones a `Button`; la lógica de validación actual se conserva. El error de
  envío (prop `error`) se sigue mostrando, pero el éxito lo notifica la página
  vía toast.

## Flujo de datos (sin cambios estructurales)

La capa `api/` y el hook `useList` no cambian. El cambio es de presentación e
interacción. Las páginas ganan estado de UI nuevo: `modalOpen`, `confirming`
(elemento pendiente de borrar) y acceso a `useToast()`.

## Manejo de errores (mejorado)

- Errores de carga de lista: banner/EmptyState con mensaje + botón "Reintentar"
  (llama a `reload`).
- Errores de envío de formulario: se muestran dentro del modal (junto al
  formulario) y/o como toast de error.
- Errores de borrado: toast de error (antes era `alert`).

## Pruebas

Nuevas (Vitest + Testing Library):
- `Button`: aplica la clase de la variante indicada y respeta `disabled`.
- `Modal`: no renderiza nada si `open=false`; al estar abierto, `Esc` y clic en
  overlay llaman a `onClose`.
- `ConfirmDialog`: el botón de confirmar llama a `onConfirm`; cancelar llama a
  `onCancel`.
- `ToastProvider`/`useToast`: lanzar un toast lo muestra; se auto-cierra tras el
  tiempo configurado (usar timers falsos de Vitest).
- `EmptyState`: renderiza título, mensaje y acción.

Actualizadas:
- Las pruebas de `ResourceForm` siguen válidas (el formulario en sí no cambia su
  lógica). Si algún selector cambia por usar `Field`, se ajusta el test.

Toda la suite debe seguir en verde, además de `tsc --noEmit` y `npm run build`.

## Organización de archivos (nuevos/modificados)

```
web/src/
  styles/
    tokens.css            (nuevo) variables + reset mínimo
  ui/                      (nuevo)
    Button.tsx / .module.css / Button.test.tsx
    Modal.tsx / .module.css / Modal.test.tsx
    ConfirmDialog.tsx / .module.css / ConfirmDialog.test.tsx
    Toast.tsx / .module.css / Toast.test.tsx   (ToastProvider + useToast)
    Skeleton.tsx / .module.css
    EmptyState.tsx / .module.css / EmptyState.test.tsx
    Field.tsx / .module.css
    Badge.tsx / .module.css
    PageHeader.tsx / .module.css
  components/
    Layout.tsx / Layout.module.css   (rediseño)
    DataTable.tsx / DataTable.module.css   (reestilizado)
    ResourcePage.tsx   (modal + confirm + toast + skeleton + empty)
    ResourceForm.tsx   (usa Field + Button)
    SceneForm.tsx      (usa Field + Button)
  pages/
    ScenesPage.tsx     (modal + confirm + toast + skeleton + empty + badges)
  main.tsx             (envolver App en ToastProvider; importar tokens.css)
```

## Fuera de alcance (YAGNI)

- Autenticación / cuentas de usuario.
- Búsqueda, filtros y paginación.
- Modo oscuro (los tokens dejarán la puerta abierta, pero no se implementa).
- Despliegue a producción.
- Librerías de UI o de iconos (todo a mano con CSS Modules + SVG inline).
- Trampa de foco completa en el modal (foco inicial + Esc + overlay bastan).

## Criterios de éxito

1. La app presenta un aspecto coherente y profesional basado en tokens.
2. Crear/editar ocurre en un modal; borrar pide confirmación con `ConfirmDialog`.
3. Éxitos y errores se comunican con toasts (sin `alert`/`confirm` nativos).
4. Las listas muestran skeleton al cargar y EmptyState cuando están vacías.
5. Las escenas muestran sus personajes y lugares como badges.
6. El layout es responsive (usable en móvil).
7. Todas las pruebas pasan, `tsc --noEmit` limpio y `npm run build` correcto.
